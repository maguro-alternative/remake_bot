package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/maguro-alternative/remake_bot/web/shared/ctxvalue"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
	"github.com/maguro-alternative/remake_bot/web/shared/session"

	"github.com/bwmarrin/discordgo"
)

type Repository interface {
	GetPermissionCodeByGuildIDAndType(ctx context.Context, guildID string, permissionType string) (int64, error)
	GetPermissionUserIDsByGuildIDAndType(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionUserID, error)
	GetPermissionRoleIDsByGuildIDAndType(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionRoleID, error)
	GetAllColumnsLineBotByGuildID(ctx context.Context, guildID string) (repository.LineBot, error)
	GetLineBotNotClientByGuildID(ctx context.Context, guildID string) (repository.LineBotNotClient, error)
	GetLineBotIvNotClientByGuildID(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error)
	InsertLineBot(ctx context.Context, lineBot *repository.LineBot) error
	InsertLineBotIvByGuildID(ctx context.Context, guildId string) error
}

var (
	_ Repository = (*repository.Repository)(nil)
)

func DiscordOAuthCheckMiddleware(
	indexService service.IndexService,
	repo Repository,
	loginRequiredFlag bool,
) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userPermissionCode int64
			userPermissionCode = 0
			permissionData := &model.DiscordPermissionData{
				PermissionCode: 0,
				User:           model.DiscordUser{},
				Permission:     "",
			}
			ctx := r.Context()
			pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
			sessionStore, err := session.NewSessionStore(r, indexService.CookieStore, config.SessionSecret())
			if err != nil {
				slog.ErrorContext(r.Context(), "sessionの取得に失敗しました。", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			discordUser, err := sessionStore.GetDiscordUser()
			if err != nil && loginRequiredFlag {
				slog.WarnContext(ctx, "ログインしていないユーザーがアクセスしました。")
				http.Redirect(w, r, "/login/discord", http.StatusFound)
				return
			}
			discordOAuthToken, err := sessionStore.GetDiscordOAuthToken()
			if err != nil && loginRequiredFlag {
				slog.WarnContext(ctx, "ログインしていないユーザーがアクセスしました。")
				http.Redirect(w, r, "/login/discord", http.StatusFound)
				return
			}
			req, err := http.NewRequestWithContext(ctx, "GET", "https://discord.com/api/users/@me", nil)
			if err != nil {
				http.Error(w, "Not get user", http.StatusInternalServerError)
				return
			}

			req.Header.Set("Authorization", "Bearer "+discordOAuthToken)
			resp, err := indexService.Client.Do(req)
			if (err != nil || resp.StatusCode != http.StatusOK) && loginRequiredFlag {
				slog.WarnContext(ctx, "ユーザー情報に問題があります。", "エラー:", err, "ステータスコード:", resp.StatusCode)
				http.Redirect(w, r, "/login/discord", http.StatusFound)
				return
			}
			defer resp.Body.Close()

			discordLoginUser := &model.DiscordOAuthSession{
				User:  *discordUser,
				Token: discordOAuthToken,
			}

			ctx = ctxvalue.ContextWithDiscordUser(ctx, discordLoginUser)
			guildId := r.PathValue("guildId")
			// 特定のサーバーのページでない場合はアクセスを許可
			if guildId == "" {
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			guild, err := indexService.DiscordBotState.Guild(guildId)
			if err != nil {
				slog.WarnContext(ctx, "ギルド情報の取得に失敗しました。", "guildId", guildId)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			member, err := indexService.DiscordBotState.Member(guildId, discordLoginUser.User.ID)
			if err != nil {
				slog.WarnContext(ctx, "メンバー情報の取得に失敗しました。", "guildId", guildId, "userId", discordLoginUser.User.ID)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			userPermissionCode = getUserRolePermissionCode(member, guild)

			memberPermission, err := indexService.DiscordSession.UserChannelPermissions(discordLoginUser.User.ID, guild.Channels[0].ID)
			if err != nil {
				slog.WarnContext(ctx, "メンバー権限の取得に失敗しました。", "guildId", guildId, "userId", discordLoginUser.User.ID)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// 設定ページの場合所属していればアクセスを許可
			permissionData.User = discordLoginUser.User
			permissionData.PermissionCode = memberPermission | userPermissionCode

			ctx = ctxvalue.ContextWithDiscordPermission(ctx, permissionData)
			// 特定の設定ページ以外はアクセスを許可
			if len(pathParts) > 0 && len(pathParts) < 3 || (pathParts[0] != "api" && pathParts[0] != "guild") {
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			permissionType := pathParts[2]
			switch permissionType {
			case "linetoken":
				permissionType = "lineBot"
			case "line-post-discord-channel":
				permissionType = "linePostDiscordChannel"
			case "vc-signal":
				permissionType = "vcSignal"
			case "webhook":
				permissionType = "webhook"
			default:
				slog.InfoContext(ctx, "権限チャンネル以外", "permissionType", permissionType)
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			permissionCode, err := repo.GetPermissionCodeByGuildIDAndType(ctx, guildId, permissionType)
			if err != nil {
				slog.WarnContext(ctx, "権限コードの取得に失敗しました。", "guildId", guildId, "permissionType", permissionType)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			permissionUserIDs, err := repo.GetPermissionUserIDsByGuildIDAndType(ctx, guildId, permissionType)
			if err != nil {
				slog.WarnContext(ctx, "権限IDの取得に失敗しました。", "guildId", guildId, "permissionType", permissionType)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			permissionRoleIDs, err := repo.GetPermissionRoleIDsByGuildIDAndType(ctx, guildId, permissionType)
			if err != nil {
				slog.WarnContext(ctx, "権限IDの取得に失敗しました。", "guildId", guildId, "permissionType", permissionType)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if (permissionCode & permissionData.PermissionCode) == 0 {
				permissionFlag := isUserAccessPermission(
					permissionUserIDs,
					permissionRoleIDs,
					permissionData,
					discordLoginUser,
					member,
				)
				if !permissionFlag {
					http.Error(w, "Forbidden", http.StatusForbidden)
					slog.WarnContext(ctx, "権限のないアクセスがありました。")
					return
				}
			}
			if permissionData.Permission == "" {
				permissionData.Permission = "all"
			}
			slog.InfoContext(ctx, "権限チェック成功", "アクセスユーザー", permissionData.User.Username, "権限", permissionData.Permission, "権限コード", permissionData.PermissionCode)
			ctx = ctxvalue.ContextWithDiscordPermission(ctx, permissionData)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getUserRolePermissionCode(
	member *discordgo.Member,
	guild *discordgo.Guild,
) int64 {
	var userPermissionCode int64
	userPermissionCode = 0
	for _, role := range member.Roles {
		for _, guildRole := range guild.Roles {
			if role == guildRole.ID {
				userPermissionCode |= guildRole.Permissions
			}
		}
	}
	return userPermissionCode
}

func isUserAccessPermission(
	permissionUserIDs []repository.PermissionUserID,
	permissionRoleIDs []repository.PermissionRoleID,
	permissionData *model.DiscordPermissionData,
	discordLoginUser *model.DiscordOAuthSession,
	member *discordgo.Member,
) bool {
	for _, permissionUserId := range permissionUserIDs {
		if permissionUserId.UserID == discordLoginUser.User.ID {
			permissionData.Permission = permissionUserId.Permission
			return true
		}
	}
	for _, permissionRoleId := range permissionRoleIDs {
		if member.Roles != nil {
			for _, role := range member.Roles {
				if permissionRoleId.RoleID == role {
					permissionData.Permission = permissionRoleId.Permission
					return true
				}
			}
		}
	}
	return false
}
