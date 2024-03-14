package guildid

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/permission"
	"github.com/maguro-alternative/remake_bot/web/shared/session/getoauth"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
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
		slog.ErrorContext(ctx, "Discordサーバーの読み取りに失敗しました: ", "エラーメッセージ:", err.Error())
		return
	}
	statusCode, discordPermissionData, err := permission.CheckDiscordPermission(ctx, w, r, g.IndexService, guild, "")
	if err != nil {
		if statusCode == 302 {
			http.Redirect(w, r, "/login/discord", http.StatusFound)
			slog.InfoContext(ctx, "Redirect to /login/discord")
			return
		}
		if discordPermissionData.Permission == "" {
			http.Error(w, "Not permission", statusCode)
			slog.WarnContext(ctx, "権限のないアクセスがありました。", "エラーメッセージ:", err.Error(), "権限コード:", discordPermissionData.PermissionCode, "権限:", discordPermissionData.Permission)
			return
		}
	}
	// Lineの認証情報なしでもアクセス可能なためエラーレスポンスは出さない
	lineSession, err := getoauth.GetLineOAuth(g.IndexService.CookieStore, r, config.SessionSecret())
	if err != nil {
		lineSession = &model.LineOAuthSession{}
	}
	accountVer := strings.Builder{}
	accountVer.WriteString(components.CreateDiscordAccountVer(discordPermissionData.User))
	accountVer.WriteString(components.CreateLineAccountVer(lineSession.User))
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
		Title        string
		AccountVer   template.HTML
		JsScriptTag  template.HTML
		GuildID      string
		GuildName    string
		GuildIconUrl string
		SettingLinks template.HTML
	}{
		Title:        guild.Name + "の設定項目一覧",
		AccountVer:   template.HTML(accountVer.String()),
		GuildID:      guild.ID,
		GuildName:    guild.Name,
		GuildIconUrl: guildIconUrl,
		SettingLinks: template.HTML(settingLinks),
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "テンプレートの実行に失敗しました: ", "エラーメッセージ:", err.Error())
		return
	}
}
