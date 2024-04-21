package linelogout

import (
    "net/http"
    "net/http/httptest"
    "testing"

	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
	"github.com/maguro-alternative/remake_bot/web/shared/session"

	"github.com/gorilla/sessions"
    "github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
    t.Run("test successful logout", func(t *testing.T) {
		cookieStore := sessions.NewCookieStore([]byte("test"))
        req, err := http.NewRequest("GET", "/logout/line", nil)
        assert.NoError(t, err)

        rr := httptest.NewRecorder()
        handler := &LineLogoutHandler{
            LineLogoutService: &service.IndexService{
				CookieStore: cookieStore,
			},
        }

		sessionStore, err := session.NewSessionStore(req, cookieStore, "test")
		assert.NoError(t, err)
		sessionStore.SetLineOAuthToken("aaa")
		sessionStore.SetGuildID("123")
		sessionStore.SetLineUser(&model.LineIdTokenUser{
			Sub: "123",
			Name: "test",
		})

        handler.ServeHTTP(rr, req)

        assert.Equal(t, http.StatusSeeOther, rr.Code)
        assert.Equal(t, "/", rr.Header().Get("Location"))

		_, err = sessionStore.GetGuildID()
		assert.Error(t, err)
		_, err = sessionStore.GetLineOAuthToken()
		assert.Error(t, err)
		_, err = sessionStore.GetLineUser()
		assert.Error(t, err)
    })

    t.Run("test session store creation failure", func(t *testing.T) {
		cookieStore := sessions.NewCookieStore([]byte(""))
        req, err := http.NewRequest("GET", "/logout/line", nil)
        assert.NoError(t, err)

        rr := httptest.NewRecorder()
        handler := &LineLogoutHandler{
            LineLogoutService: &service.IndexService{
				CookieStore: cookieStore,
			},
        }

        handler.ServeHTTP(rr, req)

        assert.Equal(t, http.StatusInternalServerError, rr.Code)
    })
}
