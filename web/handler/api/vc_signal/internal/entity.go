package internal

import (
	//validation "github.com/go-ozzo/ozzo-validation"
)

type VcSignalJson struct {
	VcChannelID string `json:"vcChannelId"`
	SendSignal    bool `json:"sendSignal"`
	SendChannelId string `json:"sendChannelId"`
	JoinBot 	 bool `json:"joinBot"`
	EveryoneMention bool `json:"everyoneMention"`
}
