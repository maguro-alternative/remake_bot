package internal

import (
	"fmt"
	"strings"
	"strconv"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/bwmarrin/discordgo"
)

func CreateUpdateWebhookForm(
	webhook *repository.Webhook,
	webhookForm string,
	subscriptionSelectForm string,
	memberSelectForm string,
	roleSelectForm string,
	ngOrWordForm string,
	ngAndWordForm string,
	searchOrWordForm string,
	searchAndWordForm string,
	mentionOrWordForm string,
	mentionAndWordForm string,
) string {
	return `
		<details style="margin: 0 0 0 1em;">
            <summary>` + webhook.SubscriptionType + `:` + webhook.SubscriptionID + `</summary>
			<label for="updateWebhookType` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `">Webhook</label>
			<select name="updateWebhookType` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `" id="updateWebhookType` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `">
				` + webhookForm + `
			</select>
			<label for="updateSubscriptionName` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `">サービス名</label>
			<select name="updateSubscriptionName` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `" id="updateSubscriptionName` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `" />
				` + subscriptionSelectForm + `
			</select>
			<br/>
			<label for="updateSubscriptionId` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `">サービスID</label>
			<input type="text" name="updateSubscriptionId` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `" id="updateSubscriptionId` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `" value="` + webhook.SubscriptionID + `" />
			<br/>
			<br/>
			<label for="updateMemberMention` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `[]">メンションするユーザー</label>
			<select name="updateMemberMention` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `[]" id="updateMemberMention` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `[]" multiple>
				` + memberSelectForm + `
			</select>
			<br/>
			<label for="updateRoleMention` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `[]">メンションするロール</label>
			<select name="updateRoleMention` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `[]" id="updateRoleMention` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `[]" multiple>
				` + roleSelectForm + `
			</select>
			<br/>
			<div id="updateNgOrWords` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `">
				` + ngOrWordForm + `
				<br/>
				<button type="button" onclick="addWord('updateNgOr', ` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `)">NGワードOR検索追加</button>
			</div>
			<br/>
			<div id="updateNgAndWords` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `">
				` + ngAndWordForm + `
				<br/>
				<button type="button" onclick="addWord('updateNgAnd', ` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `)">NGワードAND検索追加</button>
			</div>
			<br/>
			<div id="updateSearchOrWords` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `">
				` + searchOrWordForm + `
				<br/>
				<button type="button" onclick="addWord('updateSearchOr', ` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `)">キーワードOR検索追加</button>
			</div>
			<br/>
			<div id="updateSearchAndWords` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `">
				` + searchAndWordForm + `
				<br/>
				<button type="button" onclick="addWord('updateSearchAnd', ` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `)">キーワードAND検索追加</button>
			</div>
			<br/>
			<div id="updateMentionOrWords` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `">
				` + mentionOrWordForm + `
				<br/>
				<button type="button" onclick="addWord('updateMentionOr', ` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `)">メンションOR検索追加</button>
			</div>
			<br/>
			<div id="updateMentionAndWords` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `">
				` + mentionAndWordForm + `
				<br/>
				<button type="button" onclick="addWord('updateMentionAnd', ` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `)">メンションAND検索追加</button>
			</div>
			<br/>
			<label for="updateDeleteFlag` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `">削除</label>
			<input type="checkbox" name="updateDeleteFlag` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `" id="updateDeleteFlag` + strconv.Itoa(int(*webhook.WebhookSerialID)) + `" />
		</details>
	`
}

func CreateNewWebhookForm(
	guildWebhooks []*discordgo.Webhook,
	guild *discordgo.Guild,
	subscriptionNames []string,
) string {
	return `
		<details style="margin: 0 0 0 1em;">
			<summary>Webhook新規作成</summary>
			<button type="button" onclick="addWebhook()">追加</button>
			<div id="newWebhook">
				<label for="newWebhookType1">Webhook</label>
				<select name="newWebhookType1" id="newWebhookType1">
					` + CreateWebhookSelectForm(guildWebhooks, "") + `
				</select>
				<br/>
				<label for="newSubscriptionName1">サービス名</label>
				<select name="newSubscriptionName1" id="newSubscriptionName1" />
					` + CreateSubscriptionsSelectForm(subscriptionNames, "") + `
				</select>
				<br/>
				<label for="newSubscriptionId1">サービスID</label>
				<input type="text" name="newSubscriptionId1" id="newSubscriptionId1" value="" />
				<br/>
				<label for="newMemberMention1[]">メンションするユーザー</label>
				<select name="newMemberMention1[]" id="newMemberMention1[]" multiple>
					` + CreateMemberSelectForm(guild, nil) + `
				</select>
				<br/>
				<label for="newRoleMention1[]">メンションするロール</label>
				<select name="newRoleMention1[]" id="newRoleMention1[]" multiple>
					` + CreateRoleSelectForm(guild, nil) + `
				</select>
				<br/>
				<br/>
				<div id="newNgOrWords1">
					<button type="button" onclick="addWord('newNgOr', 1)">NGワードOR検索追加</button>
				</div>
				<br/>
				<br/>
				<div id="newNgAndWords1">
					<button type="button" onclick="addWord('newNgAnd', 1)">NGワードAND検索追加</button>
				</div>
				<br/>
				<br/>
				<div id="newSearchOrWords1">
					<button type="button" onclick="addWord('newSearchOr', 1)">キーワードOR検索追加</button>
				</div>
				<br/>
				<br/>
				<div id="newSearchAndWords1">
					<button type="button" onclick="addWord('newSearchAnd', 1)">キーワードAND検索追加</button>
				</div>
				<br/>
				<br/>
				<div id="newMentionOrWords1">
					<button type="button" onclick="addWord('newMentionOr', 1)">メンションOR検索追加</button>
				</div>
				<br/>
				<br/>
				<div id="newMentionAndWords1">
					<button type="button" onclick="addWord('newMentionAnd', 1)">メンションAND検索追加</button>
				</div>
				<br/>
			</div>
		</details>
	`
}

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
	webhookSerialID int64,
	label string,
) string {
	wordFormBuilder := strings.Builder{}
	for _, word := range words {
		wordId := fmt.Sprintf("update%sWord%d[]", word.Condition, webhookSerialID)
		wordFormBuilder.WriteString(fmt.Sprintf(`
			<label for="%s">%s</label>
			<input type="text" id="%s" name="%s" value="%s">
			<button type="button" onclick="document.getElementById('%s').remove(); this.remove();">削除</button>
		`, wordId, label, wordId, wordId, word.Word, wordId))
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

