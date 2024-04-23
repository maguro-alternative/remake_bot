package discordlogout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
	"github.com/maguro-alternative/remake_bot/web/shared/session"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServeHTTP(t *testing.T) {
	t.Run("discordログアウトできること", func(t *testing.T) {
		cookieStore := sessions.NewCookieStore([]byte("test"))
		indexService := &service.IndexService{
			CookieStore: cookieStore,
		} // Mock this service as per your implementation
		handler := NewDiscordOAuth2Handler(indexService)

		req, err := http.NewRequest("GET", "/logout/discord", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handlerFunc := http.HandlerFunc(handler.ServeHTTP)

		sessionStore, err := session.NewSessionStore(req, cookieStore, config.SessionSecret())
		require.NoError(t, err)
		sessionStore.SetDiscordUser(&model.DiscordUser{
			ID:       "123",
			Username: "test",
		})
		sessionStore.SetDiscordOAuthToken("aaa")

		handlerFunc.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusSeeOther, rr.Code, "Expected response code to be 303")

		_, err = sessionStore.GetDiscordUser()
		assert.Error(t, err)
		_, err = sessionStore.GetDiscordOAuthToken()
		assert.Error(t, err)
	})

	t.Run("そもそもログインされてなければ500を返す", func(t *testing.T) {
		cookieStore := sessions.NewCookieStore([]byte(""))
		indexService := &service.IndexService{
			CookieStore: cookieStore,
		} // Mock this service as per your implementation
		handler := NewDiscordOAuth2Handler(indexService)

		req, err := http.NewRequest("GET", "/logout/discord", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handlerFunc := http.HandlerFunc(handler.ServeHTTP)

		handlerFunc.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	// Add more tests for different scenarios, for example when session.NewSessionStore returns an error
}
