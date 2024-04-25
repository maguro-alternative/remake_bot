package cogs

import (
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type cogHandler struct {
	db     db.Driver
	client *http.Client
}

func newCogHandler(
	db db.Driver,
	client *http.Client,
) *cogHandler {
	return &cogHandler{
		db:     db,
		client: client,
	}
}

func RegisterHandlers(
	s *discordgo.Session,
	sqlxdb db.Driver,
	client *http.Client,
) {
	cogs := newCogHandler(sqlxdb, client)
	fmt.Println(s.State.User.Username + "としてログインしました")
	//s.AddHandler(cogs.onVoiceStateUpdate)
	s.AddHandler(cogs.onMessageCreate)
}
