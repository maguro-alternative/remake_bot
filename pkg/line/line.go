package line

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
