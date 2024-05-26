package cogs

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVcSignal(t *testing.T) {
	ctx := context.Background()
	userid := "11"
	username := "testuser"
	useravater := "a_"
	beforeGuildId := "111"
	afterGuildId := "222"
	beforeChannelId := "1111"
	afterChannelId := "2222"
	beforeSendChannelId := "11111"
	afterSendChannelId := "22222"

	discordState := discordgo.NewState()
	err := discordState.GuildAdd(&discordgo.Guild{
		ID: afterGuildId,
		Channels: []*discordgo.Channel{
			{
				ID:       afterChannelId,
				Name:     "after_test_vc",
				Position: 1,
				Type:     discordgo.ChannelTypeGuildVoice,
			},
			{
				ID:       afterSendChannelId,
				Name:     "after_test_text",
				Position: 2,
				Type:     discordgo.ChannelTypeGuildText,
			},
		},
		VoiceStates: []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: &discordgo.User{
							ID:       userid,
							Username: username,
							Avatar:   useravater,
							Bot:      false,
						},
					},
					SelfStream: false,
					SelfVideo:  false,
			},
		},
	})
	require.NoError(t, err)

	err = discordState.GuildAdd(&discordgo.Guild{
		ID: beforeGuildId,
		Channels: []*discordgo.Channel{
			{
				ID:       beforeChannelId,
				Name:     "before_test_vc",
				Position: 1,
				Type:     discordgo.ChannelTypeGuildVoice,
			},
			{
				ID:       beforeSendChannelId,
				Name:     "vefore_test_text",
				Position: 2,
				Type:     discordgo.ChannelTypeGuildText,
			},
		},
	})
	require.NoError(t, err)

	t.Run("正常系", func(t *testing.T) {
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUsersByVcChannelIDAllColumnFunc: func(ctx context.Context, vcChannelID string) ([]*repository.VcSignalNgUserAllColumn, error) {
					return []*repository.VcSignalNgUserAllColumn{}, nil
				},
				GetVcSignalNgRolesByVcChannelIDAllColumnFunc: func(ctx context.Context, vcChannelID string) ([]*repository.VcSignalNgRoleAllColumn, error) {
					return []*repository.VcSignalNgRoleAllColumn{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:   vcChannelID,
						GuildID:       afterGuildId,
						SendSignal:    true,
						SendChannelID: "123",
					}, nil
				},
				GetVcSignalMentionUsersByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]*repository.VcSignalMentionUser, error) {
					return []*repository.VcSignalMentionUser{}, nil
				},
				GetVcSignalMentionRolesByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]*repository.VcSignalMentionRole, error) {
					return []*repository.VcSignalMentionRole{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: &discordgo.User{
							ID:       userid,
							Username: username,
							Avatar:   useravater,
							Bot:      false,
						},
					},
					SelfStream: false,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 2)
	})
}
