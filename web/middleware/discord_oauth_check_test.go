package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
	"github.com/maguro-alternative/remake_bot/web/shared/session"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscordOAuthCheckMiddleware(t *testing.T) {
	cookieStore := sessions.NewCookieStore([]byte(config.SessionSecret()))
	user := model.DiscordUser{
		ID:       "123",
		Username: "test",
		Avatar:   "test",
	}
	t.Run("DiscordOAuthCheckMiddlewareが'/'で正常に動作すること", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		middleware := DiscordOAuthCheckMiddleware(
			service.IndexService{
				Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`{
							"id": "123456789",
							"username": "test",
							"global_name": "test",
							"avatar": "test",
							"avatar_decoration": "test",
							"discriminator": "1234",
							"public_flags": 0,
							"flags": 0,
							"banner": "test",
							"banner_color": "test",
							"accent_color": "test",
							"locale": "test",
							"mfa_enabled": true,
							"premium_type": 0,
							"email": "test",
							"verified": true,
							"bio": "test",
						}`)),
					}
				}),
				CookieStore: sessions.NewCookieStore([]byte(config.SessionSecret())),
			},
			&repository.RepositoryFuncMock{
				GetPermissionCodeByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) (int64, error) {
					return 0, nil
				},
				GetPermissionUserIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionUserID, error) {
					return nil, nil
				},
				GetPermissionRoleIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionRoleID, error) {
					return nil, nil
				},
			},
			false,
		)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		middleware(handler).ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("'/guilds'で正常に動作すること", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		middleware := DiscordOAuthCheckMiddleware(
			service.IndexService{
				Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`{
							"id": "123456789",
							"username": "test",
							"global_name": "test",
							"avatar": "test",
							"avatar_decoration": "test",
							"discriminator": "1234",
							"public_flags": 0,
							"flags": 0,
							"banner": "test",
							"banner_color": "test",
							"accent_color": "test",
							"locale": "test",
							"mfa_enabled": true,
							"premium_type": 0,
							"email": "test",
							"verified": true,
							"bio": "test",
						}`)),
					}
				}),
				CookieStore: sessions.NewCookieStore([]byte(config.SessionSecret())),
			},
			&repository.RepositoryFuncMock{
				GetPermissionCodeByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) (int64, error) {
					return 0, nil
				},
				GetPermissionUserIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionUserID, error) {
					return nil, nil
				},
				GetPermissionRoleIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionRoleID, error) {
					return nil, nil
				},
			},
			true,
		)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guilds", nil)

		sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
		require.NoError(t, err)
		sessionStore.SetDiscordUser(&user)
		sessionStore.SetDiscordOAuthToken("test")
		sessionStore.SessionSave(r, w)

		defer sessionStore.CleanupDiscordUser()
		defer sessionStore.CleanupDiscordOAuthToken()
		defer sessionStore.SessionSave(r, w)

		middleware(handler).ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("'/guilds'で認証情報がない場合、ログインページにリダイレクトさせること", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		middleware := DiscordOAuthCheckMiddleware(
			service.IndexService{
				Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`{
							"id": "123456789",
							"username": "test",
							"global_name": "test",
							"avatar": "test",
							"avatar_decoration": "test",
							"discriminator": "1234",
							"public_flags": 0,
							"flags": 0,
							"banner": "test",
							"banner_color": "test",
							"accent_color": "test",
							"locale": "test",
							"mfa_enabled": true,
							"premium_type": 0,
							"email": "test",
							"verified": true,
							"bio": "test",
						}`)),
					}
				}),
				CookieStore: cookieStore,
			},
			&repository.RepositoryFuncMock{
				GetPermissionCodeByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) (int64, error) {
					return 0, nil
				},
				GetPermissionUserIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionUserID, error) {
					return nil, nil
				},
				GetPermissionRoleIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionRoleID, error) {
					return nil, nil
				},
			},
			true,
		)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guilds", nil)

		middleware(handler).ServeHTTP(w, r)
		assert.Equal(t, http.StatusFound, w.Code)
	})

	t.Run("'/guild/{guildid}'で正常に動作すること", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		indexService := service.IndexService{
			Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"id": "123456789",
						"username": "test",
						"global_name": "test",
						"avatar": "test",
						"avatar_decoration": "test",
						"discriminator": "1234",
						"public_flags": 0,
						"flags": 0,
						"banner": "test",
						"banner_color": "test",
						"accent_color": "test",
						"locale": "test",
						"mfa_enabled": true,
						"premium_type": 0,
						"email": "test",
						"verified": true,
						"bio": "test",
					}`)),
				}
			}),
			CookieStore: sessions.NewCookieStore([]byte(config.SessionSecret())),
			DiscordSession: &mock.SessionMock{
				UserChannelPermissionsFunc: func(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error) {
					return 0, nil
				},
			},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "111",
			Channels: []*discordgo.Channel{
				{
					ID:       "123",
					Name:     "test",
					Position: 1,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "1234",
					Name:     "test",
					Position: 2,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "12345",
					Name:     "test",
					Position: 3,
					Type:     discordgo.ChannelTypeGuildText,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID: "123",
					},
				},
			},
		})
		require.NoError(t, err)
		middleware := DiscordOAuthCheckMiddleware(
			indexService,
			&repository.RepositoryFuncMock{
				GetPermissionCodeByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) (int64, error) {
					return 0, nil
				},
				GetPermissionUserIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionUserID, error) {
					return nil, nil
				},
				GetPermissionRoleIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionRoleID, error) {
					return nil, nil
				},
			},
			true,
		)

		mux := http.NewServeMux()

		mux.HandleFunc("/guild/{guildId}", middleware(handler).ServeHTTP)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guild/111", nil)

		sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
		require.NoError(t, err)
		sessionStore.SetDiscordUser(&user)
		sessionStore.SetDiscordOAuthToken("test")
		sessionStore.SessionSave(r, w)

		defer sessionStore.CleanupDiscordUser()
		defer sessionStore.CleanupDiscordOAuthToken()
		defer sessionStore.SessionSave(r, w)

		mux.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("'/guild/{guildid}'でサーバーのユーザーではない場合500を返すこと", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		indexService := service.IndexService{
			Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"id": "123456789",
						"username": "test",
						"global_name": "test",
						"avatar": "test",
						"avatar_decoration": "test",
						"discriminator": "1234",
						"public_flags": 0,
						"flags": 0,
						"banner": "test",
						"banner_color": "test",
						"accent_color": "test",
						"locale": "test",
						"mfa_enabled": true,
						"premium_type": 0,
						"email": "test",
						"verified": true,
						"bio": "test",
					}`)),
				}
			}),
			CookieStore: sessions.NewCookieStore([]byte(config.SessionSecret())),
			DiscordSession: &mock.SessionMock{
				UserChannelPermissionsFunc: func(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error) {
					return 0, nil
				},
			},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "111",
			Channels: []*discordgo.Channel{
				{
					ID:       "123",
					Name:     "test",
					Position: 1,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "1234",
					Name:     "test",
					Position: 2,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "12345",
					Name:     "test",
					Position: 3,
					Type:     discordgo.ChannelTypeGuildText,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID: "121",
					},
				},
			},
		})
		require.NoError(t, err)
		middleware := DiscordOAuthCheckMiddleware(
			indexService,
			&repository.RepositoryFuncMock{
				GetPermissionCodeByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) (int64, error) {
					return 0, nil
				},
				GetPermissionUserIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionUserID, error) {
					return nil, nil
				},
				GetPermissionRoleIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionRoleID, error) {
					return nil, nil
				},
			},
			true,
		)

		mux := http.NewServeMux()

		mux.HandleFunc("/guild/{guildId}", middleware(handler).ServeHTTP)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guild/111", nil)

		sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
		require.NoError(t, err)
		sessionStore.SetDiscordUser(&user)
		sessionStore.SetDiscordOAuthToken("test")
		sessionStore.SessionSave(r, w)

		defer sessionStore.CleanupDiscordUser()
		defer sessionStore.CleanupDiscordOAuthToken()
		defer sessionStore.SessionSave(r, w)

		mux.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("'/guild/{guildid}'でguildIdが不正な値の場合500を返すこと", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		indexService := service.IndexService{
			Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"id": "123456789",
						"username": "test",
						"global_name": "test",
						"avatar": "test",
						"avatar_decoration": "test",
						"discriminator": "1234",
						"public_flags": 0,
						"flags": 0,
						"banner": "test",
						"banner_color": "test",
						"accent_color": "test",
						"locale": "test",
						"mfa_enabled": true,
						"premium_type": 0,
						"email": "test",
						"verified": true,
						"bio": "test",
					}`)),
				}
			}),
			CookieStore: sessions.NewCookieStore([]byte(config.SessionSecret())),
			DiscordSession: &mock.SessionMock{
				UserChannelPermissionsFunc: func(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error) {
					return 0, nil
				},
			},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "112",
			Channels: []*discordgo.Channel{
				{
					ID:       "123",
					Name:     "test",
					Position: 1,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "1234",
					Name:     "test",
					Position: 2,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "12345",
					Name:     "test",
					Position: 3,
					Type:     discordgo.ChannelTypeGuildText,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID: "123",
					},
				},
			},
		})
		require.NoError(t, err)
		middleware := DiscordOAuthCheckMiddleware(
			indexService,
			&repository.RepositoryFuncMock{
				GetPermissionCodeByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) (int64, error) {
					return 0, nil
				},
				GetPermissionUserIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionUserID, error) {
					return nil, nil
				},
				GetPermissionRoleIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionRoleID, error) {
					return nil, nil
				},
			},
			true,
		)

		mux := http.NewServeMux()

		mux.HandleFunc("/guild/{guildId}", middleware(handler).ServeHTTP)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guild/111", nil)

		sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
		require.NoError(t, err)
		sessionStore.SetDiscordUser(&user)
		sessionStore.SetDiscordOAuthToken("test")
		sessionStore.SessionSave(r, w)

		defer sessionStore.CleanupDiscordUser()
		defer sessionStore.CleanupDiscordOAuthToken()
		defer sessionStore.SessionSave(r, w)

		mux.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("'/guild/{guildid}/linetoken'で正常に動作すること(管理者権限)", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		indexService := service.IndexService{
			Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"id": "123456789",
						"username": "test",
						"global_name": "test",
						"avatar": "test",
						"avatar_decoration": "test",
						"discriminator": "1234",
						"public_flags": 0,
						"flags": 0,
						"banner": "test",
						"banner_color": "test",
						"accent_color": "test",
						"locale": "test",
						"mfa_enabled": true,
						"premium_type": 0,
						"email": "test",
						"verified": true,
						"bio": "test",
					}`)),
				}
			}),
			CookieStore: sessions.NewCookieStore([]byte(config.SessionSecret())),
			DiscordSession: &mock.SessionMock{
				UserChannelPermissionsFunc: func(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error) {
					return 8, nil
				},
			},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "111",
			Channels: []*discordgo.Channel{
				{
					ID:       "123",
					Name:     "test",
					Position: 1,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "1234",
					Name:     "test",
					Position: 2,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "12345",
					Name:     "test",
					Position: 3,
					Type:     discordgo.ChannelTypeGuildText,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID: "123",
					},
				},
			},
		})
		require.NoError(t, err)
		middleware := DiscordOAuthCheckMiddleware(
			indexService,
			&repository.RepositoryFuncMock{
				GetPermissionCodeByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) (int64, error) {
					return 8, nil
				},
				GetPermissionUserIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionUserID, error) {
					return nil, nil
				},
				GetPermissionRoleIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionRoleID, error) {
					return nil, nil
				},
			},
			true,
		)

		mux := http.NewServeMux()

		mux.HandleFunc("/guild/{guildId}/linetoken", middleware(handler).ServeHTTP)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guild/111/linetoken", nil)

		sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
		require.NoError(t, err)
		sessionStore.SetDiscordUser(&user)
		sessionStore.SetDiscordOAuthToken("test")
		sessionStore.SessionSave(r, w)

		defer sessionStore.CleanupDiscordUser()
		defer sessionStore.CleanupDiscordOAuthToken()
		defer sessionStore.SessionSave(r, w)

		mux.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("'/guild/{guildid}/linetoken'で正常に動作すること(許可メンバー)", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		indexService := service.IndexService{
			Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"id": "123456789",
						"username": "test",
						"global_name": "test",
						"avatar": "test",
						"avatar_decoration": "test",
						"discriminator": "1234",
						"public_flags": 0,
						"flags": 0,
						"banner": "test",
						"banner_color": "test",
						"accent_color": "test",
						"locale": "test",
						"mfa_enabled": true,
						"premium_type": 0,
						"email": "test",
						"verified": true,
						"bio": "test",
					}`)),
				}
			}),
			CookieStore: sessions.NewCookieStore([]byte(config.SessionSecret())),
			DiscordSession: &mock.SessionMock{
				UserChannelPermissionsFunc: func(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error) {
					return 0, nil
				},
			},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "111",
			Channels: []*discordgo.Channel{
				{
					ID:       "123",
					Name:     "test",
					Position: 1,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "1234",
					Name:     "test",
					Position: 2,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "12345",
					Name:     "test",
					Position: 3,
					Type:     discordgo.ChannelTypeGuildText,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID: "123",
					},
					Roles: []string{"123"},
				},
			},
		})
		require.NoError(t, err)
		middleware := DiscordOAuthCheckMiddleware(
			indexService,
			&repository.RepositoryFuncMock{
				GetPermissionCodeByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) (int64, error) {
					return 8, nil
				},
				GetPermissionUserIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionUserID, error) {
					return []repository.PermissionUserID{
						{
							UserID:     "123",
							Permission: "write",
						},
					}, nil
				},
				GetPermissionRoleIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionRoleID, error) {
					return nil, nil
				},
			},
			true,
		)

		mux := http.NewServeMux()

		mux.HandleFunc("/guild/{guildId}/linetoken", middleware(handler).ServeHTTP)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guild/111/linetoken", nil)

		sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
		require.NoError(t, err)
		sessionStore.SetDiscordUser(&user)
		sessionStore.SetDiscordOAuthToken("test")
		sessionStore.SessionSave(r, w)

		defer sessionStore.CleanupDiscordUser()
		defer sessionStore.CleanupDiscordOAuthToken()
		defer sessionStore.SessionSave(r, w)

		mux.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("'/guild/{guildid}/linetoken'で正常に動作すること(許可ロール)", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		indexService := service.IndexService{
			Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"id": "123456789",
						"username": "test",
						"global_name": "test",
						"avatar": "test",
						"avatar_decoration": "test",
						"discriminator": "1234",
						"public_flags": 0,
						"flags": 0,
						"banner": "test",
						"banner_color": "test",
						"accent_color": "test",
						"locale": "test",
						"mfa_enabled": true,
						"premium_type": 0,
						"email": "test",
						"verified": true,
						"bio": "test",
					}`)),
				}
			}),
			CookieStore: sessions.NewCookieStore([]byte(config.SessionSecret())),
			DiscordSession: &mock.SessionMock{
				UserChannelPermissionsFunc: func(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error) {
					return 0, nil
				},
			},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "111",
			Channels: []*discordgo.Channel{
				{
					ID:       "123",
					Name:     "test",
					Position: 1,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "1234",
					Name:     "test",
					Position: 2,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "12345",
					Name:     "test",
					Position: 3,
					Type:     discordgo.ChannelTypeGuildText,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID: "123",
					},
				},
			},
			Roles: []*discordgo.Role{
				{
					ID: "123",
				},
			},
		})
		require.NoError(t, err)
		middleware := DiscordOAuthCheckMiddleware(
			indexService,
			&repository.RepositoryFuncMock{
				GetPermissionCodeByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) (int64, error) {
					return 8, nil
				},
				GetPermissionUserIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionUserID, error) {
					return nil, nil
				},
				GetPermissionRoleIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionRoleID, error) {
					return []repository.PermissionRoleID{
						{
							RoleID:     "123",
							Permission: "write",
						},
					}, nil
				},
			},
			true,
		)

		mux := http.NewServeMux()

		mux.HandleFunc("/guild/{guildId}/linetoken", middleware(handler).ServeHTTP)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guild/111/linetoken", nil)

		sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
		require.NoError(t, err)
		sessionStore.SetDiscordUser(&user)
		sessionStore.SetDiscordOAuthToken("test")
		sessionStore.SessionSave(r, w)

		defer sessionStore.CleanupDiscordUser()
		defer sessionStore.CleanupDiscordOAuthToken()
		defer sessionStore.SessionSave(r, w)

		mux.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("'/guild/{guildid}/linetoken'で権限がない場合403を返すこと", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		indexService := service.IndexService{
			Client: mock.NewStubHttpClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"id": "123456789",
						"username": "test",
						"global_name": "test",
						"avatar": "test",
						"avatar_decoration": "test",
						"discriminator": "1234",
						"public_flags": 0,
						"flags": 0,
						"banner": "test",
						"banner_color": "test",
						"accent_color": "test",
						"locale": "test",
						"mfa_enabled": true,
						"premium_type": 0,
						"email": "test",
						"verified": true,
						"bio": "test",
					}`)),
				}
			}),
			CookieStore: sessions.NewCookieStore([]byte(config.SessionSecret())),
			DiscordSession: &mock.SessionMock{
				UserChannelPermissionsFunc: func(userID string, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error) {
					return 0, nil
				},
			},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "111",
			Channels: []*discordgo.Channel{
				{
					ID:       "123",
					Name:     "test",
					Position: 1,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "1234",
					Name:     "test",
					Position: 2,
					Type:     discordgo.ChannelTypeGuildText,
				},
				{
					ID:       "12345",
					Name:     "test",
					Position: 3,
					Type:     discordgo.ChannelTypeGuildText,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID: "123",
					},
				},
			},
		})
		require.NoError(t, err)
		middleware := DiscordOAuthCheckMiddleware(
			indexService,
			&repository.RepositoryFuncMock{
				GetPermissionCodeByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) (int64, error) {
					return 0, nil
				},
				GetPermissionUserIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionUserID, error) {
					return nil, nil
				},
				GetPermissionRoleIDsByGuildIDAndTypeFunc: func(ctx context.Context, guildID string, permissionType string) ([]repository.PermissionRoleID, error) {
					return nil, nil
				},
			},
			true,
		)

		mux := http.NewServeMux()

		mux.HandleFunc("/guild/{guildId}/linetoken", middleware(handler).ServeHTTP)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/guild/111/linetoken", nil)

		sessionStore, err := session.NewSessionStore(r, cookieStore, config.SessionSecret())
		require.NoError(t, err)
		sessionStore.SetDiscordUser(&user)
		sessionStore.SetDiscordOAuthToken("test")
		sessionStore.SessionSave(r, w)

		defer sessionStore.CleanupDiscordUser()
		defer sessionStore.CleanupDiscordOAuthToken()
		defer sessionStore.SessionSave(r, w)

		mux.ServeHTTP(w, r)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

}
