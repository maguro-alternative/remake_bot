package middleware

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session"
)

// CSRFMiddleware protects against CSRF attacks using the double-submit cookie pattern.
// On safe methods (GET/HEAD/OPTIONS), it ensures a CSRF token exists in the session
// and sets a readable cookie so JavaScript can include the token in requests.
// On state-changing methods (POST/PUT/DELETE/PATCH), it validates that the
// X-CSRF-Token request header matches the token stored in the session.
func CSRFMiddleware(indexService service.IndexService) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionStore, err := session.NewSessionStore(r, indexService.CookieStore, config.SessionName())
			if err != nil {
				slog.ErrorContext(r.Context(), "CSRFミドルウェア: sessionの取得に失敗しました。", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// For safe methods, ensure a token exists and expose it via cookie
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				token, err := sessionStore.GetCSRFToken()
				if err != nil || token == "" {
					token = uuid.New().String()
					sessionStore.SetCSRFToken(token)
					if err := sessionStore.SessionSave(r, w); err != nil {
						slog.ErrorContext(r.Context(), "CSRFミドルウェア: sessionの保存に失敗しました。", "エラー:", err.Error())
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
				}
				// Set a JS-readable cookie with the CSRF token (not HttpOnly so JS can read it).
				// The actual validation is done against the session-stored token, not this cookie.
				http.SetCookie(w, &http.Cookie{
					Name:     "csrf_token",
					Value:    token,
					Path:     "/",
					SameSite: http.SameSiteStrictMode,
					Secure:   r.TLS != nil,
				})
				h.ServeHTTP(w, r)
				return
			}

			// For state-changing methods, validate the token
			sessionToken, err := sessionStore.GetCSRFToken()
			if err != nil || sessionToken == "" {
				slog.WarnContext(r.Context(), "CSRFミドルウェア: セッションにCSRFトークンがありません。")
				http.Error(w, "Forbidden - CSRF token missing", http.StatusForbidden)
				return
			}

			requestToken := r.Header.Get("X-CSRF-Token")
			if requestToken == "" {
				slog.WarnContext(r.Context(), "CSRFミドルウェア: リクエストにCSRFトークンがありません。")
				http.Error(w, "Forbidden - CSRF token missing", http.StatusForbidden)
				return
			}

			if requestToken != sessionToken {
				slog.WarnContext(r.Context(), "CSRFミドルウェア: CSRFトークンが一致しません。")
				http.Error(w, "Forbidden - CSRF token mismatch", http.StatusForbidden)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
