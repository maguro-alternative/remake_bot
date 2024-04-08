package line

import (
	"context"
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

// LINEBotの友達数
type LineBotFriend struct {
	Status          string `json:"status"`
	Followers       int    `json:"followers"`
	TargetedReaches int    `json:"targetedReaches"`
	Blocked         int    `json:"blocked"`
	Message         string `json:"message,omitempty"`
}

func (l *LineBotFriend) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Status, validation.Required),
		validation.Field(&l.Followers, validation.Required),
		validation.Field(&l.TargetedReaches, validation.Required),
		validation.Field(&l.Blocked, validation.Required),
	)
}

// LINEの友達数を取得
func (r *LineRequest) GetFriendCount(ctx context.Context) (int, error) {
	url := "https://api.line.me/v2/bot/followers/count"
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

	var lineBotFriend LineBotFriend
	err = json.NewDecoder(resp.Body).Decode(&lineBotFriend)
	if err != nil {
		return 0, err
	}
	err = lineBotFriend.Validate()
	return lineBotFriend.Followers, err
}

// LINEのプロフィール情報をLINEグループから取得
func (r *LineRequest) GetProfileInGroup(ctx context.Context, userID string) (LineProfile, error) {
	var lineProfile LineProfile
	url := "https://api.line.me/v2/bot/group/" + r.lineGroupID + "/member/" + userID
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return lineProfile, err
	}

	req.Header.Set("Authorization", "Bearer "+r.lineBotToken)
	resp, err := r.client.Do(req)
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
