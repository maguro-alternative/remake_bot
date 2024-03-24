package middleware

import (
	"context"
	"encoding/gob"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/getoauth"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

func init() {
	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	gob.Register(&model.DiscordUser{})
}

func DiscordOAuthCheckMiddleware(indexService service.IndexService) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			client := &http.Client{}
			oauthStore := getoauth.NewOAuthStore(indexService.CookieStore, config.SessionSecret())
			discordLoginUser, err := oauthStore.GetDiscordOAuth(ctx, r)
			if err != nil {
				http.Redirect(w, r, "/login/discord", http.StatusFound)
				return
			}
			req, err := http.NewRequestWithContext(ctx, "GET", "https://discord.com/api/users/@me", nil)
			if err != nil {
				http.Error(w, "Not get user", http.StatusInternalServerError)
				return
			}
			req.Header.Set("Authorization", "Bearer "+discordLoginUser.Token)
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode != http.StatusOK {
				http.Redirect(w, r, "/login/discord", http.StatusFound)
				return
			}
			defer resp.Body.Close()
			h.ServeHTTP(w, r)
		})
	}
}
