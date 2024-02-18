package discordoauth

import (
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/service"
)

type DiscordOAuth2Handler struct {
	DiscordOAuth2Service *service.DiscordOAuth2Service
}

func NewDiscordOAuth2Handler(discordOAuth2Service *service.DiscordOAuth2Service) *DiscordOAuth2Handler {
	return &DiscordOAuth2Handler{
		DiscordOAuth2Service: discordOAuth2Service,
	}
}

func (h *DiscordOAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}