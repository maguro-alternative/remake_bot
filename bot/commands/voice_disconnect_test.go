package commands

import (
	"testing"

	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func Test_VoiceDisconnectCommand(t *testing.T) {
	session := mock.SessionMock{}
	state := discordgo.NewState()
	state.Guilds = []*discordgo.Guild{
		{
			ID: "1",
		},
	}
	t.Run("voice_disconnect成功", func(t *testing.T) {
		err := VoiceDisconnectCommand(nil, nil).Executor(&session, state, nil, &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				GuildID: "1",
				User: &discordgo.User{
					ID: "1",
				},
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "voice_disconnect",
				},
			},
		})
		assert.NoError(t, err)
	})
}
