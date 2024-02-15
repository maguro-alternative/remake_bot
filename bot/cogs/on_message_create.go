package cogs

import (
	"context"
	"encoding/hex"
	"os/exec"
	"os"
	"io"
	"path/filepath"
	"net/http"
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
	if err.Error() == "sql: no rows in result set" {
		err = repo.InsertLineChannel(ctx, vs.ChannelID, vs.GuildID)
		if err != nil {
			return
		}
		channel, err = repo.GetLineChannel(ctx, vs.ChannelID)
		if err != nil {
			return
		}
	} else if err != nil {
		return
	}
	if channel.Ng || (!channel.BotMessage && vs.Author.Bot) {
		return
	}
	lineBotApi, err := repo.GetLineBot(ctx, vs.GuildID)
	if err != nil {
		return
	}
	var lineBotDecrypt onMessageCreate.LineBotDecrypt
	// 暗号化キーのバイトへの変換
	keyBytes, err := hex.DecodeString(config.PrivateKey())
	if err != nil {
		return
	}
	lineNotifyTokenByte, err := crypto.Decrypt(lineBotApi.LineNotifyToken, keyBytes, lineBotApi.Iv)
	if err != nil {
		return
	}
	lineBotTokenByte, err := crypto.Decrypt(lineBotApi.LineBotToken, keyBytes, lineBotApi.Iv)
	if err != nil {
		return
	}
	lineGroupByte, err := crypto.Decrypt(lineBotApi.LineGroupID, keyBytes, lineBotApi.Iv)
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
		sendText = st.Name+"にて、"+vs.Message.Author.Username
	}

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
			tmpFile := os.TempDir()+"/"+attachment.Filename
			tmpFileNotExt := os.TempDir()+"/"+fileNameNoExt
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
			f, err := os.Open(tmpFileNotExt+".m4a")
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

	for _, url := range imageUrls {
		err = lineRequ.PushImageNotify(ctx, sendText, url)
		if err != nil {
			return
		}
	}
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
