package cogs

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

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
	repo repository.RepositoryFunc,
	s mock.Session,
	state *discordgo.State,
	vs *discordgo.VoiceStateUpdate,
) ([]*discordgo.Message, error) {
	var sendText, mentionText strings.Builder
	var embed *discordgo.MessageEmbed
	var sendMessages []*discordgo.Message
	embed = nil
	vcChannelID := vs.ChannelID
	guildId := vs.GuildID
	if vs.BeforeUpdate != nil {
		vcChannelID = vs.BeforeUpdate.ChannelID
		guildId = vs.BeforeUpdate.GuildID
	}
	vcChannel, err := state.Channel(vcChannelID)
	if err != nil {
		return nil, err
	}
	if vs.BeforeUpdate != nil && (vs.BeforeUpdate.SelfDeaf != vs.SelfDeaf || vs.BeforeUpdate.SelfMute != vs.SelfMute) {
		return nil, nil
	}
	ngUserIDs, err := repo.GetVcSignalNgUsersByVcChannelIDAllColumn(ctx, vcChannelID)
	if err != nil {
		return nil, err
	}
	for _, ngUser := range ngUserIDs {
		if ngUser.UserID == vs.UserID {
			return nil, nil
		}
	}
	ngRoleIDs, err := repo.GetVcSignalNgRolesByVcChannelIDAllColumn(ctx, vcChannelID)
	if err != nil {
		return nil, err
	}
	for _, ngRole := range ngRoleIDs {
		for _, roleID := range vs.Member.Roles {
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
	if vcSignalChannel.JoinBot && vs.UserID == state.User.ID {
		return nil, nil
	}
	if vcSignalChannel.EveryoneMention {
		mentionText.WriteString("@everyone ")
	}
	mentionUserIDs, err := repo.GetVcSignalMentionUsersByVcChannelID(ctx, vcChannelID)
	if err != nil {
		return nil, err
	}
	for _, mentionUser := range mentionUserIDs {
		if mentionUser.UserID == vs.UserID {
			mentionText.WriteString("<@")
			mentionText.WriteString(vs.UserID)
			mentionText.WriteString("> ")
		}
	}
	mentionRoleIDs, err := repo.GetVcSignalMentionRolesByVcChannelID(ctx, vcChannelID)
	if err != nil {
		return nil, err
	}
	for _, mentionRole := range mentionRoleIDs {
		for _, roleID := range vs.Member.Roles {
			if mentionRole.RoleID == roleID {
				mentionText.WriteString("<@&")
				mentionText.WriteString(roleID)
				mentionText.WriteString("> ")
			}
		}
	}
	//chengeVcChannelFlag := (vs.BeforeUpdate != nil) && (vs.ChannelID != "") && (vs.BeforeUpdate.ChannelID != vs.ChannelID)
	if vs.BeforeUpdate == nil || vs.ChannelID != "" && (!vs.SelfVideo == !vs.SelfStream) && (!vs.BeforeUpdate.SelfVideo == !vs.BeforeUpdate.SelfStream) {
		vcChannel, err = state.Channel(vs.ChannelID)
		if err != nil {
			return nil, err
		}
		membersCount := vcMembersCount(state, guildId, vs.ChannelID)
		sendText.WriteString(mentionText.String())
		sendText.WriteString("現在" + strconv.Itoa(membersCount) + "人 <@" + vs.Member.User.ID + "> が " + vcChannel.Name + "に入室しました。")
		if membersCount == 1 {
			embed = &discordgo.MessageEmbed{
				Title:       "通話開始",
				Description: "<#" + vs.ChannelID + ">",
				Author: &discordgo.MessageEmbedAuthor{
					Name:    vs.Member.User.Username,
					IconURL: vs.Member.AvatarURL("64"),
				},
			}
		}
		sendMessage, err := s.ChannelMessageSend(vcSignalChannel.SendChannelID, sendText.String())
		if err != nil {
			return nil, err
		}
		sendMessages = append(sendMessages, sendMessage)
		sendText.Reset()
	}
	if vs.BeforeUpdate != nil && (!vs.BeforeUpdate.SelfVideo == !vs.BeforeUpdate.SelfStream) && (!vs.SelfVideo == !vs.SelfStream) {
		vcChannel, err = state.Channel(vs.BeforeUpdate.ChannelID)
		if err != nil {
			return nil, err
		}
		membersCount := vcMembersCount(state, guildId, vs.BeforeUpdate.ChannelID)
		sendText.WriteString("現在" + strconv.Itoa(membersCount) + "人 <@" + vs.Member.User.ID + "> が " + vcChannel.Name + "から退室しました。")
		sendMessage, err := s.ChannelMessageSend(vcSignalChannel.SendChannelID, sendText.String())
		if err != nil {
			return nil, err
		}
		sendMessages = append(sendMessages, sendMessage)
		sendText.Reset()
		if membersCount == 0 {
			guildMembersCount := guildVcMembersCount(state, guildId)
			if guildMembersCount == 0 {
				embed = &discordgo.MessageEmbed{
					Title: "通話終了",
				}
				sendText.WriteString(mentionText.String())
				sendText.WriteString("通話が終了しました。")
				sendMessage, err := s.ChannelMessageSend(vcSignalChannel.SendChannelID, sendText.String())
				if err != nil {
					return nil, err
				}
				sendMessages = append(sendMessages, sendMessage)
				sendText.Reset()
			}
		}
	}
	if (vs.BeforeUpdate != nil && !vs.BeforeUpdate.SelfVideo) && vs.SelfVideo {
		embed = &discordgo.MessageEmbed{
			Title:       "カメラ配信",
			Description: vs.Member.User.Username + "\n" + "<#" + vs.ChannelID + ">",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    vs.Member.User.Username,
				IconURL: vs.Member.AvatarURL("64"),
			},
		}
		sendText.WriteString(mentionText.String())
		sendText.WriteString("<@" + vs.Member.User.ID + "> が" + vcChannel.Name + "カメラ配信を開始しました。")
		sendMessage, err := s.ChannelMessageSend(vcSignalChannel.SendChannelID, sendText.String())
		if err != nil {
			return nil, err
		}
		sendMessages = append(sendMessages, sendMessage)
		sendText.Reset()
	}
	if (vs.BeforeUpdate != nil && vs.BeforeUpdate.SelfVideo) && !vs.SelfVideo {
		sendText.WriteString("<@" + vs.Member.User.ID + "> がカメラ配信を終了しました。")
		sendMessage, err := s.ChannelMessageSend(vcSignalChannel.SendChannelID, sendText.String())
		if err != nil {
			return nil, err
		}
		sendMessages = append(sendMessages, sendMessage)
		sendText.Reset()
	}
	if (vs.BeforeUpdate != nil && !vs.BeforeUpdate.SelfStream) && vs.SelfStream {
		presence, err := state.Presence(vs.GuildID, vs.UserID)
		if err != nil {
			return nil, err
		}
		if len(presence.Activities) == 0 {
			embed = &discordgo.MessageEmbed{
				Title:       "画面共有",
				Description: vs.Member.User.Username + "\n" + "<#" + vs.ChannelID + ">",
				Author: &discordgo.MessageEmbedAuthor{
					Name:    vs.Member.User.Username,
					IconURL: vs.Member.AvatarURL("64"),
				},
			}
			sendText.WriteString(mentionText.String())
			sendText.WriteString("<@" + vs.Member.User.ID + "> が" + vcChannel.Name + "で画面共有を開始しました。")
			sendMessage, err := s.ChannelMessageSend(vcSignalChannel.SendChannelID, sendText.String())
			if err != nil {
				return nil, err
			}
			sendMessages = append(sendMessages, sendMessage)
			sendText.Reset()
		} else {
			embed = &discordgo.MessageEmbed{
				Title:       "配信タイトル:" + presence.Activities[0].Name,
				Description: vs.Member.User.Username + "\n" + "<#" + vs.ChannelID + ">",
				Image: &discordgo.MessageEmbedImage{
					URL: "https://cdn.discordapp.com/app-assets/" + presence.Activities[0].ApplicationID + "/" + presence.Activities[0].Assets.LargeImageID + ".png",
				},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    vs.Member.User.Username,
					IconURL: vs.Member.AvatarURL("64"),
				},
			}
			sendText.WriteString(mentionText.String())
			sendText.WriteString("<@" + vs.Member.User.ID + "> が" + vcChannel.Name + "で" + presence.Activities[0].Name + "を配信開始しました。")
			sendMessage, err := s.ChannelMessageSend(vcSignalChannel.SendChannelID, sendText.String())
			if err != nil {
				return nil, err
			}
			sendMessages = append(sendMessages, sendMessage)
			sendText.Reset()
		}
	}
	if (vs.BeforeUpdate != nil && vs.BeforeUpdate.SelfStream) && !vs.SelfStream {
		sendText.WriteString("<@" + vs.Member.User.ID + "> が画面共有を終了しました。")
		sendMessage, err := s.ChannelMessageSend(vcSignalChannel.SendChannelID, sendText.String())
		if err != nil {
			return nil, err
		}
		sendMessages = append(sendMessages, sendMessage)
		sendText.Reset()
	}

	if embed != nil {
		sendMessage, err := s.ChannelMessageSendEmbed(vcSignalChannel.SendChannelID, embed)
		if err != nil {
			return sendMessages, err
		}
		sendMessages = append(sendMessages, sendMessage)
		sendText.Reset()
	}
	return sendMessages, err
}

func vcMembersCount(state *discordgo.State, guildID, vcChannelID string) int {
	guild, err := state.Guild(guildID)
	if err != nil {
		return 0
	}
	memberCount := 0
	for _, voiceState := range guild.VoiceStates {
		if voiceState.ChannelID == vcChannelID {
			memberCount++
		}
	}
	return memberCount
}

func guildVcMembersCount(state *discordgo.State, guildID string) int {
	guild, err := state.Guild(guildID)
	if err != nil {
		return 0
	}
	return len(guild.VoiceStates)
}
