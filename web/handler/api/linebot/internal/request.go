package internal

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type Request struct {
	db db.Driver

	lineBotDecrypt LineBotDecrypt
}

func NewRequest(db db.Driver, lineBotDecrypt LineBotDecrypt) *Request {
	return &Request{
		db: db,

		lineBotDecrypt: lineBotDecrypt,
	}
}

// LINEのプロフィール情報を取得
func (r *Request) GetProfile(ctx context.Context, groupID, userID string) (LineProfile, error) {
	var lineProfile LineProfile
	client := &http.Client{}
	url := "https://api.line.me/v2/bot/group/" + groupID + "/member/" + userID
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return lineProfile, err
	}

	req.Header.Set("Authorization", "Bearer "+r.lineBotDecrypt.LineBotToken)
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
