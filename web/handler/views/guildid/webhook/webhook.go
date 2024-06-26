package webhook

import (
	"context"
	//"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"strconv"

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

		ngOrWords, err := h.repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *webhook.WebhookSerialID, "ng_or")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get ng_or words: "+err.Error())
			return
		}
		ngOrWordForm := internal.CreateWordWebhookForm(ngOrWords, guildId, "NGワードOR検索(いずれかの言葉が含まれている場合、送信しない)")
		ngAndWords, err := h.repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *webhook.WebhookSerialID, "ng_and")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get ng_and words: "+err.Error())
			return
		}
		ngAndWordForm := internal.CreateWordWebhookForm(ngAndWords, guildId, "NGワードAND検索(全ての言葉が含まれている場合、送信しない)")

		searchOrWords, err := h.repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *webhook.WebhookSerialID, "search_or")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get ng_or words: "+err.Error())
			return
		}
		searchOrWordForm := internal.CreateWordWebhookForm(searchOrWords, guildId, "キーワードOR検索(いずれかの言葉が含まれている場合、送信)")
		searchAndWords, err := h.repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *webhook.WebhookSerialID, "search_and")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get ng_and words: "+err.Error())
			return
		}
		searchAndWordForm := internal.CreateWordWebhookForm(searchAndWords, guildId, "キーワードAND検索(すべての単語が含まれている場合、送信)")

		mentionOrWords, err := h.repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *webhook.WebhookSerialID, "ng_or")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get ng_or words: "+err.Error())
			return
		}
		mentionOrWordForm := internal.CreateWordWebhookForm(mentionOrWords, guildId, "メンションOR検索(いずれかの言葉が含まれている場合、メンションを付けて送信)")
		mentionAndWords, err := h.repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *webhook.WebhookSerialID, "ng_and")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get ng_and words: "+err.Error())
			return
		}
		mentionAndWordForm := internal.CreateWordWebhookForm(mentionAndWords, guildId, "メンションAND検索(すべての単語が含まれている場合、メンションを付けて送信)")

		webhookFormBuilder.WriteString(`
		<details style="margin: 0 0 0 1em;">
            <summary>` + webhook.SubscriptionType + `:` + webhook.SubscriptionID + `</summary>
			<label for="subscription_name`+strconv.Itoa(int(*webhook.WebhookSerialID))+`">サービス名</label>
			<select name="subscription_name`+strconv.Itoa(int(*webhook.WebhookSerialID))+`" id="subscription_name`+strconv.Itoa(int(*webhook.WebhookSerialID))+`" />
				`+ subscriptionSelectForm +`
			</select>
			<br/>
			<label for="subscription_id`+strconv.Itoa(int(*webhook.WebhookSerialID))+`">サービスID</label>
			<input type="text" name="subscription_id`+strconv.Itoa(int(*webhook.WebhookSerialID))+`" id="subscription_id`+strconv.Itoa(int(*webhook.WebhookSerialID))+`" value="`+webhook.SubscriptionID+`" />
			<br/>
			<select name="webhook_type`+strconv.Itoa(int(*webhook.WebhookSerialID))+`" id="webhook_type`+strconv.Itoa(int(*webhook.WebhookSerialID))+`>
				`+ webhookForm +`
			</select>
			<br/>
			<label for="member_mention`+strconv.Itoa(int(*webhook.WebhookSerialID))+`[]">メンションするユーザー</label>
			<select name="member_mention`+strconv.Itoa(int(*webhook.WebhookSerialID))+`[]" id="member_mention`+strconv.Itoa(int(*webhook.WebhookSerialID))+`[]" multiple>
				`+ memberSelectForm +`
			</select>
			<br/>
			<label for="role_mention`+strconv.Itoa(int(*webhook.WebhookSerialID))+`[]">メンションするロール</label>
			<select name="role_mention`+strconv.Itoa(int(*webhook.WebhookSerialID))+`[]" id="role_mention`+strconv.Itoa(int(*webhook.WebhookSerialID))+`[]" multiple>
				`+ roleSelectForm +`
			</select>
			<br/>
			`+ ngOrWordForm +`
			<br/>
			<button type="button" onclick="addWord('ng_or', `+strconv.Itoa(int(*webhook.WebhookSerialID))+`)">NGワードOR検索追加</button>
			<br/>
			`+ ngAndWordForm +`
			<br/>
			<button type="button" onclick="addWord('ng_and', `+strconv.Itoa(int(*webhook.WebhookSerialID))+`)">NGワードAND検索追加</button>
			<br/>
			`+ searchOrWordForm +`
			<br/>
			<button type="button" onclick="addWord('search_or', `+strconv.Itoa(int(*webhook.WebhookSerialID))+`)">キーワードOR検索追加</button>
			<br/>
			`+ searchAndWordForm +`
			<br/>
			<button type="button" onclick="addWord('search_and', `+strconv.Itoa(int(*webhook.WebhookSerialID))+`)">キーワードAND検索追加</button>
			<br/>
			`+ mentionOrWordForm +`
			<br/>
			<button type="button" onclick="addWord('mention_or', `+strconv.Itoa(int(*webhook.WebhookSerialID))+`)">メンションOR検索追加</button>
			<br/>
			`+ mentionAndWordForm +`
			<br/>
			<button type="button" onclick="addWord('mention_and', `+strconv.Itoa(int(*webhook.WebhookSerialID))+`)">メンションAND検索追加</button>
			<br/>
			<button type="button" onclick="deleteWebhook(`+strconv.Itoa(int(*webhook.WebhookSerialID))+`)">削除</button>
		</details>
		`)
	}

	webhookFormBuilder.WriteString(`
		<details style="margin: 0 0 0 1em;">
			<div id="new_webhook">
				<summary>新規Webhook追加</summary>
				<label for="new_webhook_type1">Webhook</label>
				<select name="new_webhook_type1" id="new_webhook_type1">
					`+ internal.CreateWebhookSelectForm(guildWebhooks, "") +`
				</select>
				<br/>
				<label for="new_subscription_name1">サービス名</label>
				<select name="new_subscription_name1" id="new_subscription_name1" />
					`+ internal.CreateSubscriptionsSelectForm(subscriptionNames, "") +`
				</select>
				<br/>
				<label for="new_subscription_id">サービスID</label>
				<input type="text" name="new_subscription_id" id="new_subscription_id1" value="" />
				<br/>
				<label for="new_member_mention1[]">メンションするユーザー</label>
				<select name="new_member_mention1[]" id="new_member_mention1[]" multiple>
					`+ internal.CreateMemberSelectForm(guild, nil) +`
				</select>
				<br/>
				<label for="new_role_mention1[]">メンションするロール</label>
				<select name="new_role_mention1[]" id="new_role_mention1[]" multiple>
					`+ internal.CreateRoleSelectForm(guild, nil) +`
				</select>
				<br/>
				`+ internal.CreateWordWebhookForm(nil, guildId, "NGワードOR検索(いずれかの言葉が含まれている場合、送信しない)") +`
				<br/>
				<div id="new_ng_or_words1">
					<button type="button" onclick="addWord('new_ng_or', 1)">NGワードOR検索追加</button>
				</div>
				<br/>
				`+ internal.CreateWordWebhookForm(nil, guildId, "NGワードAND検索(全ての言葉が含まれている場合、送信しない)") +`
				<br/>
				<div id="new_ng_and_words1">
					<button type="button" onclick="addWord('new_ng_and', 1)">NGワードAND検索追加</button>
				</div>
				<br/>
				`+ internal.CreateWordWebhookForm(nil, guildId, "キーワードOR検索(いずれかの言葉が含まれている場合、送信)") +`
				<br/>
				<div id="new_search_or_words1">
					<button type="button" onclick="addWord('new_search_or', 1)">キーワードOR検索追加</button>
				</div>
				<br/>
				`+ internal.CreateWordWebhookForm(nil, guildId, "キーワードAND検索(すべての単語が含まれている場合、送信)") +`
				<br/>
				<div id="new_search_and_words1">
					<button type="button" onclick="addWord('new_search_and', 1)">キーワードAND検索追加</button>
				</div>
				<br/>
				`+ internal.CreateWordWebhookForm(nil, guildId, "メンションOR検索(いずれかの言葉が含まれている場合、メンションを付けて送信)") +`
				<br/>
				<div id="new_mention_or_words1">
					<button type="button" onclick="addWord('new_mention_or', 1)">メンションOR検索追加</button>
				</div>
				<br/>
				`+ internal.CreateWordWebhookForm(nil, guildId, "メンションAND検索(すべての単語が含まれている場合、メンションを付けて送信)") +`
				<br/>
				<div id="new_mention_and_words1">
					<button type="button" onclick="addWord('new_mention_and', 1)">メンションAND検索追加</button>
				</div>
				<br/>
				<button type="button" onclick="addWebhook()">追加</button>
			</div>
		</details>
	`)

	submitTag := components.CreateSubmitTag(discordPermissionData.Permission)
	accountVer := strings.Builder{}
	accountVer.WriteString(components.CreateDiscordAccountVer(discordPermissionData.User))
	accountVer.WriteString(components.CreateLineAccountVer(lineSession.User))

	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/views/guildid/webhook.html"))
	if err := tmpl.Execute(w, &struct{
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
