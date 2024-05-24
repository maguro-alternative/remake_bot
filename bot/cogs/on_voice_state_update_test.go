package cogs

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

func TestVcSignal(t *testing.T) {
	ctx := context.Background()
	userid := "11"
	username := "testuser"
	useravater := "a_"
	//beforeGuildId := "111"
	afterGuildId := "222"
	//beforeChannelId := "1111"
	afterChannelId := "2222"

	// スタブHTTPクライアントを作成
	t.Run("正常系", func(t *testing.T) {
		_, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUsersByVcChannelIDAllColumnFunc: func(ctx context.Context, vcChannelID string) ([]*repository.VcSignalNgUserAllColumn, error) {
					return []*repository.VcSignalNgUserAllColumn{}, nil
				},
				GetVcSignalNgRolesByVcChannelIDAllColumnFunc: func(ctx context.Context, vcChannelID string) ([]*repository.VcSignalNgRoleAllColumn, error) {
					return []*repository.VcSignalNgRoleAllColumn{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					return &repository.VcSignalChannelAllColumn{}, nil
				},
				GetVcSignalMentionUsersByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]*repository.VcSignalMentionUser, error) {
					return []*repository.VcSignalMentionUser{}, nil
				},
				GetVcSignalMentionRolesByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]*repository.VcSignalMentionRole, error) {
					return []*repository.VcSignalMentionRole{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc :func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{}, nil
				},
				ChannelMessageSendEmbedFunc :func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{}, nil
				},
			},
			&discordgo.State{},
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID: afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: &discordgo.User{
							ID: userid,
							Username: username,
							Avatar: useravater,
							Bot: false,
						},
					},
					SelfStream: false,
					SelfVideo: false,
				},
			},
		)
		assert.NoError(t, err)
	})
}
