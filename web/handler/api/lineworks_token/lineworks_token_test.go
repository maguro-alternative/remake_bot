package lineworkstoken

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lib/pq"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/maguro-alternative/remake_bot/web/handler/api/lineworks_token/internal"
	"github.com/maguro-alternative/remake_bot/web/service"

	"github.com/stretchr/testify/assert"
)

// 実際に使われているものではありません。
const sampleKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn
SgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i
cqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC
PUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR
ap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA
Rdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----
`

func TestLineWorksTokenHandler_ServeHTTP(t *testing.T) {
	stubClient := mock.NewStubHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"access_token":"U1234567890","refresh_token":"R1234567890","token_type":"Bearer","expires_in":3600,"scope":"bot"}`)),
		}
	})

	t.Run("MethodがPOST以外の場合、Method Not Allowedが返ること", func(t *testing.T) {
		h := NewLineWorksTokenHandler(nil, nil, nil)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/lineworks-token", nil)

		h.ServeHTTP(w, r)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("jsonの読み取りに失敗した場合、Bad Requestが返ること", func(t *testing.T) {
		h := NewLineWorksTokenHandler(nil, nil, nil)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/lineworks-token", strings.NewReader(""))

		h.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("正常系", func(t *testing.T) {
		bodyJson, err := json.Marshal(internal.LineWorksTokenJson{
			LineWorksClientID:    "testclientid",
			LineWorksClientSecret: "testclientsecret",
			LineWorksServiceAccount: "testserviceaccount",
			LineWorksPrivateKey: sampleKey,
			LineWorksDomainID: "testdomainid",
			LineWorksAdminID: "testadminid",
			DefaultChannelID: "1",
			DebugMode: false,
		})
		assert.NoError(t, err)
		h := NewLineWorksTokenHandler(
			&service.IndexService{
				Client:         stubClient,
				DiscordSession: &discordgo.Session{},
			},
			&repository.RepositoryFuncMock{
				GetLineWorksBotByGuildIDFunc: func(ctx context.Context, guildID string) (*repository.LineWorksBot, error) {
					return &repository.LineWorksBot{
						GuildID: "1",
						LineWorksBotToken: pq.ByteaArray{[]byte("test")},
						LineWorksRefreshToken: pq.ByteaArray{[]byte("test")},
						LineWorksGroupID: pq.ByteaArray{[]byte("test")},
						LineWorksBotID: pq.ByteaArray{[]byte("test")},
						LineWorksBotSecret: pq.ByteaArray{[]byte("test")},
						RefreshTokenExpiresAt: pq.NullTime{Time: time.Now(), Valid: true},
						DefaultChannelID: "1",
						DebugMode: false,
					}, nil
				},
				GetLineWorksBotIVByGuildIDFunc: func(ctx context.Context, guildID string) (*repository.LineWorksBotIV, error) {
					return &repository.LineWorksBotIV{
						GuildID: "1",
						LineWorksBotTokenIV: pq.ByteaArray{[]byte("test")},
						LineWorksRefreshTokenIV: pq.ByteaArray{[]byte("test")},
						LineWorksGroupIDIV: pq.ByteaArray{[]byte("test")},
						LineWorksBotIDIV: pq.ByteaArray{[]byte("test")},
					}, nil
				},
				GetLineWorksBotInfoByGuildIDFunc: func(ctx context.Context, guildID string) (*repository.LineWorksBotInfo, error) {
					return &repository.LineWorksBotInfo{
						GuildID: "1",
						LineWorksClientID: pq.ByteaArray{[]byte("test")},
						LineWorksClientSecret: pq.ByteaArray{[]byte("test")},
						LineWorksServiceAccount: pq.ByteaArray{[]byte("test")},
						LineWorksPrivateKey: pq.ByteaArray{[]byte("test")},
						LineWorksDomainID: pq.ByteaArray{[]byte("test")},
						LineWorksAdminID: pq.ByteaArray{[]byte("test")},
					}, nil
				},
				GetLineWorksBotInfoIVByGuildIDFunc: func(ctx context.Context, guildID string) (*repository.LineWorksBotInfoIV, error) {
					return &repository.LineWorksBotInfoIV{
						GuildID: "1",
						LineWorksClientIDIV: pq.ByteaArray{[]byte("test")},
						LineWorksClientSecretIV: pq.ByteaArray{[]byte("test")},
						LineWorksServiceAccountIV: pq.ByteaArray{[]byte("test")},
						LineWorksPrivateKeyIV: pq.ByteaArray{[]byte("test")},
						LineWorksDomainIDIV: pq.ByteaArray{[]byte("test")},
						LineWorksAdminIDIV: pq.ByteaArray{[]byte("test")},
					}, nil
				},
				UpdateLineWorksBotFunc: func(ctx context.Context, lineWorksBot *repository.LineWorksBot) error {
					return nil
				},
				UpdateLineWorksBotIVFunc: func(ctx context.Context, lineWorksBotIV *repository.LineWorksBotIV) error {
					return nil
				},
				UpdateLineWorksBotInfoFunc: func(ctx context.Context, lineWorksBotInfo *repository.LineWorksBotInfo) error {
					return nil
				},
				UpdateLineWorksBotInfoIVFunc: func(ctx context.Context, lineWorksBotInfoIV *repository.LineWorksBotInfoIV) error {
					return nil
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (encrypted []byte, iv []byte, err error) {
					return []byte("test"), []byte("test"), nil
				},
				DecryptFunc: func(data []byte, iv []byte) ([]byte, error) {
					return []byte("test"), nil
				},
			},
		)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/lineworks-token", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
