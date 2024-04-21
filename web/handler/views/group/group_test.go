package group

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/ctxvalue"
	"github.com/maguro-alternative/remake_bot/web/shared/model"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLineGroupViewHandler(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(cwd))
	})
	require.NoError(t, os.Chdir("../../../../"))
	t.Run("test new line group view handler", func(t *testing.T) {
		indexService := &service.IndexService{
			DiscordSession: &discordgo.Session{},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
				ID: "123",
				Channels: []*discordgo.Channel{
					{
						ID: "123",
						Name: "test",
						Position: 1,
						Type: discordgo.ChannelTypeGuildText,
					},
					{
						ID: "1234",
						Name: "test",
						Position: 2,
						Type: discordgo.ChannelTypeGuildText,
					},
					{
						ID: "12345",
						Name: "test",
						Position: 3,
						Type: discordgo.ChannelTypeGuildText,
					},
				},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetAllColumnsLineBotFunc: func(ctx context.Context, guildID string) (repository.LineBot, error) {
				return repository.LineBot{
					GuildID: "123",
					DefaultChannelID: "123",
				}, nil
			},
		}
		handler := NewLineGroupViewHandler(indexService, repo)

		_, err = handler.IndexService.DiscordBotState.Guild("123")
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/group/123", nil)

		mux.HandleFunc("/group/{guildId}", handler.Index)
		mux.ServeHTTP(w, setCtxValue(r))

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func setCtxValue(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = ctxvalue.ContextWithLineUser(ctx, &model.LineOAuthSession{
		User: model.LineIdTokenUser{
			Iss: "test",
			Sub: "test",
			Aud: "test",
			Exp: 1,
			Iat: 1,
			AuthTime: 1,
			Nonce: "test",
			Amr: []string{"test"},
			Name: "test",
			Picture: "test",
			Email: "test",
		},
		DiscordGuildID: "123",
		Token: "test",
	})
	return r.WithContext(ctx)
}
