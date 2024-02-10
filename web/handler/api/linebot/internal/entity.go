package internal

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type LineBot struct {
	LineNotifyToken   []byte `db:"line_notify_token"`
	LineBotToken      []byte `db:"line_bot_token"`
	LineBotSecret     []byte `db:"line_bot_secret"`
	LineGroupID       []byte `db:"line_group_id"`
	LineClientID      []byte `db:"line_client_id"`
	LineClientSercret []byte `db:"line_client_sercret"`
	Iv                []byte `db:"iv"`
	DefaultChannelID  string `db:"default_channel_id"`
	DebugMode         bool   `db:"debug_mode"`
}

type LineBotDecrypt struct {
	LineNotifyToken   string
	LineBotToken      string
	LineGroupID       string
	LineClientID      string
	LineClientSercret string
	DefaultChannelID  string
	DebugMode         bool
}

type LineProfile struct {
	DisplayName   string `json:"displayName"`
	UserID        string `json:"userId"`
	PictureUrl    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
	Message       string `json:"message"`
}

func (l *LineProfile) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.DisplayName, validation.Required),
		validation.Field(&l.UserID, validation.Required),
		validation.Field(&l.PictureUrl, validation.Required),
		validation.Field(&l.StatusMessage, validation.Required),
	)
}
