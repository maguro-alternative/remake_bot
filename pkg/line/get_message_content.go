package line

import (
	"context"
	"errors"
	"io"
	"net/http"
)

type LineMessageContent struct {
	Content       []byte
	ContentLength int64
	ContentType   string
}

// LINEメッセージ内のファイルを取得
func (r *LineRequest) GetContent(ctx context.Context, messageID string) (LineMessageContent, error) {
	/* https://developers.line.biz/ja/reference/messaging-api/#get-content */
	url := "https://api-data.line.me/v2/bot/message/" + messageID + "/content"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return LineMessageContent{}, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+r.lineBotToken)
	resp, err := r.client.Do(req)
	if err != nil {
		return LineMessageContent{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return LineMessageContent{}, errors.New("LINEメッセージの取得に失敗しました。" + resp.Status)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return LineMessageContent{}, err
	}
	content := LineMessageContent{
		Content:       b,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
	}
	return content, nil
}
