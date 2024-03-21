package permission

import (
	//"context"
	//"log/slog"
	"net/http"

	//"github.com/maguro-alternative/remake_bot/web/config"

	"github.com/maguro-alternative/remake_bot/web/service"
	//"github.com/maguro-alternative/remake_bot/web/shared/permission"
)

type PermissionViewHandler struct {
	IndexService *service.IndexService
}

func NewPermissionViewHandler(indexService *service.IndexService) *PermissionViewHandler {
	return &PermissionViewHandler{
		IndexService: indexService,
	}
}

func (h *PermissionViewHandler) Index(w http.ResponseWriter, r *http.Request) {}