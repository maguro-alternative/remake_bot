package line

import (
	"context"
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

// LINEのプロフィール
type LineProfile struct {
	DisplayName   string `json:"displayName"`
	UserID        string `json:"userId"`
	PictureUrl    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
	Message       string `json:"message,omitempty"`
}

func (l *LineProfile) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.DisplayName, validation.Required),
		validation.Field(&l.UserID, validation.Required),
		validation.Field(&l.PictureUrl, validation.Required),
		validation.Field(&l.StatusMessage, validation.Required),
	)
}

// LINEのプロフィール情報を取得
func (r *LineRequest) GetProfile(ctx context.Context, userID string) (LineProfile, error) {
	var lineProfile LineProfile
	client := &http.Client{}
	url := "https://api.line.me/v2/bot/profile/" + userID
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return lineProfile, err
	}

	req.Header.Set("Authorization", "Bearer "+r.lineBotToken)
	resp, err := client.Do(req)
	if err != nil {
		return lineProfile, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&lineProfile)
	if err != nil {
		return lineProfile, err
	}
	err = lineProfile.Validate()
	return lineProfile, err
}
