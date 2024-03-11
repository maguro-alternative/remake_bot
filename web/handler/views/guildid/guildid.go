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
	var settingLinks string
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
	statusCode, permissionCode, err := permission.CheckDiscordPermission(ctx, w, r, g.IndexService, guild, "line_bot")
	if err != nil {
		if statusCode == 302 {
			http.Redirect(w, r, "/auth/discord", http.StatusFound)
			return
		}
		http.Error(w, "Not get guild id", statusCode)
		return
	}
	if permissionCode&8 != 0 {
		settingLinks += `
			管理者です。<br/>
			<a href="/guild/` + guild.ID + `/admin" class="btn btn-primary">管理者設定</a>
			<br/>
		`
	}
	settingLinks += `
		<a href="/guild/` + guild.ID + `/line-post-discord-channel" class="btn btn-primary">LINEへの送信設定</a>
		<a href="/guild/` + guild.ID + `/linetoken" class="btn btn-primary">LINEBOTおよびグループ設定</a>
		<a href="/guild/` + guild.ID + `/vc-signal" class="btn btn-primary">ボイスチャンネルの通知設定</a>
		<a href="/guild/` + guild.ID + `/webhook" class="btn btn-primary">webhookの送信設定</a>
	`
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/views/guildid.html"))
	err = tmpl.Execute(w, struct {
		Title        string
		JsScriptTag  template.HTML
		GuildID      string
		GuildName    string
		GuildIcon    string
		SettingLinks template.HTML
	}{
		Title:        guild.Name + "の設定項目一覧",
		JsScriptTag:  template.HTML(``),
		GuildID:      guild.ID,
		GuildName:    guild.Name,
		GuildIcon:    guild.Icon,
		SettingLinks: template.HTML(settingLinks),
	})
	if err != nil {
		http.Error(w, "Template error "+err.Error(), http.StatusInternalServerError)
	}
}
