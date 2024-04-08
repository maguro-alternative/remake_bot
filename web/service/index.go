package service

import (
	"context"
	"net/http"
	"io"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/sessions"
)

// A TODOService implements CRUD of TODO entities.
type IndexService struct {
	Client          *http.Client
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


type SessionMock struct {
	ChannelFunc                    func(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error)
	ChannelMessageSendFunc         func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelFileSendWithMessageFunc func(channelID string, content string, name string, r io.Reader, options ...discordgo.RequestOption) (*discordgo.Message, error)
	GuildFunc                      func(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error)
	GuildChannelsFunc              func(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Channel, err error)
	GuildMemberFunc                func(guildID string, userID string, options ...discordgo.RequestOption) (st *discordgo.Member, err error)
	GuildMembersFunc               func(guildID string, after string, limit int, options ...discordgo.RequestOption) (st []*discordgo.Member, err error)
	GuildRolesFunc                 func(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Role, err error)
	UserChannelPermissionsFunc     func(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error)
	UserGuildsFunc                 func(limit int, beforeID string, afterID string, options ...discordgo.RequestOption) (st []*discordgo.UserGuild, err error)
}

func (s *SessionMock) Channel(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error) {
	return s.ChannelFunc(channelID, options...)
}

func (s *SessionMock) ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	return s.ChannelMessageSendFunc(channelID, content, options...)
}

func (s *SessionMock) ChannelFileSendWithMessage(channelID string, content string, name string, r io.Reader, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	return s.ChannelFileSendWithMessageFunc(channelID, content, name, r, options...)
}

func (s *SessionMock) Guild(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error) {
	return s.GuildFunc(guildID, options...)
}

func (s *SessionMock) GuildChannels(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Channel, err error) {
	return s.GuildChannelsFunc(guildID, options...)
}

func (s *SessionMock) GuildMember(guildID string, userID string, options ...discordgo.RequestOption) (st *discordgo.Member, err error) {
	return s.GuildMemberFunc(guildID, userID, options...)
}

func (s *SessionMock) GuildMembers(guildID string, after string, limit int, options ...discordgo.RequestOption) (st []*discordgo.Member, err error) {
	return s.GuildMembersFunc(guildID, after, limit, options...)
}

func (s *SessionMock) GuildRoles(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Role, err error) {
	return s.GuildRolesFunc(guildID, options...)
}

func (s *SessionMock) UserChannelPermissions(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error) {
	return s.UserChannelPermissionsFunc(userID, channelID, fetchOptions...)
}

func (s *SessionMock) UserGuilds(limit int, beforeID string, afterID string, options ...discordgo.RequestOption) (st []*discordgo.UserGuild, err error) {
	return s.UserGuildsFunc(limit, beforeID, afterID, options...)
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
	_ Session    = (*discordgo.Session)(nil)
	_ Session    = (*SessionMock)(nil)
	_ Repository = (*repository.Repository)(nil)
)

// NewTODOService returns new TODOService.
func NewIndexService(
	client *http.Client,
	cookieStore *sessions.CookieStore,
	discordSession Session,
	discordBotState *discordgo.State,
) *IndexService {
	return &IndexService{
		Client:          client,
		CookieStore:     cookieStore,
		DiscordSession:  discordSession,
		DiscordBotState: discordBotState,
	}
}
