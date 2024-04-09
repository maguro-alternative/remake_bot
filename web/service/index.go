package service

import (
	"net/http"

	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/sessions"
)

// A TODOService implements CRUD of TODO entities.
type IndexService struct {
	Client          *http.Client
	CookieStore     *sessions.CookieStore
	DiscordSession  mock.Session
	DiscordBotState *discordgo.State
}

// NewTODOService returns new TODOService.
func NewIndexService(
	client *http.Client,
	cookieStore *sessions.CookieStore,
	discordSession mock.Session,
	discordBotState *discordgo.State,
) *IndexService {
	return &IndexService{
		Client:          client,
		CookieStore:     cookieStore,
		DiscordSession:  discordSession,
		DiscordBotState: discordBotState,
	}
}
