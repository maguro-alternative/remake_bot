package internal

import (
	"fmt"
	"strings"

	"github.com/maguro-alternative/remake_bot/web/components"

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
		`, categoryChannelName))
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
			categoryComponentBuilders[categoryIndex].WriteString(`
			<details style="margin: 0 0 0 1em;">
                <summary>` + channel.Name + `</summary>
				<div style="margin: 0 0 0 1em;">
					<label for="sendSignal` + channel.ID + `">通知を送信する</label>
					<input type="checkbox" id="sendSignal` + channel.ID + `" name="sendSignal` + channel.ID + `" ` + sendSignalFlag + ` />
					<br/>
					<label for="joinBot` + channel.ID + `">Botの入退出を通知する</label>
					<input type="checkbox" id="joinBot` + channel.ID + `" name="joinBot` + channel.ID + `"` + joinBotFlag + ` />
					<br/>
					<label for="everyoneMention` + channel.ID + `">通知に@everyoneメンションをつける</label>
					<input type="checkbox" id="everyoneMention` + channel.ID + `" name="everyoneMention` + channel.ID + `"` + everyoneMentionFlag + ` />
					<br/>
					<label for="defaultChannel` + channel.ID + `">送信先チャンネル</label><br/>
					<select id="defaultChannel` + channel.ID + `" name="defaultChannelId` + channel.ID + `" >
						` + htmlSelectChannels + `
					</select>
					<br/>
					<label for="vcSignalNgUserIds` + channel.ID + `[]">NGユーザー</label><br/>
					<select id="ng_users` + channel.ID + `[]" name="vcSignalNgUserIds` + channel.ID + `[]" multiple>
						` + selectNgMemberForm + `
					</select>
					<br/>
					<label for="vcSignalNgRoleIds` + channel.ID + `[]">NGロール</label><br/>
					<select id="ng_roles` + channel.ID + `[]" name="vcSignalNgRoleIds` + channel.ID + `[]" multiple>
						` + selectNgRoleForm + `
					</select>
					<br/>
					<label for="vcSignalMentionUserIds` + channel.ID + `[]">メンションユーザー</label><br/>
					<select id="mention_users` + channel.ID + `[]" name="vcSignalMentionUsers` + channel.ID + `[]" multiple>
						` + selectMentionMemberForm + `
					</select>
					<br/>
					<label for="vcSignalMentionRoles` + channel.ID + `[]">メンションロール</label><br/>
					<select id="mention_roles` + channel.ID + `[]" name="vcSignalMentionRoles` + channel.ID + `[]" multiple>
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
		if selectedFlag {
			selectMemberFormBuilder.WriteString(fmt.Sprintf(`<option value="%s" selected>%s</option>`, member.User.ID, member.User.Username))
			continue
		}
		selectMemberFormBuilder.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, member.User.ID, member.User.Username))
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
		if selectedFlag {
			selectRoleFormBuilder.WriteString(fmt.Sprintf(`<option value="%s" selected>%s</option>`, role.ID, role.Name))
			continue
		}
		selectRoleFormBuilder.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, role.ID, role.Name))
	}
	return selectRoleFormBuilder.String()
}

func createChannelsInCategory(
	guild *discordgo.Guild,
	channel *discordgo.Channel,
	categoryPositions map[string]components.DiscordChannel,
	channelsInCategory map[string][]components.DiscordChannelSelect,
) {
	// カテゴリー、フォーラムチャンネルはスキップ
	if channel.Type == discordgo.ChannelTypeGuildForum {
		return
	}
	if channel.Type == discordgo.ChannelTypeGuildCategory {
		return
	}
	typeIcon := "🔊"
	if channel.Type == discordgo.ChannelTypeGuildText {
		typeIcon = "📝"
	}
	categoryPosition := categoryPositions[channel.ParentID]
	// まだチャンネルがない場合は初期化
	if len(channelsInCategory[categoryPosition.ID]) == 0 {
		channelsInCategory[categoryPosition.ID] = make([]components.DiscordChannelSelect, len(guild.Channels)+1)
	}
	channelsInCategory[categoryPosition.ID][channel.Position] = components.DiscordChannelSelect{
		ID:   channel.ID,
		Name: fmt.Sprintf("%s:%s:%s", categoryPosition.Name, typeIcon, channel.Name),
	}
	if categoryPosition.ID == "" {
		channelsInCategory[categoryPosition.ID][channel.Position] = components.DiscordChannelSelect{
			ID:   channel.ID,
			Name: fmt.Sprintf("カテゴリーなし:%s:%s", typeIcon, channel.Name),
		}
	}
}

