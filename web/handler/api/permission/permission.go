package permission

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/handler/api/permission/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/getoauth"
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
	Repo         Repository
	oauthStore   OAuthStore
}

func NewPermissionHandler(indexService *service.IndexService, repo service.Repository,) *PermissionHandler {
	return &PermissionHandler{
		IndexService: indexService,
		Repo:         repo,
	}
}

func (h *PermissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var userPermissionCode int64
	var repo Repository
	var oauthStore OAuthStore
	var client http.Client
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

	// パーミッションの更新
	repo = h.Repo

	oauthStore = getoauth.NewOAuthStore(h.IndexService.CookieStore, config.SessionSecret())
	// mockの場合はmockを使用
	if h.oauthStore != nil {
		oauthStore = h.oauthStore
	}
	discordSession, err := oauthStore.GetDiscordOAuth(ctx, r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Discordの認証情報の取得に失敗しました。", "エラー:", err.Error())
		return
	}

	// ギルド情報を取得
	guild, err := h.IndexService.DiscordSession.Guild(guildId, discordgo.WithClient(&client))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "ギルド情報の取得に失敗しました。", "エラー:", err.Error())
		return
	}

	// チャンネル一覧を取得
	channels, err := h.IndexService.DiscordSession.GuildChannels(guild.ID, discordgo.WithClient(&client))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "チャンネルの取得に失敗しました。", "エラー:", err.Error())
		return
	}

	discordGuildMember, err := h.IndexService.DiscordSession.GuildMember(guild.ID, discordSession.User.ID, discordgo.WithClient(&client))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "メンバーの取得に失敗しました。", "エラー:", err.Error())
		return
	}
	guildRoles, err := h.IndexService.DiscordSession.GuildRoles(guild.ID, discordgo.WithClient(&client))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "ロールの取得に失敗しました。", "エラー:", err.Error())
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
	memberPermission, err := h.IndexService.DiscordSession.UserChannelPermissions(discordSession.User.ID, channels[0].ID, discordgo.WithClient(&client))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "メンバーの権限の取得に失敗しました。", "エラー:", err.Error())
		return
	}

	if ((memberPermission | userPermissionCode) & 8) != 8 {
		fmt.Println((memberPermission | userPermissionCode) & 8)
		http.Error(w, "Not permission", http.StatusForbidden)
		slog.WarnContext(ctx, "権限のないアクセスがありました。")
		return
	}

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

	if err := repo.UpdatePermissionCodes(ctx, permissionCodes); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの更新に失敗しました。", "エラー:", err.Error())
		return
	}

	if err := repo.DeletePermissionIDs(ctx, guildId); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの削除に失敗しました。", "エラー:", err.Error())
		return
	}

	if err := repo.InsertPermissionIDs(ctx, permissionIDs); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "パーミッションの追加に失敗しました。", "エラー:", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
