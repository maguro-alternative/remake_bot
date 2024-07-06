package commands

import (
	"fmt"
	"net/http"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
)

func VoiceDisconnectCommand(repo repository.RepositoryFunc, client *http.Client) *command {
	exec := newCogHandler(repo, client)
	return &command{
		Name:        "voice_disconnect",
		Description: "ボイスチャンネルから切断します",
		Executor:    exec.handleVoiceDisconnect,
	}
}

func (h *commandHandler) handleVoiceDisconnect(
	s mock.Session,
	state *discordgo.State,
	voice map[string]*discordgo.VoiceConnection,
	i *discordgo.InteractionCreate,
) error {
	/*
		pingコマンドの実行

		コマンドの実行結果を返す
	*/
	if i.Interaction.Data.(discordgo.ApplicationCommandInteractionData).Name != "voice_disconnect" {
		return nil
	}
	if i.Interaction.GuildID != i.GuildID {
		return nil
	}

	_, err := state.VoiceState(i.GuildID, i.User.ID)
	if err != nil || voice[i.GuildID] == nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Botが参加していません",
			},
		})
		if err != nil {
			fmt.Printf("error responding to disconnect command: %v\n", err)
			return err
		}
		return nil
	}

	err = voice[i.GuildID].Disconnect()
	if err != nil {
		fmt.Printf("error disconnecting from voice channel: %v\n", err)
		return err
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "ボイスチャンネルから切断しました",
		},
	})
	return err
}
