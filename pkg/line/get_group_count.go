package line

import (
	"context"
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

// LINEグループのメンバー数
type LineGroupMemberCount struct {
	Count   int    `json:"count"`
	Message string `json:"message,omitempty"`
}

func (l *LineGroupMemberCount) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Count, validation.Required),
	)
}

// LINEのグループのユーザー数を取得
func (r *LineRequest) GetGroupUserCount(ctx context.Context) (int, error) {
	/* https://developers.line.biz/ja/reference/messaging-api/#get-members-group-count */
	url := "https://api.line.me/v2/bot/group/" + r.lineGroupID + "/members/count"
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

	var groupMemberCount LineGroupMemberCount
	err = json.NewDecoder(resp.Body).Decode(&groupMemberCount)
	if err != nil {
		return 0, err
	}
	if err = groupMemberCount.Validate(); err != nil {
		return 0, err
	}
	return groupMemberCount.Count, nil
}
