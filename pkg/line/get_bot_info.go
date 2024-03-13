package line

import (
	"context"
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

// LINEのボット情報
type LineBotInfo struct {
	BasicID        string `json:"basicId"`
	ChatMode       string `json:"chatMode"`
	MarkAsReadMode string `json:"markAsReadMode"`
	PremiumID      string `json:"premiumId,omitempty"`
	PictureURL     string `json:"pictureUrl"`
	DisplayName    string `json:"displayName"`
	UserID         string `json:"userId"`
	Message        string `json:"message"`
	Details        []struct {
		Property string `json:"property"`
		Message  string `json:"message"`
	} `json:"details,omitempty"`
}

func (l *LineBotInfo) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.BasicID, validation.Required),
		validation.Field(&l.ChatMode, validation.Required),
		validation.Field(&l.MarkAsReadMode, validation.Required),
		validation.Field(&l.PictureURL, validation.Required),
		validation.Field(&l.DisplayName, validation.Required),
		validation.Field(&l.UserID, validation.Required),
	)
}

// LINEBotのプロフィール情報を取得
func (r *LineRequest) GetBotInfo(ctx context.Context) (LineBotInfo, error) {
	var lineBotProfile LineBotInfo
	client := &http.Client{}
	url := "https://api.line.me/v2/bot/info"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return lineBotProfile, err
	}

	req.Header.Set("Authorization", "Bearer "+r.lineBotToken)
	resp, err := client.Do(req)
	if err != nil {
		return lineBotProfile, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&lineBotProfile)
	if err != nil {
		return lineBotProfile, err
	}
	err = lineBotProfile.Validate()
	return lineBotProfile, err
}
