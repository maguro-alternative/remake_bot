package permission

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/config"

	"github.com/maguro-alternative/remake_bot/web/handler/api/permission/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/getoauth"
)

type Repository interface {
	UpdatePermissionCodes(ctx context.Context, permissionsCode []internal.PermissionCode) error
	DeletePermissionIDs(ctx context.Context, guildId string) error
	InsertPermissionIDs(ctx context.Context, permissionsID []internal.PermissionID) error
}

type PermissionHandler struct {
	IndexService *service.IndexService
}

func NewPermissionHandler(indexService *service.IndexService) *PermissionHandler {
	return &PermissionHandler{
		IndexService: indexService,
	}
}

func (h *PermissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var userPermissionCode int64
	var repo Repository
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

	discordSession, err := getoauth.GetDiscordOAuth(
		ctx,
		h.IndexService.CookieStore,
		r,
		config.SessionSecret(),
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Discordの認証情報の取得に失敗しました。", "エラー:", err.Error())
		return
	}

	// ギルド情報を取得
	guild, err := h.IndexService.DiscordSession.State.Guild(guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "ギルド情報の取得に失敗しました。", "エラー:", err.Error())
		return
	}

	discordGuildMember, err := h.IndexService.DiscordSession.GuildMember(guild.ID, discordSession.User.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "メンバーの取得に失敗しました。", "エラー:", err.Error())
		return
	}
	guildRoles, err := h.IndexService.DiscordSession.GuildRoles(guild.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "ロールの取得に失敗しました。", "エラー:", err.Error())
		return
	}

	for _, role := range discordGuildMember.Roles {
		for _, guildRole := range guildRoles {
			if role == guildRole.ID {
				userPermissionCode |= guildRole.Permissions
			}
		}
	}

	// メンバーの権限を取得
	// discordgoの場合guildMemberから正しく権限を取得できないため、UserChannelPermissionsを使用
	memberPermission, err := h.IndexService.DiscordSession.UserChannelPermissions(discordSession.User.ID, guild.Channels[0].ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "メンバーの権限の取得に失敗しました。", "エラー:", err.Error())
		return
	}

	if ((memberPermission | userPermissionCode) & 8) != 8 {
		fmt.Println((memberPermission | userPermissionCode) & 8)
		http.Error(w, "Not permission", http.StatusForbidden)
		slog.WarnContext(ctx, "権限のないアクセスがありました。")
		return
	}

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

	// パーミッションの更新
	repo = internal.NewRepository(h.IndexService.DB)
	if err := repo.UpdatePermissionCodes(ctx, permissionJson.PermissionCodes); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの更新に失敗しました。", "エラー:", err.Error())
		return
	}

	if err := repo.DeletePermissionIDs(ctx, guildId); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの削除に失敗しました。", "エラー:", err.Error())
		return
	}

	if err := repo.InsertPermissionIDs(ctx, permissionJson.PermissionIDs); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの追加に失敗しました。", "エラー:", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
