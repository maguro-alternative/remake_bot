package permission

import (
	//"context"
	//"log/slog"
	"net/http"

	//"github.com/maguro-alternative/remake_bot/web/config"

	"github.com/maguro-alternative/remake_bot/web/service"
	//"github.com/maguro-alternative/remake_bot/web/shared/permission"
)

type PermissionHandler struct {
	IndexService *service.IndexService
}

func NewPermissionHandler(indexService *service.IndexService) *PermissionHandler {
	return &PermissionHandler{
		IndexService: indexService,
	}
}

func (h *PermissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

