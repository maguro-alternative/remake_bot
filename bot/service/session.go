package service

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

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
	Channel(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error)
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
	_ Session = (*discordgo.Session)(nil)
	_ Session = (*SessionMock)(nil)
)