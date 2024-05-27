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
	testUser := &discordgo.User{
		ID:       "11",
		Username: "testuser",
		Avatar:   "a_",
		Bot:      false,
	}
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

	t.Run("正常系(通話開始)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: false,
				SelfVideo:  false,
			},
		}
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
						VcChannelID:     vcChannelID,
						GuildID:         afterGuildId,
						SendSignal:      true,
						SendChannelID:   afterSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
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
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 2)
		assert.Equal(t, messages[0].Content, "現在1人 <@11> が after_test_vcに入室しました。")
		assert.Equal(t, messages[1].Embeds[0].Title, "通話開始")
		assert.Equal(t, messages[1].Embeds[0].Description, "<#2222>")
		assert.Equal(t, messages[1].Embeds[0].Author.Name, "testuser")
		assert.Equal(t, messages[1].Embeds[0].Author.IconURL, "https://cdn.discordapp.com/avatars/11/a_.gif?size=64")
	})

	t.Run("正常系(通話終了)", func(t *testing.T) {
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
						VcChannelID:     vcChannelID,
						GuildID:         beforeGuildId,
						SendSignal:      true,
						SendChannelID:   beforeSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
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
					GuildID:   "",
					ChannelID: "",
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   beforeGuildId,
					ChannelID: beforeChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 3)
		assert.Equal(t, messages[0].Content, "現在0人 <@11> が before_test_vcから退室しました。")
		assert.Equal(t, messages[1].Content, "通話が終了しました。")
		assert.Equal(t, messages[2].Embeds[0].Title, "通話終了")
	})

	t.Run("正常系(ボイスチャンネル移動で移動前サーバー通話終了と移動先サーバー通話開始)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: false,
				SelfVideo:  false,
			},
		}
		discordState.Guilds[1].VoiceStates = []*discordgo.VoiceState{}
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
					if vcChannelID == afterChannelId {
						return &repository.VcSignalChannelAllColumn{
							VcChannelID:     vcChannelID,
							GuildID:         afterGuildId,
							SendSignal:      true,
							SendChannelID:   afterSendChannelId,
							JoinBot:         false,
							EveryoneMention: false,
						}, nil
					}
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         beforeGuildId,
						SendSignal:      true,
						SendChannelID:   beforeSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
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
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   beforeGuildId,
					ChannelID: beforeChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 5)
		assert.Equal(t, messages[0].Content, "現在1人 <@11> が after_test_vcに入室しました。")
		assert.Equal(t, messages[1].Embeds[0].Title, "通話開始")
		assert.Equal(t, messages[1].Embeds[0].Description, "<#2222>")
		assert.Equal(t, messages[1].Embeds[0].Author.Name, "testuser")
		assert.Equal(t, messages[1].Embeds[0].Author.IconURL, "https://cdn.discordapp.com/avatars/11/a_.gif?size=64")
		assert.Equal(t, messages[2].Content, "現在0人 <@11> が before_test_vcから退室しました。")
		assert.Equal(t, messages[3].Content, "通話が終了しました。")
		assert.Equal(t, messages[4].Embeds[0].Title, "通話終了")
	})
}
