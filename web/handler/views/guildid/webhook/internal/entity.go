package internal

import (
	"fmt"
	"strings"

	//"github.com/maguro-alternative/remake_bot/web/components"

	"github.com/bwmarrin/discordgo"
)



func createNewWebhookSelectForm(
	guildWebhooks []*discordgo.Webhook,
) (string) {
	selectWebhookFormBuilder := strings.Builder{}
	for _, guildWebhook := range guildWebhooks {
		selectWebhookFormBuilder.WriteString(fmt.Sprintf(`
		<option value="%s">%s</option>`,
		guildWebhook.ID, guildWebhook.Name))
	}
	return selectWebhookFormBuilder.String()
}

func createWebhookSelectForm(
	guildWebhooks []*discordgo.Webhook,
	selectedWebhookID string,
) (string) {
	selectWebhookFormBuilder := strings.Builder{}
	for _, guildWebhook := range guildWebhooks {
		if guildWebhook.ID == selectedWebhookID {
			selectWebhookFormBuilder.WriteString(fmt.Sprintf(`
			<option value="%s" selected>%s</option>`,
			guildWebhook.ID, guildWebhook.Name))
			continue
		}
		selectWebhookFormBuilder.WriteString(fmt.Sprintf(`
		<option value="%s">%s</option>`,
		guildWebhook.ID, guildWebhook.Name))
	}
	return selectWebhookFormBuilder.String()
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

