package commands

import (
	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	DB db.Driver
}

func NewCogHandler(db db.Driver) *CommandHandler {
	return &CommandHandler{
		DB: db,
	}
}

/*
スラッシュコマンドのハンドラ

スラッシュコマンドのハンドラは、
discordgo.Session.AddHandler()で登録する必要があります。

discordgo.Session.AddHandler()の引数には、
discordgo.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command.Executor(s, i)
}
のように、
discordgo.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command.Executor(s, i)
}
を渡す必要があります。
*/

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	AppCommand  *discordgo.ApplicationCommand
	Executor    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (c *Command) AddApplicationCommand(appCmd *discordgo.ApplicationCommand) {
	c.AppCommand = appCmd
}
