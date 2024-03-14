package permission

import (
	"context"
	"errors"
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
) (statusCode int, discordPermissionData *model.DiscordPermissionData, err error) {
	var userPermissionCode int64
	var permissionData model.DiscordPermissionData
	permissionData.Permission = ""
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
		return http.StatusFound, nil, err
	}

	// アクセストークンの検証
	req, err := http.NewRequestWithContext(ctx, "GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return http.StatusFound, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+discordLoginUser.Token)
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusFound, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return http.StatusFound, nil, errors.New("status code is not 200")
	}
	permissionCode, err := repo.GetPermissionCode(ctx, guild.ID, permissionType)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	permissionIDs, err := repo.GetPermissionIDs(ctx, guild.ID, permissionType)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	discordGuildMember, err := indexService.DiscordSession.GuildMember(guild.ID, discordLoginUser.User.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	guildRoles, err := indexService.DiscordSession.GuildRoles(guild.ID)
	if err != nil {
		return http.StatusInternalServerError, nil, err
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
		return http.StatusInternalServerError, nil, err
	}
	// 設定ページの場合所属していればアクセスを許可
	permissionData.User = discordLoginUser.User
	permissionData.PermissionCode = memberPermission | userPermissionCode
	if permissionType == "" {
		return http.StatusOK, &permissionData, nil
	}
	// 権限のチェック
	if (permissionCode & permissionData.PermissionCode) == 0 {
		permissionFlag := false
		for _, permissionId := range permissionIDs {
			if permissionId.TargetType == "user" && permissionId.TargetID == discordLoginUser.User.ID {
				permissionFlag = true
				permissionData.Permission = permissionId.Permission
				break
			}
			if permissionId.TargetType == "role" && discordGuildMember.Roles != nil {
				for _, role := range discordGuildMember.Roles {
					if permissionId.TargetID == role {
						permissionFlag = true
						permissionData.Permission = permissionId.Permission
						break
					}
				}
			}
		}
		if !permissionFlag {
			return http.StatusForbidden, &permissionData, errors.New("permission denied")
		}
	}
	if permissionData.Permission == "" {
		permissionData.Permission = "all"
	}
	return 200, &permissionData, nil
}
