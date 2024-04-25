package permission

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/maguro-alternative/remake_bot/web/handler/api/permission/internal"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestPermissionHandler_ServeHTTP(t *testing.T) {
	bodyJson, err := json.Marshal(internal.PermissionJson{
		PermissionUserIDs: []internal.PermissionUserID{
			{
				GuildID:    "987654321",
				Type:       "line_bot",
				UserID:     "123456789",
				Permission: "all",
			},
		},
		PermissionRoleIDs: []internal.PermissionRoleID{
			{
				GuildID:    "987654321",
				Type:       "line_bot",
				RoleID:     "987654321",
				Permission: "all",
			},
		},
		PermissionCodes: []internal.PermissionCode{
			{
				GuildID: "987654321",
				Type:    "line_bot",
				Code:    int64(8),
			},
		},
	})
	assert.NoError(t, err)

	t.Run("MethodがPOST以外の場合、Method Not Allowedが返ること", func(t *testing.T) {
		h := &PermissionHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/987654321/permission", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("パーミッションの更新が成功すること", func(t *testing.T) {
		h := &PermissionHandler{
			IndexService: &service.IndexService{
				DiscordSession: &mock.SessionMock{
					GuildFunc: func(guildID string, options ...discordgo.RequestOption) (*discordgo.Guild, error) {
						return &discordgo.Guild{
							ID: "987654321",
						}, nil
					},
					GuildChannelsFunc: func(guildID string, options ...discordgo.RequestOption) ([]*discordgo.Channel, error) {
						return []*discordgo.Channel{
							{
								ID: "123456789",
							},
						}, nil
					},
					GuildMemberFunc: func(guildID string, userID string, options ...discordgo.RequestOption) (*discordgo.Member, error) {
						return &discordgo.Member{
							User: &discordgo.User{
								ID: "123456789",
							},
						}, nil
					},
					GuildRolesFunc: func(guildID string, options ...discordgo.RequestOption) ([]*discordgo.Role, error) {
						return []*discordgo.Role{}, nil
					},
					UserChannelPermissionsFunc: func(userID, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error) {
						return 8, nil
					},
				},
			},
			Repo: &repository.RepositoryFuncMock{
				UpdatePermissionCodesFunc: func(ctx context.Context, permissionsCode []repository.PermissionCode) error {
					return nil
				},
				DeletePermissionUserIDsFunc: func(ctx context.Context, guildId string) error {
					return nil
				},
				InsertPermissionUserIDsFunc: func(ctx context.Context, permissionsID []repository.PermissionUserIDAllColumns) error {
					return nil
				},
				DeletePermissionRoleIDsFunc: func(ctx context.Context, guildId string) error {
					return nil
				},
				InsertPermissionRoleIDsFunc: func(ctx context.Context, permissionsID []repository.PermissionRoleIDAllColumns) error {
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/permission", bytes.NewBuffer(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
