package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type WebhookJson struct {
	Webhooks []Webhook `json:"webhooks"`
}

type Webhook struct {
	WebhookID        string   `json:"webhookId"`
	ChannelID        string   `json:"channelId"`
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
		validation.Field(&g.Webhooks, validation.Required),
	)
}
