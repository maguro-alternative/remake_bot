package cogs

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/maguro-alternative/remake_bot/pkg/line"
	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	onMessageCreate "github.com/maguro-alternative/remake_bot/bot/cogs/on_message_create"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				GetLineNgDiscordUserIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return nil, nil
				},
				GetLineNgDiscordRoleIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
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
				ConversionAudioFileFunc: func(tmpFile, tmpFileNotExt string) error {
					return nil
				},
				GetAudioFileSecondFunc: func(tmpFile, tmpFileNotExt string) (float64, error) {
					return 0.0, nil
				},
			},
			&mock.SessionMock{
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
				GetLineNgDiscordUserIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return nil, nil
				},
				GetLineNgDiscordRoleIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
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
				ConversionAudioFileFunc: func(tmpFile, tmpFileNotExt string) error {
					return nil
				},
				GetAudioFileSecondFunc: func(tmpFile, tmpFileNotExt string) (float64, error) {
					return 0.0, nil
				},
			},
			&mock.SessionMock{
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
				GetLineNgDiscordUserIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return []string{"userID"}, nil
				},
				GetLineNgDiscordRoleIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return []string{"roleID"}, nil
				},
				GetLineBotNotClientFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return *lineBot, nil
				},
				GetLineBotIvNotClientFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return lineBotIv, nil
				},
			},
			&onMessageCreate.FfmpegMock{
				ConversionAudioFileFunc: func(tmpFile, tmpFileNotExt string) error {
					return nil
				},
				GetAudioFileSecondFunc: func(tmpFile, tmpFileNotExt string) (float64, error) {
					return 0.0, nil
				},
			},
			&mock.SessionMock{
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
		cwd, err := os.Getwd()
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, os.Chdir(cwd))
		})
		require.NoError(t, os.Chdir("../../"))

		testFilesPath, err := os.Getwd()
		require.NoError(t, err)

		srcMp3, err := os.Open(testFilesPath + "/testutil/files/yumi_dannasama.mp3")
		require.NoError(t, err)

		dstMp3, err := os.Create(os.TempDir() + "/yumi_dannasama.mp3")
		require.NoError(t, err)

		srcM4a, err := os.Open(testFilesPath + "/testutil/files/yumi_dannasama.m4a")
		require.NoError(t, err)

		dstM4a, err := os.Create(os.TempDir() + "/yumi_dannasama.m4a")
		require.NoError(t, err)

		defer func() {
			require.NoError(t, os.Remove(dstMp3.Name()))
			require.NoError(t, os.Remove(dstM4a.Name()))
		}()

		_, err = io.Copy(dstMp3, srcMp3)
		require.NoError(t, err)
		srcMp3.Close()
		dstMp3.Close()

		_, err = io.Copy(dstM4a, srcM4a)
		require.NoError(t, err)

		srcM4a.Close()
		dstM4a.Close()
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
				GetLineNgDiscordUserIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return nil, nil
				},
				GetLineNgDiscordRoleIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
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
				ConversionAudioFileFunc: func(tmpFile, tmpFileNotExt string) error {
					return nil
				},
				GetAudioFileSecondFunc: func(tmpFile, tmpFileNotExt string) (float64, error) {
					return 0.0, nil
				},
			},
			&mock.SessionMock{
				ChannelFunc: func(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error) {
					return &discordgo.Channel{
						GuildID: "guildID",
					}, nil
				},
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{}, nil
				},
				ChannelFileSendWithMessageFunc: func(channelID string, content string, name string, r io.Reader, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Attachments: []*discordgo.MessageAttachment{
							{
								URL:      "https://example.com/yumi_dannasama.mp3",
								Filename: "yumi_dannasama.mp3",
							},
						},
					}, nil
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
							URL:      "https://example.com/yumi_dannasama.mp3",
							Filename: "yumi_dannasama.mp3",
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
