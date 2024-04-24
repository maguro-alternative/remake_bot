package guildid

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
	require.NoError(t, os.Chdir("../../../../"))
	t.Run("設定ページ一覧が正常に表示される(管理者)", func(t *testing.T) {
		indexService := &service.IndexService{
			DiscordSession: &discordgo.Session{},
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

		handler := NewGuildIDViewHandler(indexService)

		mux.HandleFunc("/guild/{guildId}", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setAdminCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), `<p>Discordアカウント: test</p>`)
		assert.Contains(t, rec.Body.String(), `<img src="https://cdn.discordapp.com/avatars/123/test.webp?size=64" alt="Discordアイコン">`)
		assert.Contains(t, rec.Body.String(), `<p>LINEアカウント: 未ログイン</p>`)

		assert.Contains(t, rec.Body.String(), `<a href="/guild/123/permission" class="btn btn-primary">管理者設定</a>`)

		assert.Contains(t, rec.Body.String(), `<a href="/guild/123/line-post-discord-channel" class="btn btn-primary">LINEへの送信設定</a>`)
		assert.Contains(t, rec.Body.String(), `<a href="/guild/123/linetoken" class="btn btn-primary">LINEBOTおよびグループ設定</a>`)
		assert.Contains(t, rec.Body.String(), `<a href="/guild/123/vc-signal" class="btn btn-primary">ボイスチャンネルの通知設定</a>`)
		assert.Contains(t, rec.Body.String(), `<a href="/guild/123/webhook" class="btn btn-primary">webhookの送信設定</a>`)
	})

	t.Run("設定ページ一覧が正常に表示される(管理者じゃない)", func(t *testing.T) {
		indexService := &service.IndexService{
			DiscordSession: &discordgo.Session{},
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

		handler := NewGuildIDViewHandler(indexService)

		mux.HandleFunc("/guild/{guildId}", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setNotAdminCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), `<p>Discordアカウント: test</p>`)
		assert.Contains(t, rec.Body.String(), `<img src="https://cdn.discordapp.com/avatars/123/test.webp?size=64" alt="Discordアイコン">`)
		assert.Contains(t, rec.Body.String(), `<p>LINEアカウント: 未ログイン</p>`)

		assert.NotContains(t, rec.Body.String(), `<a href="/guild/123/permission" class="btn btn-primary">管理者設定</a>`)

		assert.Contains(t, rec.Body.String(), `<a href="/guild/123/line-post-discord-channel" class="btn btn-primary">LINEへの送信設定</a>`)
		assert.Contains(t, rec.Body.String(), `<a href="/guild/123/linetoken" class="btn btn-primary">LINEBOTおよびグループ設定</a>`)
		assert.Contains(t, rec.Body.String(), `<a href="/guild/123/vc-signal" class="btn btn-primary">ボイスチャンネルの通知設定</a>`)
		assert.Contains(t, rec.Body.String(), `<a href="/guild/123/webhook" class="btn btn-primary">webhookの送信設定</a>`)

	})

}

func setAdminCtxValue(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = ctxvalue.ContextWithDiscordPermission(ctx, &model.DiscordPermissionData{
		PermissionCode: 8,
		User: model.DiscordUser{
			ID:       "123",
			Username: "test",
			Avatar:   "test",
		},
		Permission: "all",
	})
	return r.WithContext(ctx)
}

func setNotAdminCtxValue(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = ctxvalue.ContextWithDiscordPermission(ctx, &model.DiscordPermissionData{
		PermissionCode: 0,
		User: model.DiscordUser{
			ID:       "123",
			Username: "test",
			Avatar:   "test",
		},
		Permission: "",
	})
	return r.WithContext(ctx)
}
