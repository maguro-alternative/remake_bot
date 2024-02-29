package web

import (
	"net/http"
	"strings"

	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/pkg/middleware"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linebot"
	discordOAuth "github.com/maguro-alternative/remake_bot/web/handler/auth/discordOAuth"
	discordCallback "github.com/maguro-alternative/remake_bot/web/handler/callback/discordCallback"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linetoken"
	"github.com/maguro-alternative/remake_bot/web/handler/guilds/guildid"
	"github.com/maguro-alternative/remake_bot/web/service"

	"golang.org/x/oauth2"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
)

func NewWebRouter(
	indexDB db.Driver,
	cookieStore *sessions.CookieStore,
	discordSession *discordgo.Session,
) {
	scopes := config.DiscordScopes()
	conf := &oauth2.Config{
		ClientID:     config.DiscordClientID(),
		ClientSecret: config.DiscordClientSecret(),
		Scopes:       strings.Split(scopes, "%20"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		RedirectURL: "/discord-callback/",
	}
	// create a *service.TODOService type variable using the *sql.DB type variable
	var indexService = service.NewIndexService(
		indexDB,
		cookieStore,
		discordSession,
	)
	var discordOAuth2Service = service.NewDiscordOAuth2Service(
		conf,
		cookieStore,
		discordSession,
	)

	// register routes
	mux := http.NewServeMux()
	middleChain := alice.New(middleware.LogMiddleware)
	mux.Handle("/api/line-bot", middleChain.Then(linebot.NewLineBotHandler(indexService)))
	//mux.Handle("/discord-auth-check", middleChain.Then(testRouter.NewAuthCheckHandler(indexService)))
	mux.Handle("/auth/discord", middleChain.Then(discordOAuth.NewDiscordOAuth2Handler(discordOAuth2Service)))
	mux.Handle("/callback/discord-callback/", middleChain.Then(discordCallback.NewDiscordCallbackHandler(discordOAuth2Service)))
	mux.Handle("/api/{guildId}/linetoken", middleChain.Then(linetoken.NewLineTokenHandler(indexService)))
	mux.Handle("/guild/{guildId}/linetoken/", middleChain.ThenFunc(guildid.NewGuildIdHandler(indexService).Index))
	http.ListenAndServe(":8080", mux)
}
