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

	"github.com/bwmarrin/discordgo"
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
		r.SetPathValue("guildId", "987654321")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("jsonの読み取りに失敗すると、Internal Server Errorが返ること", func(t *testing.T) {
		h := &LinePostDiscordChannelHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/line_post_discord_channel", bytes.NewReader([]byte("")))
		r.SetPathValue("guildId", "987654321")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("jsonのバリデーションに失敗すると、Internal Server Errorが返ること", func(t *testing.T) {
		h := &LinePostDiscordChannelHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/line_post_discord_channel", bytes.NewReader([]byte(`{"channel_id":"123456789"}`)))
		r.SetPathValue("guildId", "987654321")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("LinePostDiscordChannelの更新が成功すること", func(t *testing.T) {
		h := &LinePostDiscordChannelHandler{
			indexService: newIndexService(t),
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
		r.SetPathValue("guildId", "987654321")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("他のサーバーのチャンネルを指定すると、Forbiddenが返ること", func(t *testing.T) {
		otherGuildBodyJson, err := json.Marshal(internal.LinePostDiscordChannelJson{
			Channels: []struct {
				ChannelID  string   `json:"channelId"`
				Ng         bool     `json:"ng"`
				BotMessage bool     `json:"botMessage"`
				NgTypes    []int    `json:"ngTypes"`
				NgUsers    []string `json:"ngUsers"`
				NgRoles    []string `json:"ngRoles"`
			}{
				{ChannelID: "222222222", Ng: true, NgTypes: []int{}, NgUsers: []string{}, NgRoles: []string{}},
			},
		})
		assert.NoError(t, err)
		h := &LinePostDiscordChannelHandler{
			indexService: newIndexService(t),
			repo: &repository.RepositoryFuncMock{
				UpdateLinePostDiscordChannelFunc: func(ctx context.Context, lineChannel repository.LinePostDiscordChannelAllColumns) error {
					t.Fatal("UpdateLinePostDiscordChannel should not be called")
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/line_post_discord_channel", bytes.NewReader(otherGuildBodyJson))
		r.SetPathValue("guildId", "987654321")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func newIndexService(t *testing.T) *service.IndexService {
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
