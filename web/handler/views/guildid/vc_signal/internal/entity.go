package internal

import (
	"fmt"
	"strings"

	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/shared/htmlutil"

	"github.com/bwmarrin/discordgo"
)

type VcChannelSet struct {
	ID              string
	Name            string
	SendSignal      bool
	SendChannelID   string
	JoinBot         bool
	EveryoneMention bool
	NgUsers         []string
	NgRoles         []string
	MentionUsers    []string
	MentionRoles    []string
}

func CreateVcSignalForm(
	categoryIDTmps []string,
	vcChannelSets map[string][]VcChannelSet,
	channelsInCategory map[string][]components.DiscordChannelSelect,
	categoryPositions map[string]components.DiscordChannel,
	guild *discordgo.Guild,
) string {
	htmlFormBuilder := strings.Builder{}
	categoryComponentBuilders := make([]strings.Builder, len(categoryIDTmps)+1)
	var categoryIndex int
	for categoryID, vcChannels := range vcChannelSets {
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
		for _, channel := range vcChannels {
			var sendSignalFlag, joinBotFlag, everyoneMentionFlag string
			if channel.ID == "" {
				continue
			}
			if channel.SendSignal {
				sendSignalFlag = "checked"
			}
			if channel.JoinBot {
				joinBotFlag = "checked"
			}
			if channel.EveryoneMention {
				everyoneMentionFlag = "checked"
			}
			selectNgMemberForm := createMemberSelectForm(guild, channel.NgUsers)
			selectNgRoleForm := createRoleSelectForm(guild, channel.NgRoles)
			selectMentionMemberForm := createMemberSelectForm(guild, channel.MentionUsers)
			selectMentionRoleForm := createRoleSelectForm(guild, channel.MentionRoles)
			htmlSelectChannels := components.CreateSelectChennelOptions(
				categoryIDTmps,
				channel.SendChannelID,
				channelsInCategory,
				categoryPositions,
			)
			escapedName := htmlutil.EscapeString(channel.Name)
			escapedID := htmlutil.EscapeString(channel.ID)
			categoryComponentBuilders[categoryIndex].WriteString(`
			<details style="margin: 0 0 0 1em;">
                <summary>` + escapedName + `</summary>
				<div style="margin: 0 0 0 1em;">
					<label for="sendSignal` + escapedID + `">通知を送信する</label>
					<input type="checkbox" id="sendSignal` + escapedID + `" name="sendSignal` + escapedID + `" ` + sendSignalFlag + ` />
					<br/>
					<label for="joinBot` + escapedID + `">Botの入退出を通知する</label>
					<input type="checkbox" id="joinBot` + escapedID + `" name="joinBot` + escapedID + `"` + joinBotFlag + ` />
					<br/>
					<label for="everyoneMention` + escapedID + `">通知に@everyoneメンションをつける</label>
					<input type="checkbox" id="everyoneMention` + escapedID + `" name="everyoneMention` + escapedID + `"` + everyoneMentionFlag + ` />
					<br/>
					<label for="defaultChannel` + escapedID + `">送信先チャンネル</label><br/>
					<select id="defaultChannel` + escapedID + `" name="defaultChannelId` + escapedID + `" >
						` + htmlSelectChannels + `
					</select>
					<br/>
					<label for="vcSignalNgUserIds` + escapedID + `[]">NGユーザー</label><br/>
					<select id="ng_users` + escapedID + `[]" name="vcSignalNgUserIds` + escapedID + `[]" multiple>
						` + selectNgMemberForm + `
					</select>
					<br/>
					<label for="vcSignalNgRoleIds` + escapedID + `[]">NGロール</label><br/>
					<select id="ng_roles` + escapedID + `[]" name="vcSignalNgRoleIds` + escapedID + `[]" multiple>
						` + selectNgRoleForm + `
					</select>
					<br/>
					<label for="vcSignalMentionUserIds` + escapedID + `[]">メンションユーザー</label><br/>
					<select id="mention_users` + escapedID + `[]" name="vcSignalMentionUserIds` + escapedID + `[]" multiple>
						` + selectMentionMemberForm + `
					</select>
					<br/>
					<label for="vcSignalMentionRoleIds` + escapedID + `[]">メンションロール</label><br/>
					<select id="mention_roles` + escapedID + `[]" name="vcSignalMentionRoleIds` + escapedID + `[]" multiple>
						` + selectMentionRoleForm + `
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

func createMemberSelectForm(guild *discordgo.Guild, users []string) (string) {
	selectMemberFormBuilder := strings.Builder{}
	for _, member := range guild.Members {
		selectedFlag := false
		for _, ngUserID := range users {
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
	return selectMemberFormBuilder.String()
}

func createRoleSelectForm(guild *discordgo.Guild, roles []string) (string) {
	selectRoleFormBuilder := strings.Builder{}
	for _, role := range guild.Roles {
		selectedFlag := false
		for _, ngRoleID := range roles {
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
	return selectRoleFormBuilder.String()
}
