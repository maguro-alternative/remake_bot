package cogs

import (
	"context"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/line"

	onMessageCreate "github.com/maguro-alternative/remake_bot/bot/cogs/on_message_create"
	"github.com/maguro-alternative/remake_bot/bot/config"

	"github.com/bwmarrin/discordgo"
)

type Repository interface {
	GetLinePostDiscordChannel(ctx context.Context, channelID string) (onMessageCreate.LinePostDiscordChannel, error)
	InsertLinePostDiscordChannel(ctx context.Context, channelID string, guildID string) error
	GetLineNgDiscordMessageType(ctx context.Context, channelID string) ([]int, error)
	GetLineNgDiscordID(ctx context.Context, channelID string) ([]onMessageCreate.LineNgID, error)
	GetLineBot(ctx context.Context, guildID string) (onMessageCreate.LineBot, error)
	GetLineBotIv(ctx context.Context, guildID string) (onMessageCreate.LineBotIv, error)
}

func (h *CogHandler) OnMessageCreate(s *discordgo.Session, vs *discordgo.MessageCreate) {
	var channel onMessageCreate.LinePostDiscordChannel
	var lineMessageTypes []*line.LineMessageType
	var imageUrls []string
	var videoCount, voiceCount int
	var repo Repository

	sendTextBuilder := strings.Builder{}

	ctx := context.Background()
	repo = onMessageCreate.NewRepository(h.DB)
	channel, err := repo.GetLinePostDiscordChannel(ctx, vs.ChannelID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		slog.ErrorContext(ctx, "line_post_discord_channelの取得に失敗しました", "エラー:", err.Error())
		return
	} else if err != nil {
		err = repo.InsertLinePostDiscordChannel(ctx, vs.ChannelID, vs.GuildID)
		if err != nil {
			slog.ErrorContext(ctx, "line_post_discord_channelの登録に失敗しました", "エラー:", err.Error())
			return
		}
		channel = onMessageCreate.LinePostDiscordChannel{
			Ng:         false,
			BotMessage: false,
		}
	}
	ngTypes, err := repo.GetLineNgDiscordMessageType(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_message_typeの取得に失敗しました", "エラー:", err.Error())
		return
	}
	ngDiscordIDs, err := repo.GetLineNgDiscordID(ctx, vs.ChannelID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_message_typeの登録に失敗しました", "エラー:", err.Error())
		return
	}
	// メッセージの種類がNGの場合は処理を終了
	for _, ngType := range ngTypes {
		if vs.Message.Type == discordgo.MessageType(ngType) {
			slog.InfoContext(ctx, "NG Type")
			return
		}
	}
	// メッセージの送信者がNGの場合は処理を終了
	for _, ngDiscordID := range ngDiscordIDs {
		if ngDiscordID.IDType == "user" && vs.Author.ID == ngDiscordID.ID {
			slog.InfoContext(ctx, "NG User")
			return
		}
		if ngDiscordID.IDType == "role" {
			for _, role := range vs.Member.Roles {
				if role == ngDiscordID.ID {
					slog.InfoContext(ctx, "NG Role")
					return
				}
			}
		}
	}
	// チャンネルがNGの場合、またはBotメッセージでない場合は処理を終了
	if channel.Ng || (!channel.BotMessage && vs.Author.Bot) {
		slog.InfoContext(ctx, "NG Channel or Bot Message")
		return
	}
	lineBotApi, err := repo.GetLineBot(ctx, vs.GuildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_botの取得に失敗しました", "エラー:", err.Error())
		return
	}
	lineBotIv, err := repo.GetLineBotIv(ctx, vs.GuildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_bot_ivの取得に失敗しました", "エラー:", err.Error())
		return
	}
	var lineBotDecrypt onMessageCreate.LineBotDecrypt
	// 暗号化キーのバイトへの変換
	keyBytes, err := hex.DecodeString(config.PrivateKey())
	if err != nil {
		slog.ErrorContext(ctx, "暗号化キーのバイト変換に失敗しました", "エラー:", err.Error())
		return
	}

	lineNotifyTokenByte, err := crypto.Decrypt(lineBotApi.LineNotifyToken[0], keyBytes, lineBotIv.LineNotifyTokenIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_notify_tokenの復号化に失敗しました", "エラー:", err.Error())
		return
	}
	lineBotTokenByte, err := crypto.Decrypt(lineBotApi.LineBotToken[0], keyBytes, lineBotIv.LineBotTokenIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_bot_tokenの復号化に失敗しました", "エラー:", err.Error())
		return
	}
	lineGroupIDByte, err := crypto.Decrypt(lineBotApi.LineGroupID[0], keyBytes, lineBotIv.LineGroupIDIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_group_idの復号化に失敗しました", "エラー:", err.Error())
		return
	}
	lineBotDecrypt.LineNotifyToken = string(lineNotifyTokenByte)
	lineBotDecrypt.LineBotToken = string(lineBotTokenByte)
	lineBotDecrypt.LineGroupID = string(lineGroupIDByte)
	lineBotDecrypt.DefaultChannelID = lineBotApi.DefaultChannelID
	lineBotDecrypt.DebugMode = lineBotApi.DebugMode

	lineRequ := line.NewLineRequest(
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
			return
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
				return
			}
			lineMessageType := lineRequ.NewLineVideoMessage(attachment.URL, st.IconURL("512"))
			lineMessageTypes = append(lineMessageTypes, lineMessageType)
			videoCount++
		case ".mp3", ".wav", ".ogg", ".m4a":
			tmpFile := os.TempDir() + "/" + attachment.Filename
			tmpFileNotExt := os.TempDir() + "/" + fileNameNoExt
			slog.InfoContext(ctx, "download:"+attachment.URL)
			err = downloadFile(tmpFile, attachment.URL)
			if err != nil {
				slog.ErrorContext(ctx, "ファイルのダウンロードに失敗しました", "エラー:", err.Error())
				return
			}
			if extension != ".m4a" {
				slog.InfoContext(ctx, "m4a変換:"+tmpFile)
				slog.InfoContext(ctx, "ffmpeg:"+tmpFile)
				err = exec.CommandContext(ctx, "ffmpeg", "-i", tmpFile, tmpFileNotExt+".m4a").Run()
				if err != nil {
					slog.ErrorContext(ctx, "ffmpegの秒数カウントに失敗しました", "エラー:", err.Error())
					return
				}
			}
			f, err := os.Open(tmpFileNotExt + ".m4a")
			if err != nil {
				slog.ErrorContext(ctx, "ファイルのオープンに失敗しました", "エラー:", err.Error())
				return
			}
			defer f.Close()
			messsage, err := s.ChannelFileSendWithMessage(
				vs.ChannelID,
				"m4aに変換します。",
				tmpFileNotExt+".m4a",
				f,
			)
			if err != nil {
				slog.ErrorContext(ctx, "Discordへのメッセージ送信に失敗しました", "エラー:", err.Error())
				return
			}
			// 音声ファイルの秒数を取得
			cmd := exec.CommandContext(
				ctx,
				"ffprobe",
				"-hide_banner",
				tmpFileNotExt+".m4a",
				"-show_entries",
				"format=duration",
			)
			slog.InfoContext(ctx, "秒数取得:"+tmpFileNotExt+".m4a")
			out, err := cmd.CombinedOutput()
			if err != nil {
				slog.ErrorContext(ctx, "ffmpegの実行に失敗しました", "エラー:", err.Error())
				return
			}
			re, err := regexp.Compile(`(\d+\.\d+)`)
			if err != nil {
				slog.ErrorContext(ctx, "正規表現のコンパイルに失敗しました", "エラー:", err.Error())
				return
			}
			match := re.FindStringSubmatch(string(out))
			slog.InfoContext(ctx, "秒数:"+match[0])
			audioLen, err := strconv.ParseFloat(match[0], 64)
			if err != nil {
				slog.ErrorContext(ctx, "音声ファイルの秒数の抽出に失敗しました", "エラー:", err.Error())
				return
			}
			audio := lineRequ.NewLineAudioMessage(
				messsage.Attachments[0].URL,
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
			return
		}
	}
	// 動画、音声を送信
	if len(lineMessageTypes) > 0 {
		err = lineRequ.PushMessageBotInGroup(ctx, lineMessageTypes)
		if err != nil {
			slog.ErrorContext(ctx, "音声、動画の送信に失敗しました", "エラー:", err.Error())
			return
		}
	}
	// 画像、動画、音声がない場合はテキストのみ送信
	if len(imageUrls) > 0 || videoCount > 0 || voiceCount > 0 {
		return
	}
	err = lineRequ.PushMessageNotify(ctx, sendTextBuilder.String())
	if err != nil {
		slog.ErrorContext(ctx, "LINE Notifyの送信に失敗しました", "エラー:", err.Error())
		return
	}
}

func downloadFile(tmpFilePath, url string) error {
	f, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}
