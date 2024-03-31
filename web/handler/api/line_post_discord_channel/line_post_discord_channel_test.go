package linepostdiscordchannel

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/line_post_discord_channel/internal"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/stretchr/testify/assert"
)

func TestLinePostDiscordChannelHandler_ServeHTTP(t *testing.T) {
	bodyJson, err := json.Marshal(internal.LinePostDiscordChannelJson{
		GuildID: "987654321",
		Channels: []struct {
			ChannelID   string   `json:"channel_id"`
			Ng          bool     `json:"ng"`
			BotMessage  bool     `json:"bot_message"`
			NgTypes     []int    `json:"ng_types"`
			NgUsers     []string `json:"ng_users"`
			NgRoles     []string `json:"ng_roles"`
		}{
			{
				ChannelID:   "123456789",
				NgTypes:     []int{},
				NgUsers:     []string{},
				NgRoles:     []string{},
			},
		},
	})
	assert.NoError(t, err)

	t.Run("MethodがPOST以外の場合、Method Not Allowedが返ること", func(t *testing.T) {
		h := &LinePostDiscordChannelHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/987654321/line_post_discord_channel", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("jsonの読み取りに失敗すると、Internal Server Errorが返ること", func(t *testing.T) {
		h := &LinePostDiscordChannelHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/line_post_discord_channel", bytes.NewReader([]byte("")))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("jsonのバリデーションに失敗すると、Internal Server Errorが返ること", func(t *testing.T) {
		h := &LinePostDiscordChannelHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/line_post_discord_channel", bytes.NewReader([]byte(`{"channel_id":"123456789"}`)))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("LinePostDiscordChannelの更新が成功すること", func(t *testing.T) {
		h := &LinePostDiscordChannelHandler{
			IndexService: &service.IndexService{},
			Repo: &repository.RepositoryFuncMock{
				UpdateLinePostDiscordChannelFunc: func(ctx context.Context, lineChannel repository.LinePostDiscordChannelAllColumns) error {
					return nil
				},
				InsertLineNgDiscordMessageTypesFunc: func(ctx context.Context, lineNgDiscordMessageTypes []repository.LineNgDiscordMessageType) error {
					return nil
				},
				DeleteNotInsertLineNgDiscordMessageTypesFunc: func(ctx context.Context, lineNgDiscordMessageTypes []repository.LineNgDiscordMessageType) error {
					return nil
				},
				InsertLineNgDiscordIDsFunc: func(ctx context.Context, lineNgDiscordIDs []repository.LineNgDiscordIDAllCoulmns) error {
					return nil
				},
				DeleteNotInsertLineNgDiscordIDsFunc: func(ctx context.Context, lineNgDiscordIDs []repository.LineNgDiscordIDAllCoulmns) error {
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/line_post_discord_channel", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
