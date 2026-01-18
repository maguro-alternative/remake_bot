package mock

import (
	"context"

	"github.com/maguro-alternative/remake_bot/repository"
)

// Repository is an interface for repository.

type Repository interface {
	InsertLineBotIvByGuildID(ctx context.Context, guildId string) error
	GetAllColumnsLineBotIvByGuildID(ctx context.Context, guildID string) (repository.LineBotIv, error)
	GetLineBotIvNotClientByGuildID(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error)
	UpdateLineBotIv(ctx context.Context, lineBotIv *repository.LineBotIv) error
	InsertLineBot(ctx context.Context, lineBot *repository.LineBot) error
	GetAllColumnsLineBots(ctx context.Context) ([]*repository.LineBot, error)
	GetAllColumnsLineBotByGuildID(ctx context.Context, guildId string) (repository.LineBot, error)
	GetLineBotDefaultChannelIDByGuildID(ctx context.Context, guildID string) (repository.LineBotDefaultChannelID, error)
	GetLineBotNotClientByGuildID(ctx context.Context, guildID string) (repository.LineBotNotClient, error)
	UpdateLineBot(ctx context.Context, lineBot *repository.LineBot) error
	GetLineNgDiscordUserIDByChannelID(ctx context.Context, channelID string) ([]string, error)
	GetLineNgDiscordRoleIDByChannelID(ctx context.Context, channelID string) ([]string, error)
	InsertLineNgDiscordUserIDs(ctx context.Context, lineNgDiscordUserIDs []repository.LineNgDiscordUserIDAllCoulmns) error
	InsertLineNgDiscordRoleIDs(ctx context.Context, lineNgDiscordRoleIDs []repository.LineNgDiscordRoleIDAllCoulmns) error
	DeleteUserIDsNotInProvidedList(ctx context.Context, guildId string, lineNgDiscordUserIDs []repository.LineNgDiscordUserIDAllCoulmns) error
	DeleteRoleIDsNotInProvidedList(ctx context.Context, guildId string, lineNgDiscordRoleIDs []repository.LineNgDiscordRoleIDAllCoulmns) error
	InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []repository.LineNgDiscordMessageType) error
	DeleteMessageTypesNotInProvidedList(ctx context.Context, guildId string, lineNgDiscordTypes []repository.LineNgDiscordMessageType) error
	GetLineNgDiscordMessageTypeByChannelID(ctx context.Context, channelID string) ([]int, error)
	GetLinePostDiscordChannelByChannelID(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error)
	UpdateLinePostDiscordChannel(ctx context.Context, lineChannel repository.LinePostDiscordChannelAllColumns) error
	InsertLinePostDiscordChannelByChannelIDAndGuildID(ctx context.Context, channelID string, guildID string) error
	GetPermissionCodeByGuildIDAndType(ctx context.Context, guildID, permissionType string) (int64, error)
	GetPermissionCodesByGuildID(ctx context.Context, guildID string) ([]repository.PermissionCode, error)
	UpdatePermissionCodes(ctx context.Context, permissionsCode []repository.PermissionCode) error
	InsertPermissionUserIDs(ctx context.Context, permissionsUserID []repository.PermissionUserIDAllColumns) error
	InsertPermissionRoleIDs(ctx context.Context, permissionsRoleID []repository.PermissionRoleIDAllColumns) error
	GetGuildPermissionUserIDsAllColumnsByGuildID(ctx context.Context, guildID string) ([]repository.PermissionUserIDAllColumns, error)
	GetGuildPermissionRoleIDsAllColumnsByGuildID(ctx context.Context, guildID string) ([]repository.PermissionRoleIDAllColumns, error)
	GetPermissionUserIDsByGuildIDAndType(ctx context.Context, guildID, permissionType string) ([]repository.PermissionUserID, error)
	GetPermissionRoleIDsByGuildIDAndType(ctx context.Context, guildID, permissionType string) ([]repository.PermissionRoleID, error)
	DeletePermissionUserIDsByGuildID(ctx context.Context, guildId string) error
	DeletePermissionRoleIDsByGuildID(ctx context.Context, guildId string) error
}

var (
	_ Repository = (*repository.Repository)(nil)
)
