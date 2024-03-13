package guildid

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Discordサーバーの読み取りに失敗しました: "+err.Error())
		return
	}
	statusCode, discordPermissionData, err := permission.CheckDiscordPermission(ctx, w, r, g.IndexService, guild, "line_bot")
	if err != nil {
		if statusCode == 302 {
			http.Redirect(w, r, "/login/discord", http.StatusFound)
			slog.InfoContext(ctx, "Redirect to /login/discord")
			return
		}
		http.Error(w, "Not get guild id", statusCode)
		slog.WarnContext(ctx, "権限のないアクセスがありました: "+err.Error())
		return
	}
	discordAccountVer := strings.Builder{}
	discordAccountVer.WriteString(fmt.Sprintf(`
	<p>Discordアカウント: %s</p>
	<img src="https://cdn.discordapp.com/avatars/%s/%s.webp?size=64" alt="Discordアイコン">
	<button type="button" id="popover-btn" class="btn btn-primary">
		<a href="/logout/discord" class="btn btn-primary">ログアウト</a>
	</button>
	`, discordPermissionData.User.Username, discordPermissionData.User.ID, discordPermissionData.User.Avatar))
	guildIconUrl := "https://cdn.discordapp.com/icons/" + guild.ID + "/" + guild.Icon + ".png"
	if guild.Icon == "" {
		guildIconUrl = "/static/img/discord-icon.jpg"
	}
	if discordPermissionData.PermissionCode&8 != 0 {
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
		Title             string
		LineAccountVer    template.HTML
		DiscordAccountVer template.HTML
		JsScriptTag       template.HTML
		GuildID           string
		GuildName         string
		GuildIconUrl      string
		SettingLinks      template.HTML
	}{
		Title:             guild.Name + "の設定項目一覧",
		DiscordAccountVer: template.HTML(discordAccountVer.String()),
		GuildID:           guild.ID,
		GuildName:         guild.Name,
		GuildIconUrl:      guildIconUrl,
		SettingLinks:      template.HTML(settingLinks),
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "テンプレートの実行に失敗しました: "+err.Error())
		return
	}
}
