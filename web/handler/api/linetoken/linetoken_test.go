package linetoken

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/lib/pq"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/maguro-alternative/remake_bot/web/handler/api/linetoken/internal"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/stretchr/testify/assert"
)

func TestLineTokenHandler_ServeHTTP(t *testing.T) {
	bodyJson, err := json.Marshal(internal.LineBotJson{
		GuildID:          "987654321",
		DefaultChannelID: "123456789",
		DebugMode:        true,
	})
	assert.NoError(t, err)

	stubClient := mock.NewStubHttpClient(func(req *http.Request) *http.Response {
		if req.URL.String() == "https://api.line.me/v2/bot/info" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"userId":"U1234567890","basicId":"U1234567890","chatMode":"text","markAsReadMode":"r","displayName":"testUser","pictureUrl":"https://profile.line-scdn.net/"}`)),
			}
		} else if req.URL.String() == "https://notify-api.line.me/api/status" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":200,"message":"ok"}`)),
			}
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"count":1}`)),
		}
	})

	t.Run("MethodがPOST以外の場合、Method Not Allowedが返ること", func(t *testing.T) {
		h := &LineTokenHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/987654321/linetoken", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("jsonの読み取りに失敗すると、Bad Requestが返ること", func(t *testing.T) {
		h := &LineTokenHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/linetoken", bytes.NewReader([]byte("")))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("jsonのバリデーションに失敗すると、Unprocessable Entityが返ること", func(t *testing.T) {
		h := &LineTokenHandler{
			indexService: &service.IndexService{},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/linetoken", bytes.NewReader([]byte(`{"default_channel_id":"aaa123456789"}`)))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("LineBotの更新が成功すること", func(t *testing.T) {
		h := &LineTokenHandler{
			indexService: &service.IndexService{
				Client:         stubClient,
				DiscordSession: &discordgo.Session{},
			},
			repo: &repository.RepositoryFuncMock{
				GetAllColumnsLineBotByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBot, error) {
					return repository.LineBot{
						GuildID:          "123",
						LineNotifyToken:  pq.ByteaArray{[]byte("lineNotifyStr")},
						LineBotToken:     pq.ByteaArray{[]byte("lineBotStr")},
						LineBotSecret:    pq.ByteaArray{[]byte("lineBotSecretStr")},
						LineGroupID:      pq.ByteaArray{[]byte("lineGroupStr")},
						LineClientID:     pq.ByteaArray{[]byte("lineClientID")},
						LineClientSecret: pq.ByteaArray{[]byte("lineClientSecret")},
					}, nil
				},
				GetAllColumnsLineBotIvByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotIv, error) {
					return repository.LineBotIv{
						LineNotifyTokenIv:  pq.ByteaArray{[]byte("decodeNotifyToken")},
						LineBotTokenIv:     pq.ByteaArray{[]byte("decodeBotToken")},
						LineBotSecretIv:    pq.ByteaArray{[]byte("decodeBotSecret")},
						LineGroupIDIv:      pq.ByteaArray{[]byte("decodeGroupID")},
						LineClientIDIv:     pq.ByteaArray{[]byte("decodeClientID")},
						LineClientSecretIv: pq.ByteaArray{[]byte("decodeClientSecret")},
					}, nil
				},
				UpdateLineBotFunc: func(ctx context.Context, lineBot *repository.LineBot) error {
					return nil
				},
				UpdateLineBotIvFunc: func(ctx context.Context, lineBotIv *repository.LineBotIv) error {
					return nil
				},
			},
			aesCrypto: &crypto.AESMock{
				EncryptFunc: func(data []byte) (encrypted []byte, iv []byte, err error) {
					return []byte("test"), []byte("test"), nil
				},
				DecryptFunc: func(data []byte, iv []byte) ([]byte, error) {
					if string(iv) == string("decodeNotifyToken") {
						return []byte("testnotifytoken"), nil
					} else if string(iv) == string("decodeBotToken") {
						return []byte("testbottoken"), nil
					} else if string(iv) == string("decodeBotSecret") {
						return []byte("testbotsecret"), nil
					} else if string(iv) == string("decodeGroupID") {
						return []byte("testgroupid"), nil
					}
					return nil, nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/linetoken", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Deleteフラグが立った場合、該当する(すべて)ものがnilになること", func(t *testing.T) {
		h := &LineTokenHandler{
			indexService: &service.IndexService{
				Client:         stubClient,
				DiscordSession: &discordgo.Session{},
			},
			repo: &repository.RepositoryFuncMock{
				GetAllColumnsLineBotByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBot, error) {
					return repository.LineBot{
						GuildID:          "123",
						LineNotifyToken:  pq.ByteaArray{[]byte("lineNotifyStr")},
						LineBotToken:     pq.ByteaArray{[]byte("lineBotStr")},
						LineBotSecret:    pq.ByteaArray{[]byte("lineBotSecretStr")},
						LineGroupID:      pq.ByteaArray{[]byte("lineGroupStr")},
						LineClientID:     pq.ByteaArray{[]byte("lineClientID")},
						LineClientSecret: pq.ByteaArray{[]byte("lineClientSecret")},
					}, nil
				},
				GetAllColumnsLineBotIvByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotIv, error) {
					return repository.LineBotIv{
						LineNotifyTokenIv:  pq.ByteaArray{[]byte("decodeNotifyToken")},
						LineBotTokenIv:     pq.ByteaArray{[]byte("decodeBotToken")},
						LineBotSecretIv:    pq.ByteaArray{[]byte("decodeBotSecret")},
						LineGroupIDIv:      pq.ByteaArray{[]byte("decodeGroupID")},
						LineClientIDIv:     pq.ByteaArray{[]byte("decodeClientID")},
						LineClientSecretIv: pq.ByteaArray{[]byte("decodeClientSecret")},
					}, nil
				},
				UpdateLineBotFunc: func(ctx context.Context, lineBot *repository.LineBot) error {
					assert.Nil(t, lineBot.LineNotifyToken)
					assert.Nil(t, lineBot.LineBotToken)
					assert.Nil(t, lineBot.LineBotSecret)
					assert.Nil(t, lineBot.LineGroupID)
					assert.Nil(t, lineBot.LineClientID)
					assert.Nil(t, lineBot.LineClientSecret)
					return nil
				},
				UpdateLineBotIvFunc: func(ctx context.Context, lineBotIv *repository.LineBotIv) error {
					return nil
				},
			},
			aesCrypto: &crypto.AESMock{
				EncryptFunc: func(data []byte) (encrypted []byte, iv []byte, err error) {
					return []byte("test"), []byte("test"), nil
				},
				DecryptFunc: func(data []byte, iv []byte) ([]byte, error) {
					if string(iv) == string("decodeNotifyToken") {
						return []byte("testnotifytoken"), nil
					} else if string(iv) == string("decodeBotToken") {
						return []byte("testbottoken"), nil
					} else if string(iv) == string("decodeBotSecret") {
						return []byte("testbotsecret"), nil
					} else if string(iv) == string("decodeGroupID") {
						return []byte("testgroupid"), nil
					}
					return nil, nil
				},
			},
		}
		w := httptest.NewRecorder()

		bodyJsonDelete, err := json.Marshal(internal.LineBotJson{
			GuildID:                "987654321",
			DefaultChannelID:       "123456789",
			DebugMode:              true,
			LineNotifyTokenDelete:  true,
			LineBotTokenDelete:     true,
			LineBotSecretDelete:    true,
			LineGroupIDDelete:      true,
			LineClientIDDelete:     true,
			LineClientSecretDelete: true,
		})
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/api/987654321/linetoken", bytes.NewReader(bodyJsonDelete))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Deleteフラグが立った場合、該当する(notifyとbotのtoken)ものがnilになること", func(t *testing.T) {
		h := &LineTokenHandler{
			indexService: &service.IndexService{
				Client:         stubClient,
				DiscordSession: &discordgo.Session{},
			},
			repo: &repository.RepositoryFuncMock{
				GetAllColumnsLineBotByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBot, error) {
					return repository.LineBot{
						GuildID:          "123",
						LineNotifyToken:  pq.ByteaArray{[]byte("lineNotifyStr")},
						LineBotToken:     pq.ByteaArray{[]byte("lineBotStr")},
						LineBotSecret:    pq.ByteaArray{[]byte("lineBotSecretStr")},
						LineGroupID:      pq.ByteaArray{[]byte("lineGroupStr")},
						LineClientID:     pq.ByteaArray{[]byte("lineClientID")},
						LineClientSecret: pq.ByteaArray{[]byte("lineClientSecret")},
					}, nil
				},
				GetAllColumnsLineBotIvByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotIv, error) {
					return repository.LineBotIv{
						LineNotifyTokenIv:  pq.ByteaArray{[]byte("decodeNotifyToken")},
						LineBotTokenIv:     pq.ByteaArray{[]byte("decodeBotToken")},
						LineBotSecretIv:    pq.ByteaArray{[]byte("decodeBotSecret")},
						LineGroupIDIv:      pq.ByteaArray{[]byte("decodeGroupID")},
						LineClientIDIv:     pq.ByteaArray{[]byte("decodeClientID")},
						LineClientSecretIv: pq.ByteaArray{[]byte("decodeClientSecret")},
					}, nil
				},
				UpdateLineBotFunc: func(ctx context.Context, lineBot *repository.LineBot) error {
					assert.Nil(t, lineBot.LineNotifyToken)
					assert.Nil(t, lineBot.LineBotToken)
					assert.NotNil(t, lineBot.LineBotSecret)
					assert.NotNil(t, lineBot.LineGroupID)
					assert.NotNil(t, lineBot.LineClientID)
					assert.NotNil(t, lineBot.LineClientSecret)
					return nil
				},
				UpdateLineBotIvFunc: func(ctx context.Context, lineBotIv *repository.LineBotIv) error {
					return nil
				},
			},
			aesCrypto: &crypto.AESMock{
				EncryptFunc: func(data []byte) (encrypted []byte, iv []byte, err error) {
					return []byte("test"), []byte("test"), nil
				},
				DecryptFunc: func(data []byte, iv []byte) ([]byte, error) {
					if string(iv) == string("decodeNotifyToken") {
						return []byte("testnotifytoken"), nil
					} else if string(iv) == string("decodeBotToken") {
						return []byte("testbottoken"), nil
					} else if string(iv) == string("decodeBotSecret") {
						return []byte("testbotsecret"), nil
					} else if string(iv) == string("decodeGroupID") {
						return []byte("testgroupid"), nil
					}
					return nil, nil
				},
			},
		}
		w := httptest.NewRecorder()

		bodyJsonDelete, err := json.Marshal(internal.LineBotJson{
			GuildID:               "987654321",
			DefaultChannelID:      "123456789",
			DebugMode:             true,
			LineNotifyTokenDelete: true,
			LineBotTokenDelete:    true,
		})
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/api/987654321/linetoken", bytes.NewReader(bodyJsonDelete))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
