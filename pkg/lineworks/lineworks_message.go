package lineworks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type lineWorksMessage struct {
	Content struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
		PreviewImageUrl string `json:"previewImageUrl,omitempty"`
		OriginalContentUrl string `json:"originalContentUrl,omitempty"`
	} `json:"content"`
}

func (l *lineWorks) NewLineWorksTextMessage(message string) *lineWorksMessage {
	return &lineWorksMessage{
		Content: struct {
			Type string `json:"type"`
			Text string `json:"text,omitempty"`
			PreviewImageUrl string `json:"previewImageUrl,omitempty"`
			OriginalContentUrl string `json:"originalContentUrl,omitempty"`
		}{
			Type: "text",
			Text: message,
		},
	}
}

func (l *lineWorks) NewLineWorksImageMessage(previewImageUrl, originalContentUrl string) *lineWorksMessage {
	return &lineWorksMessage{
		Content: struct {
			Type string `json:"type"`
			Text string `json:"text,omitempty"`
			PreviewImageUrl string `json:"previewImageUrl,omitempty"`
			OriginalContentUrl string `json:"originalContentUrl,omitempty"`
		}{
			Type: "image",
			PreviewImageUrl: previewImageUrl,
			OriginalContentUrl: originalContentUrl,
		},
	}
}

func (l *lineWorks) NewLineWorksFileMessage(previewImageUrl, originalContentUrl string) *lineWorksMessage {
	return &lineWorksMessage{
		Content: struct {
			Type string `json:"type"`
			Text string `json:"text,omitempty"`
			PreviewImageUrl string `json:"previewImageUrl,omitempty"`
			OriginalContentUrl string `json:"originalContentUrl,omitempty"`
		}{
			Type: "file",
			OriginalContentUrl: originalContentUrl,
		},
	}
}

func (l *lineWorks) PushLineWorksMessage(ctx context.Context, message *lineWorksMessage) (*http.Response, error) {
	messageJson, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://www.worksapis.com/v1.0/bots/%s/channels/%s/messages", l.lineWorksBotID, l.lineWorksGroupID), bytes.NewBuffer(messageJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+l.lineWorksToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := l.client.Do(req)
	if err != nil {
		return res, err
	}

	return res, nil
}
