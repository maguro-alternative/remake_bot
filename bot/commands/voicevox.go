package commands

import (
	"fmt"
	"net/http"
	"io"
	"os"
	"time"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

var (
	choices = []*discordgo.ApplicationCommandOptionChoice{
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
	}
)

func VoiceVoxCommand(repo repository.RepositoryFunc, client *http.Client) *command {
	/*
		pingコマンドの定義

		コマンド名: ping
		説明: Pong!
		オプション: なし
	*/
	exec := newCogHandler(repo, client)
	return &command{
		Name:        "voicevox",
		Description: "ずんだもんたちが喋るよ！",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "しゃべらせたいテキスト",
				Required:    true,
			},
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "speaker",
				Description:  "しゃべる人",
				Required:     true,
				Autocomplete: true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "pitch",
				Description: "声の高さ",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "intonation",
				Description: "声の抑揚",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "speed",
				Description: "しゃべる速さ",
				Required:    false,
			},
		},
		Executor: exec.handleVoiceVox,
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

	text := ""
	speaker := "3"
	pitch := int64(0)
	intonation := int64(1)
	speed := int64(1)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices, // オートコンプリート用の選択肢
		},
	})
	// オートコンプリートのレスポンスはエラーを無視する
	if err == nil {
		return nil
	}

	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "text":
			text = option.StringValue()
		case "speaker":
			speaker = option.StringValue()
		case "pitch":
			pitch = option.IntValue()
		case "intonation":
			intonation = option.IntValue()
		case "speed":
			speed = option.IntValue()
		}
	}

	playCh := make(chan bool)
	userVoiceState, err := state.VoiceState(i.GuildID, i.Member.User.ID)
	if err != nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "ボイスチャンネルに参加してからコマンドを実行して下さい",
			},
		})
		if err != nil {
			fmt.Printf("error responding to voicevox command: %v\n", err)
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

	filepath, err := getVoiceVoxFile(h.client, config.VoiceVoxKey(), text, speaker, pitch, intonation, speed)
	if err != nil {
		fmt.Printf("error getting voicevox file: %v\n", err)
		return err
	}
	defer os.Remove(filepath)

	dgvoice.PlayAudioFile(voice[i.GuildID], filepath, playCh)

	<-playCh

	return voice[i.GuildID].Speaking(false)
}

func getVoiceVoxFile(
	client *http.Client,
	key string,
	text string,
	speaker string,
	pitch int64,
	intonation int64,
	speed int64,
) (string, error) {
	filepath := fmt.Sprintf("%s/%s.wav", os.TempDir(), time.Now().Format("20060102150405"))
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("error creating file: %v\n", err)
		return "", err
	}
	defer file.Close()

	//"https://api.su-shiki.com/v2/voicevox/audio/?key={key}&speaker={id}&pitch={pitch}&intonationScale={intonation}&speed={speed}&text={text}"
	url := fmt.Sprintf("https://api.su-shiki.com/v2/voicevox/audio/?key=%s&speaker=%s&pitch=%d&intonationScale=%d&speed=%d&text=%s", key, speaker, pitch, intonation, speed, text)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error while getting voicevox file: %v", resp.Status)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	return filepath, err
}
