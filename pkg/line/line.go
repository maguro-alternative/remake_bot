package line

import (
	"net/http"
)

type LineRequest struct {
	client          http.Client
	lineNotifyToken string
	lineBotToken    string
	lineGroupID     string
}

func NewLineRequest(client http.Client, lineNotifyToken, lineBotToken, lineGroupID string) *LineRequest {
	return &LineRequest{
		client:          client,
		lineNotifyToken: lineNotifyToken,
		lineBotToken:    lineBotToken,
		lineGroupID:     lineGroupID,
	}
}
