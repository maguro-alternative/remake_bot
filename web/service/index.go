package service

import (
	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/sessions"
)

// A TODOService implements CRUD of TODO entities.
type IndexService struct {
	DB             db.Driver
	CookieStore    *sessions.CookieStore
	DiscordSession *discordgo.Session
}

// NewTODOService returns new TODOService.
func NewIndexService(
	db db.Driver,
	cookieStore *sessions.CookieStore,
	discordSession *discordgo.Session,
) *IndexService {
	return &IndexService{
		DB:             db,
		CookieStore:    cookieStore,
		DiscordSession: discordSession,
	}
}
