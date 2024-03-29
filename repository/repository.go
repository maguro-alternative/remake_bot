package repository

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type Repository struct {
	db db.Driver
}

func NewRepository(db db.Driver) *Repository {
	return &Repository{
		db: db,
	}
}

type RepositoryMock struct {
	InsertLineBotIvFunc                          func(ctx context.Context, guildId string) error
	GetAllColumnsLineBotIvFunc                   func(ctx context.Context, guildID string) (LineBotIv, error)
	GetLineBotIvNotClientFunc                    func(ctx context.Context, guildID string) (LineBotIvNotClient, error)
	UpdateLineBotIvFunc                          func(ctx context.Context, lineBotIv *LineBotIv) error
	InsertLineBotFunc                            func(ctx context.Context, lineBot *LineBot) error
	GetAllColumnsLineBotsFunc                    func(ctx context.Context) ([]*LineBot, error)
	GetAllColumnsLineBotFunc                     func(ctx context.Context, guildId string) (LineBot, error)
	GetLineBotDefaultChannelIDFunc               func(ctx context.Context, guildID string) (LineBotDefaultChannelID, error)
	GetLineBotNotClientFunc                      func(ctx context.Context, guildID string) (LineBotNotClient, error)
	UpdateLineBotFunc                            func(ctx context.Context, lineBot *LineBot) error
	GetLineNgDiscordIDFunc                       func(ctx context.Context, channelID string) ([]LineNgDiscordID, error)
	InsertLineNgDiscordIDsFunc                   func(ctx context.Context, lineNgDiscordIDs []LineNgDiscordIDAllCoulmns) error
	DeleteNotInsertLineNgDiscordIDsFunc          func(ctx context.Context, lineNgDiscordIDs []LineNgDiscordIDAllCoulmns) error
	InsertLineNgDiscordMessageTypesFunc          func(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error
	DeleteNotInsertLineNgDiscordMessageTypesFunc func(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error
	GetLineNgDiscordMessageTypeFunc              func(ctx context.Context, channelID string) ([]int, error)
	GetLinePostDiscordChannelFunc                func(ctx context.Context, channelID string) (LinePostDiscordChannel, error)
	UpdateLinePostDiscordChannelFunc             func(ctx context.Context, lineChannel LinePostDiscordChannelAllColumns) error
	InsertLinePostDiscordChannelFunc             func(ctx context.Context, channelID string, guildID string) error
	GetPermissionCodeFunc                        func(ctx context.Context, guildID, permissionType string) (int64, error)
	GetPermissionCodesFunc                       func(ctx context.Context, guildID string) ([]PermissionCode, error)
	UpdatePermissionCodesFunc                    func(ctx context.Context, permissionsCode []PermissionCode) error
	InsertPermissionIDsFunc                      func(ctx context.Context, permissionsID []PermissionIDAllColumns) error
	GetGuildPermissionIDsAllColumnsFunc          func(ctx context.Context, guildID string) ([]PermissionIDAllColumns, error)
	GetPermissionIDsFunc                         func(ctx context.Context, guildID, permissionType string) ([]PermissionID, error)
	DeletePermissionIDsFunc                      func(ctx context.Context, guildId string) error
}

func (r *RepositoryMock) InsertLineBotIv(ctx context.Context, guildId string) error {
	return r.InsertLineBotIvFunc(ctx, guildId)
}

func (r *RepositoryMock) GetAllColumnsLineBotIv(ctx context.Context, guildID string) (LineBotIv, error) {
	return r.GetAllColumnsLineBotIvFunc(ctx, guildID)
}

func (r *RepositoryMock) GetLineBotIvNotClient(ctx context.Context, guildID string) (LineBotIvNotClient, error) {
	return r.GetLineBotIvNotClientFunc(ctx, guildID)
}

func (r *RepositoryMock) UpdateLineBotIv(ctx context.Context, lineBotIv *LineBotIv) error {
	return r.UpdateLineBotIvFunc(ctx, lineBotIv)
}

func (r *RepositoryMock) InsertLineBot(ctx context.Context, lineBot *LineBot) error {
	return r.InsertLineBotFunc(ctx, lineBot)
}

func (r *RepositoryMock) GetAllColumnsLineBots(ctx context.Context) ([]*LineBot, error) {
	return r.GetAllColumnsLineBotsFunc(ctx)
}

func (r *RepositoryMock) GetAllColumnsLineBot(ctx context.Context, guildId string) (LineBot, error) {
	return r.GetAllColumnsLineBotFunc(ctx, guildId)
}

func (r *RepositoryMock) GetLineBotDefaultChannelID(ctx context.Context, guildID string) (LineBotDefaultChannelID, error) {
	return r.GetLineBotDefaultChannelIDFunc(ctx, guildID)
}

func (r *RepositoryMock) GetLineBotNotClient(ctx context.Context, guildID string) (LineBotNotClient, error) {
	return r.GetLineBotNotClientFunc(ctx, guildID)
}

func (r *RepositoryMock) UpdateLineBot(ctx context.Context, lineBot *LineBot) error {
	return r.UpdateLineBotFunc(ctx, lineBot)
}

func (r *RepositoryMock) GetLineNgDiscordID(ctx context.Context, channelID string) ([]LineNgDiscordID, error) {
	return r.GetLineNgDiscordIDFunc(ctx, channelID)
}

func (r *RepositoryMock) InsertLineNgDiscordIDs(ctx context.Context, lineNgDiscordIDs []LineNgDiscordIDAllCoulmns) error {
	return r.InsertLineNgDiscordIDsFunc(ctx, lineNgDiscordIDs)
}

func (r *RepositoryMock) DeleteNotInsertLineNgDiscordIDs(ctx context.Context, lineNgDiscordIDs []LineNgDiscordIDAllCoulmns) error {
	return r.DeleteNotInsertLineNgDiscordIDsFunc(ctx, lineNgDiscordIDs)
}

func (r *RepositoryMock) InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error {
	return r.InsertLineNgDiscordMessageTypesFunc(ctx, lineNgDiscordTypes)
}

func (r *RepositoryMock) DeleteNotInsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error {
	return r.DeleteNotInsertLineNgDiscordMessageTypesFunc(ctx, lineNgDiscordTypes)
}

func (r *RepositoryMock) GetLineNgDiscordMessageType(ctx context.Context, channelID string) ([]int, error) {
	return r.GetLineNgDiscordMessageTypeFunc(ctx, channelID)
}

func (r *RepositoryMock) GetLinePostDiscordChannel(ctx context.Context, channelID string) (LinePostDiscordChannel, error) {
	return r.GetLinePostDiscordChannelFunc(ctx, channelID)
}

func (r *RepositoryMock) UpdateLinePostDiscordChannel(ctx context.Context, lineChannel LinePostDiscordChannelAllColumns) error {
	return r.UpdateLinePostDiscordChannelFunc(ctx, lineChannel)
}

func (r *RepositoryMock) InsertLinePostDiscordChannel(ctx context.Context, channelID string, guildID string) error {
	return r.InsertLinePostDiscordChannelFunc(ctx, channelID, guildID)
}

func (r *RepositoryMock) GetPermissionCode(ctx context.Context, guildID, permissionType string) (int64, error) {
	return r.GetPermissionCodeFunc(ctx, guildID, permissionType)
}

func (r *RepositoryMock) GetPermissionCodes(ctx context.Context, guildID string) ([]PermissionCode, error) {
	return r.GetPermissionCodesFunc(ctx, guildID)
}

func (r *RepositoryMock) UpdatePermissionCodes(ctx context.Context, permissionsCode []PermissionCode) error {
	return r.UpdatePermissionCodesFunc(ctx, permissionsCode)
}

func (r *RepositoryMock) InsertPermissionIDs(ctx context.Context, permissionsID []PermissionIDAllColumns) error {
	return r.InsertPermissionIDsFunc(ctx, permissionsID)
}

func (r *RepositoryMock) GetGuildPermissionIDsAllColumns(ctx context.Context, guildID string) ([]PermissionIDAllColumns, error) {
	return r.GetGuildPermissionIDsAllColumnsFunc(ctx, guildID)
}

func (r *RepositoryMock) GetPermissionIDs(ctx context.Context, guildID, permissionType string) ([]PermissionID, error) {
	return r.GetPermissionIDsFunc(ctx, guildID, permissionType)
}

func (r *RepositoryMock) DeletePermissionIDs(ctx context.Context, guildId string) error {
	return r.DeletePermissionIDsFunc(ctx, guildId)
}
