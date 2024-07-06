package commands

import (
	"net/http"
	"testing"
	"io"
	"strings"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func Test_VoiceVoxCommand(t *testing.T) {
	repo := &repository.RepositoryFuncMock{}
	session := mock.SessionMock{}
	state := discordgo.NewState()
	stubClient := mock.NewStubHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}
	})
	t.Run("voicevox成功", func(t *testing.T) {
		err := VoiceVoxCommand(repo, stubClient).Executor(&session, state, nil, &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "voicevox",
				},
			},
		})
		assert.NoError(t, err)
	})
}
