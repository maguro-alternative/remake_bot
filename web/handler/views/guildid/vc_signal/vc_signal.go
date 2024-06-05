package vcsignal

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
	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/handler/views/guildid/vc_signal/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

type VcSignalViewHandler struct {
	indexService *service.IndexService
	repo         repository.RepositoryFunc
}

func NewVcSignalViewHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
) *VcSignalViewHandler {
	return &VcSignalViewHandler{
		indexService: indexService,
		repo:         repo,
	}
}

func (h *VcSignalViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	categoryPositions := make(map[string]components.DiscordChannel)
	var categoryIDTmps []string
	guildId := r.PathValue("guildId")
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	guild, err := h.indexService.DiscordBotState.Guild(guildId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Discordサーバーの読み取りに失敗しました:"+err.Error())
		return
	}

	if guild.Members == nil {
		guild.Members, err = h.indexService.DiscordSession.GuildMembers(guildId, "", 1000)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "Not get guild members: "+err.Error())
			return
		}
	}

	if guild.Channels == nil {
		guild.Channels, err = h.indexService.DiscordSession.GuildChannels(guildId)
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
	vcChannels := make(map[string][]internal.VcChannelSet)
	//[categoryID]map[channelPosition]channelName
	channelsInCategory := make(map[string][]components.DiscordChannelSelect)

	for _, channel := range guild.Channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}
		categoryIDTmps = append(categoryIDTmps, channel.ID)
		categoryPositions[channel.ID] = components.DiscordChannel{
			ID:       channel.ID,
			Name:     channel.Name,
			Position: channel.Position,
		}
	}
	// カテゴリーなしのチャンネルを追加
	//vcChannels[""] = make([]internal.DiscordChannelSelect, len(guild.Channels)-1, len(guild.Channels))
	for _, channel := range guild.Channels {
		err = createCategoryInChannels(
			ctx,
			h.repo,
			guild,
			channel,
			categoryPositions,
			vcChannels,
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

	htmlFormBuilder := internal.CreateVcSignalForm(
		categoryIDTmps,
		vcChannels,
		channelsInCategory,
		categoryPositions,
		guild,
	)

	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/views/guildid/vc_signal.html"))
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
		JsScriptTag:  template.HTML(`<script src="/static/js/vc_signal.js"></script>`),
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
	categoryPositions map[string]components.DiscordChannel,
	vcChannelSets map[string][]internal.VcChannelSet,
	channelsInCategory map[string][]components.DiscordChannelSelect,
) error {
	typeIcon := "📝"
	if channel.Type == discordgo.ChannelTypeGuildForum {
		return nil
	}
	if channel.Type == discordgo.ChannelTypeGuildCategory {
		return nil
	}
	if channel.Type == discordgo.ChannelTypeGuildVoice {
		typeIcon = "🔊"
	}
	if len(channelsInCategory[channel.ParentID]) == 0 {
		channelsInCategory[channel.ParentID] = make([]components.DiscordChannelSelect, len(guild.Channels)+1)
	}
	channelsInCategory[channel.ParentID][channel.Position] = components.DiscordChannelSelect{
		ID:   channel.ID,
		Name: fmt.Sprintf("%s:%s:%s", categoryPositions[channel.ParentID].Name, typeIcon, channel.Name),
	}
	if channel.Type == discordgo.ChannelTypeGuildText {
		return nil
	}
	categoryPosition := categoryPositions[channel.ParentID]
	if len(vcChannelSets[categoryPosition.ID]) == 0 {
		vcChannelSets[categoryPosition.ID] = make([]internal.VcChannelSet, len(guild.Channels)+1)
	}
	discordChannel, err := repo.GetVcSignalChannelAllColumnByVcChannelID(ctx, channel.ID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		slog.ErrorContext(ctx, "vc_signal_channelの読み取りに失敗しました", "エラー", err.Error())
		return err
	} else if err != nil {
		// チャンネルが存在しない場合は追加
		err = repo.InsertVcSignalChannel(ctx, channel.ID, guild.ID, guild.SystemChannelID)
		if err != nil {
			slog.ErrorContext(ctx, "vc_signal_channelの追加に失敗しました", "エラー", err.Error())
			return err
		}
		discordChannel = &repository.VcSignalChannelAllColumn{
			VcChannelID:     channel.ID,
			GuildID:         guild.ID,
			SendSignal:      true,
			SendChannelID:   guild.SystemChannelID,
			JoinBot:         false,
			EveryoneMention: true,
		}
	}
	ngDiscordUserIDs, err := repo.GetVcSignalNgUserIDsByVcChannelID(ctx, channel.ID)
	if err != nil {
		slog.ErrorContext(ctx, "vc_signal_ng_user_idの読み取りに失敗しました", "エラー", err.Error())
		return err
	}
	ngDiscordRoleIDs, err := repo.GetVcSignalNgRoleIDsByVcChannelID(ctx, channel.ID)
	if err != nil {
		slog.ErrorContext(ctx, "vc_signal_ng_role_idの読み取りに失敗しました", "エラー", err.Error())
		return err
	}
	mentionDiscordUserIDs, err := repo.GetVcSignalMentionUserIDsByVcChannelID(ctx, channel.ID)
	if err != nil {
		slog.ErrorContext(ctx, "vc_signal_ng_user_idの読み取りに失敗しました", "エラー", err.Error())
		return err
	}
	mentionDiscordRoleIDs, err := repo.GetVcSignalMentionRoleIDsByVcChannelID(ctx, channel.ID)
	if err != nil {
		slog.ErrorContext(ctx, "vc_signal_ng_role_idの読み取りに失敗しました", "エラー", err.Error())
		return err
	}
	vcChannelSets[categoryPosition.ID][channel.Position] = internal.VcChannelSet{
		ID:              channel.ID,
		Name:            fmt.Sprintf("%s %s", typeIcon, channel.Name),
		SendSignal:      discordChannel.SendSignal,
		SendChannelID:   discordChannel.SendChannelID,
		JoinBot:         discordChannel.JoinBot,
		EveryoneMention: discordChannel.EveryoneMention,
	}
	vcChannelSets[categoryPosition.ID][channel.Position].NgUsers = append(vcChannelSets[categoryPosition.ID][channel.Position].NgUsers, ngDiscordUserIDs...)
	vcChannelSets[categoryPosition.ID][channel.Position].NgRoles = append(vcChannelSets[categoryPosition.ID][channel.Position].NgRoles, ngDiscordRoleIDs...)
	vcChannelSets[categoryPosition.ID][channel.Position].MentionUsers = append(vcChannelSets[categoryPosition.ID][channel.Position].MentionUsers, mentionDiscordUserIDs...)
	vcChannelSets[categoryPosition.ID][channel.Position].MentionRoles = append(vcChannelSets[categoryPosition.ID][channel.Position].MentionRoles, mentionDiscordRoleIDs...)
	return nil
}
