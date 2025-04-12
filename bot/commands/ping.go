package commands

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"
	"github.com/maguro-alternative/remake_bot/bot/config"

	"github.com/bwmarrin/discordgo"
)

func PingCommand(repo repository.RepositoryFunc) *command {
	/*
		pingコマンドの定義

		コマンド名: ping
		説明: Pong!
		オプション: なし
	*/
	exec := newCogHandler(repo, nil)
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
	req, err := http.NewRequest(http.MethodPost, config.InternalURL(), strings.NewReader(`{"message":"Pong"}`))
	if err != nil {
		fmt.Printf("error creating request: %v\n", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ChannelNo())
	resp, err := h.client.Do(req)
	if err != nil {
		fmt.Printf("error sending request: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("error response from server: %v\n", resp.Status)
		return fmt.Errorf("error response from server: %v", resp.Status)
	}
	return nil
}
