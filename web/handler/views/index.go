package views

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/maguro-alternative/remake_bot/pkg/ctxvalue"

	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
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

	// Discordの認証情報なしでもアクセス可能なためエラーレスポンスは出さない
	discordLoginUser, err := ctxvalue.DiscordUserFromContext(ctx)
	if err != nil {
		discordLoginUser = &model.DiscordOAuthSession{
			User: model.DiscordUser{},
		}
	}
	// Lineの認証情報なしでもアクセス可能なためエラーレスポンスは出さない
	lineSession, err := ctxvalue.LineUserFromContext(ctx)
	if err != nil {
		lineSession = &model.LineOAuthSession{
			User: model.LineIdTokenUser{},
		}
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
		BotName:    g.IndexService.DiscordBotState.User.Username,
		GuildId:    lineSession.DiscordGuildID,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "template error: "+err.Error())
		return
	}
}
