package group

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/group/internal"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestLineGroupHandler_ServeHTTP(t *testing.T) {
	bodyJson, err := json.Marshal(internal.LineBotJson{
		DefaultChannelID: "123456789",
		DebugMode:        true,
	})
	assert.NoError(t, err)

	t.Run("MethodがPOST以外の場合、Method Not Allowedが返ること", func(t *testing.T) {
		h := &LineGroupHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/group", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("jsonの読み取りに失敗すると、Bad Requestが返ること", func(t *testing.T) {
		h := &LineGroupHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/group", bytes.NewReader([]byte("")))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("jsonのバリデーションに失敗すると、Unprocessable Entityが返ること", func(t *testing.T) {
		h := &LineGroupHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/group", bytes.NewReader([]byte(`{"channel_id":"123456789"}`)))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	newIndexService := func(t *testing.T) *service.IndexService {
		state := discordgo.NewState()
		err := state.GuildAdd(&discordgo.Guild{
			ID: "987654321",
			Channels: []*discordgo.Channel{
				{ID: "123456789", GuildID: "987654321", Type: discordgo.ChannelTypeGuildText},
			},
		})
		assert.NoError(t, err)
		err = state.GuildAdd(&discordgo.Guild{
			ID: "111111111",
			Channels: []*discordgo.Channel{
				{ID: "222222222", GuildID: "111111111", Type: discordgo.ChannelTypeGuildText},
			},
		})
		assert.NoError(t, err)
		return &service.IndexService{DiscordBotState: state}
	}

	t.Run("LineBotの更新が成功すること", func(t *testing.T) {
		h := &LineGroupHandler{
			indexService: newIndexService(t),
			repo: &repository.RepositoryFuncMock{
				UpdateLineBotFunc: func(ctx context.Context, lineBot *repository.LineBot) error {
					return nil
				},
			},
		}
		mux := http.NewServeMux()
		mux.Handle("/api/{guildId}/group", h)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/group", bytes.NewReader(bodyJson))
		mux.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("他のサーバーのチャンネルを指定すると、Bad Requestが返ること", func(t *testing.T) {
		otherGuildChannelJson, err := json.Marshal(internal.LineBotJson{
			DefaultChannelID: "222222222",
		})
		assert.NoError(t, err)
		h := &LineGroupHandler{
			indexService: newIndexService(t),
			repo: &repository.RepositoryFuncMock{
				UpdateLineBotFunc: func(ctx context.Context, lineBot *repository.LineBot) error {
					t.Fatal("UpdateLineBot should not be called")
					return nil
				},
			},
		}
		mux := http.NewServeMux()
		mux.Handle("/api/{guildId}/group", h)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/group", bytes.NewReader(otherGuildChannelJson))
		mux.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
