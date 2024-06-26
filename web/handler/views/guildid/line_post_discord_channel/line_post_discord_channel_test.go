package linepostdiscordchannel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/ctxvalue"
	"github.com/maguro-alternative/remake_bot/web/shared/model"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLinePostDiscordChannelViewHandler(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(cwd))
	})
	require.NoError(t, os.Chdir("../../../../../"))
	t.Run("line_post_discord_channelが正常に表示されること", func(t *testing.T) {
		indexService := &service.IndexService{
			DiscordSession: &discordgo.Session{},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "123",
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
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetLinePostDiscordChannelByChannelIDFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
				return repository.LinePostDiscordChannel{
					Ng: 		true,
					BotMessage: false,
				}, nil
			},
			GetLineNgDiscordMessageTypeByChannelIDFunc: func(ctx context.Context, channelID string) ([]int, error) {
				return []int{}, nil
			},
			GetLineNgDiscordUserIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
				return []string{}, nil
			},
			GetLineNgDiscordRoleIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
				return []string{}, nil
			},
		}
		handler := NewLinePostDiscordChannelViewHandler(indexService, repo)

		mux.HandleFunc("/guilds/{guildId}/line_post_discord_channel", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guilds/123/line_post_discord_channel", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), `<p>Discordアカウント: test</p>`)
		assert.Contains(t, rec.Body.String(), `<img src="https://cdn.discordapp.com/avatars/123/test.webp?size=64" alt="Discordアイコン">`)
		assert.Contains(t, rec.Body.String(), `<p>LINEアカウント: 未ログイン</p>`)

		assert.Contains(t, rec.Body.String(), `<option value=0>デフォルト</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=1>RecipientAdd</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=2>RecipientRemove</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=3>DM通話開始</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=4>チャンネル名変更</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=5>チャンネルアイコン変更</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=6>メッセージピン止め</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=7>サーバー参加</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=8>サーバーブースト</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=9>サーバーレベル1</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=10>サーバーレベル2</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=11>サーバーレベル3</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=12>サーバーフォロー</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=13>サーバーディスカバリー失格メッセージ</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=14>サーバーディスカバリー要件メッセージ</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=15>スレッド作成</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=16>リプライメッセージ</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=17>スラッシュコマンド</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=18>スレッドスタートメッセージ</option>`)
		assert.Contains(t, rec.Body.String(), `<option value=19>コンテンツメニュー</option>`)

		assert.Contains(t, rec.Body.String(), `<select id="ng_users_123[]" name="ng_users_123[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles_123[]" name="ng_roles_123[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users_1234[]" name="ng_users_1234[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles_1234[]" name="ng_roles_1234[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users_12345[]" name="ng_users_12345[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles_12345[]" name="ng_roles_12345[]" multiple>`)
	})

	t.Run("選択されたメッセージタイプが正常に表示されること", func(t *testing.T) {
		indexService := &service.IndexService{
			DiscordSession: &discordgo.Session{},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "123",
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
					ID: "234",
					Name: "test",
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetLinePostDiscordChannelByChannelIDFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
				return repository.LinePostDiscordChannel{
					Ng: 		true,
					BotMessage: false,
				}, nil
			},
			GetLineNgDiscordMessageTypeByChannelIDFunc: func(ctx context.Context, channelID string) ([]int, error) {
				return []int{4}, nil
			},
			GetLineNgDiscordUserIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
				return []string{}, nil
			},
			GetLineNgDiscordRoleIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
				return []string{}, nil
			},
		}
		handler := NewLinePostDiscordChannelViewHandler(indexService, repo)

		mux.HandleFunc("/guilds/{guildId}/line_post_discord_channel", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guilds/123/line_post_discord_channel", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), `<p>Discordアカウント: test</p>`)
		assert.Contains(t, rec.Body.String(), `<img src="https://cdn.discordapp.com/avatars/123/test.webp?size=64" alt="Discordアイコン">`)
		assert.Contains(t, rec.Body.String(), `<p>LINEアカウント: 未ログイン</p>`)

		assert.Contains(t, rec.Body.String(), `<select id="ng_users_123[]" name="ng_users_123[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles_123[]" name="ng_roles_123[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users_1234[]" name="ng_users_1234[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles_1234[]" name="ng_roles_1234[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users_12345[]" name="ng_users_12345[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles_12345[]" name="ng_roles_12345[]" multiple>`)

		assert.Contains(t, rec.Body.String(), `<option value=5 selected>チャンネルアイコン変更</option>`)
	})

	t.Run("選択されたユーザーが正常に表示されること", func(t *testing.T) {
		indexService := &service.IndexService{
			DiscordSession: &discordgo.Session{},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "123",
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
						Username: "test",
					},
				},
			},
			Roles: []*discordgo.Role{
				{
					ID: "234",
					Name: "test",
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetLinePostDiscordChannelByChannelIDFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
				return repository.LinePostDiscordChannel{
					Ng: 		true,
					BotMessage: false,
				}, nil
			},
			GetLineNgDiscordMessageTypeByChannelIDFunc: func(ctx context.Context, channelID string) ([]int, error) {
				return []int{}, nil
			},
			GetLineNgDiscordUserIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
				return []string{"123"}, nil
			},
			GetLineNgDiscordRoleIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
				return []string{}, nil
			},
		}
		handler := NewLinePostDiscordChannelViewHandler(indexService, repo)

		mux.HandleFunc("/guilds/{guildId}/line_post_discord_channel", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guilds/123/line_post_discord_channel", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), `<p>Discordアカウント: test</p>`)
		assert.Contains(t, rec.Body.String(), `<img src="https://cdn.discordapp.com/avatars/123/test.webp?size=64" alt="Discordアイコン">`)
		assert.Contains(t, rec.Body.String(), `<p>LINEアカウント: 未ログイン</p>`)

		assert.Contains(t, rec.Body.String(), `<select id="ng_users_123[]" name="ng_users_123[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles_123[]" name="ng_roles_123[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users_1234[]" name="ng_users_1234[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles_1234[]" name="ng_roles_1234[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users_12345[]" name="ng_users_12345[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles_12345[]" name="ng_roles_12345[]" multiple>`)

		assert.Contains(t, rec.Body.String(), `<option value="123" selected>test</option>`)
	})

	t.Run("選択されたロールが正常に表示されること", func(t *testing.T) {
		indexService := &service.IndexService{
			DiscordSession: &discordgo.Session{},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "123",
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
					ID: "234",
					Name: "test",
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetLinePostDiscordChannelByChannelIDFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
				return repository.LinePostDiscordChannel{
					Ng: 		true,
					BotMessage: false,
				}, nil
			},
			GetLineNgDiscordMessageTypeByChannelIDFunc: func(ctx context.Context, channelID string) ([]int, error) {
				return []int{}, nil
			},
			GetLineNgDiscordUserIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
				return []string{}, nil
			},
			GetLineNgDiscordRoleIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
				return []string{"234"}, nil
			},
		}
		handler := NewLinePostDiscordChannelViewHandler(indexService, repo)

		mux.HandleFunc("/guilds/{guildId}/line_post_discord_channel", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guilds/123/line_post_discord_channel", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), `<p>Discordアカウント: test</p>`)
		assert.Contains(t, rec.Body.String(), `<img src="https://cdn.discordapp.com/avatars/123/test.webp?size=64" alt="Discordアイコン">`)
		assert.Contains(t, rec.Body.String(), `<p>LINEアカウント: 未ログイン</p>`)

		assert.Contains(t, rec.Body.String(), `<select id="ng_users_123[]" name="ng_users_123[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles_123[]" name="ng_roles_123[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users_1234[]" name="ng_users_1234[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles_1234[]" name="ng_roles_1234[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users_12345[]" name="ng_users_12345[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles_12345[]" name="ng_roles_12345[]" multiple>`)

		assert.Contains(t, rec.Body.String(), `<option value="234" selected>test</option>`)
	})


	t.Run("不正なサーバーIDが指定された場合500エラーを出すこと", func(t *testing.T) {
		indexService := &service.IndexService{
			DiscordSession: &discordgo.Session{},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "123",
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
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetLinePostDiscordChannelByChannelIDFunc: func(ctx context.Context, channelID string) (repository.LinePostDiscordChannel, error) {
				return repository.LinePostDiscordChannel{
					Ng: 		true,
					BotMessage: false,
				}, nil
			},
			GetLineNgDiscordMessageTypeByChannelIDFunc: func(ctx context.Context, channelID string) ([]int, error) {
				return []int{}, nil
			},
			GetLineNgDiscordUserIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
				return []string{}, nil
			},
			GetLineNgDiscordRoleIDByChannelIDFunc: func(ctx context.Context, channelID string) ([]string, error) {
				return []string{}, nil
			},
		}
		handler := NewLinePostDiscordChannelViewHandler(indexService, repo)

		mux.HandleFunc("/guilds/{guildId}/line_post_discord_channel", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guilds/111/line_post_discord_channel", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

}

func setCtxValue(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = ctxvalue.ContextWithDiscordPermission(ctx, &model.DiscordPermissionData{
		PermissionCode: 8,
		User: 		 model.DiscordUser{
			ID: "123",
			Username: "test",
			Avatar: "test",
		},
		Permission: "all",
	})
	return r.WithContext(ctx)
}
