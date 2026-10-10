package vcsignal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/vc_signal/internal"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestVcSignalHandler_ServeHTTP(t *testing.T) {
	vcSignal := internal.VcSignalJson{
		VcSignals: []internal.VcSignal{
			{
				VcChannelID:     "987654321",
				SendSignal:      true,
				SendChannelId:   "987654321",
				JoinBot:         true,
				EveryoneMention: true,
			},
		},
	}

	t.Run("MethodがPOST以外の場合、Method Not Allowedが返ること", func(t *testing.T) {
		h := &VcSignalHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/987654321/vc-signal", nil)
		r.SetPathValue("guildId", "987654321")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("VCシグナルの更新が成功すること", func(t *testing.T) {
		bodyJson, err := json.Marshal(vcSignal)
		assert.NoError(t, err)
		h := &VcSignalHandler{
			indexService: newIndexService(t),
			repo: &repository.RepositoryFuncMock{
				UpdateVcSignalChannelFunc: func(ctx context.Context, vcSignalChannelNotGuildID repository.VcSignalChannelNotGuildID) error {
					return nil
				},
				InsertVcSignalNgUserFunc: func(ctx context.Context, vcChannelID, guildID, userID string) error {
					return nil
				},
				InsertVcSignalNgRoleFunc: func(ctx context.Context, vcChannelID, guildID, roleID string) error {
					return nil
				},
				InsertVcSignalMentionUserFunc: func(ctx context.Context, vcChannelID, guildID, userID string) error {
					return nil
				},
				InsertVcSignalMentionRoleFunc: func(ctx context.Context, vcChannelID, guildID, roleID string) error {
					return nil
				},
				DeleteVcSignalNgUsersNotInProvidedListFunc: func(ctx context.Context, vcChannelID string, userIDs []string) error {
					return nil
				},
				DeleteVcSignalNgRolesNotInProvidedListFunc: func(ctx context.Context, vcChannelID string, roleIDs []string) error {
					return nil
				},
				DeleteVcSignalMentionUsersNotInProvidedListFunc: func(ctx context.Context, vcChannelID string, userIDs []string) error {
					return nil
				},
				DeleteVcSignalMentionRolesNotInProvidedListFunc: func(ctx context.Context, vcChannelID string, roleIDs []string) error {
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/vc-signal", bytes.NewReader(bodyJson))
		r.SetPathValue("guildId", "987654321")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("jsonの読み取りに失敗した場合、Bad Requestが返ること", func(t *testing.T) {
		h := &VcSignalHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/vc-signal", bytes.NewReader([]byte("")))
		r.SetPathValue("guildId", "987654321")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("jsonのバリデーションに失敗した場合、Unprocessable Entityが返ること", func(t *testing.T) {
		h := &VcSignalHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/vc-signal", bytes.NewReader([]byte("{}")))
		r.SetPathValue("guildId", "987654321")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("VcSignalChannelの更新に失敗した場合、Internal Server Errorが返ること", func(t *testing.T) {
		bodyJson, err := json.Marshal(vcSignal)
		assert.NoError(t, err)
		h := &VcSignalHandler{
			indexService: newIndexService(t),
			repo: &repository.RepositoryFuncMock{
				UpdateVcSignalChannelFunc: func(ctx context.Context, vcSignalChannelNotGuildID repository.VcSignalChannelNotGuildID) error {
					return assert.AnError
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/vc-signal", bytes.NewReader(bodyJson))
		r.SetPathValue("guildId", "987654321")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("NgUserIDの追加に失敗した場合、Internal Server Errorが返ること", func(t *testing.T) {
		vcSignal.VcSignals[0].VcSignalNgUserIDs = []string{"123456789"}
		bodyJson, err := json.Marshal(vcSignal)
		assert.NoError(t, err)
		h := &VcSignalHandler{
			indexService: newIndexService(t),
			repo: &repository.RepositoryFuncMock{
				UpdateVcSignalChannelFunc: func(ctx context.Context, vcSignalChannelNotGuildID repository.VcSignalChannelNotGuildID) error {
					return nil
				},
				InsertVcSignalNgUserFunc: func(ctx context.Context, vcChannelID, guildID, userID string) error {
					return assert.AnError
				},
				InsertVcSignalNgRoleFunc: func(ctx context.Context, vcChannelID, guildID, roleID string) error {
					return nil
				},
				InsertVcSignalMentionUserFunc: func(ctx context.Context, vcChannelID, guildID, userID string) error {
					return nil
				},
				InsertVcSignalMentionRoleFunc: func(ctx context.Context, vcChannelID, guildID, roleID string) error {
					return nil
				},
				DeleteVcSignalNgUsersNotInProvidedListFunc: func(ctx context.Context, vcChannelID string, userIDs []string) error {
					return nil
				},
				DeleteVcSignalNgRolesNotInProvidedListFunc: func(ctx context.Context, vcChannelID string, roleIDs []string) error {
					return nil
				},
				DeleteVcSignalMentionUsersNotInProvidedListFunc: func(ctx context.Context, vcChannelID string, userIDs []string) error {
					return nil
				},
				DeleteVcSignalMentionRolesNotInProvidedListFunc: func(ctx context.Context, vcChannelID string, roleIDs []string) error {
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/vc-signal", bytes.NewReader(bodyJson))
		r.SetPathValue("guildId", "987654321")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("NgUserIDの追加が成功すること", func(t *testing.T) {
		vcSignal.VcSignals[0].VcSignalNgUserIDs = []string{"123456789"}
		bodyJson, err := json.Marshal(vcSignal)
		assert.NoError(t, err)
		h := &VcSignalHandler{
			indexService: newIndexService(t),
			repo: &repository.RepositoryFuncMock{
				UpdateVcSignalChannelFunc: func(ctx context.Context, vcSignalChannelNotGuildID repository.VcSignalChannelNotGuildID) error {
					return nil
				},
				InsertVcSignalNgUserFunc: func(ctx context.Context, vcChannelID, guildID, userID string) error {
					return nil
				},
				InsertVcSignalNgRoleFunc: func(ctx context.Context, vcChannelID, guildID, roleID string) error {
					return nil
				},
				InsertVcSignalMentionUserFunc: func(ctx context.Context, vcChannelID, guildID, userID string) error {
					return nil
				},
				InsertVcSignalMentionRoleFunc: func(ctx context.Context, vcChannelID, guildID, roleID string) error {
					return nil
				},
				DeleteVcSignalNgUsersNotInProvidedListFunc: func(ctx context.Context, vcChannelID string, userIDs []string) error {
					return nil
				},
				DeleteVcSignalNgRolesNotInProvidedListFunc: func(ctx context.Context, vcChannelID string, roleIDs []string) error {
					return nil
				},
				DeleteVcSignalMentionUsersNotInProvidedListFunc: func(ctx context.Context, vcChannelID string, userIDs []string) error {
					return nil
				},
				DeleteVcSignalMentionRolesNotInProvidedListFunc: func(ctx context.Context, vcChannelID string, roleIDs []string) error {
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/vc-signal", bytes.NewReader(bodyJson))
		r.SetPathValue("guildId", "987654321")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("他のサーバーのチャンネルを指定すると、Forbiddenが返ること", func(t *testing.T) {
		otherGuildBodyJson, err := json.Marshal(internal.VcSignalJson{
			VcSignals: []internal.VcSignal{
				{VcChannelID: "222222222", SendSignal: true, SendChannelId: "987654321"},
			},
		})
		assert.NoError(t, err)
		h := &VcSignalHandler{
			indexService: newIndexService(t),
			repo: &repository.RepositoryFuncMock{
				UpdateVcSignalChannelFunc: func(ctx context.Context, vcChannel repository.VcSignalChannelNotGuildID) error {
					t.Fatal("UpdateVcSignalChannel should not be called")
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/vc-signal", bytes.NewReader(otherGuildBodyJson))
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
			{ID: "987654321", GuildID: "987654321", Type: discordgo.ChannelTypeGuildVoice},
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
