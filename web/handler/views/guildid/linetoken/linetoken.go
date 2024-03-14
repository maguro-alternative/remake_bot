package guildid

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/web/handler/views/guildid/linetoken/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/permission"
)

type LineTokenViewHandler struct {
	IndexService *service.IndexService
}

func NewLineTokenViewHandler(indexService *service.IndexService) *LineTokenViewHandler {
	return &LineTokenViewHandler{
		IndexService: indexService,
	}
}

func (g *LineTokenViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	repo := internal.NewRepository(g.IndexService.DB)
	categoryPositions := make(map[string]internal.DiscordChannel)
	guildId := r.PathValue("guildId")
	var submitTag string
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	guild, err := g.IndexService.DiscordSession.State.Guild(guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Not get guild id: "+err.Error())
		return
	}
	statusCode, discordPermissionData, err := permission.CheckDiscordPermission(ctx, w, r, g.IndexService, guild, "line_bot")
	if err != nil {
		if statusCode == http.StatusFound {
			http.Redirect(w, r, "/login/discord", http.StatusFound)
			slog.InfoContext(ctx, "Redirect to /login/discord "+err.Error())
			return
		}
		if discordPermissionData.Permission == "" {
			http.Error(w, "Not permission", statusCode)
			slog.WarnContext(ctx, "権限のないアクセスがありました。 "+err.Error())
			return
		}
	}
	// カテゴリーのチャンネルを取得
	//[categoryID]map[channelPosition]channelName
	channelsInCategory := make(map[string][]internal.DiscordChannelSelect)
	var categoryIDTmps []string
	for _, channel := range guild.Channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}
		// カテゴリーIDの順番を一時保存(Goではmapの順番が保証されないため)
		categoryIDTmps = append(categoryIDTmps, channel.ID)
		// カテゴリーごとに連想配列を作成
		categoryPositions[channel.ID] = internal.DiscordChannel{
			ID:       channel.ID,
			Name:     channel.Name,
			Position: channel.Position,
		}
	}
	// カテゴリーなしのチャンネルを追加
	//channelsInCategory[""] = make([]internal.DiscordChannelSelect, len(guild.Channels)-1, len(guild.Channels))
	for _, channel := range guild.Channels {
		// カテゴリー、フォーラムチャンネルはスキップ
		if channel.Type == discordgo.ChannelTypeGuildForum {
			continue
		}
		if channel.Type == discordgo.ChannelTypeGuildCategory {
			continue
		}
		typeIcon := "🔊"
		if channel.Type == discordgo.ChannelTypeGuildText {
			typeIcon = "📝"
		}
		categoryPosition := categoryPositions[channel.ParentID]
		// まだチャンネルがない場合は初期化
		if len(channelsInCategory[categoryPosition.ID]) == 0 {
			channelsInCategory[categoryPosition.ID] = make([]internal.DiscordChannelSelect, len(guild.Channels)-2, len(guild.Channels))
		}
		channelsInCategory[categoryPosition.ID][channel.Position] = internal.DiscordChannelSelect{
			ID:   channel.ID,
			Name: fmt.Sprintf("%s:%s:%s", categoryPosition.Name, typeIcon, channel.Name),
		}
		if categoryPosition.ID == "" {
			channelsInCategory[categoryPosition.ID][channel.Position] = internal.DiscordChannelSelect{
				ID:   channel.ID,
				Name: fmt.Sprintf("カテゴリーなし:%s:%s", typeIcon, channel.Name),
			}
		}
	}
	var lineNotifyTokenEntered, lineBotTokenEntered, lineBotSecretEntered, lineGroupIDEntered, lineClientIDEntered, lineClientSecretEntered string
	lineBot, err := repo.GetLineBot(ctx, guildId)
	if err != nil && err.Error() == "sql: no rows in result set" {
		err = repo.InsertLineBot(ctx, &internal.LineBot{
			GuildID:          guildId,
			DefaultChannelID: guild.SystemChannelID,
			DebugMode:        false,
		})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "line_botの作成に失敗しました:"+err.Error())
			return
		}
		err = repo.InsertLineBotIv(ctx, &internal.LineBotIv{
			GuildID: guildId,
		})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "line_bot_ivの作成に失敗しました:"+err.Error())
			return
		}
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_botの取得に失敗しました:"+err.Error())
		return
	}
	if lineBot.LineNotifyToken != nil {
		lineNotifyTokenEntered = "入力済み"
	}
	if lineBot.LineBotToken != nil {
		lineBotTokenEntered = "入力済み"
	}
	if lineBot.LineBotSecret != nil {
		lineBotSecretEntered = "入力済み"
	}
	if lineBot.LineGroupID != nil {
		lineGroupIDEntered = "入力済み"
	}
	if lineBot.LineClientID != nil {
		lineClientIDEntered = "入力済み"
	}
	if lineBot.LineClientSecret != nil {
		lineClientSecretEntered = "入力済み"
	}

	if discordPermissionData.Permission == "write" || discordPermissionData.Permission == "all" {
		submitTag = `<input type="submit" value="送信">`
	}

	discordAccountVer := strings.Builder{}
	discordAccountVer.WriteString(fmt.Sprintf(`
	<p>Discordアカウント: %s</p>
	<img src="https://cdn.discordapp.com/avatars/%s/%s.webp?size=64" alt="Discordアイコン">
	<button type="button" id="popover-btn" class="btn btn-primary">
		<a href="/logout/discord" class="btn btn-primary">ログアウト</a>
	</button>
	`, discordPermissionData.User.Username, discordPermissionData.User.ID, discordPermissionData.User.Avatar))
	htmlSelectChannelBuilders := strings.Builder{}
	categoryOptions := make([]strings.Builder, len(categoryIDTmps)+1)
	var categoryIndex int
	for categoryID, channels := range channelsInCategory {
		for i, categoryIDTmp := range categoryIDTmps {
			if categoryID == "" {
				categoryIndex = len(categoryIDTmps)
				break
			}
			if categoryIDTmp == categoryID {
				categoryIndex = i
				break
			}
		}
		for _, channelSelect := range channels {
			if channelSelect.ID == "" {
				continue
			}
			if lineBot.DefaultChannelID == channelSelect.ID {
				categoryOptions[categoryIndex].WriteString(fmt.Sprintf(`<option value="%s" selected>%s</option>`, channelSelect.ID, channelSelect.Name))
				continue
			}
			categoryOptions[categoryIndex].WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, channelSelect.ID, channelSelect.Name))
		}
	}
	for _, categoryOption := range categoryOptions {
		htmlSelectChannelBuilders.WriteString(categoryOption.String())
	}
	data := struct {
		Title                   string
		LineAccountVer          template.HTML
		DiscordAccountVer       template.HTML
		JsScriptTag             template.HTML
		SubmitTag               template.HTML
		LineNotifyTokenEntered  string
		LineBotTokenEntered     string
		LineBotSecretEntered    string
		LineGroupIDEntered      string
		LineClientIDEntered     string
		LineClientSecretEntered string
		Channels                template.HTML
	}{
		Title:                   "LineBotの設定",
		JsScriptTag:             template.HTML(`<script src="/static/js/linetoken.js"></script>`),
		DiscordAccountVer:       template.HTML(discordAccountVer.String()),
		SubmitTag:               template.HTML(submitTag),
		LineNotifyTokenEntered:  lineNotifyTokenEntered,
		LineBotTokenEntered:     lineBotTokenEntered,
		LineBotSecretEntered:    lineBotSecretEntered,
		LineGroupIDEntered:      lineGroupIDEntered,
		LineClientIDEntered:     lineClientIDEntered,
		LineClientSecretEntered: lineClientSecretEntered,
		Channels:                template.HTML(htmlSelectChannelBuilders.String()),
	}
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/views/guildid/linetoken.html"))
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "テンプレートの実行に失敗しました:"+err.Error())
	}
}
