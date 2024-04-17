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
	GetLineNgDiscordUserIDFunc                   func(ctx context.Context, channelID string) ([]string, error)
	GetLineNgDiscordRoleIDFunc                   func(ctx context.Context, channelID string) ([]string, error)
	InsertLineNgDiscordUserIDsFunc               func(ctx context.Context, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error
	InsertLineNgDiscordRoleIDsFunc               func(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error
	DeleteNotInsertLineNgDiscordUserIDsFunc      func(ctx context.Context, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error
	DeleteNotInsertLineNgDiscordRoleIDsFunc      func(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error
	InsertLineNgDiscordMessageTypesFunc          func(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error
	DeleteNotInsertLineNgDiscordMessageTypesFunc func(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error
	GetLineNgDiscordMessageTypeFunc              func(ctx context.Context, channelID string) ([]int, error)
	GetLinePostDiscordChannelFunc                func(ctx context.Context, channelID string) (LinePostDiscordChannel, error)
	UpdateLinePostDiscordChannelFunc             func(ctx context.Context, lineChannel LinePostDiscordChannelAllColumns) error
	InsertLinePostDiscordChannelFunc             func(ctx context.Context, channelID string, guildID string) error
	GetPermissionCodeFunc                        func(ctx context.Context, guildID, permissionType string) (int64, error)
	GetPermissionCodesFunc                       func(ctx context.Context, guildID string) ([]PermissionCode, error)
	UpdatePermissionCodesFunc                    func(ctx context.Context, permissionsCode []PermissionCode) error
	InsertPermissionUserIDsFunc                  func(ctx context.Context, permissionsUserID []PermissionUserIDAllColumns) error
	InsertPermissionRoleIDsFunc                  func(ctx context.Context, permissionsRoleID []PermissionRoleIDAllColumns) error
	GetGuildPermissionUserIDsAllColumnsFunc      func(ctx context.Context, guildID string) ([]PermissionUserIDAllColumns, error)
	GetGuildPermissionRoleIDsAllColumnsFunc      func(ctx context.Context, guildID string) ([]PermissionRoleIDAllColumns, error)
	GetPermissionUserIDsFunc                     func(ctx context.Context, guildID, permissionType string) ([]PermissionUserID, error)
	GetPermissionRoleIDsFunc                     func(ctx context.Context, guildID, permissionType string) ([]PermissionRoleID, error)
	DeletePermissionUserIDsFunc                  func(ctx context.Context, guildId string) error
	DeletePermissionRoleIDsFunc                  func(ctx context.Context, guildId string) error
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

func (r *RepositoryFuncMock) GetLineNgDiscordUserID(ctx context.Context, channelID string) ([]string, error) {
	return r.GetLineNgDiscordUserIDFunc(ctx, channelID)
}

func (r *RepositoryFuncMock) GetLineNgDiscordRoleID(ctx context.Context, channelID string) ([]string, error) {
	return r.GetLineNgDiscordRoleIDFunc(ctx, channelID)
}

func (r *RepositoryFuncMock) InsertLineNgDiscordUserIDs(ctx context.Context, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error {
	return r.InsertLineNgDiscordUserIDsFunc(ctx, lineNgDiscordUserIDs)
}

func (r *RepositoryFuncMock) InsertLineNgDiscordRoleIDs(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error {
	return r.InsertLineNgDiscordRoleIDsFunc(ctx, lineNgDiscordRoleIDs)
}

func (r *RepositoryFuncMock) DeleteNotInsertLineNgDiscordUserIDs(ctx context.Context, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error {
	return r.DeleteNotInsertLineNgDiscordUserIDsFunc(ctx, lineNgDiscordUserIDs)
}

func (r *RepositoryFuncMock) DeleteNotInsertLineNgDiscordRoleIDs(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error {
	return r.DeleteNotInsertLineNgDiscordRoleIDsFunc(ctx, lineNgDiscordRoleIDs)
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

func (r *RepositoryFuncMock) InsertPermissionUserIDs(ctx context.Context, permissionsUserID []PermissionUserIDAllColumns) error {
	return r.InsertPermissionUserIDsFunc(ctx, permissionsUserID)
}

func (r *RepositoryFuncMock) InsertPermissionRoleIDs(ctx context.Context, permissionsRoleID []PermissionRoleIDAllColumns) error {
	return r.InsertPermissionRoleIDsFunc(ctx, permissionsRoleID)
}

func (r *RepositoryFuncMock) GetGuildPermissionUserIDsAllColumns(ctx context.Context, guildID string) ([]PermissionUserIDAllColumns, error) {
	return r.GetGuildPermissionUserIDsAllColumnsFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) GetGuildPermissionRoleIDsAllColumns(ctx context.Context, guildID string) ([]PermissionRoleIDAllColumns, error) {
	return r.GetGuildPermissionRoleIDsAllColumnsFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) GetPermissionUserIDs(ctx context.Context, guildID, permissionType string) ([]PermissionUserID, error) {
	return r.GetPermissionUserIDsFunc(ctx, guildID, permissionType)
}

func (r *RepositoryFuncMock) GetPermissionRoleIDs(ctx context.Context, guildID, permissionType string) ([]PermissionRoleID, error) {
	return r.GetPermissionRoleIDsFunc(ctx, guildID, permissionType)
}

func (r *RepositoryFuncMock) DeletePermissionUserIDs(ctx context.Context, guildId string) error {
	return r.DeletePermissionUserIDsFunc(ctx, guildId)
}

func (r *RepositoryFuncMock) DeletePermissionRoleIDs(ctx context.Context, guildId string) error {
	return r.DeletePermissionRoleIDsFunc(ctx, guildId)
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
	GetLineNgDiscordUserID(ctx context.Context, channelID string) ([]string, error)
	GetLineNgDiscordRoleID(ctx context.Context, channelID string) ([]string, error)
	InsertLineNgDiscordUserIDs(ctx context.Context, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error
	InsertLineNgDiscordRoleIDs(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error
	DeleteNotInsertLineNgDiscordUserIDs(ctx context.Context, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error
	DeleteNotInsertLineNgDiscordRoleIDs(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error
	InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error
	DeleteNotInsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error
	GetLineNgDiscordMessageType(ctx context.Context, channelID string) ([]int, error)
	GetLinePostDiscordChannel(ctx context.Context, channelID string) (LinePostDiscordChannel, error)
	UpdateLinePostDiscordChannel(ctx context.Context, lineChannel LinePostDiscordChannelAllColumns) error
	InsertLinePostDiscordChannel(ctx context.Context, channelID string, guildID string) error
	GetPermissionCode(ctx context.Context, guildID, permissionType string) (int64, error)
	GetPermissionCodes(ctx context.Context, guildID string) ([]PermissionCode, error)
	UpdatePermissionCodes(ctx context.Context, permissionsCode []PermissionCode) error
	InsertPermissionUserIDs(ctx context.Context, permissionsUserID []PermissionUserIDAllColumns) error
	InsertPermissionRoleIDs(ctx context.Context, permissionsRoleID []PermissionRoleIDAllColumns) error
	GetGuildPermissionUserIDsAllColumns(ctx context.Context, guildID string) ([]PermissionUserIDAllColumns, error)
	GetGuildPermissionRoleIDsAllColumns(ctx context.Context, guildID string) ([]PermissionRoleIDAllColumns, error)
	GetPermissionUserIDs(ctx context.Context, guildID, permissionType string) ([]PermissionUserID, error)
	GetPermissionRoleIDs(ctx context.Context, guildID, permissionType string) ([]PermissionRoleID, error)
	DeletePermissionUserIDs(ctx context.Context, guildId string) error
	DeletePermissionRoleIDs(ctx context.Context, guildId string) error
}

var (
	_ RepositoryFunc = (*Repository)(nil)
	_ RepositoryFunc = (*RepositoryFuncMock)(nil)
)
