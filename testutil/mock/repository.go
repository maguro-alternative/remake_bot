package mock

import (
	"context"

	"github.com/maguro-alternative/remake_bot/repository"
)

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
	GetLineNgDiscordUserID(ctx context.Context, channelID string) ([]string, error)
	GetLineNgDiscordRoleID(ctx context.Context, channelID string) ([]string, error)
	InsertLineNgDiscordUserIDs(ctx context.Context, lineNgDiscordUserIDs []repository.LineNgDiscordUserIDAllCoulmns) error
	InsertLineNgDiscordRoleIDs(ctx context.Context, lineNgDiscordRoleIDs []repository.LineNgDiscordRoleIDAllCoulmns) error
	DeleteNotInsertLineNgDiscordUserIDs(ctx context.Context, lineNgDiscordUserIDs []repository.LineNgDiscordUserIDAllCoulmns) error
	DeleteRoleIDsNotInProvidedList(ctx context.Context, lineNgDiscordRoleIDs []repository.LineNgDiscordRoleIDAllCoulmns) error
	InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []repository.LineNgDiscordMessageType) error
	DeleteMessageTypesNotInProvidedList(ctx context.Context, lineNgDiscordTypes []repository.LineNgDiscordMessageType) error
	GetLineNgDiscordMessageType(ctx context.Context, channelID string) ([]int, error)
	GetLinePostDiscordChannel(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error)
	UpdateLinePostDiscordChannel(ctx context.Context, lineChannel repository.LinePostDiscordChannelAllColumns) error
	InsertLinePostDiscordChannel(ctx context.Context, channelID string, guildID string) error
	GetPermissionCode(ctx context.Context, guildID, permissionType string) (int64, error)
	GetPermissionCodes(ctx context.Context, guildID string) ([]repository.PermissionCode, error)
	UpdatePermissionCodes(ctx context.Context, permissionsCode []repository.PermissionCode) error
	InsertPermissionUserIDs(ctx context.Context, permissionsUserID []repository.PermissionUserIDAllColumns) error
	InsertPermissionRoleIDs(ctx context.Context, permissionsRoleID []repository.PermissionRoleIDAllColumns) error
	GetGuildPermissionUserIDsAllColumns(ctx context.Context, guildID string) ([]repository.PermissionUserIDAllColumns, error)
	GetGuildPermissionRoleIDsAllColumns(ctx context.Context, guildID string) ([]repository.PermissionRoleIDAllColumns, error)
	GetPermissionUserIDs(ctx context.Context, guildID, permissionType string) ([]repository.PermissionUserID, error)
	GetPermissionRoleIDs(ctx context.Context, guildID, permissionType string) ([]repository.PermissionRoleID, error)
	DeletePermissionUserIDs(ctx context.Context, guildId string) error
	DeletePermissionRoleIDs(ctx context.Context, guildId string) error
}

var (
	_ Repository = (*repository.Repository)(nil)
)
