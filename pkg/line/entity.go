package line

import (
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

// LINEのボット情報
type LineBotInfo struct {
	BasicID        string `json:"basicId"`
	ChatMode       string `json:"chatMode"`
	MarkAsReadMode string `json:"markAsReadMode"`
	PremiumID      string `json:"premiumId"`
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
		validation.Field(&l.PremiumID, validation.Required),
		validation.Field(&l.PictureURL, validation.Required),
		validation.Field(&l.DisplayName, validation.Required),
		validation.Field(&l.UserID, validation.Required),
	)
}

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

// LINE Notifyのメッセージ
type LineNotifyMessage struct {
	Message string `json:"message"`
}

func (l *LineNotifyMessage) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Message, validation.Required),
	)
}

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
