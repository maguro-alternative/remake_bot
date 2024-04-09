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

var (
	_ Repository = (*repository.Repository)(nil)
)
