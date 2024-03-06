package guildid

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/web/handler/views/guildid/linetoken/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/session/getoauth"
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
	discordLoginUser, err := getoauth.GetDiscordOAuth(
		g.IndexService.CookieStore,
		r,
		config.SessionSecret(),
	)
	if err != nil {
		http.Redirect(w, r, "/auth/discord", http.StatusFound)
		return
	}
	repo := internal.NewRepository(g.IndexService.DB)
	categoryPositions := make(map[string]internal.DiscordChannel)
	guildId := r.PathValue("guildId")
	guild, err := g.IndexService.DiscordSession.State.Guild(guildId)
	if err != nil {
		http.Error(w, "Not get guild id", http.StatusInternalServerError)
		return
	}
	permissionCode, err := repo.GetPermissionCode(r.Context(), guildId, "")
	if err != nil {
		http.Error(w, "権限コードの取得に失敗しました", http.StatusInternalServerError)
		return
	}
	permissionIDs, err := repo.GetPermissionIDs(r.Context(), guildId, "")
	if err != nil {
		http.Error(w, "権限読み込みに失敗しました", http.StatusInternalServerError)
		return
	}
	discordGuildMember, err := g.IndexService.DiscordSession.GuildMember(guildId, discordLoginUser.User.ID)
	if err != nil {
		http.Error(w, "Not get discord member", http.StatusInternalServerError)
		return
	}
	// 権限のチェック
	if (permissionCode & discordGuildMember.Permissions) == 0 {
		http.Error(w, "権限がありません", http.StatusForbidden)
		return
	}
	for _, permissionId := range permissionIDs {
		if permissionId.TargetType == "user" && permissionId.TargetID == discordLoginUser.User.ID {
			http.Error(w, "権限がありません", http.StatusForbidden)
			return
		}
		if permissionId.TargetType == "role" && discordGuildMember.Roles != nil {
			for _, role := range discordGuildMember.Roles {
				if permissionId.TargetID == role {
					http.Error(w, "権限がありません", http.StatusForbidden)
					return
				}
			}
		}
	}
	// カテゴリーのチャンネルを取得
	//[categoryID]map[channelPosition]channelName
	channelsInCategory := make(map[string][]internal.DiscordChannelSelect)
	for _, channel := range guild.Channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}
		categoryPositions[channel.ID] = internal.DiscordChannel{
			ID:       channel.ID,
			Name:     channel.Name,
			Position: channel.Position,
		}
	}
	// カテゴリーなしのチャンネルを追加
	//channelsInCategory[""] = make([]internal.DiscordChannelSelect, len(guild.Channels)-1, len(guild.Channels))
	for _, channel := range guild.Channels {
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
	}
	lineBot, err := repo.GetLineBot(r.Context(), guildId)
	if err != nil && err.Error() == "sql: no rows in result set" {
		err = repo.InsertLineBot(r.Context(), &internal.LineBot{
			GuildID:          guildId,
			DefaultChannelID: guild.SystemChannelID,
			DebugMode:        false,
		})
		if err != nil {
			http.Error(w, "line_bot:"+err.Error(), http.StatusInternalServerError)
			return
		}
		err = repo.InsertLineBotIv(r.Context(), &internal.LineBotIv{
			GuildID: guildId,
		})
		if err != nil {
			http.Error(w, "line_bot_iv:"+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	htmlSelectChannels := ``
	for _, channels := range channelsInCategory {
		for _, channelSelect := range channels {
			if channelSelect.ID == "" {
				continue
			}
			if lineBot.DefaultChannelID == channelSelect.ID {
				htmlSelectChannels += fmt.Sprintf(`<option value="%s" selected>%s</option>`, channelSelect.ID, channelSelect.Name)
				continue
			}
			htmlSelectChannels += fmt.Sprintf(`<option value="%s">%s</option>`, channelSelect.ID, channelSelect.Name)
		}
	}
	data := struct {
		GuildID  string
		Channels template.HTML
	}{
		GuildID:  guildId,
		Channels: template.HTML(htmlSelectChannels),
	}
	t := template.Must(template.New("linetoken.html").ParseFiles("web/templates/views/guilds/linetoken.html"))
	err = t.ExecuteTemplate(w, "linetoken.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
