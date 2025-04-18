package commands

import (
	"fmt"
	"net/http"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
)

func PingCommand(repo repository.RepositoryFunc) *command {
	/*
		pingコマンドの定義

		コマンド名: ping
		説明: Pong!
		オプション: なし
	*/
	exec := newCogHandler(repo, &http.Client{})
	return &command{
		Name:        "ping",
		Description: "Pong!",
		Options:     []*discordgo.ApplicationCommandOption{},
		Executor:    exec.handlePing,
	}
}

func (h *commandHandler) handlePing(s mock.Session, state *discordgo.State, voice map[string]*discordgo.VoiceConnection, i *discordgo.InteractionCreate) error {
	/*
		pingコマンドの実行

		コマンドの実行結果を返す
	*/
	if i.Interaction.Data.(discordgo.ApplicationCommandInteractionData).Name != "ping" {
		return nil
	}
	if i.Interaction.GuildID != i.GuildID {
		return nil
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong",
		},
	})
	if err != nil {
		fmt.Printf("error responding to ping command: %v\n", err)
		return err
	}
	return nil
}
