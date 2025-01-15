package lineworkstoken

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/handler/views/guildid/lineworks_token/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/ctxvalue"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

type LineWorksTokenViewHandler struct {
	indexService *service.IndexService
	repo         repository.RepositoryFunc
}

func NewLineWorksTokenViewHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
) *LineWorksTokenViewHandler {
	return &LineWorksTokenViewHandler{
		indexService: indexService,
		repo:         repo,
	}
}
func (g *LineWorksTokenViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	categoryPositions := make(map[string]components.DiscordChannel)
	guildId := r.PathValue("guildId")
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	guild, err := g.indexService.DiscordBotState.Guild(guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Not get guild id: "+err.Error())
		return
	}

	if guild.Channels == nil {
		guild.Channels, err = g.indexService.DiscordSession.GuildChannels(guildId)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get guild channels: "+err.Error())
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
	lineWorksBot, err := g.repo.GetLineWorksBotByGuildID(ctx, guildId)
	if err != nil && err.Error() == "sql: no rows in result set" {
		slog.InfoContext(ctx, "line_works_botが存在しないため新規作成します。")
		err = g.repo.InsertLineWorksBot(ctx, &repository.LineWorksBot{
			GuildID: guildId,
			DefaultChannelID: guild.SystemChannelID,
			DebugMode: false,
		})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not insert line works bot: "+err.Error())
			return
		}
		err = g.repo.InsertLineWorksBotIV(ctx, &repository.LineWorksBotIV{
			GuildID: guildId,
		})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not insert line works bot iv: "+err.Error())
			return
		}
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_works_botの取得に失敗しました: "+err.Error())
		return
	}
	lineWorksBotInfo, err := g.repo.GetLineWorksBotInfoByGuildID(ctx, guildId)
	if err != nil && err.Error() == "sql: no rows in result set" {
		slog.InfoContext(ctx, "line_works_bot_infoが存在しないため新規作成します。")
		err = g.repo.InsertLineWorksBotInfo(ctx, &repository.LineWorksBotInfo{
			GuildID: guildId,
		})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not insert line works bot info: "+err.Error())
			return
		}
		err = g.repo.InsertLineWorksBotInfoIV(ctx, &repository.LineWorksBotInfoIV{
			GuildID: guildId,
		})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not insert line works bot info iv: "+err.Error())
			return
		}
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_works_bot_infoの取得に失敗しました: "+err.Error())
		return
	}

	lineWorksBotByte := internal.LineWorksBotByteEntered{
		LineWorksBotToken:    lineWorksBot.LineWorksBotToken,
		LineWorksRefreshToken: lineWorksBot.LineWorksRefreshToken,
		LineWorksGroupID:      lineWorksBot.LineWorksGroupID,
		LineWorksBotID:        lineWorksBot.LineWorksBotID,
		LineWorksBotSecret:    lineWorksBot.LineWorksBotSecret,
		DebugMode:             lineWorksBot.DebugMode,
	}
	lineWorksBotInfoByte := internal.LineWorksBotInfoEntered{
		LineWorksClientID:     lineWorksBotInfo.LineWorksClientID,
		LineWorksClientSecret: lineWorksBotInfo.LineWorksClientSecret,
		LineWorksServiceAccount: lineWorksBotInfo.LineWorksServiceAccount,
		LineWorksPrivateKey:    lineWorksBotInfo.LineWorksPrivateKey,
		LineWorksDomainID:      lineWorksBotInfo.LineWorksDomainID,
		LineWorksAdminID:       lineWorksBotInfo.LineWorksAdminID,
	}
	lineWorksBotForm := internal.EnteredLineWorksBotForm(lineWorksBotByte, lineWorksBotInfoByte)
	guildIconUrl := "https://cdn.discordapp.com/icons/" + guild.ID + "/" + guild.Icon + ".png"
	if guild.Icon == "" {
		guildIconUrl = "/static/img/discord-icon.jpg"
	}

	submitTag := components.CreateSubmitTag(discordPermissionData.Permission)
	accountVer := strings.Builder{}
	accountVer.WriteString(components.CreateDiscordAccountVer(discordPermissionData.User))
	accountVer.WriteString(components.CreateLineAccountVer(lineSession.User))
	htmlSelectChannelBuilders := components.CreateSelectChennelOptions(
		categoryIDTmps,
		lineWorksBot.DefaultChannelID,
		channelsInCategory,
		categoryPositions,
	)
	data := struct {
		Title 				 string
		GuildID 			 string
		GuildIconUrl 		 string
		GuildName 			 string
		AccountVer 			 template.HTML
		JsScriptTag 		 template.HTML
		SubmitTag 			 template.HTML
		LineWorksBotForm 	 internal.LineWorksBotInfo
		Channels 			 template.HTML
	}{
		Title: 				 "LineWorks Bot設定",
		GuildID: 			 guildId,
		GuildIconUrl: 		 guildIconUrl,
		GuildName: 			 guild.Name,
		AccountVer: 		 template.HTML(accountVer.String()),
		JsScriptTag: 		 template.HTML(`<script src="/static/js/lineworks_token.js"></script>`),
		SubmitTag: 			 template.HTML(submitTag),
		LineWorksBotForm: 	 lineWorksBotForm,
		Channels: 			 template.HTML(htmlSelectChannelBuilders),
	}
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/views/guildid/lineworks_token.html"))
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "テンプレートの実行に失敗しました: "+err.Error())
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
