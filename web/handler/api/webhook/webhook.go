package webhook

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/webhook/internal"
)

type WebhookHandler struct {
	repo repository.RepositoryFunc
}

func NewWebhookHandler(
	repo repository.RepositoryFunc,
) *WebhookHandler {
	return &WebhookHandler{
		repo: repo,
	}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	now := time.Now()
	guildId := r.PathValue("guildId")
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		slog.ErrorContext(ctx, "/api/webhook Method Not Allowed")
		return
	}
	var webhookJson internal.WebhookJson

	err := json.NewDecoder(r.Body).Decode(&webhookJson)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "jsonの読み取りに失敗しました:", "エラー:", err.Error())
		return
	}
	err = webhookJson.Validate()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "jsonのバリデーションに失敗しました:", "エラー:", err.Error())
		return
	}

	for _, webhook := range webhookJson.NewWebhooks {
		if webhook.SubscriptionId == "" {
			continue
		}
		webhookIdSplit := strings.Split(webhook.WebhookID, "-")
		webhookSerialID, err := h.repo.InsertWebhook(ctx, guildId, webhookIdSplit[0], webhook.SubscriptionType, webhook.SubscriptionId, now)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Webhookの作成に失敗しました:", "エラー:", err.Error())
			return
		}
		if len(webhookIdSplit) == 2 {
			err = h.repo.InsertWebhookThread(ctx, webhookSerialID, webhookIdSplit[1])
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "WebhookThreadの作成に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		for _, word := range webhook.MentionAndWords {
			err = h.repo.InsertWebhookWord(ctx, webhookSerialID, "MentionAnd", word)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		for _, word := range webhook.MentionOrWords {
			err = h.repo.InsertWebhookWord(ctx, webhookSerialID, "MentionOr", word)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		for _, word := range webhook.SearchAndWords {
			err = h.repo.InsertWebhookWord(ctx, webhookSerialID, "SearchAnd", word)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		for _, word := range webhook.SearchOrWords {
			err = h.repo.InsertWebhookWord(ctx, webhookSerialID, "SearchOr", word)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		for _, word := range webhook.NgAndWords {
			err = h.repo.InsertWebhookWord(ctx, webhookSerialID, "NgAnd", word)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		for _, word := range webhook.NgOrWords {
			err = h.repo.InsertWebhookWord(ctx, webhookSerialID, "NgOr", word)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		for _, roleId := range webhook.MentionRoles {
			err = h.repo.InsertWebhookRoleMention(ctx, webhookSerialID, roleId)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Roleの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		for _, userId := range webhook.MentionUsers {
			err = h.repo.InsertWebhookUserMention(ctx, webhookSerialID, userId)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Userの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
	}

	for _, webhook := range webhookJson.UpdateWebhooks {
		webhookIdSplit := strings.Split(webhook.WebhookID, "-")
		if webhook.DeleteFlag {
			err = h.repo.DeleteWebhookByWebhookSerialID(ctx, webhook.WebhookSerialID)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Webhookの削除に失敗しました:", "エラー:", err.Error())
				return
			}
			err = h.repo.DeleteWebhookThreadsNotInProvidedList(ctx, webhook.WebhookSerialID)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Webhookの削除に失敗しました:", "エラー:", err.Error())
				return
			}
			continue
		}
		err = h.repo.UpdateWebhookWithWebhookIDAndSubscription(ctx, webhook.WebhookSerialID, webhookIdSplit[0], webhook.SubscriptionId, webhook.SubscriptionType)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Webhookの更新に失敗しました:", "エラー:", err.Error())
			return
		}
		err = h.repo.DeleteWebhookWordsNotInProvidedList(ctx, webhook.WebhookSerialID, "MentionAnd", webhook.MentionAndWords)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
			return
		}
		for _, word := range webhook.MentionAndWords {
			err = h.repo.InsertWebhookWord(ctx, webhook.WebhookSerialID, "MentionAnd", word)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		err = h.repo.DeleteWebhookWordsNotInProvidedList(ctx, webhook.WebhookSerialID, "MentionOr", webhook.MentionOrWords)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
			return
		}
		for _, word := range webhook.MentionOrWords {
			err = h.repo.InsertWebhookWord(ctx, webhook.WebhookSerialID, "MentionOr", word)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		err = h.repo.DeleteWebhookWordsNotInProvidedList(ctx, webhook.WebhookSerialID, "SearchAnd", webhook.SearchAndWords)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
			return
		}
		for _, word := range webhook.SearchAndWords {
			err = h.repo.InsertWebhookWord(ctx, webhook.WebhookSerialID, "SearchAnd", word)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		err = h.repo.DeleteWebhookWordsNotInProvidedList(ctx, webhook.WebhookSerialID, "SearchOr", webhook.SearchOrWords)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
			return
		}
		for _, word := range webhook.SearchOrWords {
			err = h.repo.InsertWebhookWord(ctx, webhook.WebhookSerialID, "SearchOr", word)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		err = h.repo.DeleteWebhookWordsNotInProvidedList(ctx, webhook.WebhookSerialID, "NgAnd", webhook.NgAndWords)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
			return
		}
		for _, word := range webhook.NgAndWords {
			err = h.repo.InsertWebhookWord(ctx, webhook.WebhookSerialID, "NgAnd", word)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		err = h.repo.DeleteWebhookWordsNotInProvidedList(ctx, webhook.WebhookSerialID, "NgOr", webhook.NgOrWords)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
			return
		}
		for _, word := range webhook.NgOrWords {
			err = h.repo.InsertWebhookWord(ctx, webhook.WebhookSerialID, "NgOr", word)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Wordの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		err = h.repo.DeleteWebhookRoleMentionsNotInProvidedList(ctx, webhook.WebhookSerialID, webhook.MentionRoles)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Roleの更新に失敗しました:", "エラー:", err.Error())
			return
		}
		for _, roleId := range webhook.MentionRoles {
			err = h.repo.InsertWebhookRoleMention(ctx, webhook.WebhookSerialID, roleId)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Roleの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
		err = h.repo.DeleteWebhookUserMentionsNotInProvidedList(ctx, webhook.WebhookSerialID, webhook.MentionUsers)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Userの更新に失敗しました:", "エラー:", err.Error())
			return
		}
		for _, userId := range webhook.MentionUsers {
			err = h.repo.InsertWebhookUserMention(ctx, webhook.WebhookSerialID, userId)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				slog.ErrorContext(ctx, "Userの更新に失敗しました:", "エラー:", err.Error())
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}
