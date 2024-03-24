package views

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/getoauth"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

type IndexViewHandler struct {
	IndexService *service.IndexService
}

func NewIndexViewHandler(indexService *service.IndexService) *IndexViewHandler {
	return &IndexViewHandler{
		IndexService: indexService,
	}
}

func (g *IndexViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	oauthStore := getoauth.NewOAuthStore(g.IndexService.CookieStore, config.SessionSecret())
	// Discordの認証情報なしでもアクセス可能なためエラーレスポンスは出さない
	discordLoginUser, err := oauthStore.GetDiscordOAuth(ctx, r)
	if err != nil {
		discordLoginUser = &model.DiscordOAuthSession{}
	}
	// Lineの認証情報なしでもアクセス可能なためエラーレスポンスは出さない
	lineSession, err := oauthStore.GetLineOAuth(r)
	if err != nil {
		lineSession = &model.LineOAuthSession{}
	}
	accountVer := strings.Builder{}
	accountVer.WriteString(components.CreateDiscordAccountVer(discordLoginUser.User))
	accountVer.WriteString(components.CreateLineAccountVer(lineSession.User))
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/index.html"))
	err = tmpl.Execute(w, struct {
		Title       string
		AccountVer  template.HTML
		JsScriptTag template.HTML
		BotName     string
		GuildId     string
	}{
		Title:      "トップページ",
		AccountVer: template.HTML(accountVer.String()),
		BotName:    g.IndexService.DiscordSession.State.User.Username,
		GuildId:    lineSession.DiscordGuildID,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "template error: "+err.Error())
		return
	}
}
