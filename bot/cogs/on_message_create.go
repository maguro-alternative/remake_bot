package cogs

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/line"

	onMessageCreate "github.com/maguro-alternative/remake_bot/bot/cogs/on_message_create"
	"github.com/maguro-alternative/remake_bot/bot/config"

	"github.com/bwmarrin/discordgo"
)

type Repository interface {
	GetLineChannel(ctx context.Context, channelID string) (onMessageCreate.LineChannel, error)
	InsertLineChannel(ctx context.Context, channelID string, guildID string) error
	GetLineNgType(ctx context.Context, channelID string) ([]int, error)
	GetLineNgDiscordID(ctx context.Context, channelID string) ([]onMessageCreate.LineNgID, error)
	GetLineBot(ctx context.Context, guildID string) (onMessageCreate.LineBot, error)
	GetLineBotIv(ctx context.Context, guildID string) (onMessageCreate.LineBotIv, error)
}

func (h *CogHandler) OnMessageCreate(s *discordgo.Session, vs *discordgo.MessageCreate) {
	var channel onMessageCreate.LineChannel
	var lineMessageTypes []*line.LineMessageType
	var imageUrls []string
	var sendText string
	var videoCount, voiceCount int

	ctx := context.Background()
	repo := onMessageCreate.NewRepository(h.DB)
	channel, err := repo.GetLineChannel(ctx, vs.ChannelID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return
	} else if err != nil {
		err = repo.InsertLineChannel(ctx, vs.ChannelID, vs.GuildID)
		if err != nil {
			return
		}
		channel = onMessageCreate.LineChannel{
			Ng:          false,
			BotMessage:  false,
		}
	}
	ngTypes, err := repo.GetLineNgType(ctx, vs.ChannelID)
	if err != nil {
		return
	}
	ngDiscordIDs, err := repo.GetLineNgDiscordID(ctx, vs.ChannelID)
	if err != nil {
		return
	}
	// メッセージの種類がNGの場合は処理を終了
	for _, ngType := range ngTypes {
		if vs.Message.Type == discordgo.MessageType(ngType) {
			return
		}
	}
	// メッセージの送信者がNGの場合は処理を終了
	for _, ngDiscordID := range ngDiscordIDs {
		if ngDiscordID.IDType == "user" && vs.Author.ID == ngDiscordID.ID {
			return
		}
		if ngDiscordID.IDType == "role" {
			for _, role := range vs.Member.Roles {
				if role == ngDiscordID.ID {
					return
				}
			}
		}
	}
	// チャンネルがNGの場合、またはBotメッセージでない場合は処理を終了
	if channel.Ng || (!channel.BotMessage && vs.Author.Bot) {
		return
	}
	lineBotApi, err := repo.GetLineBot(ctx, vs.GuildID)
	if err != nil {
		return
	}
	lineBotIv, err := repo.GetLineBotIv(ctx, vs.GuildID)
	if err != nil {
		return
	}
	var lineBotDecrypt onMessageCreate.LineBotDecrypt
	// 暗号化キーのバイトへの変換
	keyBytes, err := hex.DecodeString(config.PrivateKey())
	if err != nil {
		return
	}

	decodeNotifyToken, err := hex.DecodeString(string(lineBotIv.LineNotifyTokenIv[0]))
	if err != nil {
		return
	}
	decodeBotToken, err := hex.DecodeString(string(lineBotIv.LineBotTokenIv[0]))
	if err != nil {
		return
	}
	decodeGroupID, err := hex.DecodeString(string(lineBotIv.LineGroupIDIv[0]))
	if err != nil {
		return
	}
	lineNotifyStr, err := base64.StdEncoding.DecodeString(string(lineBotApi.LineNotifyToken[0]))
	if err != nil {
		return
	}
	lineBotTokenStr, err := base64.StdEncoding.DecodeString(string(lineBotApi.LineBotToken[0]))
	if err != nil {
		return
	}
	lineGroupStr, err := base64.StdEncoding.DecodeString(string(lineBotApi.LineGroupID[0]))
	if err != nil {
		return
	}

	lineNotifyTokenByte, err := crypto.Decrypt(lineNotifyStr, keyBytes, decodeNotifyToken)
	if err != nil {
		return
	}
	lineBotTokenByte, err := crypto.Decrypt(lineBotTokenStr, keyBytes, decodeBotToken)
	if err != nil {
		return
	}
	lineGroupByte, err := crypto.Decrypt(lineGroupStr, keyBytes, decodeGroupID)
	if err != nil {
		return
	}
	lineBotDecrypt.LineNotifyToken = string(lineNotifyTokenByte)
	lineBotDecrypt.LineBotToken = string(lineBotTokenByte)
	lineBotDecrypt.LineGroupID = string(lineGroupByte)
	lineBotDecrypt.DefaultChannelID = lineBotApi.DefaultChannelID
	lineBotDecrypt.DebugMode = lineBotApi.DebugMode

	lineRequ := line.NewLineRequest(
		lineBotDecrypt.LineBotToken,
		lineBotDecrypt.LineGroupID,
		lineBotDecrypt.LineNotifyToken,
	)

	// メッセージの種類によって処理を分岐
	switch vs.ReferencedMessage.Type {
	case discordgo.MessageTypeUserPremiumGuildSubscription:
		sendText = vs.Message.Author.Username + "がサーバーブーストしました。"
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierOne:
		sendText = vs.Message.Author.Username + "がサーバーブーストし、レベル1になりました！！！！！！！！"
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierTwo:
		sendText = vs.Message.Author.Username + "がサーバーブーストし、レベル2になりました！！！！！！！！"
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierThree:
		sendText = vs.Message.Author.Username + "がサーバーブーストし、レベル3になりました！！！！！！！！"
	case discordgo.MessageTypeGuildMemberJoin:
		sendText = vs.Message.Author.Username + "が参加しました。"
	default:
		st, err := s.Channel(vs.ChannelID)
		if err != nil {
			return
		}
		sendText = st.Name + "にて、" + vs.Message.Author.Username
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
				return
			}
			lineMessageType := lineRequ.NewLineVideoMessage(attachment.URL, st.IconURL("512"))
			lineMessageTypes = append(lineMessageTypes, lineMessageType)
			videoCount++
		case ".mp3", ".wav", ".ogg", ".m4a":
			tmpFile := os.TempDir() + "/" + attachment.Filename
			tmpFileNotExt := os.TempDir() + "/" + fileNameNoExt
			downloadFilePath, err := downloadFile(tmpFile, attachment.URL)
			if err != nil {
				return
			}
			if extension != ".m4a" {
				err = exec.CommandContext(ctx, "ffmpeg", "-i", downloadFilePath, tmpFileNotExt+".m4a").Run()
				if err != nil {
					return
				}
			}
			f, err := os.Open(tmpFileNotExt + ".m4a")
			if err != nil {
				return
			}
			defer f.Close()
			messsage, err := s.ChannelFileSend(
				vs.ChannelID,
				"m4aに変換します。"+tmpFileNotExt+".m4a",
				f,
			)
			if err != nil {
				return
			}
			// 音声ファイルの秒数を取得
			cmd := exec.CommandContext(
				ctx,
				"ffprobe",
				"-hide_banner",
				tmpFileNotExt+".m4a",
				"show_entries",
				"format=duration",
			)
			out, err := cmd.CombinedOutput()
			if err != nil {
				return
			}
			match := regexp.MustCompile(`(\d+\.\d+)`).FindStringSubmatch(string(out))
			audioLen, err := strconv.Atoi(match[0])
			if err != nil {
				return
			}
			audio := lineRequ.NewLineAudioMessage(
				messsage.Attachments[0].URL,
				audioLen,
			)
			lineMessageTypes = append(lineMessageTypes, audio)
			voiceCount++
		}
	}

	if len(imageUrls) > 0 {
		sendText += " 画像を" + strconv.Itoa(len(imageUrls)) + "枚、"
	}
	if videoCount > 0 {
		sendText += " 動画を" + strconv.Itoa(videoCount) + "個、"
	}
	if voiceCount > 0 {
		sendText += " 音声を" + strconv.Itoa(voiceCount) + "個、"
	}
	if len(imageUrls) > 0 || videoCount > 0 || voiceCount > 0 {
		sendText += "送信しました。"
	}

	sendText += "「 " + vs.Message.Content + " 」"

	// LINEに送信
	for _, url := range imageUrls {
		err = lineRequ.PushImageNotify(ctx, sendText, url)
		if err != nil {
			return
		}
	}
	// 動画、音声を送信
	if len(lineMessageTypes) > 0 {
		err = lineRequ.PushMessageBotInGroup(ctx, lineMessageTypes)
		if err != nil {
			return
		}
	}
}

func downloadFile(tmpFilePath, url string) (string, error) {
	f, err := os.Create(tmpFilePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return "", err
}
