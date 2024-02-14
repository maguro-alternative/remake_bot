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
	var sendText string
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
		sendText = vs.Message.Author.Username + "「 " + vs.Message.Content + " 」"
	}

	for _, attachment := range vs.Message.Attachments {
		extension := filepath.Ext(attachment.Filename)
		fileNameNoExt := filepath.Base(attachment.Filename[:len(attachment.Filename)-len(extension)])
		switch extension {
		case ".png", ".jpg", ".jpeg", ".gif":
			err = lineRequ.PushImageNotify(ctx, sendText, attachment.URL)
			if err != nil {
				return
			}
		case ".mp4", ".mov", ".avi", ".wmv", ".flv", ".webm":
			st, err := s.Guild(vs.GuildID)
			if err != nil {
				return
			}
			lineMessageType := lineRequ.NewLineVideoMessage(attachment.URL, st.IconURL("512"))
			lineMessageTypes = append(lineMessageTypes, &lineMessageType)
		case ".mp3", ".wav", ".ogg", ".m4a":
			tmpFile := os.TempDir()+"/"+attachment.Filename
			tmpFileNotExt := os.TempDir()+"/"+fileNameNoExt
			downloadFilePath, err := downloadFile(tmpFile, attachment.URL)
			if err != nil {
				return
			}
			err = exec.CommandContext(ctx, "ffmpeg", "-i", downloadFilePath, tmpFileNotExt+".m4a").Run()
			if err != nil {
				return
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
			lineMessageTypes = append(lineMessageTypes, &audio)
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
