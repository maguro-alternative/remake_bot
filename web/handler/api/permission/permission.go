package permission

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/permission/internal"
)

type PermissionHandler struct {
	repo repository.RepositoryFunc
}

func NewPermissionHandler(
	repo repository.RepositoryFunc,
) *PermissionHandler {
	return &PermissionHandler{
		repo: repo,
	}
}

func (h *PermissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	guildId := r.PathValue("guildId")
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		slog.ErrorContext(ctx, "/api/permission Method Not Allowed")
		return
	}
	var permissionJson internal.PermissionJson
	var permissionCodes []repository.PermissionCode
	var permissionUserIDs []repository.PermissionUserIDAllColumns
	var permissionRoleIDs []repository.PermissionRoleIDAllColumns

	if err := json.NewDecoder(r.Body).Decode(&permissionJson); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "jsonの読み取りに失敗しました:", "エラー:", err.Error())
		return
	}
	if err := permissionJson.Validate(); err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		slog.ErrorContext(ctx, "jsonのバリデーションに失敗しました:", "エラー:", err.Error())
		return
	}

	for _, permissionCode := range permissionJson.PermissionCodes {
		permissionCodes = append(permissionCodes, repository.PermissionCode{
			GuildID: guildId,
			Type:    permissionCode.Type,
			Code:    permissionCode.Code,
		})
	}

	for _, permissionID := range permissionJson.PermissionUserIDs {
		permissionUserIDs = append(permissionUserIDs, repository.PermissionUserIDAllColumns{
			GuildID:    guildId,
			Type:       permissionID.Type,
			UserID:     permissionID.UserID,
			Permission: permissionID.Permission,
		})
	}

	for _, permissionID := range permissionJson.PermissionRoleIDs {
		permissionRoleIDs = append(permissionRoleIDs, repository.PermissionRoleIDAllColumns{
			GuildID:    guildId,
			Type:       permissionID.Type,
			RoleID:     permissionID.RoleID,
			Permission: permissionID.Permission,
		})
	}

	if err := h.repo.UpdatePermissionCodes(ctx, permissionCodes); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの更新に失敗しました。", "エラー:", err.Error())
		return
	}

	if err := h.repo.DeletePermissionUserIDsByGuildID(ctx, guildId); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの削除に失敗しました。", "エラー:", err.Error())
		return
	}

	if err := h.repo.InsertPermissionUserIDs(ctx, permissionUserIDs); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの追加に失敗しました。", "エラー:", err.Error())
		return
	}

	if err := h.repo.DeletePermissionRoleIDsByGuildID(ctx, guildId); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの削除に失敗しました。", "エラー:", err.Error())
		return
	}

	if err := h.repo.InsertPermissionRoleIDs(ctx, permissionRoleIDs); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの追加に失敗しました。", "エラー:", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
