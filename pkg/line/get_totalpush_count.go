package line

import (
	"context"
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

// LINEBotのメッセージ合計送信数
type LineBotConsumption struct {
	TotalUsage int    `json:"totalUsage"`
	Message    string `json:"message,omitempty"`
}

func (l *LineBotConsumption) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.TotalUsage, validation.Required),
	)
}

// LINEのプッシュ通知の使用した回数を取得
func (r *LineRequest) GetTotalPushCount(ctx context.Context) (int, error) {
	url := "https://api.line.me/v2/bot/message/quota/consumption"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+r.lineBotToken)
	resp, err := r.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var lineBotConsumption LineBotConsumption
	err = json.NewDecoder(resp.Body).Decode(&lineBotConsumption)
	if err != nil {
		return 0, err
	}
	err = lineBotConsumption.Validate()
	return lineBotConsumption.TotalUsage, err
}
