package repository

import (
	"context"
	"time"

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
	GetVcSignalNgUserIDsByVcChannelIDFunc                 func(ctx context.Context, vcChannelID string) ([]string, error)
	GetVcSignalNgRoleIDsByVcChannelIDFunc                 func(ctx context.Context, vcChannelID string) ([]string, error)
	GetVcSignalChannelAllColumnByVcChannelIDFunc          func(ctx context.Context, vcChannelID string) (*VcSignalChannelAllColumn, error)
	GetVcSignalMentionUserIDsByVcChannelIDFunc            func(ctx context.Context, vcChannelID string) ([]string, error)
	GetVcSignalMentionRoleIDsByVcChannelIDFunc            func(ctx context.Context, vcChannelID string) ([]string, error)
	UpdateVcSignalChannelFunc                             func(ctx context.Context, vcSignalChannelNotGuildID VcSignalChannelNotGuildID) error
	InsertVcSignalChannelFunc                             func(ctx context.Context, vcChannelID string, guildID, sendChannelID string) error
	InsertVcSignalNgUserFunc                              func(ctx context.Context, vcChannelID string, guildID string, userID string) error
	InsertVcSignalNgRoleFunc                              func(ctx context.Context, vcChannelID string, guildID string, roleID string) error
	InsertVcSignalMentionUserFunc                         func(ctx context.Context, vcChannelID string, guildID string, userID string) error
	InsertVcSignalMentionRoleFunc                         func(ctx context.Context, vcChannelID string, guildID string, roleID string) error
	DeleteVcSignalNgUsersNotInProvidedListFunc            func(ctx context.Context, vcChannelID string, userIDs []string) error
	DeleteVcSignalNgRolesNotInProvidedListFunc            func(ctx context.Context, vcChannelID string, roleIDs []string) error
	DeleteVcSignalMentionUsersNotInProvidedListFunc       func(ctx context.Context, vcChannelID string, userIDs []string) error
	DeleteVcSignalMentionRolesNotInProvidedListFunc       func(ctx context.Context, vcChannelID string, roleIDs []string) error
	InsertWebhookFunc                                     func(ctx context.Context, guildID string, webhookID string, subscriptionType string, subscriptionID string, lastPostedAt time.Time) (int64, error)
	GetAllColumnsWebhooksByGuildIDFunc                    func(ctx context.Context, guildID string) ([]*Webhook, error)
	UpdateWebhookWithLastPostedAtFunc                     func(ctx context.Context, webhookSerialID int64, lastPostedAt time.Time) error
	UpdateWebhookWithWebhookIDAndSubscriptionFunc         func(ctx context.Context, webhookSerialID int64, webhookID string, subscriptionID string, subscriptionType string) error
	DeleteWebhookByWebhookSerialIDFunc                    func(ctx context.Context, webhookSerialID int64) error
	InsertWebhookWordFunc                                 func(ctx context.Context, webhookSerialID int64, condition string, word string) error
	GetWebhookWordWithWebhookSerialIDAndConditionFunc     func(ctx context.Context, webhookSerialID int64, condition string) ([]*WebhookWord, error)
	//GetWebhookWordWithWebhookSerialIDsFunc                func(ctx context.Context, webhookSerialIDs []int64) ([]*WebhookWord, error)
	DeleteWebhookWordsNotInProvidedListFunc               func(ctx context.Context, webhookSerialID int64, conditions string, words []string) error
	InsertWebhookUserMentionFunc                          func(ctx context.Context, webhookSerialID int64, userID string) error
	GetWebhookUserMentionWithWebhookSerialIDFunc          func(ctx context.Context, webhookSerialID int64) ([]*WebhookUserMention, error)
	//GetWebhookUserMentionWithWebhookSerialIDsFunc         func(ctx context.Context, webhookSerialIDs []int64) ([]*WebhookUserMention, error)
	DeleteWebhookUserMentionsNotInProvidedListFunc        func(ctx context.Context, webhookSerialID int64, userIDs []string) error
	InsertWebhookRoleMentionFunc                          func(ctx context.Context, webhookSerialID int64, roleID string) error
	GetWebhookRoleMentionWithWebhookSerialIDFunc          func(ctx context.Context, webhookSerialID int64) ([]*WebhookRoleMention, error)
	//GetWebhookRoleMentionWithWebhookSerialIDsFunc         func(ctx context.Context, webhookSerialIDs []int64) ([]*WebhookRoleMention, error)
	DeleteWebhookRoleMentionsNotInProvidedListFunc        func(ctx context.Context, webhookSerialID int64, roleIDs []string) error
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

func (r *RepositoryFuncMock) GetVcSignalNgUserIDsByVcChannelID(ctx context.Context, vcChannelID string) ([]string, error) {
	return r.GetVcSignalNgUserIDsByVcChannelIDFunc(ctx, vcChannelID)
}

func (r *RepositoryFuncMock) GetVcSignalNgRoleIDsByVcChannelID(ctx context.Context, vcChannelID string) ([]string, error) {
	return r.GetVcSignalNgRoleIDsByVcChannelIDFunc(ctx, vcChannelID)
}

func (r *RepositoryFuncMock) GetVcSignalChannelAllColumnByVcChannelID(ctx context.Context, vcChannelID string) (*VcSignalChannelAllColumn, error) {
	return r.GetVcSignalChannelAllColumnByVcChannelIDFunc(ctx, vcChannelID)
}

func (r *RepositoryFuncMock) GetVcSignalMentionUserIDsByVcChannelID(ctx context.Context, vcChannelID string) ([]string, error) {
	return r.GetVcSignalMentionUserIDsByVcChannelIDFunc(ctx, vcChannelID)
}

func (r *RepositoryFuncMock) GetVcSignalMentionRoleIDsByVcChannelID(ctx context.Context, vcChannelID string) ([]string, error) {
	return r.GetVcSignalMentionRoleIDsByVcChannelIDFunc(ctx, vcChannelID)
}

func (r *RepositoryFuncMock) UpdateVcSignalChannel(ctx context.Context, vcSignalChannelNotGuildID VcSignalChannelNotGuildID) error {
	return r.UpdateVcSignalChannelFunc(ctx, vcSignalChannelNotGuildID)
}

func (r *RepositoryFuncMock) InsertVcSignalChannel(ctx context.Context, vcChannelID string, guildID, sendChannelID string) error {
	return r.InsertVcSignalChannelFunc(ctx, vcChannelID, guildID, sendChannelID)
}

func (r *RepositoryFuncMock) InsertVcSignalNgUser(ctx context.Context, vcChannelID string, guildID string, userID string) error {
	return r.InsertVcSignalNgUserFunc(ctx, vcChannelID, guildID, userID)
}

func (r *RepositoryFuncMock) InsertVcSignalNgRole(ctx context.Context, vcChannelID string, guildID string, roleID string) error {
	return r.InsertVcSignalNgRoleFunc(ctx, vcChannelID, guildID, roleID)
}

func (r *RepositoryFuncMock) InsertVcSignalMentionUser(ctx context.Context, vcChannelID string, guildID string, userID string) error {
	return r.InsertVcSignalMentionUserFunc(ctx, vcChannelID, guildID, userID)
}

func (r *RepositoryFuncMock) InsertVcSignalMentionRole(ctx context.Context, vcChannelID string, guildID string, roleID string) error {
	return r.InsertVcSignalMentionRoleFunc(ctx, vcChannelID, guildID, roleID)
}

func (r *RepositoryFuncMock) DeleteVcSignalNgUsersNotInProvidedList(ctx context.Context, vcChannelID string, userIDs []string) error {
	return r.DeleteVcSignalNgUsersNotInProvidedListFunc(ctx, vcChannelID, userIDs)
}

func (r *RepositoryFuncMock) DeleteVcSignalNgRolesNotInProvidedList(ctx context.Context, vcChannelID string, roleIDs []string) error {
	return r.DeleteVcSignalNgRolesNotInProvidedListFunc(ctx, vcChannelID, roleIDs)
}

func (r *RepositoryFuncMock) DeleteVcSignalMentionUsersNotInProvidedList(ctx context.Context, vcChannelID string, userIDs []string) error {
	return r.DeleteVcSignalMentionUsersNotInProvidedListFunc(ctx, vcChannelID, userIDs)
}

func (r *RepositoryFuncMock) DeleteVcSignalMentionRolesNotInProvidedList(ctx context.Context, vcChannelID string, roleIDs []string) error {
	return r.DeleteVcSignalMentionRolesNotInProvidedListFunc(ctx, vcChannelID, roleIDs)
}

func (r *RepositoryFuncMock) InsertWebhook(ctx context.Context, guildID string, webhookID string, subscriptionType string, subscriptionID string, lastPostedAt time.Time) (int64, error) {
	return r.InsertWebhookFunc(ctx, guildID, webhookID, subscriptionType, subscriptionID, lastPostedAt)
}

func (r *RepositoryFuncMock) GetAllColumnsWebhooksByGuildID(ctx context.Context, guildID string) ([]*Webhook, error) {
	return r.GetAllColumnsWebhooksByGuildIDFunc(ctx, guildID)
}

func (r *RepositoryFuncMock) UpdateWebhookWithLastPostedAt(ctx context.Context, webhookSerialID int64, lastPostedAt time.Time) error {
	return r.UpdateWebhookWithLastPostedAtFunc(ctx, webhookSerialID, lastPostedAt)
}

func (r *RepositoryFuncMock) UpdateWebhookWithWebhookIDAndSubscription(ctx context.Context, webhookSerialID int64, webhookID string, subscriptionID string, subscriptionType string) error {
	return r.UpdateWebhookWithWebhookIDAndSubscriptionFunc(ctx, webhookSerialID, webhookID, subscriptionID, subscriptionType)
}

func (r *RepositoryFuncMock) DeleteWebhookByWebhookSerialID(ctx context.Context, webhookSerialID int64) error {
	return r.DeleteWebhookByWebhookSerialIDFunc(ctx, webhookSerialID)
}

func (r *RepositoryFuncMock) InsertWebhookWord(ctx context.Context, webhookSerialID int64, condition string, word string) error {
	return r.InsertWebhookWordFunc(ctx, webhookSerialID, condition, word)
}

func (r *RepositoryFuncMock) GetWebhookWordWithWebhookSerialIDAndCondition(ctx context.Context, webhookSerialID int64, condition string) ([]*WebhookWord, error) {
	return r.GetWebhookWordWithWebhookSerialIDAndConditionFunc(ctx, webhookSerialID, condition)
}

/*func (r *RepositoryFuncMock) GetWebhookWordWithWebhookSerialIDs(ctx context.Context, webhookSerialIDs []int64) ([]*WebhookWord, error) {
	return r.GetWebhookWordWithWebhookSerialIDsFunc(ctx, webhookSerialIDs)
}*/

func (r *RepositoryFuncMock) DeleteWebhookWordsNotInProvidedList(ctx context.Context, webhookSerialID int64, conditions string, words []string) error {
	return r.DeleteWebhookWordsNotInProvidedListFunc(ctx, webhookSerialID, conditions, words)
}

func (r *RepositoryFuncMock) InsertWebhookUserMention(ctx context.Context, webhookSerialID int64, userID string) error {
	return r.InsertWebhookUserMentionFunc(ctx, webhookSerialID, userID)
}

func (r *RepositoryFuncMock) GetWebhookUserMentionWithWebhookSerialID(ctx context.Context, webhookSerialID int64) ([]*WebhookUserMention, error) {
	return r.GetWebhookUserMentionWithWebhookSerialIDFunc(ctx, webhookSerialID)
}

/*func (r *RepositoryFuncMock) GetWebhookUserMentionWithWebhookSerialIDs(ctx context.Context, webhookSerialIDs []int64) ([]*WebhookUserMention, error) {
	return r.GetWebhookUserMentionWithWebhookSerialIDsFunc(ctx, webhookSerialIDs)
}*/

func (r *RepositoryFuncMock) DeleteWebhookUserMentionsNotInProvidedList(ctx context.Context, webhookSerialID int64, userIDs []string) error {
	return r.DeleteWebhookUserMentionsNotInProvidedListFunc(ctx, webhookSerialID, userIDs)
}

func (r *RepositoryFuncMock) InsertWebhookRoleMention(ctx context.Context, webhookSerialID int64, roleID string) error {
	return r.InsertWebhookRoleMentionFunc(ctx, webhookSerialID, roleID)
}

func (r *RepositoryFuncMock) GetWebhookRoleMentionWithWebhookSerialID(ctx context.Context, webhookSerialID int64) ([]*WebhookRoleMention, error) {
	return r.GetWebhookRoleMentionWithWebhookSerialIDFunc(ctx, webhookSerialID)
}

/*func (r *RepositoryFuncMock) GetWebhookRoleMentionWithWebhookSerialIDs(ctx context.Context, webhookSerialIDs []int64) ([]*WebhookRoleMention, error) {
	return r.GetWebhookRoleMentionWithWebhookSerialIDsFunc(ctx, webhookSerialIDs)
}*/

func (r *RepositoryFuncMock) DeleteWebhookRoleMentionsNotInProvidedList(ctx context.Context, webhookSerialID int64, roleIDs []string) error {
	return r.DeleteWebhookRoleMentionsNotInProvidedListFunc(ctx, webhookSerialID, roleIDs)
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
	GetVcSignalNgUserIDsByVcChannelID(ctx context.Context, vcChannelID string) ([]string, error)
	GetVcSignalNgRoleIDsByVcChannelID(ctx context.Context, vcChannelID string) ([]string, error)
	GetVcSignalChannelAllColumnByVcChannelID(ctx context.Context, vcChannelID string) (*VcSignalChannelAllColumn, error)
	GetVcSignalMentionUserIDsByVcChannelID(ctx context.Context, vcChannelID string) ([]string, error)
	GetVcSignalMentionRoleIDsByVcChannelID(ctx context.Context, vcChannelID string) ([]string, error)
	UpdateVcSignalChannel(ctx context.Context, vcChannel VcSignalChannelNotGuildID) error
	InsertVcSignalChannel(ctx context.Context, vcChannelID string, guildID, sendChannelID string) error
	InsertVcSignalNgUser(ctx context.Context, vcChannelID string, guildID string, userID string) error
	InsertVcSignalNgRole(ctx context.Context, vcChannelID string, guildID string, roleID string) error
	InsertVcSignalMentionUser(ctx context.Context, vcChannelID string, guildID string, userID string) error
	InsertVcSignalMentionRole(ctx context.Context, vcChannelID string, guildID string, roleID string) error
	DeleteVcSignalNgUsersNotInProvidedList(ctx context.Context, vcChannelID string, userIDs []string) error
	DeleteVcSignalNgRolesNotInProvidedList(ctx context.Context, vcChannelID string, roleIDs []string) error
	DeleteVcSignalMentionUsersNotInProvidedList(ctx context.Context, vcChannelID string, userIDs []string) error
	DeleteVcSignalMentionRolesNotInProvidedList(ctx context.Context, vcChannelID string, roleIDs []string) error
	InsertWebhook(ctx context.Context, guildID string, webhookID string, subscriptionType string, subscriptionID string, lastPostedAt time.Time) (int64, error)
	GetAllColumnsWebhooksByGuildID(ctx context.Context, guildID string) ([]*Webhook, error)
	UpdateWebhookWithLastPostedAt(ctx context.Context, webhookSerialID int64, lastPostedAt time.Time) error
	UpdateWebhookWithWebhookIDAndSubscription(ctx context.Context, webhookSerialID int64, webhookID string, subscriptionID string, subscriptionType string) error
	DeleteWebhookByWebhookSerialID(ctx context.Context, webhookSerialID int64) error
	InsertWebhookWord(ctx context.Context, webhookSerialID int64, condition string, word string) error
	GetWebhookWordWithWebhookSerialIDAndCondition(ctx context.Context, webhookSerialID int64, condition string) ([]*WebhookWord, error)
	//GetWebhookWordWithWebhookSerialIDs(ctx context.Context, webhookSerialIDs []int64) ([]*WebhookWord, error)
	DeleteWebhookWordsNotInProvidedList(ctx context.Context, webhookSerialID int64, conditions string, words []string) error
	InsertWebhookUserMention(ctx context.Context, webhookSerialID int64, userID string) error
	GetWebhookUserMentionWithWebhookSerialID(ctx context.Context, webhookSerialID int64) ([]*WebhookUserMention, error)
	//GetWebhookUserMentionWithWebhookSerialIDs(ctx context.Context, webhookSerialIDs []int64) ([]*WebhookUserMention, error)
	DeleteWebhookUserMentionsNotInProvidedList(ctx context.Context, webhookSerialID int64, userIDs []string) error
	InsertWebhookRoleMention(ctx context.Context, webhookSerialID int64, roleID string) error
	GetWebhookRoleMentionWithWebhookSerialID(ctx context.Context, webhookSerialID int64) ([]*WebhookRoleMention, error)
	//GetWebhookRoleMentionWithWebhookSerialIDs(ctx context.Context, webhookSerialIDs []int64) ([]*WebhookRoleMention, error)
	DeleteWebhookRoleMentionsNotInProvidedList(ctx context.Context, webhookSerialID int64, roleIDs []string) error
}

var (
	_ RepositoryFunc = (*Repository)(nil)
	_ RepositoryFunc = (*RepositoryFuncMock)(nil)
)
