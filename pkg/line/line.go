package line

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type LineRequest struct {
	lineNotifyToken string
	lineBotToken    string
	lineGroupID     string
}

func NewLineRequest(lineNotifyToken, lineBotToken, lineGroupID string) *LineRequest {
	return &LineRequest{
		lineNotifyToken: lineNotifyToken,
		lineBotToken:    lineBotToken,
		lineGroupID:     lineGroupID,
	}
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

// LINEのプロフィール情報をLINEグループから取得
func (r *LineRequest) GetProfileInGroup(ctx context.Context, userID string) (LineProfile, error) {
	var lineProfile LineProfile
	client := &http.Client{}
	url := "https://api.line.me/v2/bot/group/" + r.lineGroupID + "/member/" + userID
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

// LINEBotのプロフィール情報を取得
func (r *LineRequest) GetBotInfo(ctx context.Context) (LineProfile, error) {
	var lineProfile LineProfile
	client := &http.Client{}
	url := "https://api.line.me/v2/bot/info"
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

// LINEの友達数を取得
func (r *LineRequest) GetFriendCount(ctx context.Context) (int, error) {
	client := &http.Client{}
	url := "https://api.line.me/v2/bot/followers/count"
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

	var lineBotFriend LineBotFriend
	err = json.NewDecoder(resp.Body).Decode(&lineBotFriend)
	if err != nil {
		return 0, err
	}
	err = lineBotFriend.Validate()
	return lineBotFriend.Followers, err
}

// LINEのグループのユーザー数を取得
func (r *LineRequest) GetGroupUserCount(ctx context.Context) (int, error) {
	client := &http.Client{}
	url := "https://api.line.me/v2/bot/group/" + r.lineGroupID + "/members/count"
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

	var groupMemberCount LineGroupMemberCount
	err = json.NewDecoder(resp.Body).Decode(&groupMemberCount)
	if err != nil {
		return 0, err
	}
	return groupMemberCount.Count, nil
}

// LINEのプッシュ通知の使用した回数を取得
func (r *LineRequest) GetTotalPushCount(ctx context.Context) (int, error) {
	client := &http.Client{}
	url := "https://api.line.me/v2/bot/message/quota/consumption"
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

	var lineBotConsumption LineBotConsumption
	err = json.NewDecoder(resp.Body).Decode(&lineBotConsumption)
	if err != nil {
		return 0, err
	}
	err = lineBotConsumption.Validate()
	return lineBotConsumption.TotalUsage, err
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
