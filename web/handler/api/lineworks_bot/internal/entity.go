package internal

import (
	"time"
)

type LineWorksResponses struct {
	Type   string `json:"type"`
	Source struct {
		UserID    string `json:"userId"`
		ChannelId string `json:"channelId"`
		DomainId  int64  `json:"domainId"`
	} `json:"source"`
	IssuedTime string `json:"issuedTime"`
	Content    struct {
		Type               string `json:"type"`
		Text               string `json:"text,omitempty"`
		FileId             string `json:"fileId,omitempty"`
		OriginalContentURL string `json:"originalContentUrl,omitempty"`
	} `json:"content"`
}

type LineWorksBotDecrypt struct {
	LineWorksBotToken string
	LineWorksRefreshToken string
	LineWorksGroupID string
	LineWorksBotID string
	LineWorksBotSecret string
	RefreshTokenExpiresAt time.Time
	DefaultChannelID string
	DebugMode bool
}
