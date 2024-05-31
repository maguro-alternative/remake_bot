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

	"github.com/stretchr/testify/assert"
)

func TestVcSignalHandler_ServeHTTP(t *testing.T) {
	vcSignal := internal.VcSignalJson{
		GuildID: "987654321",
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
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("VCシグナルの更新が成功すること", func(t *testing.T) {
		bodyJson, err := json.Marshal(vcSignal)
		assert.NoError(t, err)
		h := &VcSignalHandler{
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
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("jsonの読み取りに失敗した場合、Bad Requestが返ること", func(t *testing.T) {
		h := &VcSignalHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/vc-signal", bytes.NewReader([]byte("")))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("jsonのバリデーションに失敗した場合、Unprocessable Entityが返ること", func(t *testing.T) {
		h := &VcSignalHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/vc-signal", bytes.NewReader([]byte("{}")))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("VcSignalChannelの更新に失敗した場合、Internal Server Errorが返ること", func(t *testing.T) {
		bodyJson, err := json.Marshal(vcSignal)
		assert.NoError(t, err)
		h := &VcSignalHandler{
			repo: &repository.RepositoryFuncMock{
				UpdateVcSignalChannelFunc: func(ctx context.Context, vcSignalChannelNotGuildID repository.VcSignalChannelNotGuildID) error {
					return assert.AnError
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/vc-signal", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("VcSignalChannelの更新に失敗した場合、Internal Server Errorが返ること", func(t *testing.T) {
		bodyJson, err := json.Marshal(vcSignal)
		assert.NoError(t, err)
		h := &VcSignalHandler{
			repo: &repository.RepositoryFuncMock{
				UpdateVcSignalChannelFunc: func(ctx context.Context, vcSignalChannelNotGuildID repository.VcSignalChannelNotGuildID) error {
					return nil
				},
				InsertVcSignalNgUserFunc: func(ctx context.Context, vcChannelID, guildID, userID string) error {
					return assert.AnError
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/vc-signal", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("NgUserIDの追加に失敗した場合、Internal Server Errorが返ること", func(t *testing.T) {
		bodyJson, err := json.Marshal(vcSignal)
		assert.NoError(t, err)
		h := &VcSignalHandler{
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
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("NgUserIDの追加が成功すること", func(t *testing.T) {
		bodyJson, err := json.Marshal(vcSignal)
		assert.NoError(t, err)
		h := &VcSignalHandler{
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
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
