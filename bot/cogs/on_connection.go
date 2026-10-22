package cogs

import (
	"github.com/maguro-alternative/remake_bot/bot/internal/health"

	"github.com/bwmarrin/discordgo"
)

// RegisterConnectionHandlers はWebSocket接続の状態変化を監視するハンドラを登録する。
func RegisterConnectionHandlers(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, d *discordgo.Disconnect) {
		health.OnDisconnect(s.ShouldReconnectOnError)
	})

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Resumed) {
		health.OnResumed()
	})

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		health.OnReady()
	})
}
