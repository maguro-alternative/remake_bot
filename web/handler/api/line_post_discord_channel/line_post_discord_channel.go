package linepostdiscordchannel

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/line_post_discord_channel/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
)

type LinePostDiscordChannelHandler struct {
	IndexService          *service.IndexService
	Repo                  repository.RepositoryFunc
}

func NewLinePostDiscordChannelHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
) *LinePostDiscordChannelHandler {
	return &LinePostDiscordChannelHandler{
		IndexService:          indexService,
		Repo:                  repo,
	}
}

func (h *LinePostDiscordChannelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	// Post以外のリクエストは受け付けない
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		slog.ErrorContext(ctx, "Method Not Allowed")
		return
	}
	var lineChannelJson internal.LinePostDiscordChannelJson
	if err := json.NewDecoder(r.Body).Decode(&lineChannelJson); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Json読み取りに失敗しました。 "+err.Error())
		return
	}

	if err := lineChannelJson.Validate(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Jsonバリデーションに失敗しました。 "+err.Error())
		return
	}

	lineChannelJson.GuildID = r.PathValue("guildId")

	lineChannels, lineNgTypes, lineNgIDs := lineChannelJsonRead(lineChannelJson)

	for _, lineChannel := range lineChannels {
		linePostDiscordChannel := repository.NewLinePostDiscordChannel(
			lineChannel.ChannelID,
			lineChannel.GuildID,
			lineChannel.Ng,
			lineChannel.BotMessage,
		)
		if err := h.Repo.UpdateLinePostDiscordChannel(ctx, *linePostDiscordChannel); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "line_post_discord_channel更新に失敗しました。 "+err.Error())
			return
		}
	}

	if err := h.Repo.InsertLineNgDiscordMessageTypes(ctx, lineNgTypes); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_ng_discord_message_type更新に失敗しました。 "+err.Error())
		return
	}

	if err := h.Repo.DeleteNotInsertLineNgDiscordMessageTypes(ctx, lineNgTypes); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_ng_discord_message_type更新に失敗しました。 "+err.Error())
		return
	}

	if err := h.Repo.InsertLineNgDiscordIDs(ctx, lineNgIDs); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_ng_discord_id更新に失敗しました。 "+err.Error())
		return
	}

	if err := h.Repo.DeleteNotInsertLineNgDiscordIDs(ctx, lineNgIDs); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_ng_discord_id更新に失敗しました。 "+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}

func lineChannelJsonRead(lineChannelJson internal.LinePostDiscordChannelJson) (
	channels []repository.LinePostDiscordChannelAllColumns,
	ngTypes []repository.LineNgDiscordMessageType,
	ngIDs []repository.LineNgDiscordIDAllCoulmns,
) {
	var lineChannels []repository.LinePostDiscordChannelAllColumns
	var lineNgTypes []repository.LineNgDiscordMessageType
	var lineNgIDs []repository.LineNgDiscordIDAllCoulmns
	for _, lineChannel := range lineChannelJson.Channels {
		channel := repository.NewLinePostDiscordChannel(
			lineChannel.ChannelID,
			lineChannelJson.GuildID,
			lineChannel.Ng,
			lineChannel.BotMessage,
		)
		lineChannels = append(lineChannels, *channel)
		if len(lineChannel.NgTypes) > 0 {
			for _, ngType := range lineChannel.NgTypes {
				messageType := repository.NewLineNgDiscordMessageType(
					lineChannel.ChannelID,
					lineChannelJson.GuildID,
					ngType,
				)
				lineNgTypes = append(lineNgTypes, *messageType)
			}
		}
		if len(lineChannel.NgUsers) > 0 {
			for _, ngUser := range lineChannel.NgUsers {
				user := repository.NewLineNgDiscordID(
					lineChannel.ChannelID,
					lineChannelJson.GuildID,
					ngUser,
					"user",
				)
				lineNgIDs = append(lineNgIDs, *user)
			}
		}
		if len(lineChannel.NgRoles) > 0 {
			for _, ngRole := range lineChannel.NgRoles {
				role := repository.NewLineNgDiscordID(
					lineChannel.ChannelID,
					lineChannelJson.GuildID,
					ngRole,
					"role",
				)
				lineNgIDs = append(lineNgIDs, *role)
			}
		}
	}
	return lineChannels, lineNgTypes, lineNgIDs
}
