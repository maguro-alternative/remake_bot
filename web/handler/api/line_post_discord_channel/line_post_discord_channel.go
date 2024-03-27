package linepostdiscordchannel

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/line_post_discord_channel/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

//go:generate go run github.com/matryer/moq -out mock_test.go . Repository
type Repository interface {
	UpdateLinePostDiscordChannel(ctx context.Context, linePostDiscordChannel repository.LinePostDiscordChannelAllColumns) error
	InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordMessageTypes []repository.LineNgDiscordMessageType) error
	DeleteNotInsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordMessageTypes []repository.LineNgDiscordMessageType) error
	InsertLineNgDiscordIDs(ctx context.Context, lineNgDiscordIDs []repository.LineNgDiscordIDAllCoulmns) error
	DeleteNotInsertLineNgDiscordIDs(ctx context.Context, lineNgDiscordIDs []repository.LineNgDiscordIDAllCoulmns) error
}

//go:generate go run github.com/matryer/moq -out permission_mock_test.go . OAuthPermission
type OAuthPermission interface {
	CheckDiscordPermission(ctx context.Context, guild *discordgo.Guild, permissionType string) (statusCode int, discordPermissionData *model.DiscordPermissionData, err error)
}

type LinePostDiscordChannelHandler struct {
	IndexService          *service.IndexService
	Repo                  Repository
	DiscordPermissiondata *model.DiscordPermissionData
}

//go:generate go run github.com/matryer/moq -out discordsession_mock_test.go . Session
type Session interface {
	ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelFileSendWithMessage(channelID string, content string, name string, r io.Reader, options ...discordgo.RequestOption) (*discordgo.Message, error)
	Guild(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error)
	GuildChannels(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Channel, err error)
	GuildMember(guildID string, userID string, options ...discordgo.RequestOption) (st *discordgo.Member, err error)
	GuildMembers(guildID string, after string, limit int, options ...discordgo.RequestOption) (st []*discordgo.Member, err error)
	GuildRoles(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Role, err error)
	UserChannelPermissions(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error)
	UserGuilds(limit int, beforeID string, afterID string, options ...discordgo.RequestOption) (st []*discordgo.UserGuild, err error)
}

var (
	_ Session    = (*discordgo.Session)(nil)
	_ Session    = (service.Session)(nil)
	_ Repository = (*repository.Repository)(nil)
)

func NewLinePostDiscordChannelHandler(
	indexService *service.IndexService,
	repo Repository,
	DiscordPermissionData *model.DiscordPermissionData,
) *LinePostDiscordChannelHandler {
	return &LinePostDiscordChannelHandler{
		IndexService:          indexService,
		Repo:                  repo,
		DiscordPermissiondata: DiscordPermissionData,
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
	var repo Repository
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

	repo = h.Repo
	lineChannels, lineNgTypes, lineNgIDs := lineChannelJsonRead(lineChannelJson)

	for _, lineChannel := range lineChannels {
		linePostDiscordChannel := repository.NewLinePostDiscordChannel(
			lineChannel.ChannelID,
			lineChannel.GuildID,
			lineChannel.Ng,
			lineChannel.BotMessage,
		)
		if err := repo.UpdateLinePostDiscordChannel(ctx, *linePostDiscordChannel); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.ErrorContext(ctx, "line_post_discord_channel更新に失敗しました。 "+err.Error())
			return
		}
	}

	if err := repo.InsertLineNgDiscordMessageTypes(ctx, lineNgTypes); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_ng_discord_message_type更新に失敗しました。 "+err.Error())
		return
	}

	if err := repo.DeleteNotInsertLineNgDiscordMessageTypes(ctx, lineNgTypes); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_ng_discord_message_type更新に失敗しました。 "+err.Error())
		return
	}

	if err := repo.InsertLineNgDiscordIDs(ctx, lineNgIDs); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "line_ng_discord_id更新に失敗しました。 "+err.Error())
		return
	}

	if err := repo.DeleteNotInsertLineNgDiscordIDs(ctx, lineNgIDs); err != nil {
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
