package permission

import (
	"context"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/permission/internal"
	"github.com/maguro-alternative/remake_bot/web/shared/session/getoauth"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

func CheckDiscordPermission(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	indexService *service.IndexService,
	guild *discordgo.Guild,
	permissionType string,
) (statusCode int, permission int64, discordUserSession *model.DiscordUser, err error) {
	var userPermissionCode int64
	userPermissionCode = 0
	client := &http.Client{}
	repo := internal.NewRepository(indexService.DB)

	// ログインユーザーの取得
	discordLoginUser, err := getoauth.GetDiscordOAuth(
		ctx,
		indexService.CookieStore,
		r,
		config.SessionSecret(),
	)
	if err != nil {
		return http.StatusFound, 0, nil, err
	}

	// アクセストークンの検証
	req, err := http.NewRequestWithContext(ctx, "GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return http.StatusFound, 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+discordLoginUser.Token)
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusFound, 0, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return http.StatusFound, 0, nil, nil
	}
	permissionCode, err := repo.GetPermissionCode(ctx, guild.ID, permissionType)
	if err != nil {
		return http.StatusInternalServerError, userPermissionCode, nil, err
	}
	permissionIDs, err := repo.GetPermissionIDs(ctx, guild.ID, permissionType)
	if err != nil {
		return http.StatusInternalServerError, userPermissionCode, nil, err
	}
	discordGuildMember, err := indexService.DiscordSession.GuildMember(guild.ID, discordLoginUser.User.ID)
	if err != nil {
		return http.StatusInternalServerError, userPermissionCode, nil, err
	}
	guildRoles, err := indexService.DiscordSession.GuildRoles(guild.ID)
	if err != nil {
		return http.StatusInternalServerError, userPermissionCode, nil, err
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
		return http.StatusInternalServerError, userPermissionCode, nil, err
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
			return http.StatusForbidden, userPermissionCode, &discordLoginUser.User, nil
		}
	}
	return 200, memberPermission | userPermissionCode, &discordLoginUser.User, nil
}
