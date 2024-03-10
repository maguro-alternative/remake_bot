package linechannel

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/handler/api/line_post_discord_channel/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/permission"
)

type LineChannelHandler struct {
	IndexService *service.IndexService
}

func NewLineChannelHandler(indexService *service.IndexService) *LineChannelHandler {
	return &LineChannelHandler{
		IndexService: indexService,
	}
}

func (h *LineChannelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	var lineChannelJson internal.LineChannelJson
	if err := json.NewDecoder(r.Body).Decode(&lineChannelJson); err != nil {
		http.Error(w, "Json読み取りに失敗しました。", http.StatusBadRequest)
		return
	}

	if err := lineChannelJson.Validate(); err != nil {
		http.Error(w, "バリデーションチェックに失敗しました。", http.StatusBadRequest)
		return
	}

	lineChannelJson.GuildID = r.PathValue("guildId")
	guild, err := h.IndexService.DiscordSession.State.Guild(lineChannelJson.GuildID)
	if err != nil {
		http.Error(w, "Not get guild id", http.StatusInternalServerError)
		return
	}
	statusCode, _, err := permission.CheckDiscordPermission(ctx, w, r, h.IndexService, guild, "line_bot")
	if err != nil {
		if statusCode == http.StatusFound {
			http.Redirect(w, r, "/auth/discord", http.StatusFound)
			return
		}
		http.Error(w, "Not permission", statusCode)
		return
	}

	repo := internal.NewRepository(h.IndexService.DB)
	lineChannels, lineNgTypes, lineNgIDs := lineChannelJsonRead(lineChannelJson)

	for _, lineChannel := range lineChannels {
		if err := repo.UpdateLinePostDiscordChannel(ctx, lineChannel); err != nil {
			http.Error(w, "line_post_discord_channel更新に失敗しました。", http.StatusInternalServerError)
			return
		}
	}

	if err := repo.InsertLineNgDiscordMessageTypes(ctx, lineNgTypes); err != nil {
		http.Error(w, "line_ng_discord_message_type更新に失敗しました。", http.StatusInternalServerError)
		return
	}

	if err := repo.DeleteNotInsertLineNgDiscordMessageTypes(ctx, lineNgTypes); err != nil {
		http.Error(w, "line_ng_discord_message_type更新に失敗しました。", http.StatusInternalServerError)
		return
	}

	if err := repo.InsertLineNgDiscordIDs(ctx, lineNgIDs); err != nil {
		http.Error(w, "line_ng_discord_id更新(挿入)に失敗しました。", http.StatusInternalServerError)
		return
	}

	if err := repo.DeleteNotInsertLineNgDiscordIDs(ctx, lineNgIDs); err != nil {
		http.Error(w, "line_ng_discord_id更新(削除)に失敗しました。"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}

func lineChannelJsonRead(lineChannelJson internal.LineChannelJson) (channels []internal.LineChannel, ngTypes []internal.LineNgType, ngIDs []internal.LineNgID) {
	var lineChannels []internal.LineChannel
	var lineNgTypes []internal.LineNgType
	var lineNgIDs []internal.LineNgID
	for _, lineChannel := range lineChannelJson.Channels {
		lineChannels = append(lineChannels, internal.LineChannel{
			ChannelID:  lineChannel.ChannelID,
			GuildID:    lineChannelJson.GuildID,
			Ng:         lineChannel.Ng,
			BotMessage: lineChannel.BotMessage,
		})
		if len(lineChannel.NgTypes) > 0 {
			for _, ngType := range lineChannel.NgTypes {
				lineNgTypes = append(lineNgTypes, internal.LineNgType{
					ChannelID: lineChannel.ChannelID,
					GuildID:   lineChannelJson.GuildID,
					Type:      ngType,
				})
			}
		}
		if len(lineChannel.NgUsers) > 0 {
			for _, ngUser := range lineChannel.NgUsers {
				lineNgIDs = append(lineNgIDs, internal.LineNgID{
					ChannelID: lineChannel.ChannelID,
					GuildID:   lineChannelJson.GuildID,
					ID:        ngUser,
					IDType:    "user",
				})
			}
		}
		if len(lineChannel.NgRoles) > 0 {
			for _, ngRole := range lineChannel.NgRoles {
				lineNgIDs = append(lineNgIDs, internal.LineNgID{
					ChannelID: lineChannel.ChannelID,
					GuildID:   lineChannelJson.GuildID,
					ID:        ngRole,
					IDType:    "role",
				})
			}
		}
	}
	return lineChannels, lineNgTypes, lineNgIDs
}
