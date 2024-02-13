package line

import (
	"context"
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

// 当月に送信できるメッセージ数の上限目安
type LineBotQuota struct {
	Type    string `json:"type"`
	Value   int    `json:"value"`
	Message string `json:"message,omitempty"`
}

func (l *LineBotQuota) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Type, validation.Required),
		validation.Field(&l.Value, validation.Required),
	)
}

// LINEのプッシュ通知の使用可能回数を取得
func (r *LineRequest) GetPushLimit(ctx context.Context) (int, error) {
	client := &http.Client{}
	url := "https://api.line.me/v2/bot/message/quota"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+r.lineBotToken)
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var lineBotQuota LineBotQuota
	err = json.NewDecoder(resp.Body).Decode(&lineBotQuota)
	if err != nil {
		return 0, err
	}
	err = lineBotQuota.Validate()
	return lineBotQuota.Value, err
}
