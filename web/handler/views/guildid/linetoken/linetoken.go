package guildid

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/handler/views/guildid/linetoken/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
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
	var userPermissionCode int64
	userPermissionCode = 0
	repo := internal.NewRepository(g.IndexService.DB)
	categoryPositions := make(map[string]internal.DiscordChannel)
	guildId := r.PathValue("guildId")
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// ログインユーザーの取得
	discordLoginUser, err := getoauth.GetDiscordOAuth(
		ctx,
		g.IndexService.CookieStore,
		r,
		config.SessionSecret(),
	)
	if err != nil {
		http.Redirect(w, r, "/auth/discord", http.StatusFound)
		return
	}
	guild, err := g.IndexService.DiscordSession.State.Guild(guildId)
	if err != nil {
		http.Error(w, "Not get guild id", http.StatusInternalServerError)
		return
	}
	permissionCode, err := repo.GetPermissionCode(ctx, guildId, "line_bot")
	if err != nil {
		http.Error(w, "権限コードの取得に失敗しました", http.StatusInternalServerError)
		return
	}
	permissionIDs, err := repo.GetPermissionIDs(ctx, guildId, "line_bot")
	if err != nil {
		http.Error(w, "権限読み込みに失敗しました", http.StatusInternalServerError)
		return
	}
	discordGuildMember, err := g.IndexService.DiscordSession.GuildMember(guildId, discordLoginUser.User.ID)
	if err != nil {
		http.Error(w, "Not get discord member", http.StatusInternalServerError)
		return
	}
	guildRoles, err := g.IndexService.DiscordSession.GuildRoles(guildId)
	if err != nil {
		http.Error(w, "Not get guild roles", http.StatusInternalServerError)
		return
	}

	for _, role := range discordGuildMember.Roles {
		for _, guildRole := range guildRoles {
			if role == guildRole.ID {
				userPermissionCode |= guildRole.Permissions
			}
		}
	}
	// メンバーの権限を取得
	// discordgoの場合guildMemberから正しく権限を取得できないため、UserChannelPermissionsを使用
	memberPermission, err := g.IndexService.DiscordSession.UserChannelPermissions(discordLoginUser.User.ID, guild.Channels[0].ID)
	if err != nil {
		http.Error(w, "Not get member permission", http.StatusInternalServerError)
		return
	}
	// 権限のチェック
	if (permissionCode & (memberPermission | userPermissionCode)) == 0 {
		permissionFlag := false
		for _, permissionId := range permissionIDs {
			if permissionId.TargetType == "user" && permissionId.TargetID == discordLoginUser.User.ID {
				permissionFlag = true
				break
			}
			if permissionId.TargetType == "role" && discordGuildMember.Roles != nil {
				for _, role := range discordGuildMember.Roles {
					if permissionId.TargetID == role {
						permissionFlag = true
						break
					}
				}
			}
		}
		if !permissionFlag {
			http.Error(w, "権限がありません", http.StatusForbidden)
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
	lineBot, err := repo.GetLineBot(ctx, guildId)
	if err != nil && err.Error() == "sql: no rows in result set" {
		err = repo.InsertLineBot(ctx, &internal.LineBot{
			GuildID:          guildId,
			DefaultChannelID: guild.SystemChannelID,
			DebugMode:        false,
		})
		if err != nil {
			http.Error(w, "line_bot:"+err.Error(), http.StatusInternalServerError)
			return
		}
		err = repo.InsertLineBotIv(ctx, &internal.LineBotIv{
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
	categoryOptions := make([]string, len(categoryIDTmps)+1)
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
				categoryOptions[categoryIndex] += fmt.Sprintf(`<option value="%s" selected>%s</option>`, channelSelect.ID, channelSelect.Name)
				continue
			}
			categoryOptions[categoryIndex] += fmt.Sprintf(`<option value="%s">%s</option>`, channelSelect.ID, channelSelect.Name)
		}
	}
	for _, categoryOption := range categoryOptions {
		htmlSelectChannels += categoryOption
	}
	data := struct {
		GuildID  string
		Channels template.HTML
	}{
		GuildID:  guildId,
		Channels: template.HTML(htmlSelectChannels),
	}
	t := template.Must(template.New("linetoken.html").ParseFiles("web/templates/views/guildid/linetoken.html"))
	err = t.ExecuteTemplate(w, "linetoken.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
