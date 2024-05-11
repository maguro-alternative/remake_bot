package linelogin

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(cwd))
	})
	require.NoError(t, os.Chdir("../../../../"))

	t.Run("Getではない場合400エラーを返す", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/login/line", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := &LineLoginHandler{}

		handler.Index(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Lineログイン選択画面に遷移すること", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/login/line", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := NewLineLoginHandler(
			&service.IndexService{
				Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`{
							"basicId": "basicId",
							"chatMode": "chatMode",
							"markAsReadMode": "markAsReadMode",
							"premiumId": "premiumId",
							"pictureUrl": "pictureUrl",
							"displayName": "displayName",
							"userId": "userId",
							"message": "message"
						}`)),
					}
				}),
			},
			&repository.RepositoryFuncMock{
				GetAllColumnsLineBotsFunc: func(ctx context.Context) ([]*repository.LineBot, error) {
					return []*repository.LineBot{
						{
							GuildID:          "123",
							LineNotifyToken:  pq.ByteaArray{[]byte("lineNotifyStr")},
							LineBotToken:     pq.ByteaArray{[]byte("lineBotStr")},
							LineBotSecret:    pq.ByteaArray{[]byte("lineBotSecretStr")},
							LineGroupID:      pq.ByteaArray{[]byte("lineGroupStr")},
							LineClientID:     pq.ByteaArray{[]byte("lineClientID")},
							LineClientSecret: pq.ByteaArray{[]byte("lineClientSecret")},
						},
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
			},
			&crypto.AESMock{
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
				},
			},
		)

		handler.Index(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		assert.Contains(t, rr.Body.String(), "<title>LINEログイン選択</title>")

		assert.Contains(t, rr.Body.String(), `<a href="/login/line/123">`)
		assert.Contains(t, rr.Body.String(), `<img src="pictureUrl"/>`)
		assert.Contains(t, rr.Body.String(), `<li>displayName</li>`)
	})

	t.Run("Lineログイン選択画面に遷移すること(Botなしで何も表示されない)", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/login/line", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := NewLineLoginHandler(
			&service.IndexService{
				Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`{
							"basicId": "basicId",
							"chatMode": "chatMode",
							"markAsReadMode": "markAsReadMode",
							"premiumId": "premiumId",
							"pictureUrl": "pictureUrl",
							"displayName": "displayName",
							"userId": "userId",
							"message": "message"
						}`)),
					}
				}),
			},
			&repository.RepositoryFuncMock{
				GetAllColumnsLineBotsFunc: func(ctx context.Context) ([]*repository.LineBot, error) {
					return nil, nil
				},
				GetAllColumnsLineBotIvByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotIv, error) {
					return repository.LineBotIv{}, nil
				},
			},
			&crypto.AESMock{
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
				},
			},
		)

		handler.Index(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		assert.Contains(t, rr.Body.String(), "<title>LINEログイン選択</title>")

		assert.NotContains(t, rr.Body.String(), `<a href="/login/line/123">`)
		assert.NotContains(t, rr.Body.String(), `<img src="pictureUrl"/>`)
		assert.NotContains(t, rr.Body.String(), `<li>displayName</li>`)
	})

	t.Run("lineBotの読み込みに失敗した場合、500エラーを返す", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/login/line", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := NewLineLoginHandler(
			&service.IndexService{},
			&repository.RepositoryFuncMock{
				GetAllColumnsLineBotsFunc: func(ctx context.Context) ([]*repository.LineBot, error) {
					return nil, errors.New("failed to get all columns of line bots")
				},
			},
			&crypto.AESMock{
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
				},
			},
		)

		handler.Index(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	// Add more tests here for other error conditions
}

func TestLineLogin(t *testing.T) {
	cookieStore := sessions.NewCookieStore([]byte(config.SessionSecret()))

	t.Run("Lineログインのリダイレクトに成功する", func(t *testing.T) {
		// Mocking the necessary dependencies
		h := NewLineLoginHandler(
			&service.IndexService{
				CookieStore: cookieStore,
			},
			&repository.RepositoryFuncMock{
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
			},
			&crypto.AESMock{
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
				},
			},
		)

		req, err := http.NewRequest("GET", "/login/line/111", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.LineLogin)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusSeeOther, rr.Code)
	})

	t.Run("セッションの読み込みに失敗した場合500を返す", func(t *testing.T) {
		// Mocking the necessary dependencies
		h := NewLineLoginHandler(
			&service.IndexService{
				CookieStore: sessions.NewCookieStore([]byte("")),
			},
			&repository.RepositoryFuncMock{
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
			},
			&crypto.AESMock{
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
				},
			},
		)

		req, err := http.NewRequest("GET", "/login/line/111", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.LineLogin)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("データベースがエラーを返した場合500エラーを返す", func(t *testing.T) {
		// Mocking the necessary dependencies
		h := NewLineLoginHandler(
			&service.IndexService{
				CookieStore: cookieStore,
			},
			&repository.RepositoryFuncMock{
				GetAllColumnsLineBotByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBot, error) {
					return repository.LineBot{}, errors.New("failed to get all columns of line bots")
				},
			},
			&crypto.AESMock{
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
				},
			},
		)

		req, err := http.NewRequest("GET", "/login/line/111", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.LineLogin)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
