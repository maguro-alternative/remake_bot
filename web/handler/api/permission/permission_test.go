package permission

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/permission/internal"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestPermissionHandler_ServeHTTP(t *testing.T) {
	bodyJson, err := json.Marshal(internal.PermissionJson{
		PermissionIDs: []internal.PermissionID{
			{
				GuildID:    "987654321",
				Type:       "line_bot",
				TargetType: "user",
				TargetID:   "123456789",
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
				DiscordSession: &SessionMock{
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
				DeletePermissionIDsFunc: func(ctx context.Context, guildId string) error {
					return nil
				},
				InsertPermissionIDsFunc: func(ctx context.Context, permissionsID []repository.PermissionIDAllColumns) error {
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
