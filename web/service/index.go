package service

import (
	"context"
	"io"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/sessions"
)

// A TODOService implements CRUD of TODO entities.
type IndexService struct {
	DB              db.Driver
	CookieStore     *sessions.CookieStore
	DiscordSession  Session
	DiscordBotState *discordgo.State
}

// Repository is an interface for repository.
type Repository interface {
	InsertLineBotIv(ctx context.Context, guildId string) error
	GetAllColumnsLineBotIv(ctx context.Context, guildID string) (repository.LineBotIv, error)
	GetLineBotIvNotClient(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error)
	UpdateLineBotIv(ctx context.Context, lineBotIv *repository.LineBotIv) error
	InsertLineBot(ctx context.Context, lineBot *repository.LineBot) error
	GetAllColumnsLineBots(ctx context.Context) ([]*repository.LineBot, error)
	GetAllColumnsLineBot(ctx context.Context, guildId string) (repository.LineBot, error)
	GetLineBotDefaultChannelID(ctx context.Context, guildID string) (repository.LineBotDefaultChannelID, error)
	GetLineBotNotClient(ctx context.Context, guildID string) (repository.LineBotNotClient, error)
	UpdateLineBot(ctx context.Context, lineBot *repository.LineBot) error
	GetLineNgDiscordID(ctx context.Context, channelID string) ([]repository.LineNgDiscordID, error)
	InsertLineNgDiscordIDs(ctx context.Context, lineNgDiscordIDs []repository.LineNgDiscordIDAllCoulmns) error
	DeleteNotInsertLineNgDiscordIDs(ctx context.Context, lineNgDiscordIDs []repository.LineNgDiscordIDAllCoulmns) error
	InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []repository.LineNgDiscordMessageType) error
	DeleteNotInsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []repository.LineNgDiscordMessageType) error
	GetLineNgDiscordMessageType(ctx context.Context, channelID string) ([]int, error)
	GetLinePostDiscordChannel(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error)
	UpdateLinePostDiscordChannel(ctx context.Context, lineChannel repository.LinePostDiscordChannelAllColumns) error
	InsertLinePostDiscordChannel(ctx context.Context, channelID string, guildID string) error
	GetPermissionCode(ctx context.Context, guildID, permissionType string) (int64, error)
	GetPermissionCodes(ctx context.Context, guildID string) ([]repository.PermissionCode, error)
	UpdatePermissionCodes(ctx context.Context, permissionsCode []repository.PermissionCode) error
	InsertPermissionIDs(ctx context.Context, permissionsID []repository.PermissionIDAllColumns) error
	GetGuildPermissionIDsAllColumns(ctx context.Context, guildID string) ([]repository.PermissionIDAllColumns, error)
	GetPermissionIDs(ctx context.Context, guildID, permissionType string) ([]repository.PermissionID, error)
	DeletePermissionIDs(ctx context.Context, guildId string) error
}

// Session is an interface for discordgo.Session.
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
	_ Repository = (*repository.Repository)(nil)
)

// NewTODOService returns new TODOService.
func NewIndexService(
	db db.Driver,
	cookieStore *sessions.CookieStore,
	discordSession Session,
	discordBotState *discordgo.State,
) *IndexService {
	return &IndexService{
		DB:              db,
		CookieStore:     cookieStore,
		DiscordSession:  discordSession,
		DiscordBotState: discordBotState,
	}
}
