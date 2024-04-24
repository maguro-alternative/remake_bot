package views

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/ctxvalue"
	"github.com/maguro-alternative/remake_bot/web/shared/model"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLinePostDiscordChannelViewHandler(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(cwd))
	})
	require.NoError(t, os.Chdir("../../../"))
	t.Run("ルートページが正常に表示される(discordログイン中)", func(t *testing.T) {
		indexService := &service.IndexService{
			DiscordSession: &discordgo.Session{},
		}
		indexService.DiscordBotState = discordgo.NewState()
		indexService.DiscordBotState.User = &discordgo.User{
			Username: "test",
		}
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "123",
			Channels: []*discordgo.Channel{
				{
					ID:       "123",
					Name:     "test",
					Position: 1,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "1234",
					Name:     "test",
					Position: 2,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "12345",
					Name:     "test",
					Position: 3,
					Type:     discordgo.ChannelTypeGuildText,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID: "123",
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		handler := NewIndexViewHandler(indexService)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.Index(rec, setCtxDiscordUserValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), `<p>Discordアカウント: test</p>`)
		assert.Contains(t, rec.Body.String(), `<img src="https://cdn.discordapp.com/avatars/123/test.webp?size=64" alt="Discordアイコン">`)
		assert.Contains(t, rec.Body.String(), `<p>LINEアカウント: 未ログイン</p>`)
	})

	t.Run("ルートページが正常に表示される(lineログイン中)", func(t *testing.T) {
		indexService := &service.IndexService{
			DiscordSession: &discordgo.Session{},
		}
		indexService.DiscordBotState = discordgo.NewState()
		indexService.DiscordBotState.User = &discordgo.User{
			Username: "test",
		}
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "123",
			Channels: []*discordgo.Channel{
				{
					ID:       "123",
					Name:     "test",
					Position: 1,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "1234",
					Name:     "test",
					Position: 2,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "12345",
					Name:     "test",
					Position: 3,
					Type:     discordgo.ChannelTypeGuildText,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID: "123",
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		handler := NewIndexViewHandler(indexService)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.Index(rec, setCtxLineUserValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), `<p>Discordアカウント: 未ログイン</p>`)
		assert.Contains(t, rec.Body.String(), `<p>LINEアカウント: test</p>`)
		assert.Contains(t, rec.Body.String(), `<img src="test" style="height:64px;" alt="LINEアイコン">`)
	})

}

func setCtxDiscordUserValue(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = ctxvalue.ContextWithDiscordUser(ctx, &model.DiscordOAuthSession{
		Token: "test",
		User: model.DiscordUser{
			ID:       "123",
			Username: "test",
			Avatar:   "test",
		},
	},
	)
	return r.WithContext(ctx)
}

func setCtxLineUserValue(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = ctxvalue.ContextWithLineUser(ctx, &model.LineOAuthSession{
		Token:          "test",
		DiscordGuildID: "123",
		User: model.LineIdTokenUser{
			Iss:     "https://access.line.me",
			Sub:     "U123456789abcdef123456789abcdef12",
			Aud:     "1234567890",
			Exp:     1504169092,
			Iat:     1504168492,
			Nonce:   "0987654321fedcba987654321fedcba9",
			Amr:     []string{"pwd"},
			Name:    "test",
			Picture: "test",
			Email:   "test",
		},
	})
	return r.WithContext(ctx)
}
