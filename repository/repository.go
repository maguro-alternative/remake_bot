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
	InsertLineBotIvByGuildIDFunc                          func(ctx context.Context, guildId string) error
	GetAllColumnsLineBotIvByGuildIDFunc                   func(ctx context.Context, guildID string) (LineBotIv, error)
	GetLineBotIvNotClientByGuildIDFunc                    func(ctx context.Context, guildID string) (LineBotIvNotClient, error)
	UpdateLineBotIvFunc                                   func(ctx context.Context, lineBotIv *LineBotIv) error
	InsertLineBotFunc                                     func(ctx context.Context, lineBot *LineBot) error
	GetAllColumnsLineBotsFunc                             func(ctx context.Context) ([]*LineBot, error)
	GetAllColumnsLineBotByGuildIDFunc                     func(ctx context.Context, guildId string) (LineBot, error)
	GetLineBotDefaultChannelIDByGuildIDFunc               func(ctx context.Context, guildID string) (LineBotDefaultChannelID, error)
	GetLineBotNotClientByGuildIDFunc                      func(ctx context.Context, guildID string) (LineBotNotClient, error)
	UpdateLineBotFunc                                     func(ctx context.Context, lineBot *LineBot) error
	GetLineNgDiscordUserIDByChannelIDFunc                 func(ctx context.Context, channelID string) ([]string, error)
	GetLineNgDiscordRoleIDByChannelIDFunc                 func(ctx context.Context, channelID string) ([]string, error)
	InsertLineNgDiscordUserIDsFunc                        func(ctx context.Context, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error
	InsertLineNgDiscordRoleIDsFunc                        func(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error
	DeleteUserIDsNotInProvidedListFunc                    func(ctx context.Context, guildId string, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error
	DeleteRoleIDsNotInProvidedListFunc                    func(ctx context.Context, guildId string, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error
	InsertLineNgDiscordMessageTypesFunc                   func(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error
	DeleteMessageTypesNotInProvidedListFunc               func(ctx context.Context, guildId string, lineNgDiscordTypes []LineNgDiscordMessageType) error
	GetLineNgDiscordMessageTypeByChannelIDFunc            func(ctx context.Context, channelID string) ([]int, error)
	GetLinePostDiscordChannelByChannelIDFunc              func(ctx context.Context, channelID string) (LinePostDiscordChannel, error)
	UpdateLinePostDiscordChannelFunc                      func(ctx context.Context, lineChannel LinePostDiscordChannelAllColumns) error
	InsertLinePostDiscordChannelByChannelIDAndGuildIDFunc func(ctx context.Context, channelID string, guildID string) error
	GetPermissionCodeByGuildIDAndTypeFunc                 func(ctx context.Context, guildID, permissionType string) (int64, error)
	GetPermissionCodesByGuildIDFunc                       func(ctx context.Context, guildID string) ([]PermissionCode, error)
	UpdatePermissionCodesFunc                             func(ctx context.Context, permissionsCode []PermissionCode) error
	InsertPermissionUserIDsFunc                           func(ctx context.Context, permissionsUserID []PermissionUserIDAllColumns) error
	InsertPermissionRoleIDsFunc                           func(ctx context.Context, permissionsRoleID []PermissionRoleIDAllColumns) error
	GetGuildPermissionUserIDsAllColumnsByGuildIDFunc      func(ctx context.Context, guildID string) ([]PermissionUserIDAllColumns, error)
	GetGuildPermissionRoleIDsAllColumnsByGuildIDFunc      func(ctx context.Context, guildID string) ([]PermissionRoleIDAllColumns, error)
	GetPermissionUserIDsByGuildIDAndTypeFunc              func(ctx context.Context, guildID, permissionType string) ([]PermissionUserID, error)
	GetPermissionRoleIDsByGuildIDAndTypeFunc              func(ctx context.Context, guildID, permissionType string) ([]PermissionRoleID, error)
	DeletePermissionUserIDsByGuildIDFunc                  func(ctx context.Context, guildId string) error
	DeletePermissionRoleIDsByGuildIDFunc                  func(ctx context.Context, guildId string) error
	GetVcSignalNgUsersByVcChannelIDAllColumnFunc          func(ctx context.Context, vcChannelID string) ([]*VcSignalNgUserAllColumn, error)
	GetVcSignalNgRolesByVcChannelIDAllColumnFunc          func(ctx context.Context, vcChannelID string) ([]*VcSignalNgRoleAllColumn, error)
	GetVcSignalChannelAllColumnByVcChannelIDFunc          func(ctx context.Context, vcChannelID string) (*VcSignalChannelAllColumn,error)
	GetVcSignalMentionUsersByVcChannelIDFunc              func(ctx context.Context, vcChannelID string) ([]*VcSignalMentionUser,error)
	GetVcSignalMentionRolesByVcChannelIDFunc              func(ctx context.Context, vcChannelID string) ([]*VcSignalMentionRole,error)
}

func (r *RepositoryFuncMock) InsertLineBotIvByGuildID(ctx context.Context, guildId string) error {
	return r.InsertLineBotIvByGuildIDFunc(ctx, guildId)
}

func (r *RepositoryFuncMock) GetAllColumnsLineBotIvByGuildID(ctx context.Context, guildID string) (LineBotIv, error) {
	return r.GetAllColumnsLineBotIvByGuildIDFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) GetLineBotIvNotClientByGuildID(ctx context.Context, guildID string) (LineBotIvNotClient, error) {
	return r.GetLineBotIvNotClientByGuildIDFunc(ctx, guildID)
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

func (r *RepositoryFuncMock) GetAllColumnsLineBotByGuildID(ctx context.Context, guildId string) (LineBot, error) {
	return r.GetAllColumnsLineBotByGuildIDFunc(ctx, guildId)
}

func (r *RepositoryFuncMock) GetLineBotDefaultChannelIDByGuildID(ctx context.Context, guildID string) (LineBotDefaultChannelID, error) {
	return r.GetLineBotDefaultChannelIDByGuildIDFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) GetLineBotNotClientByGuildID(ctx context.Context, guildID string) (LineBotNotClient, error) {
	return r.GetLineBotNotClientByGuildIDFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) UpdateLineBot(ctx context.Context, lineBot *LineBot) error {
	return r.UpdateLineBotFunc(ctx, lineBot)
}

func (r *RepositoryFuncMock) GetLineNgDiscordUserIDByChannelID(ctx context.Context, channelID string) ([]string, error) {
	return r.GetLineNgDiscordUserIDByChannelIDFunc(ctx, channelID)
}

func (r *RepositoryFuncMock) GetLineNgDiscordRoleIDByChannelID(ctx context.Context, channelID string) ([]string, error) {
	return r.GetLineNgDiscordRoleIDByChannelIDFunc(ctx, channelID)
}

func (r *RepositoryFuncMock) InsertLineNgDiscordUserIDs(ctx context.Context, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error {
	return r.InsertLineNgDiscordUserIDsFunc(ctx, lineNgDiscordUserIDs)
}

func (r *RepositoryFuncMock) InsertLineNgDiscordRoleIDs(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error {
	return r.InsertLineNgDiscordRoleIDsFunc(ctx, lineNgDiscordRoleIDs)
}

func (r *RepositoryFuncMock) DeleteUserIDsNotInProvidedList(ctx context.Context, guildId string, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error {
	return r.DeleteUserIDsNotInProvidedListFunc(ctx, guildId, lineNgDiscordUserIDs)
}

func (r *RepositoryFuncMock) DeleteRoleIDsNotInProvidedList(ctx context.Context, guildId string, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error {
	return r.DeleteRoleIDsNotInProvidedListFunc(ctx, guildId, lineNgDiscordRoleIDs)
}

func (r *RepositoryFuncMock) InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error {
	return r.InsertLineNgDiscordMessageTypesFunc(ctx, lineNgDiscordTypes)
}

func (r *RepositoryFuncMock) DeleteMessageTypesNotInProvidedList(ctx context.Context, guildId string, lineNgDiscordTypes []LineNgDiscordMessageType) error {
	return r.DeleteMessageTypesNotInProvidedListFunc(ctx, guildId, lineNgDiscordTypes)
}

func (r *RepositoryFuncMock) GetLineNgDiscordMessageTypeByChannelID(ctx context.Context, channelID string) ([]int, error) {
	return r.GetLineNgDiscordMessageTypeByChannelIDFunc(ctx, channelID)
}

func (r *RepositoryFuncMock) GetLinePostDiscordChannelByChannelID(ctx context.Context, channelID string) (LinePostDiscordChannel, error) {
	return r.GetLinePostDiscordChannelByChannelIDFunc(ctx, channelID)
}

func (r *RepositoryFuncMock) UpdateLinePostDiscordChannel(ctx context.Context, lineChannel LinePostDiscordChannelAllColumns) error {
	return r.UpdateLinePostDiscordChannelFunc(ctx, lineChannel)
}

func (r *RepositoryFuncMock) InsertLinePostDiscordChannelByChannelIDAndGuildID(ctx context.Context, channelID string, guildID string) error {
	return r.InsertLinePostDiscordChannelByChannelIDAndGuildIDFunc(ctx, channelID, guildID)
}

func (r *RepositoryFuncMock) GetPermissionCodeByGuildIDAndType(ctx context.Context, guildID, permissionType string) (int64, error) {
	return r.GetPermissionCodeByGuildIDAndTypeFunc(ctx, guildID, permissionType)
}

func (r *RepositoryFuncMock) GetPermissionCodesByGuildID(ctx context.Context, guildID string) ([]PermissionCode, error) {
	return r.GetPermissionCodesByGuildIDFunc(ctx, guildID)
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

func (r *RepositoryFuncMock) GetGuildPermissionUserIDsAllColumnsByGuildID(ctx context.Context, guildID string) ([]PermissionUserIDAllColumns, error) {
	return r.GetGuildPermissionUserIDsAllColumnsByGuildIDFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) GetGuildPermissionRoleIDsAllColumnsByGuildID(ctx context.Context, guildID string) ([]PermissionRoleIDAllColumns, error) {
	return r.GetGuildPermissionRoleIDsAllColumnsByGuildIDFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) GetPermissionUserIDsByGuildIDAndType(ctx context.Context, guildID, permissionType string) ([]PermissionUserID, error) {
	return r.GetPermissionUserIDsByGuildIDAndTypeFunc(ctx, guildID, permissionType)
}

func (r *RepositoryFuncMock) GetPermissionRoleIDsByGuildIDAndType(ctx context.Context, guildID, permissionType string) ([]PermissionRoleID, error) {
	return r.GetPermissionRoleIDsByGuildIDAndTypeFunc(ctx, guildID, permissionType)
}

func (r *RepositoryFuncMock) DeletePermissionUserIDsByGuildID(ctx context.Context, guildId string) error {
	return r.DeletePermissionUserIDsByGuildIDFunc(ctx, guildId)
}

func (r *RepositoryFuncMock) DeletePermissionRoleIDsByGuildID(ctx context.Context, guildId string) error {
	return r.DeletePermissionRoleIDsByGuildIDFunc(ctx, guildId)
}

func (r *RepositoryFuncMock) GetVcSignalNgUsersByVcChannelIDAllColumn(ctx context.Context, vcChannelID string) ([]*VcSignalNgUserAllColumn, error) {
	return r.GetVcSignalNgUsersByVcChannelIDAllColumnFunc(ctx, vcChannelID)
}

func (r *RepositoryFuncMock) GetVcSignalNgRolesByVcChannelIDAllColumn(ctx context.Context, vcChannelID string) ([]*VcSignalNgRoleAllColumn, error) {
	return r.GetVcSignalNgRolesByVcChannelIDAllColumnFunc(ctx, vcChannelID)
}

func (r *RepositoryFuncMock) GetVcSignalChannelAllColumnByVcChannelID(ctx context.Context, vcChannelID string) (*VcSignalChannelAllColumn, error) {
	return r.GetVcSignalChannelAllColumnByVcChannelIDFunc(ctx, vcChannelID)
}

func (r *RepositoryFuncMock) GetVcSignalMentionUsersByVcChannelID(ctx context.Context, vcChannelID string) ([]*VcSignalMentionUser,error) {
	return r.GetVcSignalMentionUsersByVcChannelIDFunc(ctx, vcChannelID)
}

func (r *RepositoryFuncMock) GetVcSignalMentionRolesByVcChannelID(ctx context.Context, vcChannelID string) ([]*VcSignalMentionRole,error) {
	return r.GetVcSignalMentionRolesByVcChannelIDFunc(ctx, vcChannelID)
}

// Repository is an interface for repository.
type RepositoryFunc interface {
	InsertLineBotIvByGuildID(ctx context.Context, guildId string) error
	GetAllColumnsLineBotIvByGuildID(ctx context.Context, guildID string) (LineBotIv, error)
	GetLineBotIvNotClientByGuildID(ctx context.Context, guildID string) (LineBotIvNotClient, error)
	UpdateLineBotIv(ctx context.Context, lineBotIv *LineBotIv) error
	InsertLineBot(ctx context.Context, lineBot *LineBot) error
	GetAllColumnsLineBots(ctx context.Context) ([]*LineBot, error)
	GetAllColumnsLineBotByGuildID(ctx context.Context, guildId string) (LineBot, error)
	GetLineBotDefaultChannelIDByGuildID(ctx context.Context, guildID string) (LineBotDefaultChannelID, error)
	GetLineBotNotClientByGuildID(ctx context.Context, guildID string) (LineBotNotClient, error)
	UpdateLineBot(ctx context.Context, lineBot *LineBot) error
	GetLineNgDiscordUserIDByChannelID(ctx context.Context, channelID string) ([]string, error)
	GetLineNgDiscordRoleIDByChannelID(ctx context.Context, channelID string) ([]string, error)
	InsertLineNgDiscordUserIDs(ctx context.Context, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error
	InsertLineNgDiscordRoleIDs(ctx context.Context, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error
	DeleteUserIDsNotInProvidedList(ctx context.Context, guildId string, lineNgDiscordUserIDs []LineNgDiscordUserIDAllCoulmns) error
	DeleteRoleIDsNotInProvidedList(ctx context.Context, guildId string, lineNgDiscordRoleIDs []LineNgDiscordRoleIDAllCoulmns) error
	InsertLineNgDiscordMessageTypes(ctx context.Context, lineNgDiscordTypes []LineNgDiscordMessageType) error
	DeleteMessageTypesNotInProvidedList(ctx context.Context, guildId string, lineNgDiscordTypes []LineNgDiscordMessageType) error
	GetLineNgDiscordMessageTypeByChannelID(ctx context.Context, channelID string) ([]int, error)
	GetLinePostDiscordChannelByChannelID(ctx context.Context, channelID string) (LinePostDiscordChannel, error)
	UpdateLinePostDiscordChannel(ctx context.Context, lineChannel LinePostDiscordChannelAllColumns) error
	InsertLinePostDiscordChannelByChannelIDAndGuildID(ctx context.Context, channelID string, guildID string) error
	GetPermissionCodeByGuildIDAndType(ctx context.Context, guildID, permissionType string) (int64, error)
	GetPermissionCodesByGuildID(ctx context.Context, guildID string) ([]PermissionCode, error)
	UpdatePermissionCodes(ctx context.Context, permissionsCode []PermissionCode) error
	InsertPermissionUserIDs(ctx context.Context, permissionsUserID []PermissionUserIDAllColumns) error
	InsertPermissionRoleIDs(ctx context.Context, permissionsRoleID []PermissionRoleIDAllColumns) error
	GetGuildPermissionUserIDsAllColumnsByGuildID(ctx context.Context, guildID string) ([]PermissionUserIDAllColumns, error)
	GetGuildPermissionRoleIDsAllColumnsByGuildID(ctx context.Context, guildID string) ([]PermissionRoleIDAllColumns, error)
	GetPermissionUserIDsByGuildIDAndType(ctx context.Context, guildID, permissionType string) ([]PermissionUserID, error)
	GetPermissionRoleIDsByGuildIDAndType(ctx context.Context, guildID, permissionType string) ([]PermissionRoleID, error)
	DeletePermissionUserIDsByGuildID(ctx context.Context, guildId string) error
	DeletePermissionRoleIDsByGuildID(ctx context.Context, guildId string) error
	GetVcSignalNgUsersByVcChannelIDAllColumn(ctx context.Context, vcChannelID string) ([]*VcSignalNgUserAllColumn, error)
	GetVcSignalNgRolesByVcChannelIDAllColumn(ctx context.Context, vcChannelID string) ([]*VcSignalNgRoleAllColumn, error)
	GetVcSignalChannelAllColumnByVcChannelID(ctx context.Context, vcChannelID string) (*VcSignalChannelAllColumn,error)
	GetVcSignalMentionUsersByVcChannelID(ctx context.Context, vcChannelID string) ([]*VcSignalMentionUser,error)
	GetVcSignalMentionRolesByVcChannelID(ctx context.Context, vcChannelID string) ([]*VcSignalMentionRole,error)
}

var (
	_ RepositoryFunc = (*Repository)(nil)
	_ RepositoryFunc = (*RepositoryFuncMock)(nil)
)
