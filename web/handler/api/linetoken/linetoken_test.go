package linetoken

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/linetoken/internal"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/stretchr/testify/assert"
)

func TestLineTokenHandler_ServeHTTP(t *testing.T) {
	bodyJson, err := json.Marshal(internal.LineBotJson{
		GuildID:          "987654321",
		DefaultChannelID: "123456789",
		DebugMode:        true,
	})
	assert.NoError(t, err)

	t.Run("MethodがPOST以外の場合、Method Not Allowedが返ること", func(t *testing.T) {
		h := &LineTokenHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/987654321/linetoken", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("jsonの読み取りに失敗すると、Bad Requestが返ること", func(t *testing.T) {
		h := &LineTokenHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/linetoken", bytes.NewReader([]byte("")))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("jsonのバリデーションに失敗すると、Unprocessable Entityが返ること", func(t *testing.T) {
		h := &LineTokenHandler{
			IndexService: &service.IndexService{},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/linetoken", bytes.NewReader([]byte(`{"channel_id":"123456789"}`)))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("LineBotの更新が成功すること", func(t *testing.T) {
		h := &LineTokenHandler{
			IndexService: &service.IndexService{
				DiscordSession: &discordgo.Session{},
			},
			Repo: &repository.RepositoryFuncMock{
				UpdateLineBotFunc: func(ctx context.Context, lineBot *repository.LineBot) error {
					return nil
				},
				UpdateLineBotIvFunc: func(ctx context.Context, lineBotIv *repository.LineBotIv) error {
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/linetoken", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
