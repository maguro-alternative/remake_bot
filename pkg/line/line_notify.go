package line

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

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
		bodyStr := fmt.Sprintf("%+v\n", resp)
		return errors.New(bodyStr)
	}
	defer resp.Body.Close()

	return nil
}

// LINE Notifyで画像を送信
func (r *LineRequest) PushImageNotify(ctx context.Context, message, image string) error {
	client := &http.Client{}
	notifyUrl := "https://notify-api.line.me/api/notify"
	u, err := url.ParseRequestURI(notifyUrl)
	if err != nil {
		return err
	}
	form := url.Values{}
	form.Add("message", message)
	form.Add("imageThumbnail", image)
	form.Add("imageFullsize", image)

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
		bodyStr := fmt.Sprintf("%+v\n", resp)
		return errors.New(bodyStr)
	}
	defer resp.Body.Close()

	return nil
}

// LINE Notifyでスタンプを送信
func (r *LineRequest) PushStampNotify(ctx context.Context, message, stickerPackageID, stickerID string) error {
	client := &http.Client{}
	notifyUrl := "https://notify-api.line.me/api/notify"
	u, err := url.ParseRequestURI(notifyUrl)
	if err != nil {
		return err
	}
	form := url.Values{}
	form.Add("message", message)
	form.Add("stickerPackageId", stickerPackageID)
	form.Add("stickerId", stickerID)

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
		bodyStr := fmt.Sprintf("%+v\n", resp)
		return errors.New(bodyStr)
	}
	defer resp.Body.Close()

	return nil
}

// LINE Notifyのメッセージを送信するためのリクエストを作成(実行はDo()で行う)
func (r *LineRequest) PushMessageNotifyCall(ctx context.Context, message string) (*LineNotifyCall, error) {
	client := &http.Client{}
	notifyUrl := "https://notify-api.line.me/api/notify"
	u, err := url.ParseRequestURI(notifyUrl)
	if err != nil {
		return nil, err
	}
	form := url.Values{}
	form.Add("message", message)

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
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
	notifyUrl := "https://notify-api.line.me/api/notify"
	u, err := url.ParseRequestURI(notifyUrl)
	if err != nil {
		return nil, err
	}
	form := url.Values{}
	form.Add("message", message)
	form.Add("imageThumbnail", image)
	form.Add("imageFullsize", image)

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
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

// LINE Notifyのスタンプを送信するためのリクエストを作成(実行はDo()で行う)
func (r *LineRequest) PushStampNotifyCall(ctx context.Context, message, stickerPackageID, stickerID string) (*LineNotifyCall, error) {
	client := &http.Client{}
	notifyUrl := "https://notify-api.line.me/api/notify"
	u, err := url.ParseRequestURI(notifyUrl)
	if err != nil {
		return nil, err
	}
	form := url.Values{}
	form.Add("message", message)
	form.Add("stickerPackageId", stickerPackageID)
	form.Add("stickerId", stickerID)

	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
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
