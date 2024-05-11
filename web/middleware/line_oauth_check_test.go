package middleware

import (
	"context"
	"encoding/gob"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
	"github.com/maguro-alternative/remake_bot/web/shared/session"

	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineOAuthCheckMiddleware(t *testing.T) {
	cookieStore := sessions.NewCookieStore([]byte(config.SessionSecret()))
	// セッションに保存する構造体の型を登録
	// これがない場合、エラーが発生する
	gob.Register(&model.LineIdTokenUser{})
	gob.Register(&model.LineOAuthSession{})
	t.Run("LineOAuthCheckMiddlewareが正常に動作すること(ルートページのため認証なしでも可)", func(t *testing.T) {
		middlewareStartFixture := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user := &model.LineIdTokenUser{
					Iss:     "https://access.line.me",
					Sub:     "U123456789abcdef123456789abcdef12",
					Aud:     "1234567890",
					Exp:     1504169092,
					Iat:     1504168492,
					Nonce:   "0987654321fedcba987654321fedcba9",
					Amr:     []string{"pwd"},
					Name:    "test",
					Picture: "test",
					Email:   "test",
				}
				sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
				require.NoError(t, err)
				sessionStore.SetLineUser(user)
				sessionStore.SetLineOAuthToken("test")
				sessionStore.SetGuildID("111")
				err = sessionStore.SessionSave(r, w)
				require.NoError(t, err)
				h.ServeHTTP(w, r)
			})
		}
		middlewareEndFixture := func(r *http.Request) {
			sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
			require.NoError(t, err)
			sessionStore.CleanupLineUser()
			sessionStore.CleanupLineOAuthToken()
			sessionStore.CleanupGuildID()
		}
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		middleware := LineOAuthCheckMiddleware(
			service.IndexService{
				Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`{
							"userId": "U123456789",
							"displayName": "test",
							"pictureUrl": "test",
							"statusMessage": "test"
						}`)),
					}
				}),
			},
			&repository.RepositoryFuncMock{
				GetLineBotNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return repository.LineBotNotClient{
						LineNotifyToken:  pq.ByteaArray{[]byte("lineNotifyStr")},
						LineBotToken:     pq.ByteaArray{[]byte("lineBotStr")},
						LineBotSecret:    pq.ByteaArray{[]byte("lineBotSecretStr")},
						LineGroupID:      pq.ByteaArray{[]byte("lineGroupStr")},
					}, nil
				},
				GetLineBotIvNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return repository.LineBotIvNotClient{
						LineNotifyTokenIv:  pq.ByteaArray{[]byte("decodeNotifyToken")},
						LineBotTokenIv:     pq.ByteaArray{[]byte("decodeBotToken")},
						LineBotSecretIv:    pq.ByteaArray{[]byte("decodeBotSecret")},
						LineGroupIDIv:      pq.ByteaArray{[]byte("decodeGroupID")},
					}, nil
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
				},
			},
			false,
		)(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer test")
		w := httptest.NewRecorder()
		middlewareStartFixture(middleware).ServeHTTP(w, req)
		middlewareEndFixture(req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("認証ありのページで正常に動作すること", func(t *testing.T) {
		user := &model.LineIdTokenUser{
			Iss:     "https://access.line.me",
			Sub:     "U123456789abcdef123456789abcdef12",
			Aud:     "1234567890",
			Exp:     1504169092,
			Iat:     1504168492,
			Nonce:   "0987654321fedcba987654321fedcba9",
			Amr:     []string{"pwd"},
			Name:    "test",
			Picture: "test",
			Email:   "test",
		}
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		middleware := LineOAuthCheckMiddleware(
			service.IndexService{
				Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`{
							"userId": "U123456789",
							"displayName": "test",
							"pictureUrl": "test",
							"statusMessage": "test"
						}`)),
					}
				}),
				CookieStore: cookieStore,
			},
			&repository.RepositoryFuncMock{
				GetLineBotNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return repository.LineBotNotClient{
						LineNotifyToken:  pq.ByteaArray{[]byte("lineNotifyStr")},
						LineBotToken:     pq.ByteaArray{[]byte("lineBotStr")},
						LineBotSecret:    pq.ByteaArray{[]byte("lineBotSecretStr")},
						LineGroupID:      pq.ByteaArray{[]byte("lineGroupStr")},
					}, nil
				},
				GetLineBotIvNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return repository.LineBotIvNotClient{
						LineNotifyTokenIv:  pq.ByteaArray{[]byte("decodeNotifyToken")},
						LineBotTokenIv:     pq.ByteaArray{[]byte("decodeBotToken")},
						LineBotSecretIv:    pq.ByteaArray{[]byte("decodeBotSecret")},
						LineGroupIDIv:      pq.ByteaArray{[]byte("decodeGroupID")},
					}, nil
				},
			},
			&crypto.AESMock{
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
				},
			},
			true,
		)(handler)

		mux := http.NewServeMux()

		req := httptest.NewRequest(http.MethodGet, "/group/111", nil)
		req.Header.Set("Authorization", "Bearer test")
		w := httptest.NewRecorder()

		sessionStore, err := session.NewSessionStore(req, cookieStore, config.SessionSecret())
		require.NoError(t, err)
		sessionStore.SetLineUser(user)
		sessionStore.SetLineOAuthToken("test")
		sessionStore.SetGuildID("111")
		err = sessionStore.SessionSave(req, w)
		require.NoError(t, err)

		mux.HandleFunc("/group/{guildId}", middleware.ServeHTTP)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("認証情報がない場合リダイレクトすること", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		middleware := LineOAuthCheckMiddleware(
			service.IndexService{
				Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`{
							"userId": "U123456789",
							"displayName": "test",
							"pictureUrl": "test",
							"statusMessage": "test"
						}`)),
					}
				}),
				CookieStore: cookieStore,
			},
			&repository.RepositoryFuncMock{
				GetLineBotNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return repository.LineBotNotClient{
						LineNotifyToken:  pq.ByteaArray{[]byte("lineNotifyStr")},
						LineBotToken:     pq.ByteaArray{[]byte("lineBotStr")},
						LineBotSecret:    pq.ByteaArray{[]byte("lineBotSecretStr")},
						LineGroupID:      pq.ByteaArray{[]byte("lineGroupStr")},
					}, nil
				},
				GetLineBotIvNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return repository.LineBotIvNotClient{
						LineNotifyTokenIv:  pq.ByteaArray{[]byte("decodeNotifyToken")},
						LineBotTokenIv:     pq.ByteaArray{[]byte("decodeBotToken")},
						LineBotSecretIv:    pq.ByteaArray{[]byte("decodeBotSecret")},
						LineGroupIDIv:      pq.ByteaArray{[]byte("decodeGroupID")},
					}, nil
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
				},
			},
			true,
		)(handler)

		mux := http.NewServeMux()

		req := httptest.NewRequest(http.MethodGet, "/group/111", nil)
		req.Header.Set("Authorization", "Bearer test")
		w := httptest.NewRecorder()

		mux.HandleFunc("/group/{guildId}", middleware.ServeHTTP)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
	})

	t.Run("グループに所属していない場合へージを表示しないこと", func(t *testing.T) {
		user := &model.LineIdTokenUser{
			Iss:     "https://access.line.me",
			Sub:     "U123456789abcdef123456789abcdef12",
			Aud:     "1234567890",
			Exp:     1504169092,
			Iat:     1504168492,
			Nonce:   "0987654321fedcba987654321fedcba9",
			Amr:     []string{"pwd"},
			Name:    "test",
			Picture: "test",
			Email:   "test",
		}
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		middleware := LineOAuthCheckMiddleware(
			service.IndexService{
				Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body: io.NopCloser(strings.NewReader(`{
							"message": "Not found"
						}`)),
					}
				}),
				CookieStore: cookieStore,
			},
			&repository.RepositoryFuncMock{
				GetLineBotNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotNotClient, error) {
					return repository.LineBotNotClient{
						LineNotifyToken:  pq.ByteaArray{[]byte("lineNotifyStr")},
						LineBotToken:     pq.ByteaArray{[]byte("lineBotStr")},
						LineBotSecret:    pq.ByteaArray{[]byte("lineBotSecretStr")},
						LineGroupID:      pq.ByteaArray{[]byte("lineGroupStr")},
					}, nil
				},
				GetLineBotIvNotClientByGuildIDFunc: func(ctx context.Context, guildID string) (repository.LineBotIvNotClient, error) {
					return repository.LineBotIvNotClient{
						LineNotifyTokenIv:  pq.ByteaArray{[]byte("decodeNotifyToken")},
						LineBotTokenIv:     pq.ByteaArray{[]byte("decodeBotToken")},
						LineBotSecretIv:    pq.ByteaArray{[]byte("decodeBotSecret")},
						LineGroupIDIv:      pq.ByteaArray{[]byte("decodeGroupID")},
					}, nil
				},
			},
			&crypto.AESMock{
				EncryptFunc: func(data []byte) (iv []byte, encrypted []byte, err error) {
					return nil, nil, nil
				},
				DecryptFunc: func(data []byte, iv []byte) (decrypted []byte, err error) {
					return nil, nil
				},
			},
			true,
		)(handler)

		mux := http.NewServeMux()

		req := httptest.NewRequest(http.MethodGet, "/group/111", nil)
		req.Header.Set("Authorization", "Bearer test")
		w := httptest.NewRecorder()

		sessionStore, err := session.NewSessionStore(req, cookieStore, config.SessionSecret())
		require.NoError(t, err)
		sessionStore.SetLineUser(user)
		sessionStore.SetLineOAuthToken("test")
		sessionStore.SetGuildID("111")
		err = sessionStore.SessionSave(req, w)
		require.NoError(t, err)

		mux.HandleFunc("/group/{guildId}", middleware.ServeHTTP)
		mux.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

}
