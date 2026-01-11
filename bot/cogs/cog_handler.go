package cogs

import (
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/pkg/lineworks_service"
)

type cogHandler struct {
	db               db.Driver
	client           *http.Client
	lineWorksService *lineworks_service.LineWorksService
}

func newCogHandler(
	db db.Driver,
	client *http.Client,
	lineWorksService *lineworks_service.LineWorksService,
) *cogHandler {
	return &cogHandler{
		db:               db,
		client:           client,
		lineWorksService: lineWorksService,
	}
}

func RegisterHandlers(
	s *discordgo.Session,
	sqlxdb db.Driver,
	client *http.Client,
	lineWorksService *lineworks_service.LineWorksService,
) {
	cogs := newCogHandler(sqlxdb, client, lineWorksService)
	fmt.Println(s.State.User.Username + "としてログインしました")
	s.AddHandler(cogs.onVoiceStateUpdate)
	s.AddHandler(cogs.onMessageCreate)
}
