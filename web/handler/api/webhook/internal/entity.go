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
	MentionRoles     []int    `json:"mentionRoles"`
	MentionUsers     []int    `json:"mentionUsers"`
	NgOrWords        []string `json:"ngOrWords"`
	NgAndWords       []string `json:"ngAndWords"`
}

func (g WebhookJson) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.Webhooks, validation.Required),
	)
}
