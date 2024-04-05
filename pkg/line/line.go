package line

import (
	"net/http"
)

type RoundTripFn func(req *http.Request) *http.Response

func (f RoundTripFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewStubHttpClient(fn RoundTripFn) *http.Client {
	return &http.Client{
		Transport: RoundTripFn(fn),
	}
}

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
