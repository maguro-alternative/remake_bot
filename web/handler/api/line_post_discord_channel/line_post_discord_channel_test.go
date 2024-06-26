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

	"github.com/stretchr/testify/assert"
)

func TestLinePostDiscordChannelHandler_ServeHTTP(t *testing.T) {
	bodyJson, err := json.Marshal(internal.LinePostDiscordChannelJson{
		GuildID: "987654321",
		Channels: []struct {
			ChannelID  string   `json:"channelId"`
			Ng         bool     `json:"ng"`
			BotMessage bool     `json:"botMessage"`
			NgTypes    []int    `json:"ngTypes"`
			NgUsers    []string `json:"ngUsers"`
			NgRoles    []string `json:"ngRoles"`
		}{
			{
				ChannelID: "123456789",
				NgTypes:   []int{},
				NgUsers:   []string{},
				NgRoles:   []string{},
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
			repo: &repository.RepositoryFuncMock{
				UpdateLinePostDiscordChannelFunc: func(ctx context.Context, lineChannel repository.LinePostDiscordChannelAllColumns) error {
					return nil
				},
				InsertLineNgDiscordMessageTypesFunc: func(ctx context.Context, lineNgDiscordMessageTypes []repository.LineNgDiscordMessageType) error {
					return nil
				},
				DeleteMessageTypesNotInProvidedListFunc: func(ctx context.Context, guildId string, lineNgDiscordMessageTypes []repository.LineNgDiscordMessageType) error {
					return nil
				},
				InsertLineNgDiscordUserIDsFunc: func(ctx context.Context, lineNgDiscordIDs []repository.LineNgDiscordUserIDAllCoulmns) error {
					return nil
				},
				InsertLineNgDiscordRoleIDsFunc: func(ctx context.Context, lineNgDiscordIDs []repository.LineNgDiscordRoleIDAllCoulmns) error {
					return nil
				},
				DeleteUserIDsNotInProvidedListFunc: func(ctx context.Context, guildId string, lineNgDiscordIDs []repository.LineNgDiscordUserIDAllCoulmns) error {
					return nil
				},
				DeleteRoleIDsNotInProvidedListFunc: func(ctx context.Context, guildId string, lineNgDiscordIDs []repository.LineNgDiscordRoleIDAllCoulmns) error {
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
