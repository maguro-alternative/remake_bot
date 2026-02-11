package commands

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/sharedtime"
	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	voicevoxcorego "github.com/sh1ma/voicevoxcore.go"
)

const (
	// voicevox_core local API endpoint
	VoiceVoxLocalAPIURL = "http://localhost:50021"
	// voicevox_core library path
	VoiceVoxCoreDictPath = "/voicevox_core_files/open_jtalk_dic_utf_8-1.11"
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
				Required:     false,
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

func (h *commandHandler) handleVoiceVox(
	s mock.Session,
	state *discordgo.State,
	voice map[string]*discordgo.VoiceConnection,
	i *discordgo.InteractionCreate,
) error {
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
	speacker := "ずんだもん"
	speakerId := "3"
	pitch := int64(100)      // 100 = 1.0x speed
	intonation := int64(100) // 100 = 1.0x intonation
	speed := int64(100)      // 100 = 1.0x speed

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
			speakerId = option.StringValue()
			for _, choice := range choices {
				if choice.Value == speakerId {
					speacker = choice.Name
					break
				}
			}
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
			Content: fmt.Sprintf("%s「 %s 」", speacker, text),
		},
	})
	if err != nil {
		fmt.Printf("error responding to voicevox command: %v\n", err)
		return err
	}

	sharedtime.SetSharedTime(i.GuildID, time.Now())

	err = voice[i.GuildID].Speaking(true)
	if err != nil {
		fmt.Printf("error speaking: %v\n", err)
		return err
	}

	filepath, err := getVoiceVoxFile(h.client, config.VoiceVoxKey(), text, speakerId, pitch, intonation, speed)
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
	speakerId string,
	pitch int64,
	intonation int64,
	speed int64,
) (string, error) {
	// Try voicevoxcore.go first
	audioData, err := getVoiceVoxFileFromLibrary(text, speakerId, pitch, intonation, speed)
	if err == nil {
		return audioData, nil
	}

	// Fallback to external API if library is unavailable
	fmt.Printf("voicevoxcore.go not available, falling back to external API: %v\n", err)
	return getVoiceVoxFileFromExternalAPI(client, key, text, speakerId, pitch, intonation, speed)
}

func getVoiceVoxFileFromLibrary(
	text string,
	speakerId string,
	pitch int64,
	intonation int64,
	speed int64,
) (string, error) {
	// Parse speaker ID to int
	var speakerIDInt int
	_, err := fmt.Sscanf(speakerId, "%d", &speakerIDInt)
	if err != nil {
		return "", fmt.Errorf("invalid speaker ID: %v", err)
	}

	// Check if voicevox_core dictionary exists
	if _, err := os.Stat(VoiceVoxCoreDictPath); os.IsNotExist(err) {
		return "", fmt.Errorf("voicevox_core dictionary not found at %s", VoiceVoxCoreDictPath)
	}

	// Initialize voicevoxcore.go
	core := voicevoxcorego.New()
	defer core.Finalize()

	// Set up initialization options
	// accelerationMode: 0 (auto), cpuNumThreads: 0 (auto), 
	// loadAllModels: false (load on demand), dictPath
	initOpts := voicevoxcorego.NewVoicevoxInitializeOptions(0, 0, false, VoiceVoxCoreDictPath)
	if err := core.Initialize(initOpts); err != nil {
		return "", fmt.Errorf("failed to initialize voicevoxcore: %w", err)
	}

	// Load the speaker model
	if err := core.LoadModel(uint(speakerIDInt)); err != nil {
		return "", fmt.Errorf("failed to load model for speaker %d: %w", speakerIDInt, err)
	}

	// Generate AudioQuery first
	audioQueryOpts := voicevoxcorego.NewVoicevoxAudioQueryOptions(false)
	query, err := core.AudioQuery(text, uint(speakerIDInt), audioQueryOpts)
	if err != nil {
		return "", fmt.Errorf("failed to generate audio query: %w", err)
	}

	// Set voice parameters (convert from percentage to fractional scale)
	// pitch: 100 = 1.0x
	// speed: 100 = 1.0x
	// intonation: 100 = 1.0x
	query.PitchScale = float32(pitch) / 100.0
	query.SpeedScale = float32(speed) / 100.0
	query.IntonationScale = float32(intonation) / 100.0

	// Generate audio using modified AudioQuery
	synthesisOpts := voicevoxcorego.NewVoicevoxSynthesisOptions(false)
	audioData, err := core.Synthesis(query, speakerIDInt, synthesisOpts)
	if err != nil {
		return "", fmt.Errorf("failed to synthesize audio: %w", err)
	}

	// Save to temporary file
	temppath := filepath.Join(os.TempDir(), fmt.Sprintf("voicevox_%d.wav", time.Now().UnixNano()))
	if err := os.WriteFile(temppath, audioData, 0600); err != nil {
		return "", fmt.Errorf("failed to write audio file: %w", err)
	}

	return temppath, nil
}

func getVoiceVoxFileFromExternalAPI(
	client *http.Client,
	key string,
	text string,
	speakerId string,
	pitch int64,
	intonation int64,
	speed int64,
) (string, error) {
	if key == "" {
		return "", fmt.Errorf("voicevox_core API unavailable and no VOICEVOX_KEY configured for external API")
	}

	filepath := fmt.Sprintf("%s/%s.wav", os.TempDir(), time.Now().Format("20060102150405"))
	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// External API parameters - convert to the correct format (0-100 scale to -0.15~0.15 for pitch, 0.5~2.0 for speed)
	// pitch: -100~100 (percentage) → -0.15~0.15 (absolute pitch change in semitones)
	// speed: 0~200 (percentage) → 0.5~2.0 (multiplier)
	// intonation: 0~200 (percentage) → 0.0~2.0 (multiplier)

	pitchValue := (float64(pitch) - 100.0) / 100.0 * 0.15 // convert to ±0.15 range
	speedValue := float64(speed) / 100.0                  // convert to 0.5~2.0 range (default 1.0)
	intonationValue := float64(intonation) / 100.0        // convert to 0.5~2.0 range (default 1.0)

	externalURL := fmt.Sprintf(
		"https://api.su-shiki.com/v2/voicevox/audio/?key=%s&speaker=%s&pitch=%.2f&intonationScale=%.2f&speed=%.2f&text=%s",
		key, speakerId, pitchValue, intonationValue, speedValue, url.QueryEscape(text),
	)

	fmt.Printf("External API URL: %s\n", externalURL)

	req, err := http.NewRequest("GET", externalURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to external voicevox API: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return "", fmt.Errorf("external voicevox API error: %d - %s", resp.StatusCode, string(bodyBytes))
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	return filepath, err
}
