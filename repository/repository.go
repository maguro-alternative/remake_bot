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

type RepositoryFuncMock struct {
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

func (r *RepositoryFuncMock) InsertLineBotIv(ctx context.Context, guildId string) error {
	return r.InsertLineBotIvFunc(ctx, guildId)
}

func (r *RepositoryFuncMock) GetAllColumnsLineBotIv(ctx context.Context, guildID string) (LineBotIv, error) {
	return r.GetAllColumnsLineBotIvFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) GetLineBotIvNotClient(ctx context.Context, guildID string) (LineBotIvNotClient, error) {
	return r.GetLineBotIvNotClientFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) UpdateLineBotIv(ctx context.Context, lineBotIv *LineBotIv) error {
	return r.UpdateLineBotIvFunc(ctx, lineBotIv)
}

func (r *RepositoryFuncMock) InsertLineBot(ctx context.Context, lineBot *LineBot) error {
	return r.InsertLineBotFunc(ctx, lineBot)
}

func (r *RepositoryFuncMock) GetAllColumnsLineBots(ctx context.Context) ([]*LineBot, error) {
	return r.GetAllColumnsLineBotsFunc(ctx)
}

func (r *RepositoryFuncMock) GetAllColumnsLineBot(ctx context.Context, guildId string) (LineBot, error) {
	return r.GetAllColumnsLineBotFunc(ctx, guildId)
}

func (r *RepositoryFuncMock) GetLineBotDefaultChannelID(ctx context.Context, guildID string) (LineBotDefaultChannelID, error) {
	return r.GetLineBotDefaultChannelIDFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) GetLineBotNotClient(ctx context.Context, guildID string) (LineBotNotClient, error) {
	return r.GetLineBotNotClientFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) UpdateLineBot(ctx context.Context, lineBot *LineBot) error {
	return r.UpdateLineBotFunc(ctx, lineBot)
}

func (r *RepositoryFuncMock) GetLineNgDiscordID(ctx context.Context, channelID string) ([]LineNgDiscordID, error) {
	return r.GetLineNgDiscordIDFunc(ctx, channelID)
}

func (r *RepositoryFuncMock) InsertLineNgDiscordIDs(ctx context.Context, lineNgDiscordIDs []LineNgDiscordIDAllCoulmns) error {
	return r.InsertLineNgDiscordIDsFunc(ctx, lineNgDiscordIDs)
}

func (r *RepositoryFuncMock) DeleteNotInsertLineNgDiscordIDs(ctx context.Context, lineNgDiscordIDs []LineNgDiscordIDAllCoulmns) error {
	return r.DeleteNotInsertLineNgDiscordIDsFunc(ctx, lineNgDiscordIDs)
}

func (r *RepositoryFuncMock) InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error {
	return r.InsertLineNgDiscordMessageTypesFunc(ctx, lineNgDiscordTypes)
}

func (r *RepositoryFuncMock) DeleteNotInsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error {
	return r.DeleteNotInsertLineNgDiscordMessageTypesFunc(ctx, lineNgDiscordTypes)
}

func (r *RepositoryFuncMock) GetLineNgDiscordMessageType(ctx context.Context, channelID string) ([]int, error) {
	return r.GetLineNgDiscordMessageTypeFunc(ctx, channelID)
}

func (r *RepositoryFuncMock) GetLinePostDiscordChannel(ctx context.Context, channelID string) (LinePostDiscordChannel, error) {
	return r.GetLinePostDiscordChannelFunc(ctx, channelID)
}

func (r *RepositoryFuncMock) UpdateLinePostDiscordChannel(ctx context.Context, lineChannel LinePostDiscordChannelAllColumns) error {
	return r.UpdateLinePostDiscordChannelFunc(ctx, lineChannel)
}

func (r *RepositoryFuncMock) InsertLinePostDiscordChannel(ctx context.Context, channelID string, guildID string) error {
	return r.InsertLinePostDiscordChannelFunc(ctx, channelID, guildID)
}

func (r *RepositoryFuncMock) GetPermissionCode(ctx context.Context, guildID, permissionType string) (int64, error) {
	return r.GetPermissionCodeFunc(ctx, guildID, permissionType)
}

func (r *RepositoryFuncMock) GetPermissionCodes(ctx context.Context, guildID string) ([]PermissionCode, error) {
	return r.GetPermissionCodesFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) UpdatePermissionCodes(ctx context.Context, permissionsCode []PermissionCode) error {
	return r.UpdatePermissionCodesFunc(ctx, permissionsCode)
}

func (r *RepositoryFuncMock) InsertPermissionIDs(ctx context.Context, permissionsID []PermissionIDAllColumns) error {
	return r.InsertPermissionIDsFunc(ctx, permissionsID)
}

func (r *RepositoryFuncMock) GetGuildPermissionIDsAllColumns(ctx context.Context, guildID string) ([]PermissionIDAllColumns, error) {
	return r.GetGuildPermissionIDsAllColumnsFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) GetPermissionIDs(ctx context.Context, guildID, permissionType string) ([]PermissionID, error) {
	return r.GetPermissionIDsFunc(ctx, guildID, permissionType)
}

func (r *RepositoryFuncMock) DeletePermissionIDs(ctx context.Context, guildId string) error {
	return r.DeletePermissionIDsFunc(ctx, guildId)
}

// Repository is an interface for repository.
type RepositoryFunc interface {
	InsertLineBotIv(ctx context.Context, guildId string) error
	GetAllColumnsLineBotIv(ctx context.Context, guildID string) (LineBotIv, error)
	GetLineBotIvNotClient(ctx context.Context, guildID string) (LineBotIvNotClient, error)
	UpdateLineBotIv(ctx context.Context, lineBotIv *LineBotIv) error
	InsertLineBot(ctx context.Context, lineBot *LineBot) error
	GetAllColumnsLineBots(ctx context.Context) ([]*LineBot, error)
	GetAllColumnsLineBot(ctx context.Context, guildId string) (LineBot, error)
	GetLineBotDefaultChannelID(ctx context.Context, guildID string) (LineBotDefaultChannelID, error)
	GetLineBotNotClient(ctx context.Context, guildID string) (LineBotNotClient, error)
	UpdateLineBot(ctx context.Context, lineBot *LineBot) error
	GetLineNgDiscordID(ctx context.Context, channelID string) ([]LineNgDiscordID, error)
	InsertLineNgDiscordIDs(ctx context.Context, lineNgDiscordIDs []LineNgDiscordIDAllCoulmns) error
	DeleteNotInsertLineNgDiscordIDs(ctx context.Context, lineNgDiscordIDs []LineNgDiscordIDAllCoulmns) error
	InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error
	DeleteNotInsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error
	GetLineNgDiscordMessageType(ctx context.Context, channelID string) ([]int, error)
	GetLinePostDiscordChannel(ctx context.Context, channelID string) (LinePostDiscordChannel, error)
	UpdateLinePostDiscordChannel(ctx context.Context, lineChannel LinePostDiscordChannelAllColumns) error
	InsertLinePostDiscordChannel(ctx context.Context, channelID string, guildID string) error
	GetPermissionCode(ctx context.Context, guildID, permissionType string) (int64, error)
	GetPermissionCodes(ctx context.Context, guildID string) ([]PermissionCode, error)
	UpdatePermissionCodes(ctx context.Context, permissionsCode []PermissionCode) error
	InsertPermissionIDs(ctx context.Context, permissionsID []PermissionIDAllColumns) error
	GetGuildPermissionIDsAllColumns(ctx context.Context, guildID string) ([]PermissionIDAllColumns, error)
	GetPermissionIDs(ctx context.Context, guildID, permissionType string) ([]PermissionID, error)
	DeletePermissionIDs(ctx context.Context, guildId string) error
}

var (
	_ RepositoryFunc = (*Repository)(nil)
	_ RepositoryFunc = (*RepositoryFuncMock)(nil)
)
