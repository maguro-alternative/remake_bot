package guilds

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/maguro-alternative/remake_bot/testutil/mock"

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
	require.NoError(t, os.Chdir("../../../../"))
	t.Run("サーバー一覧が正常に表示される(1件のみ)", func(t *testing.T) {
		indexService := &service.IndexService{
			Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`
					[
						{
							"id": "123",
							"name": "test",
							"icon": "test",
							"owner": true,
							"permissions": 1,
							"features": ["test"]
						}
					]
					`)),
				}
			}),
			DiscordSession: &mock.SessionMock{
				UserGuildsFunc: func(limit int, beforeID string, afterID string, options ...discordgo.RequestOption) ([]*discordgo.UserGuild, error) {
					return []*discordgo.UserGuild{
						{
							ID:          "123",
							Name:        "test",
							Icon:        "test",
							Owner:       true,
							Permissions: 1,
							Features:    []discordgo.GuildFeature{"test"},
						},
					}, nil
				},
			},
		}
		indexService.DiscordBotState = discordgo.NewState()
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

		mux := http.NewServeMux()

		handler := NewGuildsViewHandler(indexService)

		mux.HandleFunc("/guilds", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guilds", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), `<p>Discordアカウント: test</p>`)
		assert.Contains(t, rec.Body.String(), `<img src="https://cdn.discordapp.com/avatars/123/test.webp?size=64" alt="Discordアイコン">`)
		assert.Contains(t, rec.Body.String(), `<p>LINEアカウント: 未ログイン</p>`)

		assert.Contains(t, rec.Body.String(), `<a href="/guild/123">`)
		assert.Contains(t, rec.Body.String(), `<img src="https://cdn.discordapp.com/icons/123/test.png" alt="test">`)
		assert.Contains(t, rec.Body.String(), `<li>test</li>`)
	})
}

func setCtxValue(r *http.Request) *http.Request {
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
