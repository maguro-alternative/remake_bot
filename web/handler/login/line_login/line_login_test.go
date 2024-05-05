package linelogin

import (
	"context"
	"errors"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"
	"github.com/maguro-alternative/remake_bot/pkg/crypto"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
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

	decodeClientID, err := hex.DecodeString("aa7c5fe80002633327f0fefe67a565de")
	assert.NoError(t, err)
	lineClientID, err := base64.StdEncoding.DecodeString(string([]byte("X+P6kmO6DnEjM3TVqXkwNA==")))
	assert.NoError(t, err)

	decodeClientSecret, err := hex.DecodeString("baeff317cb83ef55b193b6d3de194124")
	assert.NoError(t, err)
	lineClientSecret, err := base64.StdEncoding.DecodeString(string([]byte("uy2qtvYTnSoB5qIntwUdVQ==")))
	assert.NoError(t, err)

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
							LineNotifyToken:  pq.ByteaArray{lineNotifyStr},
							LineBotToken:     pq.ByteaArray{lineBotStr},
							LineBotSecret:    pq.ByteaArray{lineBotSecretStr},
							LineGroupID:      pq.ByteaArray{lineGroupStr},
							LineClientID:     pq.ByteaArray{lineClientID},
							LineClientSecret: pq.ByteaArray{lineClientSecret},
						},
					}, nil
				},
				GetAllColumnsLineBotIvFunc: func(ctx context.Context, guildID string) (repository.LineBotIv, error) {
					return repository.LineBotIv{
						LineNotifyTokenIv:  pq.ByteaArray{decodeNotifyToken},
						LineBotTokenIv:     pq.ByteaArray{decodeBotToken},
						LineBotSecretIv:    pq.ByteaArray{decodeBotSecret},
						LineGroupIDIv:      pq.ByteaArray{decodeGroupID},
						LineClientIDIv:     pq.ByteaArray{decodeClientID},
						LineClientSecretIv: pq.ByteaArray{decodeClientSecret},
					}, nil
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
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
				GetAllColumnsLineBotIvFunc: func(ctx context.Context, guildID string) (repository.LineBotIv, error) {
					return repository.LineBotIv{}, nil
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
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
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
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

	decodeClientID, err := hex.DecodeString("aa7c5fe80002633327f0fefe67a565de")
	assert.NoError(t, err)
	lineClientID, err := base64.StdEncoding.DecodeString(string([]byte("X+P6kmO6DnEjM3TVqXkwNA==")))
	assert.NoError(t, err)

	decodeClientSecret, err := hex.DecodeString("baeff317cb83ef55b193b6d3de194124")
	assert.NoError(t, err)
	lineClientSecret, err := base64.StdEncoding.DecodeString(string([]byte("uy2qtvYTnSoB5qIntwUdVQ==")))
	assert.NoError(t, err)

	t.Run("Lineログインのリダイレクトに成功する", func(t *testing.T) {
		// Mocking the necessary dependencies
		h := NewLineLoginHandler(
			&service.IndexService{
				CookieStore: cookieStore,
			},
			&repository.RepositoryFuncMock{
				GetAllColumnsLineBotFunc: func(ctx context.Context, guildID string) (repository.LineBot, error) {
					return repository.LineBot{
							GuildID:          "111",
							LineNotifyToken:  pq.ByteaArray{lineNotifyStr},
							LineBotToken:     pq.ByteaArray{lineBotStr},
							LineBotSecret:    pq.ByteaArray{lineBotSecretStr},
							LineGroupID:      pq.ByteaArray{lineGroupStr},
							LineClientID:     pq.ByteaArray{lineClientID},
							LineClientSecret: pq.ByteaArray{lineClientSecret},
					}, nil
				},
				GetAllColumnsLineBotIvFunc: func(ctx context.Context, guildID string) (repository.LineBotIv, error) {
					return repository.LineBotIv{
						LineNotifyTokenIv:  pq.ByteaArray{decodeNotifyToken},
						LineBotTokenIv:     pq.ByteaArray{decodeBotToken},
						LineBotSecretIv:    pq.ByteaArray{decodeBotSecret},
						LineGroupIDIv:      pq.ByteaArray{decodeGroupID},
						LineClientIDIv:     pq.ByteaArray{decodeClientID},
						LineClientSecretIv: pq.ByteaArray{decodeClientSecret},
					}, nil
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
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
				GetAllColumnsLineBotFunc: func(ctx context.Context, guildID string) (repository.LineBot, error) {
					return repository.LineBot{
							GuildID:          "111",
							LineNotifyToken:  pq.ByteaArray{lineNotifyStr},
							LineBotToken:     pq.ByteaArray{lineBotStr},
							LineBotSecret:    pq.ByteaArray{lineBotSecretStr},
							LineGroupID:      pq.ByteaArray{lineGroupStr},
							LineClientID:     pq.ByteaArray{lineClientID},
							LineClientSecret: pq.ByteaArray{lineClientSecret},
					}, nil
				},
				GetAllColumnsLineBotIvFunc: func(ctx context.Context, guildID string) (repository.LineBotIv, error) {
					return repository.LineBotIv{
						LineNotifyTokenIv:  pq.ByteaArray{decodeNotifyToken},
						LineBotTokenIv:     pq.ByteaArray{decodeBotToken},
						LineBotSecretIv:    pq.ByteaArray{decodeBotSecret},
						LineGroupIDIv:      pq.ByteaArray{decodeGroupID},
						LineClientIDIv:     pq.ByteaArray{decodeClientID},
						LineClientSecretIv: pq.ByteaArray{decodeClientSecret},
					}, nil
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
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
				GetAllColumnsLineBotFunc: func(ctx context.Context, guildID string) (repository.LineBot, error) {
					return repository.LineBot{}, errors.New("failed to get all columns of line bots")
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
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
