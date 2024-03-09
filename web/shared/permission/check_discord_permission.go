package permission

import (
	"context"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/getoauth"
	"github.com/maguro-alternative/remake_bot/web/shared/permission/internal"
)

func CheckDiscordPermission(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	indexService *service.IndexService,
	guild *discordgo.Guild,
	permissionType string,
) (int, error) {
	var userPermissionCode int64
	userPermissionCode = 0
	repo := internal.NewRepository(indexService.DB)

	// ログインユーザーの取得
	discordLoginUser, err := getoauth.GetDiscordOAuth(
		ctx,
		indexService.CookieStore,
		r,
		config.SessionSecret(),
	)
	if err != nil {
		return http.StatusFound, err
	}
	permissionCode, err := repo.GetPermissionCode(ctx, guild.ID, permissionType)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	permissionIDs, err := repo.GetPermissionIDs(ctx, guild.ID, permissionType)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	discordGuildMember, err := indexService.DiscordSession.GuildMember(guild.ID, discordLoginUser.User.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	guildRoles, err := indexService.DiscordSession.GuildRoles(guild.ID)
	if err != nil {
		return http.StatusInternalServerError, err
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
	memberPermission, err := indexService.DiscordSession.UserChannelPermissions(discordLoginUser.User.ID, guild.Channels[0].ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// 権限のチェック
	if (permissionCode & (memberPermission | userPermissionCode)) == 0 {
		permissionFlag := false
		for _, permissionId := range permissionIDs {
			if permissionId.TargetType == "user" && permissionId.TargetID == discordLoginUser.User.ID {
				permissionFlag = true
				break
			}
			if permissionId.TargetType == "role" && discordGuildMember.Roles != nil {
				for _, role := range discordGuildMember.Roles {
					if permissionId.TargetID == role {
						permissionFlag = true
						break
					}
				}
			}
		}
		if !permissionFlag {
			return http.StatusForbidden, nil
		}
	}
	return 200, nil
}
