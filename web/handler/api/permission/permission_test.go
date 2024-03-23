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
}
