package internal

import (
	"fmt"
	"strings"

	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/shared/htmlutil"

	"github.com/bwmarrin/discordgo"
)

type DiscordChannelSet struct {
	ID         string
	Name       string
	Ng         bool
	BotMessage bool
	NgTypes    []int
	NgUsers    []string
	NgRoles    []string
}

func CreateLinePostDiscordChannelForm(
	categoryIDTmps []string,
	channelsInCategory map[string][]DiscordChannelSet,
	categoryPositions map[string]components.DiscordChannel,
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
		`, htmlutil.EscapeString(categoryChannelName)))
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
			escapedName := htmlutil.EscapeString(channel.Name)
			escapedID := htmlutil.EscapeString(channel.ID)
			categoryComponentBuilders[categoryIndex].WriteString(`
			<details style="margin: 0 0 0 1em;">
                <summary>` + escapedName + `</summary>
				<div style="margin: 0 0 0 1em;">
					<label for="ng` + escapedID + `">LINEへ送信しない</label>
					<input type="checkbox" id="ng` + escapedID + `" name="ng` + escapedID + `" ` + messageNgFlag + ` />
					<br/>
					<label for="botMessage` + escapedID + `">Botのメッセージを送信する</label>
					<input type="checkbox" id="botMessage` + escapedID + `" name="botMessage` + escapedID + `"` + botNgFlag + ` />
					<br/>
					<label for="ngTypes` + escapedID + `[]">NGタイプ</label><br/>
					<select id="ngTypes` + escapedID + `[]" name="ngTypes` + escapedID + `[]" multiple>
						` + selectMessageTypeForm + `
					</select>
					<br/>
					<label for="ngUsers` + escapedID + `[]">NGユーザー</label><br/>
					<select id="ng_users` + escapedID + `[]" name="ngUsers` + escapedID + `[]" multiple>
						` + selectMemberForm + `
					</select>
					<br/>
					<label for="ngRoles` + escapedID + `[]">NGロール</label><br/>
					<select id="ngRoles` + escapedID + `[]" name="ngRoles` + escapedID + `[]" multiple>
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
		escapedID := htmlutil.EscapeString(member.User.ID)
		escapedName := htmlutil.EscapeString(member.User.Username)
		if selectedFlag {
			selectMemberFormBuilder.WriteString(fmt.Sprintf(`<option value="%s" selected>%s</option>`, escapedID, escapedName))
			continue
		}
		selectMemberFormBuilder.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, escapedID, escapedName))
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
		escapedID := htmlutil.EscapeString(role.ID)
		escapedName := htmlutil.EscapeString(role.Name)
		if selectedFlag {
			selectRoleFormBuilder.WriteString(fmt.Sprintf(`<option value="%s" selected>%s</option>`, escapedID, escapedName))
			continue
		}
		selectRoleFormBuilder.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, escapedID, escapedName))
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
		escapedType := htmlutil.EscapeString(messageType)
		if selectedFlag {
			selectMessageTypeFormBuilder.WriteString(fmt.Sprintf(`<option value=%d selected>%s</option>`, i, escapedType))
			continue
		}
		selectMessageTypeFormBuilder.WriteString(fmt.Sprintf(`<option value=%d>%s</option>`, i, escapedType))
	}
	return selectMemberFormBuilder.String(), selectRoleFormBuilder.String(), selectMessageTypeFormBuilder.String()
}

