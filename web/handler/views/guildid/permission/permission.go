package permission

import (
	"context"
	"log/slog"
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

func (h *PermissionViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		slog.ErrorContext(ctx, "/api/line-bot Method Not Allowed")
		return
	}
}