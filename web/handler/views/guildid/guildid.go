package guildid

import (
	"context"
	"html/template"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/permission"
)

type GuildIDViewHandler struct {
	IndexService *service.IndexService
}

func NewGuildIDViewHandler(indexService *service.IndexService) *GuildIDViewHandler {
	return &GuildIDViewHandler{
		IndexService: indexService,
	}
}

func (g *GuildIDViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	guildId := r.PathValue("guildId")
	guild, err := g.IndexService.DiscordSession.State.Guild(guildId)
	if err != nil {
		http.Error(w, "Not get guild id", http.StatusInternalServerError)
		return
	}
	statusCode, err := permission.CheckDiscordPermission(ctx, w, r, g.IndexService, guild, "line_bot")
	if err != nil {
		if statusCode == 302 {
			http.Redirect(w, r, "/auth/discord", http.StatusFound)
			return
		}
		if statusCode != 200 {
			http.Error(w, "Not get guild id", statusCode)
			return
		}
		http.Error(w, "Not get guild id", http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(template.ParseFiles("web/templates/views/guildid.html"))
	tmpl.Execute(w, nil)
}
