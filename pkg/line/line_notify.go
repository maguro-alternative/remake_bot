package line

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

// LINE Notifyのリクエスト
type LineNotifyCall struct {
	c *http.Client
	r *http.Request
}

// LINE Notifyでメッセージを送信
func (r *LineRequest) PushMessageNotify(ctx context.Context, message string) error {
	client := &http.Client{}
	notifyUrl := "https://notify-api.line.me/api/notify"
	u, err := url.ParseRequestURI(notifyUrl)
	if err != nil {
		return err
	}
	form := url.Values{}
	form.Add("message", message)

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+r.lineNotifyToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintln("%+v\n", resp))
	}
	defer resp.Body.Close()

	return nil
}

// LINE Notifyで画像を送信
func (r *LineRequest) PushImageNotify(ctx context.Context, message, image string) error {
	client := &http.Client{}
	url := "https://notify-api.line.me/api/notify"
	notifyMessage := LineNotifyImage{
		Message:        message,
		ImageThumbnail: image,
		ImageFullsize:  image,
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

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+r.lineNotifyToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// LINE Notifyのメッセージを送信するためのリクエストを作成(実行はDo()で行う)
func (r *LineRequest) PushMessageNotifyCall(ctx context.Context, message string) (*LineNotifyCall, error) {
	client := &http.Client{}
	url := "https://notify-api.line.me/api/notify"
	notifyMessage := LineNotifyMessage{
		Message: message,
	}
	err := notifyMessage.Validate()
	if err != nil {
		return nil, err
	}
	notifyMessageByte, err := json.Marshal(notifyMessage)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(notifyMessageByte))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+r.lineNotifyToken)
	return &LineNotifyCall{
		c: client,
		r: req,
	}, nil
}

// LINE Notifyの画像を送信するためのリクエストを作成(実行はDo()で行う)
func (r *LineRequest) PushImageNotifyCall(ctx context.Context, message, image string) (*LineNotifyCall, error) {
	client := &http.Client{}
	url := "https://notify-api.line.me/api/notify"
	notifyMessage := LineNotifyImage{
		Message:        message,
		ImageThumbnail: image,
		ImageFullsize:  image,
	}
	err := notifyMessage.Validate()
	if err != nil {
		return nil, err
	}
	notifyMessageByte, err := json.Marshal(notifyMessage)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(notifyMessageByte))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+r.lineNotifyToken)
	return &LineNotifyCall{
		c: client,
		r: req,
	}, nil
}

// LINE Notifyのスタンプを送信するためのリクエストを作成(実行はDo()で行う)
func (r *LineRequest) PushStampNotifyCall(ctx context.Context, message, stickerPackageID, stickerID string) (*LineNotifyCall, error) {
	client := &http.Client{}
	url := "https://notify-api.line.me/api/notify"
	notifyMessage := LineNotifySticker{
		Message:          message,
		StickerPackageID: stickerPackageID,
		StickerID:        stickerID,
	}
	err := notifyMessage.Validate()
	if err != nil {
		return nil, err
	}

	notifyMessageByte, err := json.Marshal(notifyMessage)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(notifyMessageByte))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+r.lineNotifyToken)
	return &LineNotifyCall{
		c: client,
		r: req,
	}, nil
}

// LINE Notifyのメッセージを送信するためのリクエストを実行
func (l *LineNotifyCall) Do() (*http.Response, error) {
	return l.c.Do(l.r)
}
