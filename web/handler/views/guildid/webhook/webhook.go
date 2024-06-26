package webhook

import (
	"context"
	//"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	//"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/web/shared/ctxvalue"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/handler/views/guildid/webhook/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

var (
	subscriptionNames = []string{"youtube", "niconico"}
)

type WebhookViewHandler struct {
	indexService *service.IndexService
	repo         repository.RepositoryFunc
}

func NewWebhookViewHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
) *WebhookViewHandler {
	return &WebhookViewHandler{
		indexService: indexService,
		repo:         repo,
	}
}

func (h *WebhookViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	//var webhookSerialIDs []int64
	var webhookFormBuilder strings.Builder
	var mentionIds []string
	guildId := r.PathValue("guildId")
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	guild, err := h.indexService.DiscordBotState.Guild(guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Discordサーバーの読み取りに失敗しました:"+err.Error())
		return
	}

	if guild.Members == nil {
		guild.Members, err = h.indexService.DiscordSession.GuildMembers(guildId, "", 1000)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get guild members: "+err.Error())
			return
		}
	}

	if guild.Roles == nil {
		guild.Roles, err = h.indexService.DiscordSession.GuildRoles(guildId)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get guild roles: "+err.Error())
			return
		}
	}

	discordPermissionData, err := ctxvalue.DiscordPermissionFromContext(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Discord認証情報の取得に失敗しました: ", "エラーメッセージ:", err.Error())
		return
	}
	// Lineの認証情報なしでもアクセス可能なためエラーレスポンスは出さない
	lineSession, err := ctxvalue.LineUserFromContext(ctx)
	if err != nil {
		lineSession = &model.LineOAuthSession{}
	}

	guildWebhooks, err := h.indexService.DiscordSession.GuildWebhooks(guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Not get guild webhooks: "+err.Error())
		return
	}

	webhooks, err := h.repo.GetAllColumnsWebhooksByGuildID(ctx, guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Not get webhooks: "+err.Error())
		return
	}

	for _, webhook := range webhooks {
		webhookForm := internal.CreateWebhookSelectForm(guildWebhooks, webhook.WebhookID)

		subscriptionSelectForm := internal.CreateSubscriptionsSelectForm(subscriptionNames, webhook.SubscriptionType)

		userMentions, err := h.repo.GetWebhookUserMentionWithWebhookSerialID(ctx, *webhook.WebhookSerialID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get user mention: "+err.Error())
			return
		}
		for _, userMention := range userMentions {
			mentionIds = append(mentionIds, userMention.UserID)
		}
		memberSelectForm := internal.CreateMemberSelectForm(guild, mentionIds)
		mentionIds = nil

		roleMentions, err := h.repo.GetWebhookRoleMentionWithWebhookSerialID(ctx, *webhook.WebhookSerialID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get role mention: "+err.Error())
			return
		}
		for _, roleMention := range roleMentions {
			mentionIds = append(mentionIds, roleMention.RoleID)
		}
		roleSelectForm := internal.CreateRoleSelectForm(guild, mentionIds)
		mentionIds = nil

		ngOrWords, err := h.repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *webhook.WebhookSerialID, "NgOr")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get NgOr words: "+err.Error())
			return
		}
		ngOrWordForm := internal.CreateWordWebhookForm(ngOrWords, *webhook.WebhookSerialID, "NGワードOR検索(いずれかの言葉が含まれている場合、送信しない)")
		ngAndWords, err := h.repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *webhook.WebhookSerialID, "NgAnd")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get NgAnd words: "+err.Error())
			return
		}
		ngAndWordForm := internal.CreateWordWebhookForm(ngAndWords, *webhook.WebhookSerialID, "NGワードAND検索(全ての言葉が含まれている場合、送信しない)")

		searchOrWords, err := h.repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *webhook.WebhookSerialID, "SearchOr")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get NgOr words: "+err.Error())
			return
		}
		searchOrWordForm := internal.CreateWordWebhookForm(searchOrWords, *webhook.WebhookSerialID, "キーワードOR検索(いずれかの言葉が含まれている場合、送信)")
		searchAndWords, err := h.repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *webhook.WebhookSerialID, "SearchAnd")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get NgAnd words: "+err.Error())
			return
		}
		searchAndWordForm := internal.CreateWordWebhookForm(searchAndWords, *webhook.WebhookSerialID, "キーワードAND検索(すべての単語が含まれている場合、送信)")

		mentionOrWords, err := h.repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *webhook.WebhookSerialID, "MentionOr")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get NgOr words: "+err.Error())
			return
		}
		mentionOrWordForm := internal.CreateWordWebhookForm(mentionOrWords, *webhook.WebhookSerialID, "メンションOR検索(いずれかの言葉が含まれている場合、メンションを付けて送信)")
		mentionAndWords, err := h.repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *webhook.WebhookSerialID, "MentionAnd")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get NgAnd words: "+err.Error())
			return
		}
		mentionAndWordForm := internal.CreateWordWebhookForm(mentionAndWords, *webhook.WebhookSerialID, "メンションAND検索(すべての単語が含まれている場合、メンションを付けて送信)")

		webhookFormBuilder.WriteString(`
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
		`)
	}

	webhookFormBuilder.WriteString(`
		<details style="margin: 0 0 0 1em;">
			<summary>Webhook新規作成</summary>
			<button type="button" onclick="addWebhook()">追加</button>
			<div id="newWebhook">
				<label for="newWebhookType1">Webhook</label>
				<select name="newWebhookType1" id="newWebhookType1">
					` + internal.CreateWebhookSelectForm(guildWebhooks, "") + `
				</select>
				<br/>
				<label for="newSubscriptionName1">サービス名</label>
				<select name="newSubscriptionName1" id="newSubscriptionName1" />
					` + internal.CreateSubscriptionsSelectForm(subscriptionNames, "") + `
				</select>
				<br/>
				<label for="newSubscriptionId1">サービスID</label>
				<input type="text" name="newSubscriptionId1" id="newSubscriptionId1" value="" />
				<br/>
				<label for="newMemberMention1[]">メンションするユーザー</label>
				<select name="newMemberMention1[]" id="newMemberMention1[]" multiple>
					` + internal.CreateMemberSelectForm(guild, nil) + `
				</select>
				<br/>
				<label for="newRoleMention1[]">メンションするロール</label>
				<select name="newRoleMention1[]" id="newRoleMention1[]" multiple>
					` + internal.CreateRoleSelectForm(guild, nil) + `
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
	`)

	submitTag := components.CreateSubmitTag(discordPermissionData.Permission)
	accountVer := strings.Builder{}
	accountVer.WriteString(components.CreateDiscordAccountVer(discordPermissionData.User))
	accountVer.WriteString(components.CreateLineAccountVer(lineSession.User))

	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/views/guildid/webhook.html"))
	if err := tmpl.Execute(w, &struct {
		Title        string
		AccountVer   template.HTML
		JsScriptTag  template.HTML
		SubmitTag    template.HTML
		GuildName    string
		GuildIconUrl string
		GuildID      string
		HTMLForm     template.HTML
	}{
		Title:        "Webhook設定",
		AccountVer:   template.HTML(accountVer.String()),
		JsScriptTag:  template.HTML(`<script src="/static/js/webhook.js"></script>`),
		SubmitTag:    template.HTML(submitTag),
		GuildName:    guild.Name,
		GuildIconUrl: guild.IconURL("64"),
		GuildID:      guild.ID,
		HTMLForm:     template.HTML(webhookFormBuilder.String()),
	}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Failed to execute template: "+err.Error())
		return
	}
}
