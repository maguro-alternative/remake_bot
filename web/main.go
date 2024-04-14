package web

import (
	"encoding/gob"
	"net/http"
	"strings"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/middleware"

	"github.com/maguro-alternative/remake_bot/web/handler/api/group"
	linePostDiscordChannel "github.com/maguro-alternative/remake_bot/web/handler/api/line_post_discord_channel"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linebot"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linetoken"
	"github.com/maguro-alternative/remake_bot/web/handler/api/permission"
	discordCallback "github.com/maguro-alternative/remake_bot/web/handler/callback/discord_callback"
	lineCallback "github.com/maguro-alternative/remake_bot/web/handler/callback/line_callback"
	discordLogin "github.com/maguro-alternative/remake_bot/web/handler/login/discord_login"
	lineLogin "github.com/maguro-alternative/remake_bot/web/handler/login/line_login"
	discordLogout "github.com/maguro-alternative/remake_bot/web/handler/logout/discord_logout"
	lineLogout "github.com/maguro-alternative/remake_bot/web/handler/logout/line_logout"

	indexView "github.com/maguro-alternative/remake_bot/web/handler/views"
	groupView "github.com/maguro-alternative/remake_bot/web/handler/views/group"
	guildIdView "github.com/maguro-alternative/remake_bot/web/handler/views/guildid"
	linePostDiscordChannelView "github.com/maguro-alternative/remake_bot/web/handler/views/guildid/line_post_discord_channel"
	linetokenView "github.com/maguro-alternative/remake_bot/web/handler/views/guildid/linetoken"
	permissionView "github.com/maguro-alternative/remake_bot/web/handler/views/guildid/permission"
	guildsView "github.com/maguro-alternative/remake_bot/web/handler/views/guilds"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"

	"golang.org/x/oauth2"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
)

func init() {
	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	gob.Register(&model.DiscordUser{})
	gob.Register(&model.LineIdTokenUser{})
	gob.Register(&model.LineOAuthSession{})
}

func NewWebRouter(
	indexDB db.Driver,
	client *http.Client,
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

	repo := repository.NewRepository(indexDB)

	// セッションストアを作成します。
	cookieStore := sessions.NewCookieStore([]byte(config.SessionSecret()))

	// create a *service.TODOService type variable using the *sql.DB type variable
	var indexService = service.NewIndexService(
		client,
		cookieStore,
		discordSession,
		discordSession.State,
	)
	var discordOAuth2Service = service.NewDiscordOAuth2Service(
		conf,
		cookieStore,
		discordSession,
	)

	// register routes
	mux := http.NewServeMux()
	middleChain := alice.New(middleware.LogMiddleware)
	discordMiddleChain := alice.New(middleware.DiscordOAuthCheckMiddleware(*indexService, repo, true), middleware.LogMiddleware)
	lineMiddleChain := alice.New(middleware.LineOAuthCheckMiddleware(*indexService, repo, true), middleware.LogMiddleware)
	loginRequiredChain := alice.New(
		middleware.DiscordOAuthCheckMiddleware(*indexService, repo, false),
		middleware.LineOAuthCheckMiddleware(*indexService, repo, false),
		middleware.LogMiddleware,
	)
	discordLoginRequiredChain := alice.New(
		middleware.DiscordOAuthCheckMiddleware(*indexService, repo, true),
		middleware.LineOAuthCheckMiddleware(*indexService, repo, false),
		middleware.LogMiddleware,
	)
	lineLoginRequiredChain := alice.New(
		middleware.DiscordOAuthCheckMiddleware(*indexService, repo, false),
		middleware.LineOAuthCheckMiddleware(*indexService, repo, true),
		middleware.LogMiddleware,
	)

	// 静的ファイルのハンドリング
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/templates/static/"))))

	mux.Handle("/", loginRequiredChain.ThenFunc(indexView.NewIndexViewHandler(indexService).Index))
	mux.Handle("/login/line", loginRequiredChain.ThenFunc(lineLogin.NewLineLoginHandler(indexService, repo).Index))
	mux.Handle("/guilds", discordLoginRequiredChain.ThenFunc(guildsView.NewGuildsViewHandler(indexService).Index))
	mux.Handle("/guild/{guildId}", discordLoginRequiredChain.ThenFunc(guildIdView.NewGuildIDViewHandler(indexService).Index))
	mux.Handle("/guild/{guildId}/permission", discordLoginRequiredChain.ThenFunc(permissionView.NewPermissionViewHandler(indexService, repo).Index))
	mux.Handle("/guild/{guildId}/linetoken", discordLoginRequiredChain.ThenFunc(linetokenView.NewLineTokenViewHandler(indexService, repo).Index))
	mux.Handle("/guild/{guildId}/line-post-discord-channel", discordLoginRequiredChain.ThenFunc(linePostDiscordChannelView.NewLinePostDiscordChannelViewHandler(indexService, repo).Index))
	mux.Handle("/group/{guildId}", lineLoginRequiredChain.ThenFunc(groupView.NewLineGroupViewHandler(indexService, repo).Index))

	mux.Handle("/api/line-bot", middleChain.Then(linebot.NewLineBotHandler(indexService, repo)))
	mux.Handle("/login/discord", middleChain.Then(discordLogin.NewDiscordOAuth2Handler(discordOAuth2Service)))
	mux.Handle("/logout/discord", middleChain.Then(discordLogout.NewDiscordOAuth2Handler(discordOAuth2Service)))
	mux.Handle("/login/line/{guildId}", middleChain.ThenFunc(lineLogin.NewLineLoginHandler(indexService, repo).LineLogin))
	mux.Handle("/logout/line", middleChain.Then(lineLogout.NewLineLogoutHandler(indexService)))
	mux.Handle("/callback/discord-callback/", middleChain.Then(discordCallback.NewDiscordCallbackHandler(discordOAuth2Service)))
	mux.Handle("/callback/line-callback/", middleChain.Then(lineCallback.NewLineCallbackHandler(indexService, repo)))
	mux.Handle("/api/{guildId}/group", lineMiddleChain.Then(group.NewLineGroupHandler(indexService, repo)))
	mux.Handle("/api/{guildId}/permission", discordMiddleChain.Then(permission.NewPermissionHandler(indexService, repo)))
	mux.Handle("/api/{guildId}/linetoken", discordMiddleChain.Then(linetoken.NewLineTokenHandler(indexService, repo)))
	mux.Handle("/api/{guildId}/line-post-discord-channel", discordMiddleChain.Then(linePostDiscordChannel.NewLinePostDiscordChannelHandler(indexService, repo)))

	http.ListenAndServe(":"+config.Port(), mux)
}
