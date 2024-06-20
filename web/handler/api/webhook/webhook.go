package webhook


import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/webhook/internal"
)

type WebhookHandler struct {
	repo repository.RepositoryFunc
}

func NewVcSignalHandler(
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

	for _, webhook := range webhookJson.NewWebhooks {
		err = h.repo.InsertWebhook(ctx, guildId, webhook.WebhookID, webhook.SubscriptionType, webhook.SubscriptionId, now)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Webhookの更新に失敗しました:", "エラー:", err.Error())
			return
		}
		
	}
}
