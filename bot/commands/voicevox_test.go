package commands


import (
	"net/http"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func Test_VoiceVoxCommand(t *testing.T) {
	repo := &repository.RepositoryFuncMock{}
	session := mock.SessionMock{}
	state := discordgo.NewState()
	client := &http.Client{}
	t.Run("voicevox成功", func(t *testing.T){
		err := VoiceVoxCommand(repo, client).Executor(&session, state, nil, &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "voicevox",
				},
			},
		})
		assert.NoError(t, err)
	})
}

