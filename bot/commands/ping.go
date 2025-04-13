package commands

import (
	"fmt"
	"net/http"
	"strings"
	"encoding/json"
	"net/url"

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
	form := url.Values{}
	form.Add("message", "test")
	req, err := http.NewRequest(http.MethodPost, config.InternalURL(), strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Printf("error creating request: %v\n", err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+config.ChannelNo())
	resp, err := h.client.Do(req)
	if err != nil {
		fmt.Printf("error sending request: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var body map[string]interface{}
		fmt.Printf("error response from server: %v\n", resp.Status)
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			fmt.Printf("error decoding response body: %v\n", err)
		}
		fmt.Printf("response body: %v\n", body)
		return fmt.Errorf("error response from server: %v", resp.Status)
	}
	return nil
}
