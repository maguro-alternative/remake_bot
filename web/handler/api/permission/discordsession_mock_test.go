// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package permission

import (
	"github.com/bwmarrin/discordgo"
	"io"
	"sync"
)

// Ensure, that SessionMock does implement Session.
// If this is not the case, regenerate this file with moq.
var _ Session = &SessionMock{}

// SessionMock is a mock implementation of Session.
//
//	func TestSomethingThatUsesSession(t *testing.T) {
//
//		// make and configure a mocked Session
//		mockedSession := &SessionMock{
//			ChannelFileSendWithMessageFunc: func(channelID string, content string, name string, r io.Reader, options ...discordgo.RequestOption) (*discordgo.Message, error) {
//				panic("mock out the ChannelFileSendWithMessage method")
//			},
//			ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
//				panic("mock out the ChannelMessageSend method")
//			},
//			GuildFunc: func(guildID string, options ...discordgo.RequestOption) (*discordgo.Guild, error) {
//				panic("mock out the Guild method")
//			},
//			GuildChannelsFunc: func(guildID string, options ...discordgo.RequestOption) ([]*discordgo.Channel, error) {
//				panic("mock out the GuildChannels method")
//			},
//			GuildMemberFunc: func(guildID string, userID string, options ...discordgo.RequestOption) (*discordgo.Member, error) {
//				panic("mock out the GuildMember method")
//			},
//			GuildMembersFunc: func(guildID string, after string, limit int, options ...discordgo.RequestOption) ([]*discordgo.Member, error) {
//				panic("mock out the GuildMembers method")
//			},
//			GuildRolesFunc: func(guildID string, options ...discordgo.RequestOption) ([]*discordgo.Role, error) {
//				panic("mock out the GuildRoles method")
//			},
//			UserChannelPermissionsFunc: func(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (int64, error) {
//				panic("mock out the UserChannelPermissions method")
//			},
//			UserGuildsFunc: func(limit int, beforeID string, afterID string, options ...discordgo.RequestOption) ([]*discordgo.UserGuild, error) {
//				panic("mock out the UserGuilds method")
//			},
//		}
//
//		// use mockedSession in code that requires Session
//		// and then make assertions.
//
//	}
type SessionMock struct {
	// ChannelFileSendWithMessageFunc mocks the ChannelFileSendWithMessage method.
	ChannelFileSendWithMessageFunc func(channelID string, content string, name string, r io.Reader, options ...discordgo.RequestOption) (*discordgo.Message, error)

	// ChannelMessageSendFunc mocks the ChannelMessageSend method.
	ChannelMessageSendFunc func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)

	// GuildFunc mocks the Guild method.
	GuildFunc func(guildID string, options ...discordgo.RequestOption) (*discordgo.Guild, error)

	// GuildChannelsFunc mocks the GuildChannels method.
	GuildChannelsFunc func(guildID string, options ...discordgo.RequestOption) ([]*discordgo.Channel, error)

	// GuildMemberFunc mocks the GuildMember method.
	GuildMemberFunc func(guildID string, userID string, options ...discordgo.RequestOption) (*discordgo.Member, error)

	// GuildMembersFunc mocks the GuildMembers method.
	GuildMembersFunc func(guildID string, after string, limit int, options ...discordgo.RequestOption) ([]*discordgo.Member, error)

	// GuildRolesFunc mocks the GuildRoles method.
	GuildRolesFunc func(guildID string, options ...discordgo.RequestOption) ([]*discordgo.Role, error)

	// UserChannelPermissionsFunc mocks the UserChannelPermissions method.
	UserChannelPermissionsFunc func(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (int64, error)

	// UserGuildsFunc mocks the UserGuilds method.
	UserGuildsFunc func(limit int, beforeID string, afterID string, options ...discordgo.RequestOption) ([]*discordgo.UserGuild, error)

	// calls tracks calls to the methods.
	calls struct {
		// ChannelFileSendWithMessage holds details about calls to the ChannelFileSendWithMessage method.
		ChannelFileSendWithMessage []struct {
			// ChannelID is the channelID argument value.
			ChannelID string
			// Content is the content argument value.
			Content string
			// Name is the name argument value.
			Name string
			// R is the r argument value.
			R io.Reader
			// Options is the options argument value.
			Options []discordgo.RequestOption
		}
		// ChannelMessageSend holds details about calls to the ChannelMessageSend method.
		ChannelMessageSend []struct {
			// ChannelID is the channelID argument value.
			ChannelID string
			// Content is the content argument value.
			Content string
			// Options is the options argument value.
			Options []discordgo.RequestOption
		}
		// Guild holds details about calls to the Guild method.
		Guild []struct {
			// GuildID is the guildID argument value.
			GuildID string
			// Options is the options argument value.
			Options []discordgo.RequestOption
		}
		// GuildChannels holds details about calls to the GuildChannels method.
		GuildChannels []struct {
			// GuildID is the guildID argument value.
			GuildID string
			// Options is the options argument value.
			Options []discordgo.RequestOption
		}
		// GuildMember holds details about calls to the GuildMember method.
		GuildMember []struct {
			// GuildID is the guildID argument value.
			GuildID string
			// UserID is the userID argument value.
			UserID string
			// Options is the options argument value.
			Options []discordgo.RequestOption
		}
		// GuildMembers holds details about calls to the GuildMembers method.
		GuildMembers []struct {
			// GuildID is the guildID argument value.
			GuildID string
			// After is the after argument value.
			After string
			// Limit is the limit argument value.
			Limit int
			// Options is the options argument value.
			Options []discordgo.RequestOption
		}
		// GuildRoles holds details about calls to the GuildRoles method.
		GuildRoles []struct {
			// GuildID is the guildID argument value.
			GuildID string
			// Options is the options argument value.
			Options []discordgo.RequestOption
		}
		// UserChannelPermissions holds details about calls to the UserChannelPermissions method.
		UserChannelPermissions []struct {
			// UserID is the userID argument value.
			UserID string
			// ChannelID is the channelID argument value.
			ChannelID string
			// FetchOptions is the fetchOptions argument value.
			FetchOptions []discordgo.RequestOption
		}
		// UserGuilds holds details about calls to the UserGuilds method.
		UserGuilds []struct {
			// Limit is the limit argument value.
			Limit int
			// BeforeID is the beforeID argument value.
			BeforeID string
			// AfterID is the afterID argument value.
			AfterID string
			// Options is the options argument value.
			Options []discordgo.RequestOption
		}
	}
	lockChannelFileSendWithMessage sync.RWMutex
	lockChannelMessageSend         sync.RWMutex
	lockGuild                      sync.RWMutex
	lockGuildChannels              sync.RWMutex
	lockGuildMember                sync.RWMutex
	lockGuildMembers               sync.RWMutex
	lockGuildRoles                 sync.RWMutex
	lockUserChannelPermissions     sync.RWMutex
	lockUserGuilds                 sync.RWMutex
}

// ChannelFileSendWithMessage calls ChannelFileSendWithMessageFunc.
func (mock *SessionMock) ChannelFileSendWithMessage(channelID string, content string, name string, r io.Reader, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	if mock.ChannelFileSendWithMessageFunc == nil {
		panic("SessionMock.ChannelFileSendWithMessageFunc: method is nil but Session.ChannelFileSendWithMessage was just called")
	}
	callInfo := struct {
		ChannelID string
		Content   string
		Name      string
		R         io.Reader
		Options   []discordgo.RequestOption
	}{
		ChannelID: channelID,
		Content:   content,
		Name:      name,
		R:         r,
		Options:   options,
	}
	mock.lockChannelFileSendWithMessage.Lock()
	mock.calls.ChannelFileSendWithMessage = append(mock.calls.ChannelFileSendWithMessage, callInfo)
	mock.lockChannelFileSendWithMessage.Unlock()
	return mock.ChannelFileSendWithMessageFunc(channelID, content, name, r, options...)
}

// ChannelFileSendWithMessageCalls gets all the calls that were made to ChannelFileSendWithMessage.
// Check the length with:
//
//	len(mockedSession.ChannelFileSendWithMessageCalls())
func (mock *SessionMock) ChannelFileSendWithMessageCalls() []struct {
	ChannelID string
	Content   string
	Name      string
	R         io.Reader
	Options   []discordgo.RequestOption
} {
	var calls []struct {
		ChannelID string
		Content   string
		Name      string
		R         io.Reader
		Options   []discordgo.RequestOption
	}
	mock.lockChannelFileSendWithMessage.RLock()
	calls = mock.calls.ChannelFileSendWithMessage
	mock.lockChannelFileSendWithMessage.RUnlock()
	return calls
}

// ChannelMessageSend calls ChannelMessageSendFunc.
func (mock *SessionMock) ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	if mock.ChannelMessageSendFunc == nil {
		panic("SessionMock.ChannelMessageSendFunc: method is nil but Session.ChannelMessageSend was just called")
	}
	callInfo := struct {
		ChannelID string
		Content   string
		Options   []discordgo.RequestOption
	}{
		ChannelID: channelID,
		Content:   content,
		Options:   options,
	}
	mock.lockChannelMessageSend.Lock()
	mock.calls.ChannelMessageSend = append(mock.calls.ChannelMessageSend, callInfo)
	mock.lockChannelMessageSend.Unlock()
	return mock.ChannelMessageSendFunc(channelID, content, options...)
}

// ChannelMessageSendCalls gets all the calls that were made to ChannelMessageSend.
// Check the length with:
//
//	len(mockedSession.ChannelMessageSendCalls())
func (mock *SessionMock) ChannelMessageSendCalls() []struct {
	ChannelID string
	Content   string
	Options   []discordgo.RequestOption
} {
	var calls []struct {
		ChannelID string
		Content   string
		Options   []discordgo.RequestOption
	}
	mock.lockChannelMessageSend.RLock()
	calls = mock.calls.ChannelMessageSend
	mock.lockChannelMessageSend.RUnlock()
	return calls
}

// Guild calls GuildFunc.
func (mock *SessionMock) Guild(guildID string, options ...discordgo.RequestOption) (*discordgo.Guild, error) {
	if mock.GuildFunc == nil {
		panic("SessionMock.GuildFunc: method is nil but Session.Guild was just called")
	}
	callInfo := struct {
		GuildID string
		Options []discordgo.RequestOption
	}{
		GuildID: guildID,
		Options: options,
	}
	mock.lockGuild.Lock()
	mock.calls.Guild = append(mock.calls.Guild, callInfo)
	mock.lockGuild.Unlock()
	return mock.GuildFunc(guildID, options...)
}

// GuildCalls gets all the calls that were made to Guild.
// Check the length with:
//
//	len(mockedSession.GuildCalls())
func (mock *SessionMock) GuildCalls() []struct {
	GuildID string
	Options []discordgo.RequestOption
} {
	var calls []struct {
		GuildID string
		Options []discordgo.RequestOption
	}
	mock.lockGuild.RLock()
	calls = mock.calls.Guild
	mock.lockGuild.RUnlock()
	return calls
}

// GuildChannels calls GuildChannelsFunc.
func (mock *SessionMock) GuildChannels(guildID string, options ...discordgo.RequestOption) ([]*discordgo.Channel, error) {
	if mock.GuildChannelsFunc == nil {
		panic("SessionMock.GuildChannelsFunc: method is nil but Session.GuildChannels was just called")
	}
	callInfo := struct {
		GuildID string
		Options []discordgo.RequestOption
	}{
		GuildID: guildID,
		Options: options,
	}
	mock.lockGuildChannels.Lock()
	mock.calls.GuildChannels = append(mock.calls.GuildChannels, callInfo)
	mock.lockGuildChannels.Unlock()
	return mock.GuildChannelsFunc(guildID, options...)
}

// GuildChannelsCalls gets all the calls that were made to GuildChannels.
// Check the length with:
//
//	len(mockedSession.GuildChannelsCalls())
func (mock *SessionMock) GuildChannelsCalls() []struct {
	GuildID string
	Options []discordgo.RequestOption
} {
	var calls []struct {
		GuildID string
		Options []discordgo.RequestOption
	}
	mock.lockGuildChannels.RLock()
	calls = mock.calls.GuildChannels
	mock.lockGuildChannels.RUnlock()
	return calls
}

// GuildMember calls GuildMemberFunc.
func (mock *SessionMock) GuildMember(guildID string, userID string, options ...discordgo.RequestOption) (*discordgo.Member, error) {
	if mock.GuildMemberFunc == nil {
		panic("SessionMock.GuildMemberFunc: method is nil but Session.GuildMember was just called")
	}
	callInfo := struct {
		GuildID string
		UserID  string
		Options []discordgo.RequestOption
	}{
		GuildID: guildID,
		UserID:  userID,
		Options: options,
	}
	mock.lockGuildMember.Lock()
	mock.calls.GuildMember = append(mock.calls.GuildMember, callInfo)
	mock.lockGuildMember.Unlock()
	return mock.GuildMemberFunc(guildID, userID, options...)
}

// GuildMemberCalls gets all the calls that were made to GuildMember.
// Check the length with:
//
//	len(mockedSession.GuildMemberCalls())
func (mock *SessionMock) GuildMemberCalls() []struct {
	GuildID string
	UserID  string
	Options []discordgo.RequestOption
} {
	var calls []struct {
		GuildID string
		UserID  string
		Options []discordgo.RequestOption
	}
	mock.lockGuildMember.RLock()
	calls = mock.calls.GuildMember
	mock.lockGuildMember.RUnlock()
	return calls
}

// GuildMembers calls GuildMembersFunc.
func (mock *SessionMock) GuildMembers(guildID string, after string, limit int, options ...discordgo.RequestOption) ([]*discordgo.Member, error) {
	if mock.GuildMembersFunc == nil {
		panic("SessionMock.GuildMembersFunc: method is nil but Session.GuildMembers was just called")
	}
	callInfo := struct {
		GuildID string
		After   string
		Limit   int
		Options []discordgo.RequestOption
	}{
		GuildID: guildID,
		After:   after,
		Limit:   limit,
		Options: options,
	}
	mock.lockGuildMembers.Lock()
	mock.calls.GuildMembers = append(mock.calls.GuildMembers, callInfo)
	mock.lockGuildMembers.Unlock()
	return mock.GuildMembersFunc(guildID, after, limit, options...)
}

// GuildMembersCalls gets all the calls that were made to GuildMembers.
// Check the length with:
//
//	len(mockedSession.GuildMembersCalls())
func (mock *SessionMock) GuildMembersCalls() []struct {
	GuildID string
	After   string
	Limit   int
	Options []discordgo.RequestOption
} {
	var calls []struct {
		GuildID string
		After   string
		Limit   int
		Options []discordgo.RequestOption
	}
	mock.lockGuildMembers.RLock()
	calls = mock.calls.GuildMembers
	mock.lockGuildMembers.RUnlock()
	return calls
}

// GuildRoles calls GuildRolesFunc.
func (mock *SessionMock) GuildRoles(guildID string, options ...discordgo.RequestOption) ([]*discordgo.Role, error) {
	if mock.GuildRolesFunc == nil {
		panic("SessionMock.GuildRolesFunc: method is nil but Session.GuildRoles was just called")
	}
	callInfo := struct {
		GuildID string
		Options []discordgo.RequestOption
	}{
		GuildID: guildID,
		Options: options,
	}
	mock.lockGuildRoles.Lock()
	mock.calls.GuildRoles = append(mock.calls.GuildRoles, callInfo)
	mock.lockGuildRoles.Unlock()
	return mock.GuildRolesFunc(guildID, options...)
}

// GuildRolesCalls gets all the calls that were made to GuildRoles.
// Check the length with:
//
//	len(mockedSession.GuildRolesCalls())
func (mock *SessionMock) GuildRolesCalls() []struct {
	GuildID string
	Options []discordgo.RequestOption
} {
	var calls []struct {
		GuildID string
		Options []discordgo.RequestOption
	}
	mock.lockGuildRoles.RLock()
	calls = mock.calls.GuildRoles
	mock.lockGuildRoles.RUnlock()
	return calls
}

// UserChannelPermissions calls UserChannelPermissionsFunc.
func (mock *SessionMock) UserChannelPermissions(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (int64, error) {
	if mock.UserChannelPermissionsFunc == nil {
		panic("SessionMock.UserChannelPermissionsFunc: method is nil but Session.UserChannelPermissions was just called")
	}
	callInfo := struct {
		UserID       string
		ChannelID    string
		FetchOptions []discordgo.RequestOption
	}{
		UserID:       userID,
		ChannelID:    channelID,
		FetchOptions: fetchOptions,
	}
	mock.lockUserChannelPermissions.Lock()
	mock.calls.UserChannelPermissions = append(mock.calls.UserChannelPermissions, callInfo)
	mock.lockUserChannelPermissions.Unlock()
	return mock.UserChannelPermissionsFunc(userID, channelID, fetchOptions...)
}

// UserChannelPermissionsCalls gets all the calls that were made to UserChannelPermissions.
// Check the length with:
//
//	len(mockedSession.UserChannelPermissionsCalls())
func (mock *SessionMock) UserChannelPermissionsCalls() []struct {
	UserID       string
	ChannelID    string
	FetchOptions []discordgo.RequestOption
} {
	var calls []struct {
		UserID       string
		ChannelID    string
		FetchOptions []discordgo.RequestOption
	}
	mock.lockUserChannelPermissions.RLock()
	calls = mock.calls.UserChannelPermissions
	mock.lockUserChannelPermissions.RUnlock()
	return calls
}

// UserGuilds calls UserGuildsFunc.
func (mock *SessionMock) UserGuilds(limit int, beforeID string, afterID string, options ...discordgo.RequestOption) ([]*discordgo.UserGuild, error) {
	if mock.UserGuildsFunc == nil {
		panic("SessionMock.UserGuildsFunc: method is nil but Session.UserGuilds was just called")
	}
	callInfo := struct {
		Limit    int
		BeforeID string
		AfterID  string
		Options  []discordgo.RequestOption
	}{
		Limit:    limit,
		BeforeID: beforeID,
		AfterID:  afterID,
		Options:  options,
	}
	mock.lockUserGuilds.Lock()
	mock.calls.UserGuilds = append(mock.calls.UserGuilds, callInfo)
	mock.lockUserGuilds.Unlock()
	return mock.UserGuildsFunc(limit, beforeID, afterID, options...)
}

// UserGuildsCalls gets all the calls that were made to UserGuilds.
// Check the length with:
//
//	len(mockedSession.UserGuildsCalls())
func (mock *SessionMock) UserGuildsCalls() []struct {
	Limit    int
	BeforeID string
	AfterID  string
	Options  []discordgo.RequestOption
} {
	var calls []struct {
		Limit    int
		BeforeID string
		AfterID  string
		Options  []discordgo.RequestOption
	}
	mock.lockUserGuilds.RLock()
	calls = mock.calls.UserGuilds
	mock.lockUserGuilds.RUnlock()
	return calls
}