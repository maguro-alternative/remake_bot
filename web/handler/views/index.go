package views

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/getoauth"
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
	discordAccountVer := strings.Builder{}
	discordLoginUser, err := getoauth.GetDiscordOAuth(
		ctx,
		g.IndexService.CookieStore,
		r,
		config.SessionSecret(),
	)
	if err != nil {
		discordAccountVer.WriteString(`
		<p>Discordアカウント</p>
		<button type="button" id="popover-btn" class="btn btn-primary">
			<a href="/" class="btn btn-primary">ログイン</a>
		</button>
		`)
	} else {
		discordAccountVer.WriteString(fmt.Sprintf(`
		<p>Discordアカウント: %s</p>
		<img src="https://cdn.discordapp.com/avatars/%s/%s.webp?size=64" alt="Discordアイコン">
		<button type="button" id="popover-btn" class="btn btn-primary">
			<a href="/" class="btn btn-primary">ログアウト</a>
		</button>
		`, discordLoginUser.User.Username, discordLoginUser.User.ID, discordLoginUser.User.Avatar))
	}
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/index.html"))
	err = tmpl.Execute(w, struct {
		Title             string
		LineAccountVer    template.HTML
		DiscordAccountVer template.HTML
		JsScriptTag       template.HTML
		BotName           string
		GuildId           string
	}{
		Title:             "Remake Bot",
		DiscordAccountVer: template.HTML(discordAccountVer.String()),
		BotName:           g.IndexService.DiscordSession.State.User.Username,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "template error: "+err.Error())
		return
	}
}
