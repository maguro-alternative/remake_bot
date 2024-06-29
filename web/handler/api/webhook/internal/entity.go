package internal

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

type WebhookJson struct {
	NewWebhooks    []*NewWebhook    `json:"newWebhooks"`
	UpdateWebhooks []*UpdateWebhook `json:"updateWebhooks"`
}

type NewWebhook struct {
	WebhookID        string   `json:"webhookId"`
	SubscriptionType string   `json:"subscriptionType"`
	SubscriptionId   string   `json:"subscriptionId"`
	MentionRoles     []string `json:"mentionRoles"`
	MentionUsers     []string `json:"mentionUsers"`
	NgOrWords        []string `json:"ngOrWords"`
	NgAndWords       []string `json:"ngAndWords"`
	SearchOrWords    []string `json:"searchOrWords"`
	SearchAndWords   []string `json:"searchAndWords"`
	MentionOrWords   []string `json:"mentionOrWords"`
	MentionAndWords  []string `json:"mentionAndWords"`
}

type UpdateWebhook struct {
	WebhookSerialID  int64    `json:"webhookSerialId"`
	WebhookID        string   `json:"webhookId"`
	SubscriptionType string   `json:"subscriptionType"`
	SubscriptionId   string   `json:"subscriptionId"`
	MentionRoles     []string `json:"mentionRoles"`
	MentionUsers     []string `json:"mentionUsers"`
	NgOrWords        []string `json:"ngOrWords"`
	NgAndWords       []string `json:"ngAndWords"`
	SearchOrWords    []string `json:"searchOrWords"`
	SearchAndWords   []string `json:"searchAndWords"`
	MentionOrWords   []string `json:"mentionOrWords"`
	MentionAndWords  []string `json:"mentionAndWords"`
	DeleteFlag       bool     `json:"deleteFlag"`
}

func (g WebhookJson) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.NewWebhooks,
			validation.Length(0, 1000),
			validation.Each(validation.By(func(value interface{}) error {
				// valueをNewWebhook型へキャストし、そのアドレスを取得
				webhookValue, ok := value.(NewWebhook)
				if !ok {
					return errors.New("type assertion to NewWebhook failed")
				}
				return validation.ValidateStruct(&webhookValue,
					validation.Field(&webhookValue.WebhookID, validation.Required),
					validation.Field(&webhookValue.SubscriptionType, validation.In("youtube", "niconico")),
				)
			})),
		),
		validation.Field(&g.UpdateWebhooks,
			validation.Length(0, 1000),
			validation.Each(validation.By(func(value interface{}) error {
				// valueをUpdateWebhook型へキャストし、そのアドレスを取得
				webhookValue, ok := value.(UpdateWebhook)
				if !ok {
					return errors.New("type assertion to UpdateWebhook failed")
				}
				return validation.ValidateStruct(&webhookValue,
					validation.Field(&webhookValue.WebhookID, validation.Required),
					validation.Field(&webhookValue.SubscriptionType, validation.In("youtube", "niconico")),
				)
			})),
		),
	)
}
