package discordcallback

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"io"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestDiscordCallbackHandler_ServeHTTP(t *testing.T) {
	cookieStore := sessions.NewCookieStore([]byte(config.SessionSecret()))
	t.Run("successful callback", func(t *testing.T) {
		middlewareStartFixture := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
				require.NoError(t, err)
				sessionStore.SetDiscordState("123")
				err = sessionStore.SessionSave(r, w)
				require.NoError(t, err)
				h.ServeHTTP(w, r)
			})
		}
		middlewareEndFixture := func(r *http.Request) {
			sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
			require.NoError(t, err)
			sessionStore.CleanupDiscordOAuthToken()
			sessionStore.CleanupDiscordUser()
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

		// Create a new handler with the mocked service
		handler := NewDiscordCallbackHandler(svc)

		// Create a new HTTP request
		req, err := http.NewRequest(http.MethodGet, "/callback?state=123&code=abc", nil)
		require.NoError(t, err)

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Call ServeHTTP directly
		middlewareStartFixture(handler).ServeHTTP(rr, req)
		middlewareEndFixture(req)

		// Check the status code
		assert.Equal(t, http.StatusFound, rr.Code)

		// Check the response body
		assert.Contains(t, rr.Body.String(), "guilds")
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

		// Create a new handler with the mocked service
		handler := NewDiscordCallbackHandler(svc)

		// Create a new HTTP request with a different state
		req, err := http.NewRequest(http.MethodGet, "/callback?state=456&code=abc", nil)
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
