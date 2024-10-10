package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/webhook/internal"

	"github.com/stretchr/testify/assert"
)

func TestWebhookHandler_ServeHTTP(t *testing.T) {
	webhook := internal.WebhookJson{
		NewWebhooks: []*internal.NewWebhook{
			{
				WebhookID:        "987654321",
				SubscriptionType: "youtube",
				SubscriptionId:   "987654321",
				MentionAndWords:  []string{"word1", "word2"},
			},
		},
	}

	t.Run("MethodがPOST以外の場合、Method Not Allowedが返ること", func(t *testing.T) {
		h := &WebhookHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/987654321/webhook", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("バリデーションチェックに失敗した場合、BadRequestが返ること", func(t *testing.T) {
		h := &WebhookHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/webhook", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Webhookの更新が成功すること", func(t *testing.T) {
		bodyJson, err := json.Marshal(webhook)
		assert.NoError(t, err)
		h := &WebhookHandler{
			repo: &repository.RepositoryFuncMock{
				InsertWebhookFunc: func(ctx context.Context, guildID, webhookID, subscriptionType, subscriptionID string, lastPostedAt time.Time) (int64, error) {
					return 1, nil
				},
				InsertWebhookWordFunc: func(ctx context.Context, webhookSerialID int64, mentionAndWordType, word string) error {
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/webhook", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Webhookの更新が失敗すること(新規のwebhookのinsert失敗)", func(t *testing.T) {
		bodyJson, err := json.Marshal(webhook)
		assert.NoError(t, err)
		h := &WebhookHandler{
			repo: &repository.RepositoryFuncMock{
				InsertWebhookFunc: func(ctx context.Context, guildID, webhookID, subscriptionType, subscriptionID string, lastPostedAt time.Time) (int64, error) {
					return 0, assert.AnError
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/webhook", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Webhookの更新が失敗すること(mentionAndWordのinsert失敗)", func(t *testing.T) {
		bodyJson, err := json.Marshal(webhook)
		assert.NoError(t, err)
		h := &WebhookHandler{
			repo: &repository.RepositoryFuncMock{
				InsertWebhookFunc: func(ctx context.Context, guildID, webhookID, subscriptionType, subscriptionID string, lastPostedAt time.Time) (int64, error) {
					return 1, nil
				},
				InsertWebhookWordFunc: func(ctx context.Context, webhookSerialID int64, mentionAndWordType, word string) error {
					return assert.AnError
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/webhook", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Webhookの更新が失敗すること(mentionOrWordのinsert失敗)", func(t *testing.T) {
		webhook.NewWebhooks[0].MentionAndWords = []string{"word1", "word2", "word3"}
		bodyJson, err := json.Marshal(webhook)
		assert.NoError(t, err)
		h := &WebhookHandler{
			repo: &repository.RepositoryFuncMock{
				InsertWebhookFunc: func(ctx context.Context, guildID, webhookID, subscriptionType, subscriptionID string, lastPostedAt time.Time) (int64, error) {
					return 1, nil
				},
				InsertWebhookWordFunc: func(ctx context.Context, webhookSerialID int64, mentionAndWordType, word string) error {
					if word == "word3" {
						return assert.AnError
					}
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/webhook", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Webhookの更新が失敗すること(SearchAndWordのinsert失敗)", func(t *testing.T) {
		webhook.NewWebhooks[0].MentionAndWords = []string{"word1", "word2", "word3"}
		bodyJson, err := json.Marshal(webhook)
		assert.NoError(t, err)
		h := &WebhookHandler{
			repo: &repository.RepositoryFuncMock{
				InsertWebhookFunc: func(ctx context.Context, guildID, webhookID, subscriptionType, subscriptionID string, lastPostedAt time.Time) (int64, error) {
					return 1, nil
				},
				InsertWebhookWordFunc: func(ctx context.Context, webhookSerialID int64, mentionAndWordType, word string) error {
					if word == "word3" {
						return assert.AnError
					}
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/webhook", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Webhookの更新が失敗すること(mentionAndWordのinsert失敗)", func(t *testing.T) {
		webhook.NewWebhooks[0].MentionAndWords = []string{"word1", "word2", "word3"}
		bodyJson, err := json.Marshal(webhook)
		assert.NoError(t, err)
		h := &WebhookHandler{
			repo: &repository.RepositoryFuncMock{
				InsertWebhookFunc: func(ctx context.Context, guildID, webhookID, subscriptionType, subscriptionID string, lastPostedAt time.Time) (int64, error) {
					return 1, nil
				},
				InsertWebhookWordFunc: func(ctx context.Context, webhookSerialID int64, mentionAndWordType, word string) error {
					if word == "word3" {
						return assert.AnError
					}
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/webhook", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Webhookの更新が失敗すること(リクエストのパース失敗)", func(t *testing.T) {
		h := &WebhookHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/webhook", bytes.NewReader([]byte("invalid json")))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Threadに対するWebhookの登録に成功すること", func(t *testing.T) {
		webhook := internal.WebhookJson{
			NewWebhooks: []*internal.NewWebhook{
				{
					WebhookID:        "987654321-123456789",
					SubscriptionType: "youtube",
					SubscriptionId:   "987654321",
					MentionAndWords:  []string{"word1", "word2"},
				},
			},
		}
		bodyJson, err := json.Marshal(webhook)
		assert.NoError(t, err)
		h := WebhookHandler{
			repo: &repository.RepositoryFuncMock{
				InsertWebhookFunc: func(ctx context.Context, guildID, webhookID, subscriptionType, subscriptionID string, lastPostedAt time.Time) (int64, error) {
					assert.Equal(t, webhookID, "987654321")
					return 1, nil
				},
				InsertWebhookWordFunc: func(ctx context.Context, webhookSerialID int64, mentionAndWordType, word string) error {
					return nil
				},
				InsertWebhookThreadFunc: func(ctx context.Context, webhookSerialID int64, threadID string) error {
					assert.Equal(t, threadID, "123456789")
					return nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/webhook", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
