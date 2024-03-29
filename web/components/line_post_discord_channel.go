package components

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func CreateLinePostDiscordChannelForm(
	categoryIDTmps []string,
	channelsInCategory map[string][]DiscordChannelSet,
	categoryPositions map[string]DiscordChannel,
	guild *discordgo.Guild,
	messageTypes []string,
) string {
	htmlFormBuilder := strings.Builder{}
	categoryComponentBuilders := make([]strings.Builder, len(categoryIDTmps)+1)
	var categoryIndex int
	for categoryID, channels := range channelsInCategory {
		for i, categoryIDTmp := range categoryIDTmps {
			if categoryID == "" {
				categoryIndex = len(categoryIDTmps)
				break
			}
			if categoryIDTmp == categoryID {
				categoryIndex = i
				break
			}
		}
		categoryChannelName := categoryPositions[categoryID].Name
		if categoryID == "" {
			categoryChannelName = "カテゴリーなし"
		}
		categoryComponentBuilders[categoryIndex].WriteString(fmt.Sprintf(`
		<details>
            <summary>%s</summary>
		`, categoryChannelName))
		for _, channel := range channels {
			if channel.ID == "" {
				continue
			}
			messageNgFlag, botNgFlag := "", ""
			selectMemberForm, selectRoleForm, selectMessageTypeForm := createSelectForm(guild, channel, messageTypes)
			if channel.Ng {
				messageNgFlag = "checked"
			}
			if channel.BotMessage {
				botNgFlag = "checked"
			}
			categoryComponentBuilders[categoryIndex].WriteString(`
			<details style="margin: 0 0 0 1em;">
                <summary>` + channel.Name + `</summary>
				<div style="margin: 0 0 0 1em;">
					<label for="ng_` + channel.ID + `">LINEへ送信しない</label>
					<input type="checkbox" id="ng_` + channel.ID + `" name="ng_` + channel.ID + `" ` + messageNgFlag + ` />
					<br/>
					<label for="bot_message_` + channel.ID + `">Botのメッセージを送信しない</label>
					<input type="checkbox" id="bot_message_` + channel.ID + `" name="bot_message_` + channel.ID + `"` + botNgFlag + ` />
					<br/>
					<label for="ng_types_` + channel.ID + `[]">NGタイプ</label><br/>
					<select id="ng_types_` + channel.ID + `[]" name="ng_types_` + channel.ID + `[]" multiple>
						` + selectMessageTypeForm + `
					</select>
					<br/>
					<label for="ng_users_` + channel.ID + `[]">NGユーザー</label><br/>
					<select id="ng_users_` + channel.ID + `[]" name="ng_users_` + channel.ID + `[]" multiple>
						` + selectMemberForm + `
					</select>
					<br/>
					<label for="ng_roles_` + channel.ID + `[]">NGロール</label><br/>
					<select id="ng_roles_` + channel.ID + `[]" name="ng_roles_` + channel.ID + `[]" multiple>
						` + selectRoleForm + `
					</select>
					<br/>
				</div>
			</details>
			`)
		}
		categoryComponentBuilders[categoryIndex].WriteString(`
		</details>`)
	}
	for _, categoryComponent := range categoryComponentBuilders {
		htmlFormBuilder.WriteString(categoryComponent.String())
	}
	return htmlFormBuilder.String()
}

func createSelectForm(guild *discordgo.Guild, channel DiscordChannelSet, messageTypes []string) (member string, role string, messageType string) {
	selectMemberFormBuilder := strings.Builder{}
	for _, member := range guild.Members {
		selectedFlag := false
		for _, ngUserID := range channel.NgUsers {
			if ngUserID == member.User.ID {
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
		for _, ngRoleID := range channel.NgRoles {
			if ngRoleID == role.ID {
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
	selectMessageTypeFormBuilder := strings.Builder{}
	for i, messageType := range messageTypes {
		selectedFlag := false
		for _, ngType := range channel.NgTypes {
			if ngType == i {
				selectedFlag = true
				break
			}
		}
		if selectedFlag {
			selectMessageTypeFormBuilder.WriteString(fmt.Sprintf(`<option value=%d selected>%s</option>`, i, messageType))
			continue
		}
		selectMessageTypeFormBuilder.WriteString(fmt.Sprintf(`<option value=%d>%s</option>`, i, messageType))
	}
	return selectMemberFormBuilder.String(), selectRoleFormBuilder.String(), selectMessageTypeFormBuilder.String()
}

