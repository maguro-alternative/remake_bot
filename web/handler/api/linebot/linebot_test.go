package linebot

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/maguro-alternative/remake_bot/pkg/line"
	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/maguro-alternative/remake_bot/web/handler/api/linebot/internal"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/bwmarrin/discordgo"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineBotHandler_ServeHTTP(t *testing.T) {
	bodyJson, err := json.Marshal(internal.LineResponses{
		Events: []struct {
			ReplyToken string `json:"replyToken"`
			Type       string `json:"type"`
			Source     struct {
				GroupID string `json:"groupId"`
				UserID  string `json:"userId"`
				Type    string `json:"type"`
			} `json:"source"`
			Timestamp float64 `json:"timestamp"`
			Message   struct {
				ID                  string   `json:"id"`
				Text                string   `json:"text"`
				Type                string   `json:"type"`
				Duration            int64    `json:"duration"`
				FileName            string   `json:"fileName"`
				FileSize            int64    `json:"fileSize"`
				Title               string   `json:"title"`
				Address             string   `json:"address"`
				Latitude            float64  `json:"latitude"`
				Longitude           float64  `json:"longitude"`
				PackageID           string   `json:"packageId"`
				StickerID           string   `json:"stickerId"`
				StickerResourceType string   `json:"stickerResourceType"`
				Keywords            []string `json:"keywords"`
				ImageSet            struct {
					ID    string  `json:"id"`
					Index float64 `json:"index"`
					Total float64 `json:"total"`
				} `json:"imageSet"`
				ContentProvider struct {
					Type               string `json:"type"`
					OriginalContentURL string `json:"originalContentUrl"`
					PreviewImageURL    string `json:"previewImageUrl"`
				} `json:"contentProvider"`
			} `json:"message"`
			Mode            string `json:"mode"`
			WebhookEventID  string `json:"webhookEventId"`
			DeliveryContext struct {
				IsRedelivery bool `json:"isRedelivery"`
			} `json:"deliveryContext"`
		}{
			{
				ReplyToken: "123456789",
				Type:       "message",
				Source: struct {
					GroupID string `json:"groupId"`
					UserID  string `json:"userId"`
					Type    string `json:"type"`
				}{
					GroupID: "123456789",
					UserID:  "123456789",
					Type:    "group",
				},
				Timestamp: 123456789,
				Message: struct {
					ID                  string   `json:"id"`
					Text                string   `json:"text"`
					Type                string   `json:"type"`
					Duration            int64    `json:"duration"`
					FileName            string   `json:"fileName"`
					FileSize            int64    `json:"fileSize"`
					Title               string   `json:"title"`
					Address             string   `json:"address"`
					Latitude            float64  `json:"latitude"`
					Longitude           float64  `json:"longitude"`
					PackageID           string   `json:"packageId"`
					StickerID           string   `json:"stickerId"`
					StickerResourceType string   `json:"stickerResourceType"`
					Keywords            []string `json:"keywords"`
					ImageSet            struct {
						ID    string  `json:"id"`
						Index float64 `json:"index"`
						Total float64 `json:"total"`
					} `json:"imageSet"`
					ContentProvider struct {
						Type               string `json:"type"`
						OriginalContentURL string `json:"originalContentUrl"`
						PreviewImageURL    string `json:"previewImageUrl"`
					} `json:"contentProvider"`
				}{
					ID:                  "123456789",
					Text:                "test",
					Type:                "text",
					Duration:            123456789,
					FileName:            "test",
					FileSize:            123456789,
					Title:               "test",
					Address:             "test",
					Latitude:            123.456789,
					Longitude:           123.456789,
					PackageID:           "123456789",
					StickerID:           "123456789",
					StickerResourceType: "test",
					Keywords:            []string{"test"},
					ImageSet: struct {
						ID    string  `json:"id"`
						Index float64 `json:"index"`
						Total float64 `json:"total"`
					}{
						ID:    "123456789",
						Index: 123.456789,
						Total: 123.456789,
					},
					ContentProvider: struct {
						Type               string `json:"type"`
						OriginalContentURL string `json:"originalContentUrl"`
						PreviewImageURL    string `json:"previewImageUrl"`
					}{
						Type:               "test",
						OriginalContentURL: "test",
						PreviewImageURL:    "test",
					},
				},
				Mode:           "active",
				WebhookEventID: "123456789",
				DeliveryContext: struct {
					IsRedelivery bool `json:"isRedelivery"`
				}{
					IsRedelivery: true,
				},
			},
		},
	})
	require.NoError(t, err)

	decodeNotifyToken, err := hex.DecodeString("aa7c5fe80002633327f0fefe67a565de")
	assert.NoError(t, err)
	lineNotifyStr, err := base64.StdEncoding.DecodeString(string([]byte("X+P6kmO6DnEjM3TVqXkwNA==")))
	assert.NoError(t, err)

	decodeBotToken, err := hex.DecodeString("baeff317cb83ef55b193b6d3de194124")
	assert.NoError(t, err)
	lineBotStr, err := base64.StdEncoding.DecodeString(string([]byte("uy2qtvYTnSoB5qIntwUdVQ==")))
	assert.NoError(t, err)

	decodeBotSecret, err := hex.DecodeString("0ffa8ed72efcb5f1d834e4ce8463a62c")
	assert.NoError(t, err)
	lineBotSecretStr, err := base64.StdEncoding.DecodeString(string([]byte("i2uHQCyn58wRR/b03fRw6w==")))
	assert.NoError(t, err)

	decodeGroupID, err := hex.DecodeString("e14db710b23520766fd652c0f19d437a")
	assert.NoError(t, err)
	lineGroupStr, err := base64.StdEncoding.DecodeString(string([]byte("YgexFQQlLcaXmsw9mFN35Q==")))
	assert.NoError(t, err)

	// スタブHTTPクライアントを作成
	stubClient := line.NewStubHttpClient(func(req *http.Request) *http.Response {
		if strings.Contains(req.URL.String(), "https://api.line.me/v2/bot/profile/") {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"displayName": "test",
					"userId": "123456789",
					"pictureUrl": "test",
					"statusMessage": "test"
				}`)),
			}
		} else if strings.Contains(req.URL.String(), "https://api-data.line.me/v2/bot/message/") {
			cwd, err := os.Getwd()
			require.NoError(t, err)
			t.Cleanup(func() {
				require.NoError(t, os.Chdir(cwd))
			})
			require.NoError(t, os.Chdir("../.."))
			srcMp3, err := os.Open(cwd + "/on_message_create/yumi_dannasama.mp3")
			require.NoError(t, err)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(srcMp3),
			}
		}
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("")),
		}
	})

	t.Run("MethodがPOST以外の場合、Method Not Allowedが返ること", func(t *testing.T) {
		h := &LineBotHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/linebot", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("jsonの読み取りに失敗すると、Internal Server Errorが返ること", func(t *testing.T) {
		h := &LineBotHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/linebot", bytes.NewReader([]byte("")))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("jsonのバリデーションに失敗すると、Internal Server Errorが返ること", func(t *testing.T) {
		h := &LineBotHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/linebot", bytes.NewReader([]byte(`{"line_bot_secret":"123456789"}`)))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("LineBotがない場合Bad Requestを返すこと", func(t *testing.T) {
		h := &LineBotHandler{
			IndexService: &service.IndexService{
				Client:         stubClient,
				DiscordSession: &discordgo.Session{},
			},
			Repo: &repository.RepositoryFuncMock{
				GetAllColumnsLineBotsFunc: func(ctx context.Context) ([]*repository.LineBot, error) {
					return []*repository.LineBot{}, nil
				},
				GetLineBotIvNotClientFunc: func(ctx context.Context, guildId string) (repository.LineBotIvNotClient, error) {
					return repository.LineBotIvNotClient{}, nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/linebot", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("LineBotが正常に送信できること", func(t *testing.T) {
		h := &LineBotHandler{
			IndexService: &service.IndexService{
				Client:         stubClient,
				DiscordSession: &mock.SessionMock{
					ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
						return &discordgo.Message{}, nil
					},
				},
			},
			Repo: &repository.RepositoryFuncMock{
				GetAllColumnsLineBotsFunc: func(ctx context.Context) ([]*repository.LineBot, error) {
					return []*repository.LineBot{
						{
							GuildID:         "123456789",
							LineNotifyToken: pq.ByteaArray{lineNotifyStr},
							LineBotToken:    pq.ByteaArray{lineBotStr},
							LineBotSecret:   pq.ByteaArray{lineBotSecretStr},
							LineGroupID:     pq.ByteaArray{lineGroupStr},
						},
					}, nil
				},
				GetLineBotIvNotClientFunc: func(ctx context.Context, guildId string) (repository.LineBotIvNotClient, error) {
					return repository.LineBotIvNotClient{
						LineNotifyTokenIv: pq.ByteaArray{decodeNotifyToken},
						LineBotTokenIv:    pq.ByteaArray{decodeBotToken},
						LineBotSecretIv:   pq.ByteaArray{decodeBotSecret},
						LineGroupIDIv:     pq.ByteaArray{decodeGroupID},
					}, nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/linebot", bytes.NewReader([]byte(`{"events":[{"replyToken":"","type":"message","timestamp":0,"source":{"userId":"Udeadbw00dbaadbeefdeadbeefdeadbeef","type":"user"},"message":{"type":"text","id":"1234567890","text":"Hello, world"}}]}`)))
		r.Header.Set("X-Line-Signature", "6eMInZT4CEsIf/P5Iv+9VmezoOPqXs1il6R4QjtUG4o=")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
