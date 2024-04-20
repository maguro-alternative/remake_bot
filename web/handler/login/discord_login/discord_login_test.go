package discordlogin

import (
    "net/http"
    "net/http/httptest"
    "testing"

	"github.com/maguro-alternative/remake_bot/web/service"

    "github.com/gorilla/sessions"
    "github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
    t.Run("test successful ServeHTTP", func(t *testing.T) {
        // Mocking the necessary dependencies
        h := NewDiscordOAuth2Handler(
            &service.IndexService{
                CookieStore: sessions.NewCookieStore([]byte("test")),
            },
		)

        req, err := http.NewRequest("GET", "/login/discord", nil)
        assert.NoError(t, err)

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(h.ServeHTTP)

        handler.ServeHTTP(rr, req)

        assert.Equal(t, http.StatusSeeOther, rr.Code)
    })

    t.Run("test ServeHTTP with session store creation error", func(t *testing.T) {
        // Mocking the necessary dependencies
        h := NewDiscordOAuth2Handler(
            &service.IndexService{
                CookieStore: sessions.NewCookieStore([]byte("")),
            },
		)

        req, err := http.NewRequest("GET", "/login/discord", nil)
        assert.NoError(t, err)

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(h.ServeHTTP)

        handler.ServeHTTP(rr, req)

        assert.Equal(t, http.StatusInternalServerError, rr.Code)
    })
}
