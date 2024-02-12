package line

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

// LINE Notifyのメッセージ
type LineNotifyMessage struct {
	Message string `json:"message"`
}

func (l *LineNotifyMessage) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Message, validation.Required),
	)
}

// LINE Notifyの画像
type LineNotifyImage struct {
	ImageThumbnail string `json:"imageThumbnail"`
	ImageFullsize  string `json:"imageFullsize"`
	Message        string `json:"message"`
}

func (l *LineNotifyImage) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ImageThumbnail, validation.Required),
		validation.Field(&l.ImageFullsize, validation.Required),
		validation.Field(&l.Message, validation.Required),
	)
}

// LINE Notifyのスタンプ
type LineNotifySticker struct {
	StickerPackageID string `json:"stickerPackageId"`
	StickerID        string `json:"stickerId"`
	Message          string `json:"message"`
}

func (l *LineNotifySticker) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.StickerPackageID, validation.Required),
		validation.Field(&l.StickerID, validation.Required),
		validation.Field(&l.Message, validation.Required),
	)
}

// LINE Notifyでメッセージを送信
func (r *LineRequest) PushMessageNotify(ctx context.Context, message string) error {
	client := &http.Client{}
	url := "https://notify-api.line.me/api/notify"
	notifyMessage := LineNotifyMessage{
		Message: "message: " + message,
	}
	err := notifyMessage.Validate()
	if err != nil {
		return err
	}
	notifyMessageByte, err := json.Marshal(notifyMessage)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(notifyMessageByte))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+r.lineNotifyToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// LINE Notifyで画像を送信
func (r *LineRequest) PushImageNotify(ctx context.Context, message, image string) error {
	client := &http.Client{}
	url := "https://notify-api.line.me/api/notify"
	notifyMessage := LineNotifyImage{
		Message: message,
	}
	err := notifyMessage.Validate()
	if err != nil {
		return err
	}
	notifyMessageByte, err := json.Marshal(notifyMessage)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(notifyMessageByte))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+r.lineNotifyToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// LINE Notifyでスタンプを送信
func (r *LineRequest) PushStampNotify(ctx context.Context, message, stickerPackageID, stickerID string) error {
	client := &http.Client{}
	url := "https://notify-api.line.me/api/notify"
	notifyMessage := LineNotifySticker{
		Message:          message,
		StickerPackageID: stickerPackageID,
		StickerID:        stickerID,
	}
	err := notifyMessage.Validate()
	if err != nil {
		return err
	}
	notifyMessageByte, err := json.Marshal(notifyMessage)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(notifyMessageByte))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+r.lineNotifyToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
