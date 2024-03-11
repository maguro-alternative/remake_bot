package linepostdiscordchannel

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/web/handler/views/guildid/line_post_discord_channel/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/permission"
)

type LinePostDiscordChannelViewHandler struct {
	IndexService *service.IndexService
}

func NewLinePostDiscordChannelViewHandler(indexService *service.IndexService) *LinePostDiscordChannelViewHandler {
	return &LinePostDiscordChannelViewHandler{
		IndexService: indexService,
	}
}

func (g *LinePostDiscordChannelViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	messageTypes := []string{
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
	categoryPositions := make(map[string]internal.DiscordChannel)
	var categoryIDTmps []string
	guildId := r.PathValue("guildId")
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	guild, err := g.IndexService.DiscordSession.State.Guild(guildId)
	if err != nil {
		http.Error(w, "Not get guild id", http.StatusInternalServerError)
		return
	}
	statusCode, _, err := permission.CheckDiscordPermission(ctx, w, r, g.IndexService, guild, "line_bot")
	if err != nil {
		if statusCode == http.StatusFound {
			http.Redirect(w, r, "/auth/discord", http.StatusFound)
			return
		}
		http.Error(w, "Not permission", statusCode)
		return
	}
	//[categoryID]map[channelPosition]channelName
	channelsInCategory := make(map[string][]internal.DiscordChannelSet)
	repo := internal.NewRepository(g.IndexService.DB)
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
		if len(channelsInCategory[categoryPosition.ID]) == 0 {
			channelsInCategory[categoryPosition.ID] = make([]internal.DiscordChannelSet, len(guild.Channels)-2, len(guild.Channels))
		}
		discordChannel, err := repo.GetLineChannel(ctx, channel.ID)
		if err != nil && err.Error() != "sql: no rows in result set" {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else if err != nil {
			if err := repo.InsertLineChannel(ctx, channel.ID, guildId); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			discordChannel = internal.LineChannel{
				ChannelID:  channel.ID,
				GuildID:    guildId,
				Ng:         false,
				BotMessage: false,
			}
		}
		ngTypes, err := repo.GetLineNgType(ctx, channel.ID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		ngDiscordIDs, err := repo.GetLineNgDiscordID(ctx, channel.ID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		channelsInCategory[categoryPosition.ID][channel.Position] = internal.DiscordChannelSet{
			ID:         channel.ID,
			Name:       fmt.Sprintf("%s %s", typeIcon, channel.Name),
			Ng:         discordChannel.Ng,
			BotMessage: discordChannel.BotMessage,
			NgTypes:    ngTypes,
		}
		for _, ngDiscordID := range ngDiscordIDs {
			if ngDiscordID.IDType == "user" {
				channelsInCategory[categoryPosition.ID][channel.Position].NgUsers = append(channelsInCategory[categoryPosition.ID][channel.Position].NgUsers, ngDiscordID.ID)
				continue
			}
			channelsInCategory[categoryPosition.ID][channel.Position].NgRoles = append(channelsInCategory[categoryPosition.ID][channel.Position].NgRoles, ngDiscordID.ID)
		}
	}

	htmlForm := ``
	categoryComponents := make([]string, len(categoryIDTmps)+1)
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
		categoryChannelName := categoryPositions[categoryID].Name
		if categoryID == "" {
			categoryChannelName = "カテゴリーなし"
		}
		categoryComponents[categoryIndex] = fmt.Sprintf(`
		<details>
            <summary>%s</summary>
		`, categoryChannelName)
		for _, channel := range channels {
			if channel.ID == "" {
				continue
			}
			messageNgFlag, botNgFlag := "", ""
			selectMemberForm, selectRoleForm, selectMessageTypeForm := createSelectForm(guild, channel, messageTypes)
			if channel.Ng {
				messageNgFlag = "checked"
			}
			if channel.BotMessage {
				botNgFlag = "checked"
			}
			categoryComponents[categoryIndex] += `
			<details>
                <summary>` + channel.Name + `</summary>
				<label for="ng_` + channel.ID + `">LINEへ送信しない</label>
				<input type="checkbox" id="ng_` + channel.ID + `" name="ng_` + channel.ID + `" ` + messageNgFlag + ` />
				<br/>
				<label for="bot_message_` + channel.ID + `">Botのメッセージを送信しない</label>
				<input type="checkbox" id="bot_message_` + channel.ID + `" name="bot_message_` + channel.ID + `"` + botNgFlag + ` />
				<br/>
				<label for="ng_types_` + channel.ID + `[]">NGタイプ</label>
				<select id="ng_types_` + channel.ID + `[]" name="ng_types_` + channel.ID + `[]" multiple>
					` + selectMessageTypeForm + `
				</select>
				<br/>
				<label for="ng_users_` + channel.ID + `[]">NGユーザー</label>
				<select id="ng_users_` + channel.ID + `[]" name="ng_users_` + channel.ID + `[]" multiple>
					` + selectMemberForm + `
				</select>
				<br/>
				<label for="ng_roles_` + channel.ID + `[]">NGロール</label>
				<select id="ng_roles_` + channel.ID + `[]" name="ng_roles_` + channel.ID + `[]" multiple>
					` + selectRoleForm + `
				</select>
				<br/>
			</details>
			`
		}
		categoryComponents[categoryIndex] += `
		</details>`
	}
	for _, categoryComponent := range categoryComponents {
		htmlForm += categoryComponent
	}

	tmpl := template.Must(template.ParseFiles("web/templates/layout.html","web/templates/views/guildid/line_post_discord_channel.html"))
	if err := tmpl.Execute(w, struct {
		Title     string
		GuildName string
		GuildID   string
		HTMLForm  template.HTML
	}{
		Title:     "DiscordからLINEへの送信設定",
		GuildName: guild.Name,
		GuildID:   guildId,
		HTMLForm:  template.HTML(htmlForm),
	}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func createSelectForm(guild *discordgo.Guild, channel internal.DiscordChannelSet, messageTypes []string) (member string, role string, messageType string) {
	selectMemberForm := ""
	for _, member := range guild.Members {
		selectedFlag := false
		for _, ngUserID := range channel.NgUsers {
			if ngUserID == member.User.ID {
				selectedFlag = true
				break
			}
		}
		if selectedFlag {
			selectMemberForm += fmt.Sprintf(`<option value="%s" selected>%s</option>`, member.User.ID, member.User.Username)
			continue
		}
		selectMemberForm += fmt.Sprintf(`<option value="%s">%s</option>`, member.User.ID, member.User.Username)
	}
	selectRoleForm := ""
	for _, role := range guild.Roles {
		selectedFlag := false
		for _, ngRoleID := range channel.NgRoles {
			if ngRoleID == role.ID {
				selectedFlag = true
				break
			}
		}
		if selectedFlag {
			selectRoleForm += fmt.Sprintf(`<option value="%s" selected>%s</option>`, role.ID, role.Name)
			continue
		}
		selectRoleForm += fmt.Sprintf(`<option value="%s">%s</option>`, role.ID, role.Name)
	}
	selectMessageTypeForm := ""
	for i, messageType := range messageTypes {
		selectedFlag := false
		for _, ngType := range channel.NgTypes {
			if ngType == i {
				selectedFlag = true
				break
			}
		}
		if selectedFlag {
			selectMessageTypeForm += fmt.Sprintf(`<option value=%d selected>%s</option>`, i, messageType)
			continue
		}
		selectMessageTypeForm += fmt.Sprintf(`<option value=%d>%s</option>`, i, messageType)
	}
	return selectMemberForm, selectRoleForm, selectMessageTypeForm
}
