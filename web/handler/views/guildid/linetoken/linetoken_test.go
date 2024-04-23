package linetoken

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/maguro-alternative/remake_bot/repository"

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
	require.NoError(t, os.Chdir("../../../../../"))
	t.Run("正常に表示されるか", func(t *testing.T) {
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

		repo := &repository.RepositoryFuncMock{
			GetAllColumnsLineBotFunc: func(ctx context.Context, guildId string) (repository.LineBot, error) {
				return repository.LineBot{
					GuildID:          "111",
					LineNotifyToken:  pq.ByteaArray{[]byte("test")},
					LineBotToken:     pq.ByteaArray{[]byte("test")},
					LineBotSecret:    pq.ByteaArray{[]byte("test")},
					LineGroupID:      pq.ByteaArray{[]byte("test")},
					LineClientID:     pq.ByteaArray{[]byte("test")},
					LineClientSecret: pq.ByteaArray{[]byte("test")},
				}, nil
			},
		}
		handler := NewLineTokenViewHandler(indexService, repo)

		mux.HandleFunc("/guilds/{guildId}/linetoken", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guilds/123/linetoken", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), `<p>Discordアカウント: test</p>`)
		assert.Contains(t, rec.Body.String(), `<img src="https://cdn.discordapp.com/avatars/123/test.webp?size=64" alt="Discordアイコン">`)
		assert.Contains(t, rec.Body.String(), `<p>LINEアカウント: 未ログイン</p>`)

		assert.Contains(t, rec.Body.String(), `<input id="lineNotifyToken" type="password" name="line_notify_token">`)
		assert.Contains(t, rec.Body.String(), `<input id="lineBotToken" type="password" name="line_bot_token">`)
		assert.Contains(t, rec.Body.String(), `<input id="lineBotSecret" type="password" name="line_bot_secret">`)
		assert.Contains(t, rec.Body.String(), `<input id="lineGroupId" type="password" name="line_group_id">`)
		assert.Contains(t, rec.Body.String(), `<input id="lineClientId" type="password" name="line_client_id">`)
		assert.Contains(t, rec.Body.String(), `<input id="lineClientSecret" type="password" name="line_client_secret">`)
		assert.Equal(t, strings.Count(rec.Body.String(),"入力済み"), 6)

		assert.Contains(t, rec.Body.String(), `<option value="123">カテゴリーなし:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="1234">カテゴリーなし:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="12345">カテゴリーなし:📝:test</option>`)
	})
}

func setCtxValue(r *http.Request) *http.Request {
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
