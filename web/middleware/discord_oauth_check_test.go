package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

type roundTripFn func(req *http.Request) *http.Response

func newStubHttpClient(fn roundTripFn) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func (f roundTripFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestDiscordOAuthCheckMiddleware(t *testing.T) {
	t.Run("DiscordOAuthCheckMiddlewareが正常に動作すること", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		middleware := DiscordOAuthCheckMiddleware(
			service.IndexService{
				Client: newStubHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`{
							"id": "123456789",
							"username": "test",
							"global_name": "test",
							"avatar": "test",
							"avatar_decoration": "test",
							"discriminator": "1234",
							"public_flags": 0,
							"flags": 0,
							"banner": "test",
							"banner_color": "test",
							"accent_color": "test",
							"locale": "test",
							"mfa_enabled": true,
							"premium_type": 0,
							"email": "test",
							"verified": true,
							"bio": "test",
						}`)),
					}
				}),
				CookieStore: sessions.NewCookieStore([]byte(config.SessionSecret())),
			},
			&repository.RepositoryFuncMock{
				GetPermissionCodeFunc: func(ctx context.Context, guildID string, permissionType string) (int64, error) {
					return 0, nil
				},
				GetPermissionIDsFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionID, error) {
					return nil, nil
				},
			},
			false,
		)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		middleware(handler).ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
