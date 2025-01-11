package lineworksbot

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/maguro-alternative/remake_bot/web/handler/api/lineworks_bot/internal"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/bwmarrin/discordgo"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineWorksBotHandler_ServeHTTP(t *testing.T) {
	bodyJson, err := json.Marshal(internal.LineWorksResponses{
		Type: "text",
		Source: struct {
			UserID    string `json:"userId"`
			ChannelId string `json:"channelId"`
			DomainId  int64  `json:"domainId"`
		}{
			UserID:    "test_user_id",
			ChannelId: "test_channel_id",
			DomainId:  1234567890,
		},
		IssuedTime: "2021-01-01T00:00:00Z",
		Content: struct {
			Type               string `json:"type"`
			Text               string `json:"text,omitempty"`
			FileId             string `json:"fileId,omitempty"`
			OriginalContentURL string `json:"originalContentUrl,omitempty"`
		}{
			Type:               "text",
			Text:               "test_text",
		},
	})
	require.NoError(t, err)

	stubClient := mock.NewStubHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"userId": "user123","email": "user@example.com","userName": {"lastName": "Doe","firstName": "John","phoneticLastName": "Doe","phoneticFirstName": "John"},"i18nNames": [{"language": "en","firstName": "John","lastName": "Doe"}],"organizations": [{"domainId": 1,"primary": true,"userExternalKey": "key123","email": "org@example.com","levelId": "level1","levelExternalKey": "levelKey1","levelName": "Level 1","executive": false,"organizationName": "Org Name","orgUnits": [{"orgUnitId": "unit1","orgUnitExternalKey": "unitKey1","orgUnitName": "Unit 1","orgUnitEmail": "unit@example.com","primary": true,"positionId": "pos1","positionExternalKey": "posKey1","positionName": "Position 1","isManager": false,"visible": true,"useTeamFeature": false}]}],"telephone": "123-456-7890","cellPhone": "098-765-4321","location": "Somewhere"}`)),
			Header:     make(http.Header),
		}
	})

	t.Run("MethodがPOST以外の場合、Method Not Allowedが返ること", func(t *testing.T) {
		h := &LineWorksHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/lineworks-bot", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("line_works_botの取得に失敗した場合、Internal Server Errorが返ること", func(t *testing.T) {
		h := &LineWorksHandler{
			indexService: &service.IndexService{
				Client:         stubClient,
				DiscordSession: &discordgo.Session{},
			},
			repo: &repository.RepositoryFuncMock{
				GetAllLineWorksBotsFunc: func(ctx context.Context) ([]*repository.LineWorksBot, error) {
					return nil, os.ErrNotExist
				},
			},
			aesCrypto: &crypto.AES{},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/lineworks-bot", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("line_works_botが存在しない場合、Internal Server Errorが返ること", func(t *testing.T) {
		h := &LineWorksHandler{
			indexService: &service.IndexService{
				Client:         stubClient,
				DiscordSession: &discordgo.Session{},
			},
			repo: &repository.RepositoryFuncMock{
				GetAllLineWorksBotsFunc: func(ctx context.Context) ([]*repository.LineWorksBot, error) {
					return []*repository.LineWorksBot{}, nil
				},
			},
			aesCrypto: &crypto.AES{},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/lineworks-bot", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("正常系", func(t *testing.T) {
		h := &LineWorksHandler{
			indexService: &service.IndexService{
				Client:         stubClient,
				DiscordSession: &mock.SessionMock{
					ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
						return &discordgo.Message{}, nil
					},
				},
			},
			repo: &repository.RepositoryFuncMock{
				GetAllLineWorksBotsFunc: func(ctx context.Context) ([]*repository.LineWorksBot, error) {
					return []*repository.LineWorksBot{
						{
							GuildID:               "1",
							LineWorksBotToken:     pq.ByteaArray{[]byte("test")},
							LineWorksRefreshToken: pq.ByteaArray{[]byte("test")},
							LineWorksGroupID:      pq.ByteaArray{[]byte("test")},
							LineWorksBotID:        pq.ByteaArray{[]byte("test")},
							LineWorksBotSecret:    pq.ByteaArray{[]byte("test")},
							RefreshTokenExpiresAt: pq.NullTime{Valid: true},
							DefaultChannelID:      "1",
							DebugMode:             false,
						},
					}, nil
				},
				GetLineWorksBotIVByGuildIDFunc: func(ctx context.Context, guildID string) (*repository.LineWorksBotIV, error) {
					return &repository.LineWorksBotIV{
						GuildID:               "1",
						LineWorksBotTokenIV:    pq.ByteaArray{[]byte("test")},
						LineWorksRefreshTokenIV: pq.ByteaArray{[]byte("test")},
						LineWorksGroupIDIV:     pq.ByteaArray{[]byte("test")},
						LineWorksBotIDIV:       pq.ByteaArray{[]byte("test")},
						LineWorksBotSecretIV:  pq.ByteaArray{[]byte("test")},
					}, nil
				},
			},
			aesCrypto: &crypto.AESMock{
				EncryptFunc: func(data []byte) (encrypted []byte, iv []byte, err error) {
					return []byte("test"), []byte("test"), nil
				},
				DecryptFunc: func(data []byte, iv []byte) ([]byte, error) {
					return []byte("test"), nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/lineworks-bot", strings.NewReader(string(bodyJson)))
		r.Header.Set("X-WORKS-Signature", "KgInvTpsXJ2kfNsk7NRN3Ky/vrje6RG8FC5ZhVlUOiE=")
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
