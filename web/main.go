package web

import (
	"net/http"
	"strings"

	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/web/middleware"
	"github.com/maguro-alternative/remake_bot/web/config"

	linePostDiscordChannel "github.com/maguro-alternative/remake_bot/web/handler/api/line_post_discord_channel"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linebot"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linetoken"
	discordLogin "github.com/maguro-alternative/remake_bot/web/handler/login/discord_login"
	discordLogout "github.com/maguro-alternative/remake_bot/web/handler/logout/discord_logout"
	discordCallback "github.com/maguro-alternative/remake_bot/web/handler/callback/discord_callback"

	indexView "github.com/maguro-alternative/remake_bot/web/handler/views"
	linePostDiscordChannelView "github.com/maguro-alternative/remake_bot/web/handler/views/guildid/line_post_discord_channel"
	linetokenView "github.com/maguro-alternative/remake_bot/web/handler/views/guildid/linetoken"
	guildsView "github.com/maguro-alternative/remake_bot/web/handler/views/guilds"
	guildIdView "github.com/maguro-alternative/remake_bot/web/handler/views/guildid"
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
		RedirectURL: config.ServerUrl() + "/callback/discord-callback/",
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
	// 静的ファイルのハンドリング
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/templates/static/"))))

	mux.Handle("/", middleChain.ThenFunc(indexView.NewIndexViewHandler(indexService).Index))
	mux.Handle("/guilds", middleChain.ThenFunc(guildsView.NewGuildsViewHandler(indexService).Index))
	mux.Handle("/guild/{guildId}", middleChain.ThenFunc(guildIdView.NewGuildIDViewHandler(indexService).Index))
	mux.Handle("/guild/{guildId}/linetoken", middleChain.ThenFunc(linetokenView.NewLineTokenViewHandler(indexService).Index))
	mux.Handle("/guild/{guildId}/line-post-discord-channel", middleChain.ThenFunc(linePostDiscordChannelView.NewLinePostDiscordChannelViewHandler(indexService).Index))

	mux.Handle("/api/line-bot", middleChain.Then(linebot.NewLineBotHandler(indexService)))
	mux.Handle("/login/discord", middleChain.Then(discordLogin.NewDiscordOAuth2Handler(discordOAuth2Service)))
	mux.Handle("/logout/discord", middleChain.Then(discordLogout.NewDiscordOAuth2Handler(discordOAuth2Service)))
	mux.Handle("/callback/discord-callback/", middleChain.Then(discordCallback.NewDiscordCallbackHandler(discordOAuth2Service)))
	mux.Handle("/api/{guildId}/linetoken", middleChain.Then(linetoken.NewLineTokenHandler(indexService)))
	mux.Handle("/api/{guildId}/line-post-discord-channel", middleChain.Then(linePostDiscordChannel.NewLineChannelHandler(indexService)))

	http.ListenAndServe(":"+config.Port(), mux)
}
