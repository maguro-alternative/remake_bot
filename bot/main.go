package bot

import (
	"fmt"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/bot/cogs"
	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/bwmarrin/discordgo"
	"github.com/cockroachdb/errors"
)

/*
スラッシュコマンドとハンドラの登録

スラッシュコマンドとハンドラの登録は、
discordgo.Session.ApplicationCommandCreate()と
discordgo.Session.AddHandler()を使って行います。
*/

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	AppCommand  *discordgo.ApplicationCommand
	Executor    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type Handler struct {
	session  *discordgo.Session
	commands map[string]*Command
	guild    string
}

func BotOnReady(indexDB db.Driver) (*discordgo.Session, error) {
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
	discordSession, err := discordgo.New("Bot " + config.Token())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	discordSession.Identify.Intents = discordgo.IntentsAll
	discordSession.Token = config.Token()
	err = discordSession.Open()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	registerHandlers(discordSession, indexDB)
	return discordSession, nil
}

func registerHandlers(
	s *discordgo.Session,
	sqlxdb db.Driver,
) {
	cogs := cogs.NewCogHandler(sqlxdb)
	fmt.Println(s.State.User.Username + "としてログインしました")
	s.AddHandler(cogs.OnVoiceStateUpdate)
}
