package commands

import (
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func Test_PingCommand(t *testing.T) {
	repo := &repository.RepositoryFuncMock{}
	session := mock.SessionMock{}
	state := discordgo.NewState()
	t.Run("ping成功", func(t *testing.T){
		err := PingCommand(repo).Executor(&session, state, nil, &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "ping",
				},
			},
		})
		assert.NoError(t, err)
	})
}
