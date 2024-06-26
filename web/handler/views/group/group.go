package group

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/shared/ctxvalue"

	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

type LineGroupViewHandler struct {
	indexService *service.IndexService
	repo         repository.RepositoryFunc
}

func NewLineGroupViewHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
) *LineGroupViewHandler {
	return &LineGroupViewHandler{
		indexService: indexService,
		repo:         repo,
	}
}

func (g *LineGroupViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	guildId := r.PathValue("guildId")
	categoryPositions := make(map[string]components.DiscordChannel)
	guild, err := g.indexService.DiscordBotState.Guild(guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Discordサーバーの読み取りに失敗しました: ", "エラーメッセージ:", err.Error())
		return
	}

	if guild.Channels == nil {
		guild.Channels, err = g.indexService.DiscordSession.GuildChannels(guildId)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "サーバー内のチャンネルの取得に失敗しました。 ", "エラー", err.Error())
			return
		}
	}

	// Discordの認証情報なしでもアクセス可能なためエラーレスポンスは出さない
	discordPermissionData, err := ctxvalue.DiscordPermissionFromContext(ctx)
	if err != nil {
		discordPermissionData = &model.DiscordPermissionData{}
	}
	lineSession, err := ctxvalue.LineUserFromContext(ctx)
	if err != nil {
		http.Redirect(w, r, "/login/line", http.StatusFound)
		slog.InfoContext(ctx, "Redirect to /login/line")
		return
	}

	// カテゴリーのチャンネルを取得
	//[categoryID]map[channelPosition]channelName
	channelsInCategory := make(map[string][]components.DiscordChannelSelect)
	var categoryIDTmps []string
	for _, channel := range guild.Channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}
		// カテゴリーIDの順番を一時保存(Goではmapの順番が保証されないため)
		categoryIDTmps = append(categoryIDTmps, channel.ID)
		// カテゴリーごとに連想配列を作成
		categoryPositions[channel.ID] = components.DiscordChannel{
			ID:       channel.ID,
			Name:     channel.Name,
			Position: channel.Position,
		}
	}
	// カテゴリーなしのチャンネルを追加
	//channelsInCategory[""] = make([]components.DiscordChannelSelect, len(guild.Channels)-1, len(guild.Channels))
	for _, channel := range guild.Channels {
		createChannelsInCategory(
			guild,
			channel,
			categoryPositions,
			channelsInCategory,
		)
	}

	lineBot, err := g.repo.GetAllColumnsLineBotByGuildID(ctx, guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_botの取得に失敗しました:"+err.Error())
		return
	}
	htmlSelectChannels := components.CreateSelectChennelOptions(
		categoryIDTmps,
		lineBot.DefaultChannelID,
		channelsInCategory,
		categoryPositions,
	)

	accountVer := strings.Builder{}
	accountVer.WriteString(components.CreateDiscordAccountVer(discordPermissionData.User))
	accountVer.WriteString(components.CreateLineAccountVer(lineSession.User))
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/views/group/group.html"))
	err = tmpl.Execute(w, struct {
		Title       string
		AccountVer  template.HTML
		JsScriptTag template.HTML
		Channels    template.HTML
	}{
		Title:       "グループ",
		AccountVer:  template.HTML(accountVer.String()),
		JsScriptTag: template.HTML(`<script src="/static/js/group.js"></script>`),
		Channels:    template.HTML(htmlSelectChannels),
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "テンプレートの描画に失敗しました: ", "エラーメッセージ:", err.Error())
	}
}

func createChannelsInCategory(
	guild *discordgo.Guild,
	channel *discordgo.Channel,
	categoryPositions map[string]components.DiscordChannel,
	channelsInCategory map[string][]components.DiscordChannelSelect,
) {
	// カテゴリー、フォーラムチャンネルはスキップ
	if channel.Type == discordgo.ChannelTypeGuildForum {
		return
	}
	if channel.Type == discordgo.ChannelTypeGuildCategory {
		return
	}
	typeIcon := "🔊"
	if channel.Type == discordgo.ChannelTypeGuildText {
		typeIcon = "📝"
	}
	categoryPosition := categoryPositions[channel.ParentID]
	// まだチャンネルがない場合は初期化
	if len(channelsInCategory[categoryPosition.ID]) == 0 {
		channelsInCategory[categoryPosition.ID] = make([]components.DiscordChannelSelect, len(guild.Channels)+1)
	}
	channelsInCategory[categoryPosition.ID][channel.Position] = components.DiscordChannelSelect{
		ID:   channel.ID,
		Name: fmt.Sprintf("%s:%s:%s", categoryPosition.Name, typeIcon, channel.Name),
	}
	if categoryPosition.ID == "" {
		channelsInCategory[categoryPosition.ID][channel.Position] = components.DiscordChannelSelect{
			ID:   channel.ID,
			Name: fmt.Sprintf("カテゴリーなし:%s:%s", typeIcon, channel.Name),
		}
	}
}
