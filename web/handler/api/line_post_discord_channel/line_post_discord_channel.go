package linepostdiscordchannel

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/web/handler/api/line_post_discord_channel/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/permission"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

//go:generate go run github.com/matryer/moq -out mock_test.go . Repository
type Repository interface {
	UpdateLinePostDiscordChannel(ctx context.Context, linePostDiscordChannel internal.LinePostDiscordChannel) error
	InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordMessageTypes []internal.LineNgDiscordMessageType) error
	DeleteNotInsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordMessageTypes []internal.LineNgDiscordMessageType) error
	InsertLineNgDiscordIDs(ctx context.Context, lineNgDiscordIDs []internal.LineNgID) error
	DeleteNotInsertLineNgDiscordIDs(ctx context.Context, lineNgDiscordIDs []internal.LineNgID) error
}

//go:generate go run github.com/matryer/moq -out permission_mock_test.go . OAuthPermission
type OAuthPermission interface {
	CheckDiscordPermission(ctx context.Context, guild *discordgo.Guild, permissionType string) (statusCode int, discordPermissionData *model.DiscordPermissionData, err error)
}

type LinePostDiscordChannelHandler struct {
	IndexService    *service.IndexService
	repo            Repository
	oauthPermission OAuthPermission
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
	_ Session = (*discordgo.Session)(nil)
	_ Session = (service.Session)(nil)
)

func NewLinePostDiscordChannelHandler(indexService *service.IndexService) *LinePostDiscordChannelHandler {
	return &LinePostDiscordChannelHandler{
		IndexService: indexService,
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
	var oauthPermission OAuthPermission
	var client http.Client
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
	guild, err := h.IndexService.DiscordSession.Guild(lineChannelJson.GuildID, discordgo.WithClient(&client))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Guild情報取得に失敗しました。 "+err.Error())
		return
	}

	oauthPermission = permission.NewPermissionHandler(r, &client, h.IndexService)
	// テスト用
	if h.oauthPermission != nil {
		oauthPermission = h.oauthPermission
	}
	statusCode, discordPermissionData, err := oauthPermission.CheckDiscordPermission(ctx, guild, "line_post_discord_channel")
	if err != nil {
		if statusCode == http.StatusFound {
			slog.InfoContext(ctx, "Redirect to /login/discord")
			http.Redirect(w, r, "/login/discord", http.StatusFound)
			return
		}
		if discordPermissionData.Permission == "" {
			http.Error(w, "Not permission", statusCode)
			slog.WarnContext(ctx, "権限のないアクセスがありました。 "+err.Error())
			return
		}
	}

	repo = internal.NewRepository(h.IndexService.DB)
	// mockの場合はmockを使用
	if h.repo != nil {
		repo = h.repo
	}
	lineChannels, lineNgTypes, lineNgIDs := lineChannelJsonRead(lineChannelJson)

	for _, lineChannel := range lineChannels {
		if err := repo.UpdateLinePostDiscordChannel(ctx, lineChannel); err != nil {
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

func lineChannelJsonRead(lineChannelJson internal.LinePostDiscordChannelJson) (channels []internal.LinePostDiscordChannel, ngTypes []internal.LineNgDiscordMessageType, ngIDs []internal.LineNgID) {
	var lineChannels []internal.LinePostDiscordChannel
	var lineNgTypes []internal.LineNgDiscordMessageType
	var lineNgIDs []internal.LineNgID
	for _, lineChannel := range lineChannelJson.Channels {
		lineChannels = append(lineChannels, internal.LinePostDiscordChannel{
			ChannelID:  lineChannel.ChannelID,
			GuildID:    lineChannelJson.GuildID,
			Ng:         lineChannel.Ng,
			BotMessage: lineChannel.BotMessage,
		})
		if len(lineChannel.NgTypes) > 0 {
			for _, ngType := range lineChannel.NgTypes {
				lineNgTypes = append(lineNgTypes, internal.LineNgDiscordMessageType{
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
