package middleware

import (
	"context"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
	"github.com/maguro-alternative/remake_bot/web/shared/session"

	"github.com/lib/pq"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineOAuthCheckMiddleware(t *testing.T) {
	cookieStore := sessions.NewCookieStore([]byte(config.SessionSecret()))
	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	gob.Register(&model.LineIdTokenUser{})
	gob.Register(&model.LineOAuthSession{})
	decodeNotifyToken, err := hex.DecodeString("aa7c5fe80002633327f0fefe67a565de")
	require.NoError(t, err)
	lineNotifyStr, err := base64.StdEncoding.DecodeString(string([]byte("X+P6kmO6DnEjM3TVqXkwNA==")))
	require.NoError(t, err)

	decodeBotToken, err := hex.DecodeString("baeff317cb83ef55b193b6d3de194124")
	require.NoError(t, err)
	lineBotStr, err := base64.StdEncoding.DecodeString(string([]byte("uy2qtvYTnSoB5qIntwUdVQ==")))
	require.NoError(t, err)

	decodeBotSecret, err := hex.DecodeString("0ffa8ed72efcb5f1d834e4ce8463a62c")
	require.NoError(t, err)
	lineBotSecretStr, err := base64.StdEncoding.DecodeString(string([]byte("i2uHQCyn58wRR/b03fRw6w==")))
	require.NoError(t, err)

	decodeGroupID, err := hex.DecodeString("e14db710b23520766fd652c0f19d437a")
	require.NoError(t, err)
	lineGroupStr, err := base64.StdEncoding.DecodeString(string([]byte("YgexFQQlLcaXmsw9mFN35Q==")))
	require.NoError(t, err)
	t.Run("LineOAuthCheckMiddlewareが正常に動作すること", func(t *testing.T) {
		middlewareStartFixture := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user := &model.LineIdTokenUser{
					Iss:     "https://access.line.me",
					Sub:     "U123456789abcdef123456789abcdef12",
					Aud:     "1234567890",
					Exp:     1504169092,
					Iat:     1504168492,
					Nonce:   "0987654321fedcba987654321fedcba9",
					Amr:     []string{"pwd"},
					Name:    "test",
					Picture: "test",
					Email:   "test",
				}
				sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
				require.NoError(t, err)
				sessionStore.SetLineUser(user)
				sessionStore.SetLineOAuthToken("test")
				sessionStore.SetGuildID("test")
				err = sessionStore.SessionSave(r, w)
				require.NoError(t, err)
				h.ServeHTTP(w, r)
			})
		}
		middlewareEndFixture := func(r *http.Request) {
				sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
				require.NoError(t, err)
				sessionStore.CleanupLineUser()
				sessionStore.CleanupLineOAuthToken()
				sessionStore.CleanupGuildID()
		}
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		middleware := LineOAuthCheckMiddleware(
			service.IndexService{
				Client: newStubHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`{
							"userId": "U123456789",
							"displayName": "test",
							"pictureUrl": "test",
							"statusMessage": "test"
						}`)),
					}
				}),
			},
			&repository.RepositoryFuncMock{
				GetLineBotNotClientFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return repository.LineBotNotClient{
							LineNotifyToken: pq.ByteaArray{lineNotifyStr},
							LineBotToken:    pq.ByteaArray{lineBotStr},
							LineBotSecret:   pq.ByteaArray{lineBotSecretStr},
							LineGroupID:     pq.ByteaArray{lineGroupStr},
					}, nil
				},
				GetLineBotIvNotClientFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return repository.LineBotIvNotClient{
						LineNotifyTokenIv: pq.ByteaArray{decodeNotifyToken},
						LineBotTokenIv:    pq.ByteaArray{decodeBotToken},
						LineBotSecretIv:   pq.ByteaArray{decodeBotSecret},
						LineGroupIDIv:     pq.ByteaArray{decodeGroupID},
					}, nil
				},
			},
			true,
		)(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer test")
		w := httptest.NewRecorder()
		middlewareStartFixture(middleware).ServeHTTP(w, req)
		middlewareEndFixture(req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
