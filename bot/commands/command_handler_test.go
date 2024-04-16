package commands

import (
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func Test_PingCommandRegister(t *testing.T) {
	repo := &repository.RepositoryFuncMock{}
	session := mock.SessionMock{
		AddHandlerFunc: func(interface{}) func() {
			return func() {}
		},
		ApplicationCommandCreateFunc: func(guildID string, appID string, appCmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (ccmd *discordgo.ApplicationCommand, err error) {
			return appCmd, nil
		},
	}
	t.Run("ping登録成功", func(t *testing.T){
		state := discordgo.NewState()
		state.User = &discordgo.User{
			ID: "1234567890",
		}
		h := newCommandHandler(&session, state, "")
		err := h.commandRegister(PingCommand(repo))
		assert.NoError(t, err)
		c := h.getCommands()
		assert.NotNil(t, c)
		assert.Equal(t, 1, len(c))
		assert.Equal(t, "ping", c[0].Name)
	})
}

func Test_PingCommandRemove(t *testing.T) {
	repo := &repository.RepositoryFuncMock{}
	session := mock.SessionMock{
		AddHandlerFunc: func(interface{}) func() {
			return func() {}
		},
		ApplicationCommandCreateFunc: func(guildID string, appID string, appCmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (ccmd *discordgo.ApplicationCommand, err error) {
			return appCmd, nil
		},
	}
	t.Run("ping削除成功", func(t *testing.T){
		state := discordgo.NewState()
		state.User = &discordgo.User{
			ID: "1234567890",
		}
		h := newCommandHandler(&session, state, "")
		err := h.commandRegister(PingCommand(repo))
		assert.NoError(t, err)
		c := h.getCommands()
		assert.NotNil(t, c)
		assert.Equal(t, 1, len(c))
		assert.Equal(t, "ping", c[0].Name)

		err = h.commandRemove(c[0])
		assert.NoError(t, err)

		c = h.getCommands()
		assert.Len(t, c, 0)
	})
}
