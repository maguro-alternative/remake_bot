package line

import (
	"context"
	"io"
	"net/http"
)

type LineMessageContent struct {
	Content       io.ReadCloser
	ContentLength int64
	ContentType   string
}

// LINEメッセージ内のファイルを取得
func (r *LineRequest) GetContent(ctx context.Context, messageID string) (LineMessageContent, error) {
	//var buf bytes.Buffer
	client := &http.Client{}
	url := "https://api.line.me/v2/bot/message/" + messageID + "/content"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return LineMessageContent{}, err
	}

	req.Header.Set("Authorization", "Bearer "+r.lineBotToken)
	resp, err := client.Do(req)
	if err != nil {
		return LineMessageContent{}, err
	}
	defer resp.Body.Close()
	content := LineMessageContent{
		Content: resp.Body,
		ContentType: resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
	}
	return content, nil
}
