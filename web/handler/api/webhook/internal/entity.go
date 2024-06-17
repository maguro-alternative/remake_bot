package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type WebhookJson struct {
	NewWebhooks    []Webhook `json:"newWebhooks"`
	UpdateWebhooks []Webhook `json:"updateWebhooks"`
}

type Webhook struct {
	WebhookID        string   `json:"webhookId"`
	WebhookURL       string   `json:"webhookUrl"`
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
		validation.Field(&g.NewWebhooks, validation.Required),
	)
}
