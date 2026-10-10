package service

import (
	"net/http"

	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/sessions"
)

// A TODOService implements CRUD of TODO entities.
type IndexService struct {
	Client          *http.Client
	CookieStore     *sessions.CookieStore
	DiscordSession  mock.Session
	DiscordBotState *discordgo.State
}

// NewTODOService returns new TODOService.
func NewIndexService(
	client *http.Client,
	cookieStore *sessions.CookieStore,
	discordSession mock.Session,
	discordBotState *discordgo.State,
) *IndexService {
	return &IndexService{
		Client:          client,
		CookieStore:     cookieStore,
		DiscordSession:  discordSession,
		DiscordBotState: discordBotState,
	}
}

// IsGuildChannel はチャンネルが指定したサーバーに属しているかを返します。
// 他のサーバーのチャンネルIDを設定値として保存されないよう、APIで受け取ったチャンネルIDの検証に使います。
// チャンネルIDが空の場合は未設定とみなし、trueを返します。
func (s *IndexService) IsGuildChannel(guildID, channelID string) bool {
	if channelID == "" {
		return true
	}
	if s.DiscordBotState == nil {
		return false
	}
	channel, err := s.DiscordBotState.Channel(channelID)
	if err != nil {
		return false
	}
	return channel.GuildID == guildID
}
