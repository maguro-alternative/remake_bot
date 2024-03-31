package permission

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/permission/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

//go:generate go run github.com/matryer/moq -out mock_test.go . Repository
type Repository interface {
	UpdatePermissionCodes(ctx context.Context, permissionsCode []repository.PermissionCode) error
	DeletePermissionIDs(ctx context.Context, guildId string) error
	InsertPermissionIDs(ctx context.Context, permissionsID []repository.PermissionIDAllColumns) error
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

//go:generate go run github.com/matryer/moq -out oauth_mock_test.go . OAuthStore
type OAuthStore interface {
	GetDiscordOAuth(ctx context.Context, r *http.Request) (*model.DiscordOAuthSession, error)
	GetLineOAuth(r *http.Request) (*model.LineOAuthSession, error)
}

type PermissionHandler struct {
	IndexService *service.IndexService
	Repo         repository.RepositoryFunc
}

func NewPermissionHandler(
	indexService *service.IndexService,
	repo service.Repository,
) *PermissionHandler {
	return &PermissionHandler{
		IndexService: indexService,
		Repo:         repo,
	}
}

func (h *PermissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	guildId := r.PathValue("guildId")
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		slog.ErrorContext(ctx, "/api/permission Method Not Allowed")
		return
	}
	var permissionJson internal.PermissionJson
	var permissionCodes []repository.PermissionCode
	var permissionIDs []repository.PermissionIDAllColumns

	if err := json.NewDecoder(r.Body).Decode(&permissionJson); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "jsonの読み取りに失敗しました:", "エラー:", err.Error())
		return
	}
	if err := permissionJson.Validate(); err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		slog.ErrorContext(ctx, "jsonのバリデーションに失敗しました:", "エラー:", err.Error())
		return
	}

	for _, permissionCode := range permissionJson.PermissionCodes {
		permissionCodes = append(permissionCodes, repository.PermissionCode{
			GuildID: guildId,
			Type:    permissionCode.Type,
			Code:    permissionCode.Code,
		})
	}

	for _, permissionID := range permissionJson.PermissionIDs {
		permissionIDs = append(permissionIDs, repository.PermissionIDAllColumns{
			GuildID:    guildId,
			Type:       permissionID.Type,
			TargetType: permissionID.TargetType,
			TargetID:   permissionID.TargetID,
			Permission: permissionID.Permission,
		})
	}

	if err := h.Repo.UpdatePermissionCodes(ctx, permissionCodes); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの更新に失敗しました。", "エラー:", err.Error())
		return
	}

	if err := h.Repo.DeletePermissionIDs(ctx, guildId); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの削除に失敗しました。", "エラー:", err.Error())
		return
	}

	if err := h.Repo.InsertPermissionIDs(ctx, permissionIDs); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの追加に失敗しました。", "エラー:", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
