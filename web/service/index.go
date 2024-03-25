package service

import (
	"io"

	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/sessions"
)

// A TODOService implements CRUD of TODO entities.
type IndexService struct {
	DB              db.Driver
	CookieStore     *sessions.CookieStore
	DiscordSession  Session
	DiscordBotState *discordgo.State
}

// Session is an interface for discordgo.Session.
type Session interface {
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
)

// NewTODOService returns new TODOService.
func NewIndexService(
	db db.Driver,
	cookieStore *sessions.CookieStore,
	discordSession Session,
	discordBotState *discordgo.State,
) *IndexService {
	return &IndexService{
		DB:              db,
		CookieStore:     cookieStore,
		DiscordSession:  discordSession,
		DiscordBotState: discordBotState,
	}
}
