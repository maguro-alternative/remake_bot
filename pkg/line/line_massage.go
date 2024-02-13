package line

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

// LINEのメッセージ
type LineMessage struct {
	To       string            `json:"to"`
	Messages []LineMessageType `json:"messages"`
}

func (l *LineMessage) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.To, validation.Required),
		validation.Field(&l.Messages, validation.Required),
	)
}

// LINEのメッセージタイプ
type LineMessageType struct {
	Type               string `json:"type"`
	Text               string `json:"text,omitempty"`
	ImageThumbnail     string `json:"imageThumbnail,omitempty"`
	ImageFullsize      string `json:"imageFullsize,omitempty"`
	OriginalContentUrl string `json:"originalContentUrl,omitempty"`
	PreviewImageUrl    string `json:"previewImageUrl,omitempty"`
	Duration           int    `json:"duration,omitempty"`
}

// テキストメッセージを作成
func (r *LineRequest) NewLineTextMessage(message string) LineMessageType {
	return LineMessageType{
		Type: "text",
		Text: message,
	}
}

// 画像メッセージを作成
func (r *LineRequest) NewLineImageMessage(imageThumbnail, imageFullsize string) LineMessageType {
	return LineMessageType{
		Type:           "image",
		ImageThumbnail: imageThumbnail,
		ImageFullsize:  imageFullsize,
	}
}

// 動画メッセージを作成
func (r *LineRequest) NewLineVideoMessage(originalContentUrl, previewImageUrl string) LineMessageType {
	return LineMessageType{
		Type:               "video",
		OriginalContentUrl: originalContentUrl,
		PreviewImageUrl:    previewImageUrl,
	}
}

// 音声メッセージを作成
func (r *LineRequest) NewLineAudioMessage(originalContentUrl string, duration int) LineMessageType {
	return LineMessageType{
		Type:               "audio",
		OriginalContentUrl: originalContentUrl,
		Duration:           duration * 1000,
	}
}

// LINEにメッセージを送信
func (r *LineRequest) PushMessageBotInGroup(ctx context.Context, messages []LineMessageType) error {
	client := &http.Client{}
	url := "https://api.line.me/v2/bot/message/push"
	lineMessage := LineMessage{
		To:       r.lineGroupID,
		Messages: messages,
	}
	err := lineMessage.Validate()
	if err != nil {
		return err
	}
	lineMessageByte, err := json.Marshal(lineMessage)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(lineMessageByte))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+r.lineBotToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
