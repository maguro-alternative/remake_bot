package cogs

import (
	"context"
	"log/slog"
	"strings"

	"github.com/maguro-alternative/remake_bot/repository"
	//"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
)

func (h *cogHandler) onVoiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	ctx := context.Background()
	repo := repository.NewRepository(h.db)
	slog.InfoContext(ctx, "OnVoiceStateUpdate")
	_, err := onVoiceStateUpdateFunc(ctx, repo, s, s.State, vs)
	if err != nil {
		slog.ErrorContext(ctx, "", "", err.Error())
	}
}

func onVoiceStateUpdateFunc(
	ctx context.Context,
	//repo repository.RepositoryFunc,
	repo *repository.Repository,
	//s mock.Session,
	s *discordgo.Session,
	state *discordgo.State,
	m *discordgo.VoiceStateUpdate,
) ([]*discordgo.Message, error) {
	slog.InfoContext(ctx, "OnVoiceStateUpdateFunc")
	var vcChannelID string
	var sendText strings.Builder
	var embed *discordgo.MessageEmbed
	embed = nil
	vcChannelID = m.ChannelID
	if m.BeforeUpdate != nil {
		vcChannelID = m.BeforeUpdate.ChannelID
	}
	ngUserIDs, err := repo.GetVcSignalNgUsersByVcChannelIDAllColumn(ctx, vcChannelID)
	if err != nil {
		return nil, err
	}
	for _, ngUser := range ngUserIDs {
		if ngUser.UserID == m.UserID {
			return nil, nil
		}
	}
	ngRoleIDs, err := repo.GetVcSignalNgRolesByVcChannelIDAllColumn(ctx, vcChannelID)
	if err != nil {
		return nil, err
	}
	for _, ngRole := range ngRoleIDs {
		for _, roleID := range m.Member.Roles {
			if ngRole.RoleID == roleID {
				return nil, nil
			}
		}
	}
	vcSignalChannel, err := repo.GetVcSignalChannelAllColumnByVcChannelID(ctx, vcChannelID)
	if err != nil {
		return nil, err
	}
	if !vcSignalChannel.SendSignal {
		return nil, nil
	}
	if vcSignalChannel.SendChannelID == "" {
		return nil, nil
	}
	if vcSignalChannel.JoinBot && m.UserID == state.User.ID {
		return nil, nil
	}
	if vcSignalChannel.EveryoneMention {
		sendText.WriteString("@everyone ")
	}
	mentionUserIDs, err := repo.GetVcSignalMentionUsersByVcChannelID(ctx, vcChannelID)
	if err != nil {
		return nil, err
	}
	for _, mentionUser := range mentionUserIDs {
		if mentionUser.UserID == m.UserID {
			sendText.WriteString("<@!")
			sendText.WriteString(m.UserID)
			sendText.WriteString("> ")
		}
	}
	mentionRoleIDs, err := repo.GetVcSignalMentionRolesByVcChannelID(ctx, vcChannelID)
	if err != nil {
		return nil, err
	}
	for _, mentionRole := range mentionRoleIDs {
		for _, roleID := range m.Member.Roles {
			if mentionRole.RoleID == roleID {
				sendText.WriteString("<@&")
				sendText.WriteString(roleID)
				sendText.WriteString("> ")
			}
		}
	}
	// 画像がアニメーションの場合はa_が先頭につく
	//strings.HasPrefix(m.Member.Avatar, "a_")
	//chengeVcChannelFlag := (m.BeforeUpdate != nil) && (m.ChannelID != "") && (m.BeforeUpdate.ChannelID != m.ChannelID)
	if m.BeforeUpdate == nil || m.ChannelID != ""  && (!m.SelfVideo == !m.SelfStream) && (m.BeforeUpdate != nil && (!m.BeforeUpdate.SelfVideo == !m.BeforeUpdate.SelfStream)) {
		sendText.WriteString("入室")
	}
	if m.BeforeUpdate != nil && (!m.BeforeUpdate.SelfVideo == !m.BeforeUpdate.SelfStream) && (!m.SelfVideo == !m.SelfStream) {
		sendText.WriteString("退出")
	}
	if (m.BeforeUpdate != nil && !m.BeforeUpdate.SelfVideo) && m.SelfVideo {
		embed = &discordgo.MessageEmbed{
			Title:       "カメラ配信",
			Description: m.Member.User.Username + "\n" + "<#" + m.ChannelID + ">",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    m.Member.User.Username,
				IconURL: "https://cdn.discordapp.com/avatars/" + m.Member.User.ID + "/" + m.Member.Avatar + ".webp?size=64",
			},
		}
		sendText.WriteString("ビデオON")
	}
	if (m.BeforeUpdate != nil && m.BeforeUpdate.SelfVideo) && !m.SelfVideo {
		sendText.WriteString("ビデオOFF")
	}
	if (m.BeforeUpdate != nil && !m.BeforeUpdate.SelfStream) && m.SelfStream {
		presence, err := state.Presence(m.GuildID, m.UserID)
		if err != nil {
			return nil, err
		}
		if len(presence.Activities) == 0 {
			embed = &discordgo.MessageEmbed{
				Title:       "画面共有",
				Description: m.Member.User.Username + "\n" + "<#" + m.ChannelID + ">",
				Author: &discordgo.MessageEmbedAuthor{
					Name:    m.Member.User.Username,
					IconURL: "https://cdn.discordapp.com/avatars/" + m.Member.User.ID + "/" + m.Member.Avatar + ".webp?size=64",
				},
			}
			sendText.WriteString("配信ON")
		} else {
			embed = &discordgo.MessageEmbed{
				Title:       "配信タイトル:" + presence.Activities[0].Name,
				Description: m.Member.User.Username + "\n" + "<#" + m.ChannelID + ">",
				Image: &discordgo.MessageEmbedImage{
					URL: "https://cdn.discordapp.com/app-assets/" + presence.Activities[0].ApplicationID + "/" + presence.Activities[0].Assets.LargeImageID + ".png",
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    m.Member.User.Username,
					IconURL: "https://cdn.discordapp.com/avatars/" + m.Member.User.ID + "/" + m.Member.Avatar + ".webp?size=64",
				},
			}
			sendText.WriteString("配信ON")
		}
	}
	if (m.BeforeUpdate != nil && m.BeforeUpdate.SelfStream) && !m.SelfStream {
		sendText.WriteString("配信OFF")
	}

	_, err = s.ChannelMessageSend(vcSignalChannel.SendChannelID, sendText.String())
	if err == nil && embed != nil {
		_, err = s.ChannelMessageSendEmbed(vcSignalChannel.SendChannelID, embed)
		return nil, err
	}
	return nil, err
}
