package linebot

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type LineResponses struct {
	Events []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Source     struct {
			GroupID string `json:"groupId"`
			UserID  string `json:"userId"`
			Type    string `json:"type"`
		} `json:"source"`
		Timestamp float64 `json:"timestamp"`
		Message   struct {
			ID                  string   `json:"id"`
			Text                string   `json:"text"`
			Duration            int64    `json:"duration"`
			FileName            string   `json:"fileName"`
			FileSize            int64    `json:"fileSize"`
			Title               string   `json:"title"`
			Address             string   `json:"address"`
			Latitude            float64  `json:"latitude"`
			Longitude           float64  `json:"longitude"`
			PackageID           string   `json:"packageId"`
			StickerID           string   `json:"stickerId"`
			StickerResourceType string   `json:"stickerResourceType"`
			Keywords            []string `json:"keywords"`
			ImageSet            struct {
				ID    string  `json:"id"`
				Index float64 `json:"index"`
				Total float64 `json:"total"`
			} `json:"imageSet"`
			ContentProvider struct {
				Type               string `json:"type"`
				OriginalContentURL string `json:"originalContentUrl"`
				PreviewImageURL    string `json:"previewImageUrl"`
			} `json:"contentProvider"`
		} `json:"message"`
		Mode            string `json:"mode"`
		WebhookEventID  string `json:"webhookEventId"`
		DeliveryContext struct {
			IsRedelivery bool `json:"isRedelivery"`
		} `json:"deliveryContext"`
	} `json:"events"`
}

// Validate validates the LineResponses.
func (h *LineResponses) Validate() error {
	return validation.ValidateStruct(h,
		validation.Field(&h.Events, validation.Required),
	)
}