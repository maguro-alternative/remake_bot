package commands

import (
	"fmt"

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

type Handler struct {
	session  *discordgo.Session
	commands map[string]*Command
	guild    string
}

// スラッシュコマンドの作成
func NewCommandHandler(
	session *discordgo.Session,
	guildID string,
) *Handler {
	return &Handler{
		session:  session,
		commands: make(map[string]*Command),
		guild:    guildID,
	}
}

// スラッシュコマンドの登録
func (h *Handler) CommandRegister(command *Command) error {
	if _, exists := h.commands[command.Name]; exists {
		return fmt.Errorf("command with name `%s` already exists", command.Name)
	}

	appCmd, err := h.session.ApplicationCommandCreate(
		h.session.State.User.ID,
		h.guild,
		&discordgo.ApplicationCommand{
			ApplicationID: h.session.State.User.ID,
			Name:          command.Name,
			Description:   command.Description,
			Options:       command.Options,
		},
	)
	if err != nil {
		return err
	}

	command.AddApplicationCommand(appCmd)

	h.commands[command.Name] = command

	h.session.AddHandler(
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			command.Executor(s, i)
		},
	)

	return nil
}

// スラッシュコマンドの削除
func (h *Handler) CommandRemove(command *Command) error {
	err := h.session.ApplicationCommandDelete(h.session.State.User.ID, h.guild, command.AppCommand.ID)
	if err != nil {
		return fmt.Errorf("error while deleting application command: %v", err)
	}

	delete(h.commands, command.Name)

	return nil
}

// スラッシュコマンドの取得
func (h *Handler) GetCommands() []*Command {
	var commands []*Command

	for _, v := range h.commands {
		commands = append(commands, v)
	}

	return commands
}

func RegisterCommands(discordSession *discordgo.Session, db db.Driver) func() {
	var commandHandlers []*Handler
	// 所属しているサーバすべてにスラッシュコマンドを追加する
	// NewCommandHandlerの第二引数を空にすることで、グローバルでの使用を許可する
	commandHandler := NewCommandHandler(discordSession, "")
	// 追加したいコマンドをここに追加
	commandHandler.CommandRegister(PingCommand(db))
	commandHandlers = append(commandHandlers, commandHandler)
	cleanupCommandHandlers := func() {
		for _, handler := range commandHandlers {
			for _, command := range handler.GetCommands() {
				err := handler.CommandRemove(command)
				if err != nil {
					fmt.Printf("error while removing command: %v\n", err)
				}
			}
		}
	}
	return cleanupCommandHandlers
}
