package cogs

import (
	"context"
	"encoding/hex"
	"path/filepath"

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
		sendText = vs.Message.Author.Username+"がサーバーブーストしました。"
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierOne:
		sendText = vs.Message.Author.Username+"がサーバーブーストし、レベル1になりました！！！！！！！！"
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierTwo:
		sendText = vs.Message.Author.Username+"がサーバーブーストし、レベル2になりました！！！！！！！！"
	case discordgo.MessageTypeUserPremiumGuildSubscriptionTierThree:
		sendText = vs.Message.Author.Username+"がサーバーブーストし、レベル3になりました！！！！！！！！"
	case discordgo.MessageTypeGuildMemberJoin:
		sendText = vs.Message.Author.Username+"が参加しました。"
	default:
		sendText = vs.Message.Author.Username+"「 "+vs.Message.Content+" 」"
	}

	for _, attachment := range vs.Message.Attachments {
		extension := filepath.Ext(attachment.Filename)
		switch extension{
		case ".png":
		case ".jpg":
		case ".jpeg":
		case ".gif":
			err = lineRequ.PushImageNotify(ctx, sendText, attachment.URL)
			if err != nil {
				return
			}
		case ".mp4":
		case ".mov":
			st, err := s.Guild(vs.GuildID)
			if err != nil {
				return
			}
			lineMessageType := lineRequ.NewLineVideoMessage(attachment.URL, st.IconURL("512"))
			lineMessageTypes = append(lineMessageTypes, &lineMessageType)
		case ".mp3":
		case ".wav":
		case ".ogg":
		case ".m4a":

		}
	}
}
