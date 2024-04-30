package cogs

import (
	"context"
	"log/slog"

	"github.com/maguro-alternative/remake_bot/repository"
	//"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
)

func (h *cogHandler) onVoiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	ctx := context.Background()
	repo := repository.NewRepository(h.db)
	slog.InfoContext(ctx, "OnVoiceStateUpdate")
	_, err := h.onVoiceStateUpdateFunc(ctx, repo, s, vs)
	if err != nil {
		slog.ErrorContext(ctx, "", "", err.Error())
	}
}

func (h *cogHandler) onVoiceStateUpdateFunc(
	ctx context.Context,
	repo repository.RepositoryFunc,
	//s mock.Session,
	s *discordgo.Session,
	m *discordgo.VoiceStateUpdate,
) (*discordgo.Message, error) {
	slog.InfoContext(ctx, "OnVoiceStateUpdateFunc")
	//fmt.Println(m.ChannelID)				// After
	//fmt.Println(m.BeforeUpdate.ChannelID)	// Before
	var vcChannelID string
	vcChannelID = m.ChannelID
	if m.BeforeUpdate != nil {
		vcChannelID = m.BeforeUpdate.ChannelID
	}
	// ngUserIDs, err := repo.GetVcSignalNgUserIDs(ctx, vcChannelID)
	// ngRoleIDs, err := repo.GetVcSignalNgRoleIDs(ctx, vcChannelID)
	// vcSignalChannel, err := repo.GetVcSignalChannel(ctx, vcChannelID)
	// mentionUserIDs, err := repo.GetMentionUserIDs(ctx, vcChannelID)
	// mentionRoleIDs, err := repo.GetMentionRoleIDs(ctx, vcChannelID)
	// chengeVcChannelFlag := (m.BeforeUpdate != nil) && (m.ChannelID != "") && (m.BeforeUpdate.ChannelID != m.ChannelID)
	// chengeVcChannelFlag || m.ChannelID != "" //taisyutsu
	// chengeVcChannelFlag || m.BeforeUpdate != nil //nyuusitsu
	// m.BeforeUpdate.SelfVideo == false && m.SelfVideo //kamera start
	// m.BeforeUpdate.SelfStream == false && m.SelfStream // haishin start

	_, err := s.ChannelMessageSend(vcChannelID, "Hello")
	return nil, err
}
