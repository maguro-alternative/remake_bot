package linepostdiscordchannel

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/web/shared/ctxvalue"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/views/guildid/line_post_discord_channel/internal"
	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

var (
	messageTypes = []string{
		"デフォルト",
		"RecipientAdd",
		"RecipientRemove",
		"DM通話開始",
		"チャンネル名変更",
		"チャンネルアイコン変更",
		"メッセージピン止め",
		"サーバー参加",
		"サーバーブースト",
		"サーバーレベル1",
		"サーバーレベル2",
		"サーバーレベル3",
		"サーバーフォロー",
		"サーバーディスカバリー失格メッセージ",
		"サーバーディスカバリー要件メッセージ",
		"スレッド作成",
		"リプライメッセージ",
		"スラッシュコマンド",
		"スレッドスタートメッセージ",
		"コンテンツメニュー",
	}
)

type LinePostDiscordChannelViewHandler struct {
	IndexService *service.IndexService
	Repo         repository.RepositoryFunc
}

func NewLinePostDiscordChannelViewHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
) *LinePostDiscordChannelViewHandler {
	return &LinePostDiscordChannelViewHandler{
		IndexService: indexService,
		Repo:         repo,
	}
}

func (g *LinePostDiscordChannelViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	categoryPositions := make(map[string]internal.DiscordChannel)
	var categoryIDTmps []string
	guildId := r.PathValue("guildId")
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	guild, err := g.IndexService.DiscordBotState.Guild(guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Discordサーバーの読み取りに失敗しました:"+err.Error())
		return
	}

	if guild.Members == nil {
		guild.Members, err = g.IndexService.DiscordSession.GuildMembers(guildId, "", 1000)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get guild members: "+err.Error())
			return
		}
	}

	if guild.Channels == nil {
		guild.Channels, err = g.IndexService.DiscordSession.GuildChannels(guildId)
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
	//[categoryID]map[channelPosition]channelName
	channelsInCategory := make(map[string][]internal.DiscordChannelSet)

	for _, channel := range guild.Channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}
		categoryIDTmps = append(categoryIDTmps, channel.ID)
		categoryPositions[channel.ID] = internal.DiscordChannel{
			ID:       channel.ID,
			Name:     channel.Name,
			Position: channel.Position,
		}
	}
	// カテゴリーなしのチャンネルを追加
	//channelsInCategory[""] = make([]internal.DiscordChannelSelect, len(guild.Channels)-1, len(guild.Channels))
	for _, channel := range guild.Channels {
		err = createCategoryInChannels(
			ctx,
			g.Repo,
			guild,
			channel,
			categoryPositions,
			channelsInCategory,
		)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "カテゴリーの読み取りに失敗しました:"+err.Error())
			return
		}
	}

	guildIconUrl := "https://cdn.discordapp.com/icons/" + guild.ID + "/" + guild.Icon + ".png"
	if guild.Icon == "" {
		guildIconUrl = "/static/img/discord-icon.jpg"
	}

	submitTag := components.CreateSubmitTag(discordPermissionData.Permission)
	accountVer := strings.Builder{}
	accountVer.WriteString(components.CreateDiscordAccountVer(discordPermissionData.User))
	accountVer.WriteString(components.CreateLineAccountVer(lineSession.User))

	htmlFormBuilder := internal.CreateLinePostDiscordChannelForm(
		categoryIDTmps,
		channelsInCategory,
		categoryPositions,
		guild,
		messageTypes,
	)

	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/views/guildid/line_post_discord_channel.html"))
	if err := tmpl.Execute(w, struct {
		Title        string
		AccountVer   template.HTML
		JsScriptTag  template.HTML
		SubmitTag    template.HTML
		GuildName    string
		GuildIconUrl string
		GuildID      string
		HTMLForm     template.HTML
	}{
		Title:        "DiscordからLINEへの送信設定",
		AccountVer:   template.HTML(accountVer.String()),
		JsScriptTag:  template.HTML(`<script src="/static/js/line_post_discord_channel.js"></script>`),
		SubmitTag:    template.HTML(submitTag),
		GuildName:    guild.Name,
		GuildIconUrl: guildIconUrl,
		GuildID:      guild.ID,
		HTMLForm:     template.HTML(htmlFormBuilder),
	}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "テンプレートの実行に失敗しました:"+err.Error())
	}
}

func createCategoryInChannels(
	ctx context.Context,
	repo repository.RepositoryFunc,
	guild *discordgo.Guild,
	channel *discordgo.Channel,
	categoryPositions map[string]internal.DiscordChannel,
	channelsInCategory map[string][]internal.DiscordChannelSet,
) error {
	if channel.Type == discordgo.ChannelTypeGuildForum {
		return nil
	}
	if channel.Type == discordgo.ChannelTypeGuildCategory {
		return nil
	}
	typeIcon := "🔊"
	if channel.Type == discordgo.ChannelTypeGuildText {
		typeIcon = "📝"
	}
	categoryPosition := categoryPositions[channel.ParentID]
	if len(channelsInCategory[categoryPosition.ID]) == 0 {
		channelsInCategory[categoryPosition.ID] = make([]internal.DiscordChannelSet, len(guild.Channels)+1)
	}
	discordChannel, err := repo.GetLinePostDiscordChannel(ctx, channel.ID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		slog.ErrorContext(ctx, "line_post_discord_channelの読み取りに失敗しました:"+err.Error())
		return err
	} else if err != nil {
		// チャンネルが存在しない場合は追加
		if err := repo.InsertLinePostDiscordChannel(ctx, channel.ID, guild.ID); err != nil {
			slog.ErrorContext(ctx, "line_post_discord_channelの追加に失敗しました:"+err.Error())
			return err
		}
		discordChannel = repository.LinePostDiscordChannel{
			Ng:         false,
			BotMessage: false,
		}
	}
	ngTypes, err := repo.GetLineNgDiscordMessageType(ctx, channel.ID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_typeの読み取りに失敗しました:"+err.Error())
		return err
	}
	ngDiscordUserIDs, err := repo.GetLineNgDiscordUserID(ctx, channel.ID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_user_idの読み取りに失敗しました:"+err.Error())
		return err
	}
	ngDiscordRoleIDs, err := repo.GetLineNgDiscordRoleID(ctx, channel.ID)
	if err != nil {
		slog.ErrorContext(ctx, "line_ng_discord_role_idの読み取りに失敗しました:"+err.Error())
		return err
	}
	channelsInCategory[categoryPosition.ID][channel.Position] = internal.DiscordChannelSet{
		ID:         channel.ID,
		Name:       fmt.Sprintf("%s %s", typeIcon, channel.Name),
		Ng:         discordChannel.Ng,
		BotMessage: discordChannel.BotMessage,
		NgTypes:    ngTypes,
	}
	channelsInCategory[categoryPosition.ID][channel.Position].NgUsers = append(channelsInCategory[categoryPosition.ID][channel.Position].NgUsers, ngDiscordUserIDs...)
	channelsInCategory[categoryPosition.ID][channel.Position].NgRoles = append(channelsInCategory[categoryPosition.ID][channel.Position].NgRoles, ngDiscordRoleIDs...)
	return nil
}
