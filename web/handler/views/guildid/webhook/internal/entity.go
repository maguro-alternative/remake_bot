package internal

import (
	"fmt"
	"strings"

	//"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/bwmarrin/discordgo"
)

func CreateSubscriptionsSelectForm(subscriptionNames []string, selectedSucscriptionName string) string {
	selectSubscriptionsFormBuilder := strings.Builder{}
	for _, subscriptionName := range subscriptionNames {
		if subscriptionName == selectedSucscriptionName {
			selectSubscriptionsFormBuilder.WriteString(fmt.Sprintf(`<option value="%s" selected>%s</option>`, subscriptionName, subscriptionName))
			continue
		}
		selectSubscriptionsFormBuilder.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, subscriptionName, subscriptionName))
	}
	return selectSubscriptionsFormBuilder.String()
}

func CreateWordWebhookForm(
	words []*repository.WebhookWord,
	guildId string,
	label string,
) string {
	wordFormBuilder := strings.Builder{}
	for i, word := range words {
		wordId := fmt.Sprintf("update%sWord%s_%d[]", word.Condition, guildId, i)
		wordFormBuilder.WriteString(fmt.Sprintf(`
			<label for="%s">%s</label>
			<input type="text" name="%s" value="%s">
		`, wordId, label, wordId, word.Word))
	}
	return wordFormBuilder.String()
}

func CreateNewWebhookSelectForm(
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

func CreateWebhookSelectForm(
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

func CreateMemberSelectForm(guild *discordgo.Guild, users []string) (string) {
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

func CreateRoleSelectForm(guild *discordgo.Guild, roles []string) (string) {
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

