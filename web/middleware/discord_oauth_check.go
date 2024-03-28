package middleware

import (
	"context"
	"encoding/gob"
	"net/http"
	"strings"
	"log/slog"

	"github.com/maguro-alternative/remake_bot/pkg/ctxvalue"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/getoauth"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

type Repository interface {
	GetPermissionCode(ctx context.Context, guildID string, permissionType string) (int64, error)
	GetPermissionIDs(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionID, error)
	GetAllColumnsLineBot(ctx context.Context, guildID string) (repository.LineBot, error)
	GetLineBotNotClient(ctx context.Context, guildID string) (repository.LineBotNotClient, error)
	GetLineBotIvNotClient(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error)
	InsertLineBot(ctx context.Context, lineBot *repository.LineBot) error
	InsertLineBotIv(ctx context.Context, guildId string) error
}

var (
	_ Repository = (*repository.Repository)(nil)
)

func init() {
	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	gob.Register(&model.DiscordUser{})
}


func DiscordOAuthCheckMiddleware(
	indexService service.IndexService,
	repo Repository,
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
			client := &http.Client{}
			guildId := r.PathValue("guildId")
			oauthStore := getoauth.NewOAuthStore(indexService.CookieStore, config.SessionSecret())

			discordLoginUser, err := oauthStore.GetDiscordOAuth(ctx, r)
			if err != nil {
				http.Redirect(w, r, "/login/discord", http.StatusFound)
				return
			}
			req, err := http.NewRequestWithContext(ctx, "GET", "https://discord.com/api/users/@me", nil)
			if err != nil {
				http.Error(w, "Not get user", http.StatusInternalServerError)
				return
			}
			req.Header.Set("Authorization", "Bearer "+discordLoginUser.Token)
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode != http.StatusOK {
				http.Redirect(w, r, "/login/discord", http.StatusFound)
				return
			}
			defer resp.Body.Close()

			ctx = ctxvalue.ContextWithDiscordUser(ctx, &discordLoginUser.User)
			// 特定の設定ページ以外はアクセスを許可
			if len(pathParts) > 0 && len(pathParts) < 3 || (pathParts[0] != "api" && pathParts[0] != "guild") {
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			permissionType := pathParts[2]
			switch permissionType {
			case "linetoken":
				permissionType = "line_bot"
			case "line-post-discord-channel":
				permissionType = "line_post_discord_channel"
			default:
				slog.InfoContext(ctx, "権限チャンネル以外", "permissionType", permissionType)
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			permissionCode, err := repo.GetPermissionCode(ctx, guildId, permissionType)
			if err != nil {
				return
			}
			permissionIDs, err := repo.GetPermissionIDs(ctx, guildId, permissionType)
			if err != nil {
				return
			}

			guild, err := indexService.DiscordBotState.Guild(guildId)
			if err != nil {
				return
			}
			member, err := indexService.DiscordBotState.Member(guildId, discordLoginUser.User.ID)
			if err != nil {
				return
			}
			for _, role := range member.Roles {
				for _, guildRole := range guild.Roles {
					if role == guildRole.ID {
						userPermissionCode |= guildRole.Permissions
					}
				}
			}
			memberPermission, err := indexService.DiscordSession.UserChannelPermissions(discordLoginUser.User.ID, guild.Channels[0].ID)
			if err != nil {
				return
			}
			// 設定ページの場合所属していればアクセスを許可
			permissionData.User = discordLoginUser.User
			permissionData.PermissionCode = memberPermission | userPermissionCode
			if (permissionCode&permissionData.PermissionCode) == 0 {
				permissionFlag := false
				for _, permissionId := range permissionIDs {
					if permissionId.TargetType == "user" && permissionId.TargetID == discordLoginUser.User.ID {
						permissionFlag = true
						permissionData.Permission = permissionId.Permission
						break
					}
					if permissionId.TargetType == "role" && member.Roles != nil {
						for _, role := range member.Roles {
							if permissionId.TargetID == role {
								permissionFlag = true
								permissionData.Permission = permissionId.Permission
								break
							}
						}
					}
				}
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
