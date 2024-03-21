package components

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func CreatePermissionCodeForm(guildID string, permissionCode PermissionCode) string {
	return fmt.Sprintf(`
	<h3>%s</h3>
	<h6>編集を許可する権限コード</h6>
	<input type="number" name="%s_permission_code" value="%d" min="0" max="1099511627775"/>
	`, permissionCode.Type, permissionCode.Type, permissionCode.Code)
}

func CreatePermissionSelectForm(guild *discordgo.Guild, permissionIDs []PermissionID, permission string) string {
	selectMemberFormBuilder := strings.Builder{}
	for _, member := range guild.Members {
		selectedFlag := false
		for _, permissionID := range permissionIDs {
			if permissionID.TargetID == member.User.ID && permissionID.TargetType == "user" && permissionID.Permission == permission {
				selectedFlag = true
				break
			}
		}
		if selectedFlag {
			selectMemberFormBuilder.WriteString(fmt.Sprintf(`<option value="%s" selected>%s</option>`, member.User.ID, member.User.Username))
			continue
		}
		selectMemberFormBuilder.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, member.User.ID, member.User.Username))
	}
	selectRoleFormBuilder := strings.Builder{}
	for _, role := range guild.Roles {
		selectedFlag := false
		for _, permissionID := range permissionIDs {
			if permissionID.TargetID == role.ID && permissionID.TargetType == "role" && permissionID.Permission == permission {
				selectedFlag = true
				break
			}
		}
		if selectedFlag {
			selectRoleFormBuilder.WriteString(fmt.Sprintf(`<option value="%s" selected>%s</option>`, role.ID, role.Name))
			continue
		}
		selectRoleFormBuilder.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, role.ID, role.Name))
	}
	return fmt.Sprintf(`
	<h6>編集を許可するメンバー</h6>
	<select name="%s_member_permission_id" multiple>%s</select>
	<h6>編集を許可するロール</h6>
	<select name="%s_role_permission_id" multiple>%s</select>
	<br/><br/>
	`, permission, selectMemberFormBuilder.String(), permission, selectRoleFormBuilder.String())
}
