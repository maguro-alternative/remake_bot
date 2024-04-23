package linecallback

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/gob"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maguro-alternative/remake_bot/testutil/mock"
	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session"
	"github.com/maguro-alternative/remake_bot/web/shared/model"

	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	gob.Register(&model.LineIdTokenUser{})
}

func TestDiscordCallbackHandler_ServeHTTP(t *testing.T) {
	cookieStore := sessions.NewCookieStore([]byte(config.SessionSecret()))
	decodeClientID, err := hex.DecodeString("aa7c5fe80002633327f0fefe67a565de")
	assert.NoError(t, err)
	lineClientID, err := base64.StdEncoding.DecodeString(string([]byte("X+P6kmO6DnEjM3TVqXkwNA==")))
	assert.NoError(t, err)

	decodeClientSecret, err := hex.DecodeString("baeff317cb83ef55b193b6d3de194124")
	assert.NoError(t, err)
	lineClientSecret, err := base64.StdEncoding.DecodeString(string([]byte("uy2qtvYTnSoB5qIntwUdVQ==")))
	assert.NoError(t, err)
	t.Run("successful callback", func(t *testing.T) {
		middlewareStartFixture := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
				require.NoError(t, err)
				sessionStore.SetLineState("123")
				sessionStore.SetLineNonce("456")
				sessionStore.SetGuildID("111")
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
		}

		// Mock the DiscordOAuth2Service
		svc := &service.IndexService{
			Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"access_token": "abc",
						"token_type": "Bearer",
						"expires_in": 604800,
						"refresh_token": "def",
						"scope": "identify"
					}`)),
				}
			}),
			CookieStore: cookieStore,
		}
		repo := &repository.RepositoryFuncMock{
			GetAllColumnsLineBotFunc: func(ctx context.Context, guildId string) (repository.LineBot, error) {
				return repository.LineBot{
					GuildID:          "",
					LineClientID:     pq.ByteaArray{lineClientID},
					LineClientSecret: pq.ByteaArray{lineClientSecret},
				}, nil
			},
			GetAllColumnsLineBotIvFunc: func(ctx context.Context, guildID string) (repository.LineBotIv, error) {
				return repository.LineBotIv{
					LineClientIDIv:     pq.ByteaArray{decodeClientID},
					LineClientSecretIv: pq.ByteaArray{decodeClientSecret},
				}, nil
			},
		}

		// Create a new handler with the mocked service
		handler := NewLineCallbackHandler(svc, repo)

		// Create a new HTTP request
		req, err := http.NewRequest(http.MethodGet, "/callback?state=123&code=abc&nonce=456", nil)
		require.NoError(t, err)

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Call ServeHTTP directly
		middlewareStartFixture(handler).ServeHTTP(rr, req)
		middlewareEndFixture(req)

		// Check the status code
		assert.Equal(t, http.StatusFound, rr.Code)

		// Check the response body
		assert.Contains(t, rr.Body.String(), "/group/111")
	})

	t.Run("error on state mismatch", func(t *testing.T) {
		// Mock the DiscordOAuth2Service
		svc := &service.IndexService{
			Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
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
			CookieStore: cookieStore,
		}
		repo := &repository.RepositoryFuncMock{}

		// Create a new handler with the mocked service
		handler := NewLineCallbackHandler(svc, repo)

		// Create a new HTTP request with a different state
		req, err := http.NewRequest(http.MethodGet, "/callback?state=456&code=abc&nonce=456", nil)
		require.NoError(t, err)

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Call ServeHTTP directly
		handler.ServeHTTP(rr, req)

		// Check the status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check the response body
		assert.Contains(t, rr.Body.String(), "Internal Server Error")
	})
}
