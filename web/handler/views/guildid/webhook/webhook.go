package webhook

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

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
		webhookForm := internal.CreateWebhookSelectForm(h.indexService.DiscordSession, h.indexService.DiscordBotState, guildWebhooks, webhook.WebhookID)

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

		webhookFormBuilder.WriteString(internal.CreateUpdateWebhookForm(
			webhook,
			webhookForm,
			subscriptionSelectForm,
			memberSelectForm,
			roleSelectForm,
			ngOrWordForm,
			ngAndWordForm,
			searchOrWordForm,
			searchAndWordForm,
			mentionOrWordForm,
			mentionAndWordForm,
		))
	}

	webhookFormBuilder.WriteString(internal.CreateNewWebhookForm(
		h.indexService.DiscordSession,
		h.indexService.DiscordBotState,
		guildWebhooks,
		guild,
		subscriptionNames,
	))

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
