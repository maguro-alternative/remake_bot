package cogs

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"
	"github.com/maguro-alternative/remake_bot/pkg/crypto"

	"github.com/maguro-alternative/remake_bot/bot/ffmpeg"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineRequest_PushMessageNotify(t *testing.T) {
	ctx := context.Background()
	// スタブHTTPクライアントを作成
	stubClient := mock.NewStubHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}
	})

	lineBot := &repository.LineBotNotClient{
		LineNotifyToken:  pq.ByteaArray{[]byte("lineNotifyStr")},
		LineBotToken:     pq.ByteaArray{[]byte("lineBotStr")},
		LineBotSecret:    pq.ByteaArray{[]byte("lineBotSecretStr")},
		LineGroupID:      pq.ByteaArray{[]byte("lineGroupStr")},
	}
	lineBotIv := repository.LineBotIvNotClient{
		LineNotifyTokenIv:  pq.ByteaArray{[]byte("decodeNotifyToken")},
		LineBotTokenIv:     pq.ByteaArray{[]byte("decodeBotToken")},
		LineBotSecretIv:    pq.ByteaArray{[]byte("decodeBotSecret")},
		LineGroupIDIv:      pq.ByteaArray{[]byte("decodeGroupID")},
	}

	t.Run("正常系", func(t *testing.T) {
		err := onMessageCreateFunc(
			ctx,
			stubClient,
			&repository.RepositoryFuncMock{
				GetLinePostDiscordChannelByChannelIDFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
					return repository.LinePostDiscordChannel{
						Ng:         false,
						BotMessage: false,
					}, nil
				},
				InsertLinePostDiscordChannelByChannelIDAndGuildIDFunc: func(ctx context.Context, channelID string, guildID string) error {
					return nil
				},
				GetLineNgDiscordMessageTypeByChannelIDFunc: func(ctx context.Context, channelID string) ([]int, error) {
					return nil, nil
				},
				GetLineNgDiscordUserIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return nil, nil
				},
				GetLineNgDiscordRoleIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return nil, nil
				},
				GetLineBotNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return *lineBot, nil
				},
				GetLineBotIvNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return lineBotIv, nil
				},
			},
			&ffmpeg.FfmpegMock{
				ConversionAudioFileFunc: func(tmpFile, tmpFileNotExt string) error {
					return nil
				},
				GetAudioFileSecondFunc: func(tmpFile, tmpFileNotExt string) (float64, error) {
					return 0.0, nil
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
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
		err := onMessageCreateFunc(
			ctx,
			stubClient,
			&repository.RepositoryFuncMock{
				GetLinePostDiscordChannelByChannelIDFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
					return repository.LinePostDiscordChannel{
						Ng:         false,
						BotMessage: false,
					}, nil
				},
				InsertLinePostDiscordChannelByChannelIDAndGuildIDFunc: func(ctx context.Context, channelID string, guildID string) error {
					return nil
				},
				GetLineNgDiscordMessageTypeByChannelIDFunc: func(ctx context.Context, channelID string) ([]int, error) {
					return nil, nil
				},
				GetLineNgDiscordUserIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return nil, nil
				},
				GetLineNgDiscordRoleIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return nil, nil
				},
				GetLineBotNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return *lineBot, nil
				},
				GetLineBotIvNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return lineBotIv, nil
				},
			},
			&ffmpeg.FfmpegMock{
				ConversionAudioFileFunc: func(tmpFile, tmpFileNotExt string) error {
					return nil
				},
				GetAudioFileSecondFunc: func(tmpFile, tmpFileNotExt string) (float64, error) {
					return 0.0, nil
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
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
		err := onMessageCreateFunc(
			ctx,
			stubClient,
			&repository.RepositoryFuncMock{
				GetLinePostDiscordChannelByChannelIDFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
					return repository.LinePostDiscordChannel{
						Ng:         true,
						BotMessage: false,
					}, nil
				},
				InsertLinePostDiscordChannelByChannelIDAndGuildIDFunc: func(ctx context.Context, channelID string, guildID string) error {
					return nil
				},
				GetLineNgDiscordMessageTypeByChannelIDFunc: func(ctx context.Context, channelID string) ([]int, error) {
					return nil, nil
				},
				GetLineNgDiscordUserIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return []string{"userID"}, nil
				},
				GetLineNgDiscordRoleIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return []string{"roleID"}, nil
				},
				GetLineBotNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return *lineBot, nil
				},
				GetLineBotIvNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return lineBotIv, nil
				},
			},
			&ffmpeg.FfmpegMock{
				ConversionAudioFileFunc: func(tmpFile, tmpFileNotExt string) error {
					return nil
				},
				GetAudioFileSecondFunc: func(tmpFile, tmpFileNotExt string) (float64, error) {
					return 0.0, nil
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
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
				GetLinePostDiscordChannelByChannelIDFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
					return repository.LinePostDiscordChannel{
						Ng:         false,
						BotMessage: false,
					}, nil
				},
				InsertLinePostDiscordChannelByChannelIDAndGuildIDFunc: func(ctx context.Context, channelID string, guildID string) error {
					return nil
				},
				GetLineNgDiscordMessageTypeByChannelIDFunc: func(ctx context.Context, channelID string) ([]int, error) {
					return nil, nil
				},
				GetLineNgDiscordUserIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return nil, nil
				},
				GetLineNgDiscordRoleIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
					return nil, nil
				},
				GetLineBotNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return *lineBot, nil
				},
				GetLineBotIvNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return lineBotIv, nil
				},
			},
			&ffmpeg.FfmpegMock{
				ConversionAudioFileFunc: func(tmpFile, tmpFileNotExt string) error {
					return nil
				},
				GetAudioFileSecondFunc: func(tmpFile, tmpFileNotExt string) (float64, error) {
					return 0.0, nil
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return []byte("aaa"), nil
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
