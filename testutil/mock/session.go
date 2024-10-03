package mock

import (
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
)

type SessionMock struct {
	AddHandlerFunc                 func(handler interface{}) func()
	ApplicationCommandCreateFunc   func(appID string, guildID string, cmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (ccmd *discordgo.ApplicationCommand, err error)
	ApplicationCommandDeleteFunc   func(appID string, guildID string, cmdID string, options ...discordgo.RequestOption) error
	ChannelFunc                    func(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error)
	ChannelMessageSendFunc         func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelFileSendWithMessageFunc func(channelID string, content string, name string, r io.Reader, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelMessageSendEmbedFunc    func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelVoiceJoinFunc           func(gID string, cID string, mute bool, deaf bool) (voice *discordgo.VoiceConnection, err error)
	GuildFunc                      func(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error)
	GuildChannelsFunc              func(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Channel, err error)
	GuildThreadsActiveFunc func(guildID string, options ...discordgo.RequestOption) (threads *discordgo.ThreadsList, err error)
	GuildMemberFunc                func(guildID string, userID string, options ...discordgo.RequestOption) (st *discordgo.Member, err error)
	GuildMembersFunc               func(guildID string, after string, limit int, options ...discordgo.RequestOption) (st []*discordgo.Member, err error)
	GuildRolesFunc                 func(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Role, err error)
	GuildWebhooksFunc              func(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Webhook, err error)
	InteractionRespondFunc         func(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error
	UserChannelPermissionsFunc     func(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error)
	UserGuildsFunc                 func(limit int, beforeID string, afterID string, options ...discordgo.RequestOption) (st []*discordgo.UserGuild, err error)
	ThreadsActiveFunc              func(channelID string, options ...discordgo.RequestOption) (threads *discordgo.ThreadsList, err error)
	ThreadsArchivedFunc            func(channelID string, before *time.Time, limit int, options ...discordgo.RequestOption) (threads *discordgo.ThreadsList, err error)
	ThreadsPrivateArchivedFunc     func(channelID string, before *time.Time, limit int, options ...discordgo.RequestOption) (threads *discordgo.ThreadsList, err error)
	WebhookFunc                    func(webhookID string, options ...discordgo.RequestOption) (st *discordgo.Webhook, err error)
	WebhookExecuteFunc             func(webhookID string, token string, wait bool, data *discordgo.WebhookParams, options ...discordgo.RequestOption) (st *discordgo.Message, err error)
	WebhookThreadExecuteFunc       func(webhookID, token string, wait bool, threadID string, data *discordgo.WebhookParams, options ...discordgo.RequestOption) (st *discordgo.Message, err error)
}

func (s *SessionMock) AddHandler(handler interface{}) func() {
	return s.AddHandlerFunc(handler)
}

func (s *SessionMock) ApplicationCommandCreate(appID string, guildID string, cmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (ccmd *discordgo.ApplicationCommand, err error) {
	return s.ApplicationCommandCreateFunc(appID, guildID, cmd, options...)
}

func (s *SessionMock) ApplicationCommandDelete(appID string, guildID string, cmdID string, options ...discordgo.RequestOption) error {
	return s.ApplicationCommandDeleteFunc(appID, guildID, cmdID, options...)
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

func (s *SessionMock) ChannelMessageSendEmbed(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	return s.ChannelMessageSendEmbedFunc(channelID, embed, options...)
}

func (s *SessionMock) ChannelVoiceJoin(gID string, cID string, mute bool, deaf bool) (voice *discordgo.VoiceConnection, err error) {
	return s.ChannelVoiceJoinFunc(gID, cID, mute, deaf)
}

func (s *SessionMock) Guild(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error) {
	return s.GuildFunc(guildID, options...)
}

func (s *SessionMock) GuildChannels(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Channel, err error) {
	return s.GuildChannelsFunc(guildID, options...)
}

func (s *SessionMock) GuildThreadsActive(guildID string, options ...discordgo.RequestOption) (threads *discordgo.ThreadsList, err error) {
	return s.GuildThreadsActiveFunc(guildID, options...)
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

func (s *SessionMock) GuildWebhooks(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Webhook, err error) {
	return s.GuildWebhooksFunc(guildID, options...)
}

func (s *SessionMock) InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error {
	return nil
}

func (s *SessionMock) ThreadsActive(channelID string, options ...discordgo.RequestOption) (threads *discordgo.ThreadsList, err error) {
	return s.ThreadsActiveFunc(channelID, options...)
}

func (s *SessionMock) ThreadsArchived(channelID string, before *time.Time, limit int, options ...discordgo.RequestOption) (threads *discordgo.ThreadsList, err error) {
	return s.ThreadsArchivedFunc(channelID, before, limit, options...)
}

func (s *SessionMock) ThreadsPrivateArchived(channelID string, before *time.Time, limit int, options ...discordgo.RequestOption) (threads *discordgo.ThreadsList, err error) {
	return s.ThreadsPrivateArchivedFunc(channelID, before, limit, options...)
}

func (s *SessionMock) UserChannelPermissions(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error) {
	return s.UserChannelPermissionsFunc(userID, channelID, fetchOptions...)
}

func (s *SessionMock) UserGuilds(limit int, beforeID string, afterID string, withCounts bool, options ...discordgo.RequestOption) (st []*discordgo.UserGuild, err error) {
	return s.UserGuildsFunc(limit, beforeID, afterID, options...)
}

func (s *SessionMock) Webhook(webhookID string, options ...discordgo.RequestOption) (st *discordgo.Webhook, err error) {
	return s.WebhookFunc(webhookID, options...)
}

func (s *SessionMock) WebhookExecute(webhookID string, token string, wait bool, data *discordgo.WebhookParams, options ...discordgo.RequestOption) (st *discordgo.Message, err error) {
	return s.WebhookExecuteFunc(webhookID, token, wait, data, options...)
}

func (s *SessionMock) WebhookThreadExecute(webhookID, token string, wait bool, threadID string, data *discordgo.WebhookParams, options ...discordgo.RequestOption) (st *discordgo.Message, err error) {
	return s.WebhookThreadExecuteFunc(webhookID, token, wait, threadID, data, options...)
}

// Session is an interface for discordgo.Session.
type Session interface {
	AddHandler(handler interface{}) func()
	ApplicationCommandCreate(appID string, guildID string, cmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (ccmd *discordgo.ApplicationCommand, err error)
	ApplicationCommandDelete(appID string, guildID string, cmdID string, options ...discordgo.RequestOption) error
	Channel(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error)
	ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelFileSendWithMessage(channelID string, content string, name string, r io.Reader, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelMessageSendEmbed(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelVoiceJoin(gID string, cID string, mute bool, deaf bool) (voice *discordgo.VoiceConnection, err error)
	Guild(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error)
	GuildChannels(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Channel, err error)
	GuildThreadsActive(guildID string, options ...discordgo.RequestOption) (threads *discordgo.ThreadsList, err error)
	GuildMember(guildID string, userID string, options ...discordgo.RequestOption) (st *discordgo.Member, err error)
	GuildMembers(guildID string, after string, limit int, options ...discordgo.RequestOption) (st []*discordgo.Member, err error)
	GuildRoles(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Role, err error)
	GuildWebhooks(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Webhook, err error)
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error
	UserChannelPermissions(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error)
	UserGuilds(limit int, beforeID string, afterID string, withCounts bool, options ...discordgo.RequestOption) (st []*discordgo.UserGuild, err error)
	ThreadsActive(channelID string, options ...discordgo.RequestOption) (threads *discordgo.ThreadsList, err error)
	ThreadsArchived(channelID string, before *time.Time, limit int, options ...discordgo.RequestOption) (threads *discordgo.ThreadsList, err error)
	ThreadsPrivateArchived(channelID string, before *time.Time, limit int, options ...discordgo.RequestOption) (threads *discordgo.ThreadsList, err error)
	Webhook(webhookID string, options ...discordgo.RequestOption) (st *discordgo.Webhook, err error)
	WebhookExecute(webhookID string, token string, wait bool, data *discordgo.WebhookParams, options ...discordgo.RequestOption) (st *discordgo.Message, err error)
	WebhookThreadExecute(webhookID, token string, wait bool, threadID string, data *discordgo.WebhookParams, options ...discordgo.RequestOption) (st *discordgo.Message, err error)
}

var (
	_ Session = (*discordgo.Session)(nil)
	_ Session = (*SessionMock)(nil)
)
