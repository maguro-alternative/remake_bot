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
		Options:     []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "しゃべらせたいテキスト",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "speaker",
				Description: "しゃべる人",
				Autocomplete: true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "四国めたん",
						Value: "2",
					},
					{
						Name:  "四国めたんあまあま",
						Value: "0",
					},
					{
						Name:  "四国めたんツンツン",
						Value: "6",
					},
					{
						Name:  "四国めたんセクシー",
						Value: "4",
					},
					{
						Name:  "ずんだもん",
						Value: "3",
					},
					{
						Name:  "ずんだもんあまあま",
						Value: "1",
					},
					{
						Name:  "ずんだもんツンツン",
						Value: "7",
					},
					{
						Name:  "ずんだもんセクシー",
						Value: "5",
					},
					{
						Name:  "ずんだもんささやき",
						Value: "22",
					},
					{
						Name:  "春日部つむぎ",
						Value: "8",
					},
					{
						Name:  "雨晴はう",
						Value: "10",
					},
					{
						Name:  "波音リツ",
						Value: "9",
					},
					{
						Name:  "玄野武宏",
						Value: "11",
					},
					{
						Name:  "白上虎太郎",
						Value: "12",
					},
					{
						Name:  "青山龍星",
						Value: "13",
					},
					{
						Name:  "冥鳴ひまり",
						Value: "14",
					},
					{
						Name:  "九州そら",
						Value: "16",
					},
					{
						Name:  "九州そらあまあま",
						Value: "15",
					},
					{
						Name:  "九州そらツンツン",
						Value: "18",
					},
					{
						Name:  "九州そらセクシー",
						Value: "17",
					},
					{
						Name:  "九州そらささやき",
						Value: "19",
					},
					{
						Name:  "もち子さん",
						Value: "20",
					},
					{
						Name:  "剣崎雌雄",
						Value: "21",
					},
				},
			},
		},
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

	playCh := make(chan bool)
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
	i.ApplicationCommandData()

	dgvoice.PlayAudioFile(voice[i.GuildID], "testutil/files/yumi_dannasama.mp3", playCh)

	<-playCh

	return voice[i.GuildID].Speaking(false)
}
