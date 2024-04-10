package cogs

import (
	"context"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/line"

	onMessageCreate "github.com/maguro-alternative/remake_bot/bot/cogs/on_message_create"
	"github.com/maguro-alternative/remake_bot/bot/config"

	"github.com/bwmarrin/discordgo"
)

func (h *CogHandler) OnMessageCreate(s *discordgo.Session, vs *discordgo.MessageCreate) {
	ctx := context.Background()
	repo := repository.NewRepository(h.DB)
	ffmpeg := onMessageCreate.NewFfmpeg(ctx)
	err := onMessageCreateFunc(ctx, h.client, repo, *ffmpeg, s, vs)
	if err != nil {
		slog.ErrorContext(ctx, "OnMessageCreate Error", "Error:", err.Error())
	}
}

func onMessageCreateFunc(
	ctx context.Context,
	client *http.Client,
	repo repository.RepositoryFunc,
	ffmpeg onMessageCreate.FfmpegInterface,
	s mock.Session,
	vs *discordgo.MessageCreate,
) error {
	var channel repository.LinePostDiscordChannel
	var lineMessageTypes []*line.LineMessageType
	var imageUrls []string
	var videoCount, voiceCount int

	sendTextBuilder := strings.Builder{}

	channel, err := repo.GetLinePostDiscordChannel(ctx, vs.ChannelID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		slog.ErrorContext(ctx, "line_post_discord_channelの取得に失敗しました", "エラー:", err.Error())
		return err
	} else if err != nil {
		err = repo.InsertLinePostDiscordChannel(ctx, vs.ChannelID, vs.GuildID)
		if err != nil {
			slog.ErrorContext(ctx, "line_post_discord_channelの登録に失敗しました", "エラー:", err.Error())
			return err
		}
		channel = repository.LinePostDiscordChannel{
			Ng:         false,
			BotMessage: false,
		}
	}
	ngTypes, err := repo.GetLineNgDiscordMessageType(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_message_typeの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	ngDiscordIDs, err := repo.GetLineNgDiscordID(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_message_typeの登録に失敗しました", "エラー:", err.Error())
		return err
	}
	// メッセージの種類がNGの場合は処理を終了
	for _, ngType := range ngTypes {
		if vs.Message.Type == discordgo.MessageType(ngType) {
			slog.InfoContext(ctx, "NG Type")
			return err
		}
	}
	// メッセージの送信者がNGの場合は処理を終了
	for _, ngDiscordID := range ngDiscordIDs {
		if ngDiscordID.IDType == "user" && vs.Author.ID == ngDiscordID.ID {
			slog.InfoContext(ctx, "NG User")
			return err
		}
		if ngDiscordID.IDType == "role" {
			for _, role := range vs.Member.Roles {
				if role == ngDiscordID.ID {
					slog.InfoContext(ctx, "NG Role")
					return err
				}
			}
		}
	}
	// チャンネルがNGの場合、またはBotメッセージでない場合は処理を終了
	if channel.Ng || (!channel.BotMessage && vs.Author.Bot) {
		slog.InfoContext(ctx, "NG Channel or Bot Message")
		return err
	}
	lineBotApi, err := repo.GetLineBotNotClient(ctx, vs.GuildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_botの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	lineBotIv, err := repo.GetLineBotIvNotClient(ctx, vs.GuildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_bot_ivの取得に失敗しました", "エラー:", err.Error())
		return err
	}
	var lineBotDecrypt onMessageCreate.LineBotDecrypt
	// 暗号化キーのバイトへの変換
	keyBytes, err := hex.DecodeString(config.PrivateKey())
	if err != nil {
		slog.ErrorContext(ctx, "暗号化キーのバイト変換に失敗しました", "エラー:", err.Error())
		return err
	}

	lineNotifyTokenByte, err := crypto.Decrypt(lineBotApi.LineNotifyToken[0], keyBytes, lineBotIv.LineNotifyTokenIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_notify_tokenの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineBotTokenByte, err := crypto.Decrypt(lineBotApi.LineBotToken[0], keyBytes, lineBotIv.LineBotTokenIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_bot_tokenの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineGroupIDByte, err := crypto.Decrypt(lineBotApi.LineGroupID[0], keyBytes, lineBotIv.LineGroupIDIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_group_idの復号化に失敗しました", "エラー:", err.Error())
		return err
	}
	lineBotDecrypt.LineNotifyToken = string(lineNotifyTokenByte)
	lineBotDecrypt.LineBotToken = string(lineBotTokenByte)
	lineBotDecrypt.LineGroupID = string(lineGroupIDByte)
	lineBotDecrypt.DefaultChannelID = lineBotApi.DefaultChannelID
	lineBotDecrypt.DebugMode = lineBotApi.DebugMode

	lineRequ := line.NewLineRequest(
		*client,
		lineBotDecrypt.LineNotifyToken,
		lineBotDecrypt.LineBotToken,
		lineBotDecrypt.LineGroupID,
	)

	// メッセージの種類によって処理を分岐
	switch vs.Message.Type {
	case discordgo.MessageTypeUserPremiumGuildSubscription:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "がサーバーブーストしました。")
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierOne:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "がサーバーブーストし、レベル1になりました！！！！！！！！")
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierTwo:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "がサーバーブーストし、レベル2になりました！！！！！！！！")
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierThree:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "がサーバーブーストし、レベル3になりました！！！！！！！！")
	case discordgo.MessageTypeGuildMemberJoin:
		sendTextBuilder.WriteString(vs.Message.Author.Username + "が参加しました。")
	default:
		st, err := s.Channel(vs.ChannelID)
		if err != nil {
			slog.ErrorContext(ctx, "channel取得に失敗しました", "エラー:", err.Error())
			return err
		}
		sendTextBuilder.WriteString(st.Name + "にて、" + vs.Message.Author.Username)
	}

	// スタンプが送信されていた場合、画像URLを取得
	if vs.StickerItems != nil {
		for _, sticker := range vs.StickerItems {
			switch sticker.FormatType {
			case 1:
				imageUrls = append(imageUrls, "https://cdn.discordapp.com/stickers/"+sticker.ID+".png")
			case 2:
				imageUrls = append(imageUrls, "https://cdn.discordapp.com/stickers/"+sticker.ID+".apng")
			case 4:
				imageUrls = append(imageUrls, "https://cdn.discordapp.com/stickers/"+sticker.ID+".gif")
			}
		}
	}

	videoCount = 0
	voiceCount = 0
	// 添付ファイルが送信されていた場合、LINE用に変換
	for _, attachment := range vs.Message.Attachments {
		extension := filepath.Ext(attachment.Filename)
		fileNameNoExt := filepath.Base(attachment.Filename[:len(attachment.Filename)-len(extension)])
		switch extension {
		case ".png", ".jpg", ".jpeg", ".gif":
			imageUrls = append(imageUrls, attachment.URL)
		case ".mp4", ".mov", ".avi", ".wmv", ".flv", ".webm":
			st, err := s.Guild(vs.GuildID)
			if err != nil {
				slog.ErrorContext(ctx, "サーバーの取得に失敗しました", "エラー:", err.Error())
				return err
			}
			lineMessageType := lineRequ.NewLineVideoMessage(attachment.URL, st.IconURL("512"))
			lineMessageTypes = append(lineMessageTypes, lineMessageType)
			videoCount++
		case ".mp3", ".wav", ".ogg", ".m4a":
			tmpFile := os.TempDir() + "/" + attachment.Filename
			tmpFileNotExt := os.TempDir() + "/" + fileNameNoExt
			slog.InfoContext(ctx, "download:"+attachment.URL)
			err = downloadFile(client, tmpFile, attachment.URL)
			if err != nil {
				slog.ErrorContext(ctx, "ファイルのダウンロードに失敗しました", "エラー:", err.Error())
				return err
			}
			if extension != ".m4a" {
				slog.InfoContext(ctx, "m4a変換:"+tmpFile)
				slog.InfoContext(ctx, "ffmpeg:"+tmpFile)
				err = ffmpeg.ConversionAudioFile(tmpFile, tmpFileNotExt)
				if err != nil {
					slog.ErrorContext(ctx, "ffmpegの秒数カウントに失敗しました", "エラー:", err.Error())
					return err
				}
			}
			f, err := os.Open(tmpFileNotExt + ".m4a")
			if err != nil {
				slog.ErrorContext(ctx, "ファイルのオープンに失敗しました", "エラー:", err.Error())
				return err
			}
			defer f.Close()
			fileSendMesssage, err := s.ChannelFileSendWithMessage(
				vs.ChannelID,
				"m4aに変換します。",
				tmpFileNotExt+".m4a",
				f,
			)
			if err != nil {
				slog.ErrorContext(ctx, "Discordへのメッセージ送信に失敗しました", "エラー:", err.Error())
				return err
			}
			// 音声ファイルの秒数を取得
			slog.InfoContext(ctx, "秒数取得:"+tmpFileNotExt+".m4a")
			audioLen, err := ffmpeg.GetAudioFileSecond(tmpFile, tmpFileNotExt)
			if err != nil {
				slog.ErrorContext(ctx, "音声ファイルの秒数の抽出に失敗しました", "エラー:", err.Error())
				return err
			}
			slog.InfoContext(ctx, "秒数:"+strconv.FormatFloat(audioLen, 'f', -1, 64))
			audio := lineRequ.NewLineAudioMessage(
				fileSendMesssage.Attachments[0].URL,
				audioLen,
			)
			lineMessageTypes = append(lineMessageTypes, audio)
			voiceCount++
		default:
			slog.InfoContext(ctx, "未対応のファイル形式です。")
			sendTextBuilder.WriteString(attachment.URL + "\n")
		}
	}

	if len(imageUrls) > 0 {
		sendTextBuilder.WriteString(" 画像を" + strconv.Itoa(len(imageUrls)) + "枚、")
	}
	if videoCount > 0 {
		sendTextBuilder.WriteString(" 動画を" + strconv.Itoa(videoCount) + "個、")
	}
	if voiceCount > 0 {
		sendTextBuilder.WriteString(" 音声を" + strconv.Itoa(voiceCount) + "個、")
	}
	if len(imageUrls) > 0 || videoCount > 0 || voiceCount > 0 {
		sendTextBuilder.WriteString(" 送信しました。")
	}

	sendTextBuilder.WriteString("「 " + vs.Message.Content + " 」")

	// LINEに送信
	for _, url := range imageUrls {
		err = lineRequ.PushImageNotify(ctx, sendTextBuilder.String(), url)
		if err != nil {
			slog.ErrorContext(ctx, "LINE Notifyの画像送信に失敗しました", "エラー:", err.Error())
			return err
		}
	}
	// 動画、音声を送信
	if len(lineMessageTypes) > 0 {
		err = lineRequ.PushMessageBotInGroup(ctx, lineMessageTypes)
		if err != nil {
			slog.ErrorContext(ctx, "音声、動画の送信に失敗しました", "エラー:", err.Error())
			return err
		}
	}
	// 画像、動画、音声がない場合はテキストのみ送信
	if len(imageUrls) > 0 || videoCount > 0 || voiceCount > 0 {
		return err
	}
	err = lineRequ.PushMessageNotify(ctx, sendTextBuilder.String())
	if err != nil {
		slog.ErrorContext(ctx, "LINE Notifyの送信に失敗しました", "エラー:", err.Error())
		return err
	}
	return err
}

func downloadFile(client *http.Client, tmpFilePath, url string) error {
	f, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}
