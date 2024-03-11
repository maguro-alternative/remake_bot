package guilds

import (
	"context"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/getoauth"
)

type userGuild struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Icon        string                   `json:"icon"`
	Owner       bool                     `json:"owner"`
	Permissions int64                    `json:"permissions"`
	Features    []discordgo.GuildFeature `json:"features"`
}

type GuildsViewHandler struct {
	IndexService *service.IndexService
}

func NewGuildsViewHandler(indexService *service.IndexService) *GuildsViewHandler {
	return &GuildsViewHandler{
		IndexService: indexService,
	}
}

func (g *GuildsViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	discordLoginUser, err := getoauth.GetDiscordOAuth(
		ctx,
		g.IndexService.CookieStore,
		r,
		config.SessionSecret(),
	)
	if err != nil {
		http.Redirect(w, r, "/auth/discord", http.StatusFound)
		return
	}
	var matchGuilds []discordgo.UserGuild
	botGuilds, err := g.IndexService.DiscordSession.UserGuilds(100, "", "")
	if err != nil {
		http.Error(w, "Not get bot guilds", http.StatusInternalServerError)
		return
	}
	userGuilds, err := getUserGuilds(discordLoginUser.Token)
	if err != nil {
		slog.InfoContext(ctx, "user guilds error: "+err.Error())
		http.Error(w, "Not get user guilds", http.StatusInternalServerError)
		return
	}
	for _, botGuild := range botGuilds {
		for _, userGuild := range userGuilds {
			if botGuild.ID == userGuild.ID {
				matchGuilds = append(matchGuilds, userGuild)
				break
			}
		}
	}
	htmlGuilds := ``
	for _, guild := range matchGuilds {
		if guild.Icon == "" {
			htmlGuilds += `
			<a href="/guild/` + guild.ID + `">
				<li>` + guild.Name + `</li>
			</a><br>
			`
			continue
		}
		htmlGuilds += `
		<a href="/guild/` + guild.ID + `">
			<img src="https://cdn.discordapp.com/icons/` + guild.ID + `/` + guild.Icon + `.png" alt="` + guild.Name + `">
			<li>` + guild.Name + `</li>
		</a><br>
		`
	}
	data := struct {
		Title       string
		JsScriptTag template.HTML
		Guilds      template.HTML
	}{
		Title:       "サーバー一覧",
		JsScriptTag: template.HTML(``),
		Guilds:      template.HTML(htmlGuilds),
	}
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/views/guilds/guilds.html"))
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getUserGuilds(token string) ([]discordgo.UserGuild, error) {
	url := "https://discord.com/api/users/@me/guilds"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var guilds []userGuild
	if err := json.NewDecoder(resp.Body).Decode(&guilds); err != nil {
		return nil, err
	}
	var userGuilds []discordgo.UserGuild
	for _, guild := range guilds {
		userGuilds = append(userGuilds, discordgo.UserGuild{
			ID:          guild.ID,
			Name:        guild.Name,
			Icon:        guild.Icon,
			Owner:       guild.Owner,
			Permissions: guild.Permissions,
			Features:    guild.Features,
		})
	}
	return userGuilds, nil
}
