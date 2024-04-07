package cogs

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/maguro-alternative/remake_bot/pkg/line"
	"github.com/maguro-alternative/remake_bot/repository"

	onMessageCreate "github.com/maguro-alternative/remake_bot/bot/cogs/on_message_create"
	"github.com/maguro-alternative/remake_bot/bot/service"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestLineRequest_PushMessageNotify(t *testing.T) {
	ctx := context.Background()
	// スタブHTTPクライアントを作成
	stubClient := line.NewStubHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}
	})

	decodeNotifyToken, err := hex.DecodeString("aa7c5fe80002633327f0fefe67a565de")
	assert.NoError(t, err)
	lineNotifyStr, err := base64.StdEncoding.DecodeString(string([]byte("X+P6kmO6DnEjM3TVqXkwNA==")))
	assert.NoError(t, err)

	decodeBotToken, err := hex.DecodeString("baeff317cb83ef55b193b6d3de194124")
	assert.NoError(t, err)
	lineBotStr, err := base64.StdEncoding.DecodeString(string([]byte("uy2qtvYTnSoB5qIntwUdVQ==")))
	assert.NoError(t, err)

	decodeBotSecret, err := hex.DecodeString("0ffa8ed72efcb5f1d834e4ce8463a62c")
	assert.NoError(t, err)
	lineBotSecretStr, err := base64.StdEncoding.DecodeString(string([]byte("i2uHQCyn58wRR/b03fRw6w==")))
	assert.NoError(t, err)

	decodeGroupID, err := hex.DecodeString("e14db710b23520766fd652c0f19d437a")
	assert.NoError(t, err)
	lineGroupStr, err := base64.StdEncoding.DecodeString(string([]byte("YgexFQQlLcaXmsw9mFN35Q==")))
	assert.NoError(t, err)

	lineBot := &repository.LineBotNotClient{
		LineNotifyToken: pq.ByteaArray{lineNotifyStr},
		LineBotToken:    pq.ByteaArray{lineBotStr},
		LineBotSecret:   pq.ByteaArray{lineBotSecretStr},
		LineGroupID:     pq.ByteaArray{lineGroupStr},
	}
	lineBotIv := repository.LineBotIvNotClient{
		LineNotifyTokenIv: pq.ByteaArray{decodeNotifyToken},
		LineBotTokenIv:    pq.ByteaArray{decodeBotToken},
		LineBotSecretIv:   pq.ByteaArray{decodeBotSecret},
		LineGroupIDIv:     pq.ByteaArray{decodeGroupID},
	}

	t.Run("正常系", func(t *testing.T) {
		err = onMessageCreateFunc(
			ctx,
			stubClient,
			&repository.RepositoryFuncMock{
				GetLinePostDiscordChannelFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
					return repository.LinePostDiscordChannel{
						Ng:         false,
						BotMessage: false,
					}, nil
				},
				InsertLinePostDiscordChannelFunc: func(ctx context.Context, channelID string, guildID string) error {
					return nil
				},
				GetLineNgDiscordMessageTypeFunc: func(ctx context.Context, channelID string) ([]int, error) {
					return nil, nil
				},
				GetLineNgDiscordIDFunc: func(ctx context.Context, channelID string) ([]repository.LineNgDiscordID, error) {
					return nil, nil
				},
				GetLineBotNotClientFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return *lineBot, nil
				},
				GetLineBotIvNotClientFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return lineBotIv, nil
				},
			},
			&onMessageCreate.FfmpegMock{
				ConversionAudioFileFunc: func(ctx context.Context, tmpFile, tmpFileNotExt string) error {
					return nil
				},
				GetAudioFileSecondFunc: func(ctx context.Context, tmpFile, tmpFileNotExt string) (float64, error) {
					return 0.0, nil
				},
			},
			&service.SessionMock{
				ChannelFunc: func(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error) {
					return &discordgo.Channel{
						GuildID: "guildID",
					}, nil
				},
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{}, nil
				},
				GuildFunc: func(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error) {
					return &discordgo.Guild{
						ID: "guildID",
					}, nil
				},
			},
			&discordgo.MessageCreate{
				Message: &discordgo.Message{
					Content: "test",
					Author: &discordgo.User{
						Bot: false,
					},
				},
			},
		)
		assert.NoError(t, err)
	})

	t.Run("正常系(画像)", func(t *testing.T) {
		err = onMessageCreateFunc(
			ctx,
			stubClient,
			&repository.RepositoryFuncMock{
				GetLinePostDiscordChannelFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
					return repository.LinePostDiscordChannel{
						Ng:         false,
						BotMessage: false,
					}, nil
				},
				InsertLinePostDiscordChannelFunc: func(ctx context.Context, channelID string, guildID string) error {
					return nil
				},
				GetLineNgDiscordMessageTypeFunc: func(ctx context.Context, channelID string) ([]int, error) {
					return nil, nil
				},
				GetLineNgDiscordIDFunc: func(ctx context.Context, channelID string) ([]repository.LineNgDiscordID, error) {
					return nil, nil
				},
				GetLineBotNotClientFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return *lineBot, nil
				},
				GetLineBotIvNotClientFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return lineBotIv, nil
				},
			},
			&onMessageCreate.FfmpegMock{
				ConversionAudioFileFunc: func(ctx context.Context, tmpFile, tmpFileNotExt string) error {
					return nil
				},
				GetAudioFileSecondFunc: func(ctx context.Context, tmpFile, tmpFileNotExt string) (float64, error) {
					return 0.0, nil
				},
			},
			&service.SessionMock{
				ChannelFunc: func(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error) {
					return &discordgo.Channel{
						GuildID: "guildID",
					}, nil
				},
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{}, nil
				},
				GuildFunc: func(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error) {
					return &discordgo.Guild{
						ID: "guildID",
					}, nil
				},
			},
			&discordgo.MessageCreate{
				Message: &discordgo.Message{
					Content: "test",
					Attachments: []*discordgo.MessageAttachment{
						{
							URL:      "https://example.com/image.jpg",
							Filename: "image.jpg",
						},
					},
					Author: &discordgo.User{
						Bot: false,
					},
				},
			},
		)
		assert.NoError(t, err)
	})

	t.Run("正常系(NGユーザー)", func(t *testing.T) {
		err = onMessageCreateFunc(
			ctx,
			stubClient,
			&repository.RepositoryFuncMock{
				GetLinePostDiscordChannelFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
					return repository.LinePostDiscordChannel{
						Ng:         true,
						BotMessage: false,
					}, nil
				},
				InsertLinePostDiscordChannelFunc: func(ctx context.Context, channelID string, guildID string) error {
					return nil
				},
				GetLineNgDiscordMessageTypeFunc: func(ctx context.Context, channelID string) ([]int, error) {
					return nil, nil
				},
				GetLineNgDiscordIDFunc: func(ctx context.Context, channelID string) ([]repository.LineNgDiscordID, error) {
					return []repository.LineNgDiscordID{
						{
							ID:     "userID",
							IDType: "user",
						},
					}, nil
				},
				GetLineBotNotClientFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return *lineBot, nil
				},
				GetLineBotIvNotClientFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return lineBotIv, nil
				},
			},
			&onMessageCreate.FfmpegMock{
				ConversionAudioFileFunc: func(ctx context.Context, tmpFile, tmpFileNotExt string) error {
					return nil
				},
				GetAudioFileSecondFunc: func(ctx context.Context, tmpFile, tmpFileNotExt string) (float64, error) {
					return 0.0, nil
				},
			},
			&service.SessionMock{
				ChannelFunc: func(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error) {
					return &discordgo.Channel{
						GuildID: "guildID",
					}, nil
				},
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{}, nil
				},
				GuildFunc: func(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error) {
					return &discordgo.Guild{
						ID: "guildID",
					}, nil
				},
			},
			&discordgo.MessageCreate{
				Message: &discordgo.Message{
					Author: &discordgo.User{
						ID:  "userID",
						Bot: false,
					},
					Content: "test",
				},
			},
		)
		assert.NoError(t, err)
	})

	t.Run("正常系(音声)", func(t *testing.T) {
		err = onMessageCreateFunc(
			ctx,
			stubClient,
			&repository.RepositoryFuncMock{
				GetLinePostDiscordChannelFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
					return repository.LinePostDiscordChannel{
						Ng:         false,
						BotMessage: false,
					}, nil
				},
				InsertLinePostDiscordChannelFunc: func(ctx context.Context, channelID string, guildID string) error {
					return nil
				},
				GetLineNgDiscordMessageTypeFunc: func(ctx context.Context, channelID string) ([]int, error) {
					return nil, nil
				},
				GetLineNgDiscordIDFunc: func(ctx context.Context, channelID string) ([]repository.LineNgDiscordID, error) {
					return nil, nil
				},
				GetLineBotNotClientFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return *lineBot, nil
				},
				GetLineBotIvNotClientFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return lineBotIv, nil
				},
			},
			&onMessageCreate.FfmpegMock{
				ConversionAudioFileFunc: func(ctx context.Context, tmpFile, tmpFileNotExt string) error {
					return nil
				},
				GetAudioFileSecondFunc: func(ctx context.Context, tmpFile, tmpFileNotExt string) (float64, error) {
					return 0.0, nil
				},
			},
			&service.SessionMock{
				ChannelFunc: func(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error) {
					return &discordgo.Channel{
						GuildID: "guildID",
					}, nil
				},
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{}, nil
				},
				GuildFunc: func(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error) {
					return &discordgo.Guild{
						ID: "guildID",
					}, nil
				},
			},
			&discordgo.MessageCreate{
				Message: &discordgo.Message{
					Content: "test",
					Attachments: []*discordgo.MessageAttachment{
						{
							URL:      "https://example.com/voice.mp3",
							Filename: "voice.mp3",
						},
					},
					Author: &discordgo.User{
						Bot: false,
					},
				},
			},
		)
		assert.NoError(t, err)
	})

}
