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
	var sendText, beforeMentionText, afterMentionText strings.Builder
	var embed *discordgo.MessageEmbed
	var sendMessages []*discordgo.Message
	var beforeNgUserIDs, afterNgUserIDs []*repository.VcSignalNgUserAllColumn
	var beforeNgRoleIDs, afterNgRoleIDs[]*repository.VcSignalNgRoleAllColumn
	var beforeVcSignalChannel, afterVcSignalChannel *repository.VcSignalChannelAllColumn
	var beforeMentionUserIDs, afterMentionUserIDs []*repository.VcSignalMentionUser
	var beforeMentionRoleIDs, afterMentionRoleIDs []*repository.VcSignalMentionRole
	embed = nil
	if vs.BeforeUpdate != nil && (vs.BeforeUpdate.SelfDeaf != vs.SelfDeaf || vs.BeforeUpdate.SelfMute != vs.SelfMute) {
		return nil, nil
	}
	afterNgUserIDs, err := repo.GetVcSignalNgUsersByVcChannelIDAllColumn(ctx, vs.ChannelID)
	if err != nil {
		return nil, err
	}
	for _, afterNgUser := range afterNgUserIDs {
		if afterNgUser.UserID == vs.UserID {
			return nil, nil
		}
	}
	afterNgRoleIDs, err = repo.GetVcSignalNgRolesByVcChannelIDAllColumn(ctx, vs.ChannelID)
	if err != nil {
		return nil, err
	}
	for _, afterNgRole := range afterNgRoleIDs {
		for _, roleID := range vs.Member.Roles {
			if afterNgRole.RoleID == roleID {
				return nil, nil
			}
		}
	}
	afterVcSignalChannel, err = repo.GetVcSignalChannelAllColumnByVcChannelID(ctx, vs.ChannelID)
	if err != nil {
		return nil, err
	}
	if !afterVcSignalChannel.SendSignal {
		return nil, nil
	}
	if afterVcSignalChannel.SendChannelID == "" {
		return nil, nil
	}
	if afterVcSignalChannel.JoinBot && vs.UserID == state.User.ID {
		return nil, nil
	}
	if afterVcSignalChannel.EveryoneMention {
		afterMentionText.WriteString("@everyone ")
	}
	afterMentionUserIDs, err = repo.GetVcSignalMentionUsersByVcChannelID(ctx, vs.ChannelID)
	if err != nil {
		return nil, err
	}
	for _, afterMentionUser := range afterMentionUserIDs {
		if afterMentionUser.UserID == vs.UserID {
			afterMentionText.WriteString("<@")
			afterMentionText.WriteString(vs.UserID)
			afterMentionText.WriteString("> ")
		}
	}
	afterMentionRoleIDs, err = repo.GetVcSignalMentionRolesByVcChannelID(ctx, vs.ChannelID)
	if err != nil {
		return nil, err
	}
	for _, afterMentionRole := range afterMentionRoleIDs {
		for _, roleID := range vs.Member.Roles {
			if afterMentionRole.RoleID == roleID {
				afterMentionText.WriteString("<@&")
				afterMentionText.WriteString(roleID)
				afterMentionText.WriteString("> ")
			}
		}
	}
	if vs.BeforeUpdate != nil && vs.ChannelID != vs.BeforeUpdate.ChannelID {
		beforeVcSignalChannel, err = repo.GetVcSignalChannelAllColumnByVcChannelID(ctx, vs.BeforeUpdate.ChannelID)
		if err != nil {
			return nil, err
		}
		if beforeVcSignalChannel.SendSignal {
			beforeMentionUserIDs, err = repo.GetVcSignalMentionUsersByVcChannelID(ctx, vs.BeforeUpdate.ChannelID)
			if err != nil {
				return nil, err
			}
			for _, beforeMentionUser := range beforeMentionUserIDs {
				if beforeMentionUser.UserID == vs.UserID {
					beforeMentionText.WriteString("<@")
					beforeMentionText.WriteString(vs.UserID)
					beforeMentionText.WriteString("> ")
				}
			}
			beforeMentionRoleIDs, err = repo.GetVcSignalMentionRolesByVcChannelID(ctx, vs.BeforeUpdate.ChannelID)
			if err != nil {
				return nil, err
			}
			for _, beforeMentionRole := range beforeMentionRoleIDs {
				for _, roleID := range vs.Member.Roles {
					if beforeMentionRole.RoleID == roleID {
						beforeMentionText.WriteString("<@&")
						beforeMentionText.WriteString(roleID)
						beforeMentionText.WriteString("> ")
					}
				}
			}

			beforeNgUserIDs, err = repo.GetVcSignalNgUsersByVcChannelIDAllColumn(ctx, vs.BeforeUpdate.ChannelID)
			if err != nil {
				return nil, err
			}
			for _, beforeNgUser := range beforeNgUserIDs {
				if beforeNgUser.UserID == vs.UserID {
					return nil, nil
				}
			}
			beforeNgRoleIDs, err = repo.GetVcSignalNgRolesByVcChannelIDAllColumn(ctx, vs.BeforeUpdate.ChannelID)
			if err != nil {
				return nil, err
			}
			for _, beforeNgRole := range beforeNgRoleIDs {
				for _, roleID := range vs.Member.Roles {
					if beforeNgRole.RoleID == roleID {
						return nil, nil
					}
				}
			}
		}
	}
	//chengeVcChannelFlag := (vs.BeforeUpdate != nil) && (vs.ChannelID != "") && (vs.BeforeUpdate.ChannelID != vs.ChannelID)
	if vs.BeforeUpdate == nil || vs.ChannelID != "" && (!vs.SelfVideo == !vs.SelfStream) && (!vs.BeforeUpdate.SelfVideo == !vs.BeforeUpdate.SelfStream) {
		vcChannel, err := state.Channel(vs.ChannelID)
		if err != nil {
			return nil, err
		}
		membersCount := vcMembersCount(state, vs.GuildID, vs.ChannelID)
		sendText.WriteString(afterMentionText.String())
		sendText.WriteString("現在" + strconv.Itoa(membersCount) + "人 <@" + vs.Member.User.ID + "> が " + vcChannel.Name + "に入室しました。")
		sendMessage, err := s.ChannelMessageSend(afterVcSignalChannel.SendChannelID, sendText.String())
		if err != nil {
			return nil, err
		}
		sendMessages = append(sendMessages, sendMessage)
		sendText.Reset()
		if membersCount == 1 {
			embed = &discordgo.MessageEmbed{
				Title:       "通話開始",
				Description: "<#" + vs.ChannelID + ">",
				Author: &discordgo.MessageEmbedAuthor{
					Name:    vs.Member.User.Username,
					IconURL: vs.Member.AvatarURL("64"),
				},
			}
			sendMessage, err := s.ChannelMessageSendEmbed(afterVcSignalChannel.SendChannelID, embed)
			if err != nil {
				return sendMessages, err
			}
			sendMessages = append(sendMessages, sendMessage)
		}
	}
	if vs.BeforeUpdate != nil && (!vs.BeforeUpdate.SelfVideo == !vs.BeforeUpdate.SelfStream) && (!vs.SelfVideo == !vs.SelfStream) {
		vcChannel, err := state.Channel(vs.BeforeUpdate.ChannelID)
		if err != nil {
			return nil, err
		}
		membersCount := vcMembersCount(state, vs.BeforeUpdate.GuildID, vs.BeforeUpdate.ChannelID)
		sendText.WriteString("現在" + strconv.Itoa(membersCount) + "人 <@" + vs.Member.User.ID + "> が " + vcChannel.Name + "から退室しました。")
		sendMessage, err := s.ChannelMessageSend(afterVcSignalChannel.SendChannelID, sendText.String())
		if err != nil {
			return nil, err
		}
		sendMessages = append(sendMessages, sendMessage)
		sendText.Reset()
		if membersCount == 0 {
			guildMembersCount := guildVcMembersCount(state, vs.BeforeUpdate.GuildID)
			slog.InfoContext(ctx,"","",guildMembersCount)
			if guildMembersCount == 0 {
				embed = &discordgo.MessageEmbed{
					Title: "通話終了",
				}
				sendText.WriteString(beforeMentionText.String())
				sendText.WriteString("通話が終了しました。")
				sendMessage, err := s.ChannelMessageSend(afterVcSignalChannel.SendChannelID, sendText.String())
				if err != nil {
					return nil, err
				}
				sendMessages = append(sendMessages, sendMessage)
				sendText.Reset()
			}
		}
	}
	if (vs.BeforeUpdate != nil && !vs.BeforeUpdate.SelfVideo) && vs.SelfVideo {
		vcChannel, err := state.Channel(vs.BeforeUpdate.ChannelID)
		if err != nil {
			return nil, err
		}
		embed = &discordgo.MessageEmbed{
			Title:       "カメラ配信",
			Description: vs.Member.User.Username + "\n" + "<#" + vs.ChannelID + ">",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    vs.Member.User.Username,
				IconURL: vs.Member.AvatarURL("64"),
			},
		}
		sendText.WriteString(afterMentionText.String())
		sendText.WriteString("<@" + vs.Member.User.ID + "> が" + vcChannel.Name + "カメラ配信を開始しました。")
		sendMessage, err := s.ChannelMessageSend(afterVcSignalChannel.SendChannelID, sendText.String())
		if err != nil {
			return nil, err
		}
		sendMessages = append(sendMessages, sendMessage)
		sendText.Reset()
		sendMessage, err = s.ChannelMessageSendEmbed(afterVcSignalChannel.SendChannelID, embed)
		if err != nil {
			return sendMessages, err
		}
		sendMessages = append(sendMessages, sendMessage)
	}
	if (vs.BeforeUpdate != nil && vs.BeforeUpdate.SelfVideo) && !vs.SelfVideo {
		sendText.WriteString("<@" + vs.Member.User.ID + "> がカメラ配信を終了しました。")
		sendMessage, err := s.ChannelMessageSend(afterVcSignalChannel.SendChannelID, sendText.String())
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
		vcChannel, err := state.Channel(vs.BeforeUpdate.ChannelID)
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
			sendText.WriteString(afterMentionText.String())
			sendText.WriteString("<@" + vs.Member.User.ID + "> が" + vcChannel.Name + "で画面共有を開始しました。")
			sendMessage, err := s.ChannelMessageSend(afterVcSignalChannel.SendChannelID, sendText.String())
			if err != nil {
				return nil, err
			}
			sendMessages = append(sendMessages, sendMessage)
			sendText.Reset()
			sendMessage, err = s.ChannelMessageSendEmbed(afterVcSignalChannel.SendChannelID, embed)
			if err != nil {
				return sendMessages, err
			}
			sendMessages = append(sendMessages, sendMessage)
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
			sendText.WriteString(afterMentionText.String())
			sendText.WriteString("<@" + vs.Member.User.ID + "> が" + vcChannel.Name + "で" + presence.Activities[0].Name + "を配信開始しました。")
			sendMessage, err := s.ChannelMessageSend(afterVcSignalChannel.SendChannelID, sendText.String())
			if err != nil {
				return nil, err
			}
			sendMessages = append(sendMessages, sendMessage)
			sendText.Reset()
			sendMessage, err = s.ChannelMessageSendEmbed(afterVcSignalChannel.SendChannelID, embed)
			if err != nil {
				return sendMessages, err
			}
			sendMessages = append(sendMessages, sendMessage)
		}
	}
	if (vs.BeforeUpdate != nil && vs.BeforeUpdate.SelfStream) && !vs.SelfStream {
		sendText.WriteString("<@" + vs.Member.User.ID + "> が画面共有を終了しました。")
		sendMessage, err := s.ChannelMessageSend(afterVcSignalChannel.SendChannelID, sendText.String())
		if err != nil {
			return nil, err
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
