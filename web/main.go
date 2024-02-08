package web

import (
	"net/http"

	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/pkg/middleware"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linebot"

	//"golang.org/x/oauth2"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
)

func NewWebRouter(
	indexDB db.Driver,
	cookieStore *sessions.CookieStore,
	discordSession *discordgo.Session,
) {
	//conf := &oauth2.Config{
		//ClientID:     "",
		//ClientSecret: "",
		//Scopes:       []string{"identify"},
		//Endpoint: oauth2.Endpoint{
			//AuthURL:  "https://discord.com/api/oauth2/authorize",
			//TokenURL: "https://discord.com/api/oauth2/token",
		//},
		//RedirectURL: "/discord-callback/",
	//}
	// create a *service.TODOService type variable using the *sql.DB type variable
	var indexService = service.NewIndexService(
		indexDB,
		cookieStore,
		discordSession,
	)
	//var discordOAuth2Service = service.NewDiscordOAuth2Service(
		//conf,
		//cookieStore,
		//discordSession,
	//)

	// register routes
	mux := http.NewServeMux()
	middleChain := alice.New(middleware.LogMiddleware)
	mux.Handle("/api/line-bot", middleChain.Then(linebot.NewLineBotHandler(indexService)))
	//mux.Handle("/discord-auth-check", middleChain.Then(testRouter.NewAuthCheckHandler(indexService)))
	//mux.Handle("/discord/auth", middleChain.Then(controllersDiscord.NewDiscordAuthHandler(discordOAuth2Service)))
	//mux.Handle("/discord-callback/", middleChain.Then(controllersDiscord.NewDiscordCallbackHandler(discordOAuth2Service)))
	http.ListenAndServe(":8080", mux)
}