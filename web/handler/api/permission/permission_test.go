package permission

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/web/config"
	//"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/handler/api/permission/internal"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestPermissionHandler_ServeHTTP(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewPermissionsID(ctx, func(pi *fixtures.PermissionsID) {
			pi.GuildID = "987654321"
			pi.Type = "line_bot"
			pi.TargetType = "user"
			pi.TargetID = "123456789"
			pi.Permission = "all"
		}),
	)

	t.Run("MethodがPOST以外の場合、Method Not Allowedが返ること", func(t *testing.T) {
		h := &PermissionHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/987654321/permission", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("パーミッションの更新が成功すること", func(t *testing.T) {
		h := &PermissionHandler{
			IndexService: IndexService{
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
			repo: &RepositoryMock{
				UpdatePermissionCodesFunc: func(ctx context.Context, permissionsCode []internal.PermissionCode) error {
					return nil
				},
				DeletePermissionIDsFunc: func(ctx context.Context, guildId string) error {
					return nil
				},
				InsertPermissionIDsFunc: func(ctx context.Context, permissionsID []internal.PermissionID) error {
					return nil
				},
			},
			oauthStore: &OAuthStoreMock{
				GetDiscordOAuthFunc: func(ctx context.Context, r *http.Request) (*model.DiscordOAuthSession, error) {
					return &model.DiscordOAuthSession{
						User: model.DiscordUser{
							ID: "123456789",
						},
						Token: "token",
					}, nil
				},
				GetLineOAuthFunc: func(r *http.Request) (*model.LineOAuthSession, error) {
					return nil, nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/permission", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
