package commands

import (
	"fmt"

	"github.com/maguro-alternative/remake_bot/bot/ffmpeg"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

func VoiceVoxCommand(repo repository.RepositoryFunc, playFf *ffmpeg.PlayFfmpegInterface) *command {
	/*
		pingコマンドの定義

		コマンド名: ping
		説明: Pong!
		オプション: なし
	*/
	exec := newCogHandler(repo, playFf)
	return &command{
		Name:        "voicevox",
		Description: "ずんだもんたちが喋るよ！",
		Options:     []*discordgo.ApplicationCommandOption{},
		Executor:    exec.handleVoiceVox,
	}
}

func (h *commandHandler) handleVoiceVox(s mock.Session, state *discordgo.State, voice map[string]*discordgo.VoiceConnection, i *discordgo.InteractionCreate) error {
	/*
		pingコマンドの実行

		コマンドの実行結果を返す
	*/
	if i.Interaction.Data.(discordgo.ApplicationCommandInteractionData).Name != "voicevox" {
		return nil
	}
	if i.Interaction.GuildID != i.GuildID {
		return nil
	}
	userVoiceState, err := state.VoiceState(i.GuildID, i.Member.User.ID)
	if err != nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "ボイスチャンネルに参加してからコマンドを実行してね！",
			},
		})
		if err != nil {
			fmt.Printf("error responding to ping command: %v\n", err)
			return err
		}
		return nil
	}

	if voice[i.GuildID] == nil {
		_, err = s.ChannelVoiceJoin(i.GuildID, userVoiceState.ChannelID, false, true)
		if err != nil {
			fmt.Printf("error joining voice channel: %v\n", err)
			return nil
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong",
		},
	})
	if err != nil {
		fmt.Printf("error responding to voicevox command: %v\n", err)
		return err
	}

	err = voice[i.GuildID].Speaking(true)
	if err != nil {
		fmt.Printf("error speaking: %v\n", err)
		return err
	}

	dgvoice.PlayAudioFile(voice[i.GuildID], "testutil/files/yumi_dannasama.mp3", make(chan bool))

	defer func() error {
		err = voice[i.GuildID].Speaking(false)
		return err
	}()
	return err
}

