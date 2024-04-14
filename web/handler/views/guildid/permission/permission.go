package permission

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/pkg/ctxvalue"

	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

type PermissionViewHandler struct {
	IndexService *service.IndexService
	Repo         repository.RepositoryFunc
}

func NewPermissionViewHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
) *PermissionViewHandler {
	return &PermissionViewHandler{
		IndexService: indexService,
		Repo:         repo,
	}
}

func (h *PermissionViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	guildId := r.PathValue("guildId")
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		slog.ErrorContext(ctx, "/guild/permission Method Not Allowed")
		return
	}
	var componentPermissionIDs []components.PermissionID

	guild, err := h.IndexService.DiscordBotState.Guild(guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Not get guild id: "+err.Error())
		return
	}

	if guild.Members == nil {
		guild.Members, err = h.IndexService.DiscordSession.GuildMembers(guildId, "", 1000)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get guild members: "+err.Error())
			return
		}
	}

	discordPermissionData, err := ctxvalue.DiscordPermissionFromContext(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Discord認証情報の取得に失敗しました: ", "エラーメッセージ:", err.Error())
		return
	}
	// Lineの認証情報なしでもアクセス可能なためエラーレスポンスは出さない
	lineSession, err := ctxvalue.LineUserFromContext(ctx)
	if err != nil {
		lineSession = &model.LineOAuthSession{}
	}

	permissionCodes, err := h.Repo.GetPermissionCodes(ctx, guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "permissions_codeの取得に失敗しました。", "エラー:", err.Error())
		return
	}

	permissionIDs, err := h.Repo.GetGuildPermissionIDsAllColumns(ctx, guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "permissions_idの取得に失敗しました。", "エラー:", err.Error())
		return
	}

	permissionForm := strings.Builder{}
	for _, permissionID := range permissionIDs {
		componentPermissionIDs = append(
			componentPermissionIDs,
			components.PermissionID{
				GuildID:    guildId,
				Type:       permissionID.Type,
				TargetType: permissionID.TargetType,
				TargetID:   permissionID.TargetID,
				Permission: permissionID.Permission,
			},
		)
	}
	for _, permissionCode := range permissionCodes {
		permissionForm.WriteString(
			components.CreatePermissionCodeForm(guildId, components.PermissionCode{
				GuildID: guildId,
				Type:    permissionCode.Type,
				Code:    permissionCode.Code,
			}),
		)
		permissionForm.WriteString(components.CreatePermissionSelectForm(
			guild,
			componentPermissionIDs,
			permissionCode.Type,
		))
	}

	guildIconUrl := "https://cdn.discordapp.com/icons/" + guild.ID + "/" + guild.Icon + ".png"
	if guild.Icon == "" {
		guildIconUrl = "/static/img/discord-icon.jpg"
	}

	accountVer := strings.Builder{}
	accountVer.WriteString(components.CreateDiscordAccountVer(discordPermissionData.User))
	accountVer.WriteString(components.CreateLineAccountVer(lineSession.User))

	data := struct {
		Title          string
		JsScriptTag    template.HTML
		GuildIconUrl   string
		GuildName      string
		GuildID        string
		AccountVer     template.HTML
		PermissionForm template.HTML
	}{
		Title:          "権限設定",
		JsScriptTag:    template.HTML(`<script src="/static/js/permission.js"></script>`),
		GuildIconUrl:   guildIconUrl,
		GuildName:      guild.Name,
		GuildID:        guildId,
		AccountVer:     template.HTML(accountVer.String()),
		PermissionForm: template.HTML(permissionForm.String()),
	}
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/views/guildid/permission.html"))
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "テンプレートの実行に失敗しました", "エラー:", err.Error())
	}
}
