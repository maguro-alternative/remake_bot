package bot

import (
	"net/http"
	"fmt"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/bot/cogs"
	"github.com/maguro-alternative/remake_bot/bot/commands"
	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/bwmarrin/discordgo"
	"github.com/cockroachdb/errors"
)

func BotOnReady(indexDB db.Driver, client *http.Client) (*discordgo.Session, func(), error) {
	/*
		ボットの起動

		args:
		indexDB: db.Driver
		データベースのドライバー

		return:
		*discordgo.Session
		エラーがなければ、セッションを返します。
		エラーがあれば、エラーを返します。
	*/
	// セッションを作成
	discordToken := "Bot " + config.DiscordBotToken()
	discordSession, err := discordgo.New(discordToken)
	if err != nil {
		return nil, func(){}, errors.WithStack(err)
	}
	discordSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	discordSession.Token = discordToken
	err = discordSession.Open()
	if err != nil {
		return nil, func(){}, errors.WithStack(err)
	}
	registerHandlers(discordSession, indexDB, client)
	cleanupCommandHandlers, err := commands.RegisterCommands(discordSession, indexDB)
	return discordSession, cleanupCommandHandlers, err
}

func registerHandlers(
	s *discordgo.Session,
	sqlxdb db.Driver,
	client *http.Client,
) {
	cogs := cogs.NewCogHandler(sqlxdb, client)
	fmt.Println(s.State.User.Username + "としてログインしました")
	//s.AddHandler(cogs.OnVoiceStateUpdate)
	s.AddHandler(cogs.OnMessageCreate)
}
