package service

import (
	"golang.org/x/oauth2"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/sessions"
)

type DiscordOAuth2Service struct {
	OAuth2Conf     oauth2.Config
	CookieStore    *sessions.CookieStore
	DiscordSession *discordgo.Session
}

func NewDiscordOAuth2Service(
	conf *oauth2.Config,
	cookieStore *sessions.CookieStore,
	discordSession *discordgo.Session,
) *DiscordOAuth2Service {
	return &DiscordOAuth2Service{
		OAuth2Conf:     *conf,
		CookieStore:    cookieStore,
		DiscordSession: discordSession,
	}
}
