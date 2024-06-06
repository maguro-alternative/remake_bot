package vcsignal

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
	t.Run("正常に表示される", func(t *testing.T) {
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
				{
					ID:       "123456",
					Name:     "test",
					Position: 4,
					Type:     discordgo.ChannelTypeGuildVoice,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID:       "123",
						Username: "test",
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
				if vcChannelID == "123456" {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     "123456",
						GuildID:         "123",
						SendSignal:      false,
						SendChannelID:   "123456",
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				}
				return &repository.VcSignalChannelAllColumn{}, nil
			},
			InsertVcSignalChannelFunc: func(ctx context.Context, vcChannelID string, guildID string, sendChannelID string) error {
				return nil
			},
			GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
		}

		handler := NewVcSignalViewHandler(indexService, repo)

		mux.HandleFunc("/guild/{guildId}/vc-signal", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123/vc-signal", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), "<h1>ボイスチャンネルの入退出通知</h1>")
		assert.Contains(t, rec.Body.String(), `<label for="sendSignal123456">通知を送信する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="sendSignal123456" name="sendSignal123456"  />`)

		assert.Contains(t, rec.Body.String(), `<label for="joinBot123456">Botの入退出を通知する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="joinBot123456" name="joinBot123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="everyoneMention123456">通知に@everyoneメンションをつける</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="everyoneMention123456" name="everyoneMention123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="defaultChannel123456">送信先チャンネル</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="defaultChannel123456" name="defaultChannelId123456" >`)
		assert.Contains(t, rec.Body.String(), `<option value="123">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="1234">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="12345">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="123456" selected>:🔊:test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)

		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgUserIds123456[]">NGユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users123456[]" name="vcSignalNgUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgRoleIds123456[]">NGロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles123456[]" name="vcSignalNgRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionUserIds123456[]">メンションユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_users123456[]" name="vcSignalMentionUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionRoleIds123456[]">メンションロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_roles123456[]" name="vcSignalMentionRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<button type="submit" class="btn btn-primary">送信</button>`)
	})

	t.Run("正常に表示される(roleアリ)", func(t *testing.T) {
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
				{
					ID:       "123456",
					Name:     "test",
					Position: 4,
					Type:     discordgo.ChannelTypeGuildVoice,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID:       "123",
						Username: "test",
					},
				},
			},
			Roles: []*discordgo.Role{
				{
					ID:   "1235",
					Name: "test",
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
				if vcChannelID == "123456" {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     "123456",
						GuildID:         "123",
						SendSignal:      false,
						SendChannelID:   "123456",
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				}
				return &repository.VcSignalChannelAllColumn{}, nil
			},
			InsertVcSignalChannelFunc: func(ctx context.Context, vcChannelID string, guildID string, sendChannelID string) error {
				return nil
			},
			GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
		}

		handler := NewVcSignalViewHandler(indexService, repo)

		mux.HandleFunc("/guild/{guildId}/vc-signal", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123/vc-signal", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), "<h1>ボイスチャンネルの入退出通知</h1>")
		assert.Contains(t, rec.Body.String(), `<label for="sendSignal123456">通知を送信する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="sendSignal123456" name="sendSignal123456"  />`)

		assert.Contains(t, rec.Body.String(), `<label for="joinBot123456">Botの入退出を通知する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="joinBot123456" name="joinBot123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="everyoneMention123456">通知に@everyoneメンションをつける</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="everyoneMention123456" name="everyoneMention123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="defaultChannel123456">送信先チャンネル</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="defaultChannel123456" name="defaultChannelId123456" >`)
		assert.Contains(t, rec.Body.String(), `<option value="123">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="1234">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="12345">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="123456" selected>:🔊:test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)

		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgUserIds123456[]">NGユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users123456[]" name="vcSignalNgUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgRoleIds123456[]">NGロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles123456[]" name="vcSignalNgRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionUserIds123456[]">メンションユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_users123456[]" name="vcSignalMentionUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionRoleIds123456[]">メンションロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_roles123456[]" name="vcSignalMentionRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<button type="submit" class="btn btn-primary">送信</button>`)
	})

	t.Run("正常に表示される(ngUser指定アリ)", func(t *testing.T) {
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
				{
					ID:       "123456",
					Name:     "test",
					Position: 4,
					Type:     discordgo.ChannelTypeGuildVoice,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID:       "123",
						Username: "test",
					},
				},
			},
			Roles: []*discordgo.Role{
				{
					ID:   "1235",
					Name: "test",
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
				if vcChannelID == "123456" {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     "123456",
						GuildID:         "123",
						SendSignal:      false,
						SendChannelID:   "123456",
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				}
				return &repository.VcSignalChannelAllColumn{}, nil
			},
			InsertVcSignalChannelFunc: func(ctx context.Context, vcChannelID string, guildID string, sendChannelID string) error {
				return nil
			},
			GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				if vcChannelID == "123456" {
					return []string{"123"}, nil
				}
				return []string{}, nil
			},
			GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
		}

		handler := NewVcSignalViewHandler(indexService, repo)

		mux.HandleFunc("/guild/{guildId}/vc-signal", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123/vc-signal", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), "<h1>ボイスチャンネルの入退出通知</h1>")
		assert.Contains(t, rec.Body.String(), `<label for="sendSignal123456">通知を送信する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="sendSignal123456" name="sendSignal123456"  />`)

		assert.Contains(t, rec.Body.String(), `<label for="joinBot123456">Botの入退出を通知する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="joinBot123456" name="joinBot123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="everyoneMention123456">通知に@everyoneメンションをつける</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="everyoneMention123456" name="everyoneMention123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="defaultChannel123456">送信先チャンネル</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="defaultChannel123456" name="defaultChannelId123456" >`)
		assert.Contains(t, rec.Body.String(), `<option value="123">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="1234">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="12345">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="123456" selected>:🔊:test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)

		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgUserIds123456[]">NGユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users123456[]" name="vcSignalNgUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123" selected>test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgRoleIds123456[]">NGロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles123456[]" name="vcSignalNgRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionUserIds123456[]">メンションユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_users123456[]" name="vcSignalMentionUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionRoleIds123456[]">メンションロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_roles123456[]" name="vcSignalMentionRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<button type="submit" class="btn btn-primary">送信</button>`)
	})

	t.Run("正常に表示される(ngRole指定アリ)", func(t *testing.T) {
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
				{
					ID:       "123456",
					Name:     "test",
					Position: 4,
					Type:     discordgo.ChannelTypeGuildVoice,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID:       "123",
						Username: "test",
					},
				},
			},
			Roles: []*discordgo.Role{
				{
					ID:   "1235",
					Name: "test",
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
				if vcChannelID == "123456" {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     "123456",
						GuildID:         "123",
						SendSignal:      false,
						SendChannelID:   "123456",
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				}
				return &repository.VcSignalChannelAllColumn{}, nil
			},
			InsertVcSignalChannelFunc: func(ctx context.Context, vcChannelID string, guildID string, sendChannelID string) error {
				return nil
			},
			GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				if vcChannelID == "123456" {
					return []string{"1235"}, nil
				}
				return []string{}, nil
			},
			GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
		}

		handler := NewVcSignalViewHandler(indexService, repo)

		mux.HandleFunc("/guild/{guildId}/vc-signal", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123/vc-signal", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), "<h1>ボイスチャンネルの入退出通知</h1>")
		assert.Contains(t, rec.Body.String(), `<label for="sendSignal123456">通知を送信する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="sendSignal123456" name="sendSignal123456"  />`)

		assert.Contains(t, rec.Body.String(), `<label for="joinBot123456">Botの入退出を通知する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="joinBot123456" name="joinBot123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="everyoneMention123456">通知に@everyoneメンションをつける</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="everyoneMention123456" name="everyoneMention123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="defaultChannel123456">送信先チャンネル</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="defaultChannel123456" name="defaultChannelId123456" >`)
		assert.Contains(t, rec.Body.String(), `<option value="123">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="1234">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="12345">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="123456" selected>:🔊:test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)

		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgUserIds123456[]">NGユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users123456[]" name="vcSignalNgUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgRoleIds123456[]">NGロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles123456[]" name="vcSignalNgRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235" selected>test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionUserIds123456[]">メンションユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_users123456[]" name="vcSignalMentionUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionRoleIds123456[]">メンションロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_roles123456[]" name="vcSignalMentionRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<button type="submit" class="btn btn-primary">送信</button>`)
	})

	t.Run("正常に表示される(mentionUser指定アリ)", func(t *testing.T) {
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
				{
					ID:       "123456",
					Name:     "test",
					Position: 4,
					Type:     discordgo.ChannelTypeGuildVoice,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID:       "123",
						Username: "test",
					},
				},
			},
			Roles: []*discordgo.Role{
				{
					ID:   "1235",
					Name: "test",
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
				if vcChannelID == "123456" {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     "123456",
						GuildID:         "123",
						SendSignal:      false,
						SendChannelID:   "123456",
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				}
				return &repository.VcSignalChannelAllColumn{}, nil
			},
			InsertVcSignalChannelFunc: func(ctx context.Context, vcChannelID string, guildID string, sendChannelID string) error {
				return nil
			},
			GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				if vcChannelID == "123456" {
					return []string{"123"}, nil
				}
				return []string{}, nil
			},
			GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
		}

		handler := NewVcSignalViewHandler(indexService, repo)

		mux.HandleFunc("/guild/{guildId}/vc-signal", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123/vc-signal", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), "<h1>ボイスチャンネルの入退出通知</h1>")
		assert.Contains(t, rec.Body.String(), `<label for="sendSignal123456">通知を送信する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="sendSignal123456" name="sendSignal123456"  />`)

		assert.Contains(t, rec.Body.String(), `<label for="joinBot123456">Botの入退出を通知する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="joinBot123456" name="joinBot123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="everyoneMention123456">通知に@everyoneメンションをつける</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="everyoneMention123456" name="everyoneMention123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="defaultChannel123456">送信先チャンネル</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="defaultChannel123456" name="defaultChannelId123456" >`)
		assert.Contains(t, rec.Body.String(), `<option value="123">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="1234">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="12345">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="123456" selected>:🔊:test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)

		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgUserIds123456[]">NGユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users123456[]" name="vcSignalNgUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgRoleIds123456[]">NGロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles123456[]" name="vcSignalNgRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionUserIds123456[]">メンションユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_users123456[]" name="vcSignalMentionUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123" selected>test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionRoleIds123456[]">メンションロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_roles123456[]" name="vcSignalMentionRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<button type="submit" class="btn btn-primary">送信</button>`)
	})

	t.Run("正常に表示される(mentionRole指定アリ)", func(t *testing.T) {
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
				{
					ID:       "123456",
					Name:     "test",
					Position: 4,
					Type:     discordgo.ChannelTypeGuildVoice,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID:       "123",
						Username: "test",
					},
				},
			},
			Roles: []*discordgo.Role{
				{
					ID:   "1235",
					Name: "test",
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
				if vcChannelID == "123456" {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     "123456",
						GuildID:         "123",
						SendSignal:      false,
						SendChannelID:   "123456",
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				}
				return &repository.VcSignalChannelAllColumn{}, nil
			},
			InsertVcSignalChannelFunc: func(ctx context.Context, vcChannelID string, guildID string, sendChannelID string) error {
				return nil
			},
			GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				if vcChannelID == "123456" {
					return []string{"1235"}, nil
				}
				return []string{}, nil
			},
		}

		handler := NewVcSignalViewHandler(indexService, repo)

		mux.HandleFunc("/guild/{guildId}/vc-signal", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123/vc-signal", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), "<h1>ボイスチャンネルの入退出通知</h1>")
		assert.Contains(t, rec.Body.String(), `<label for="sendSignal123456">通知を送信する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="sendSignal123456" name="sendSignal123456"  />`)

		assert.Contains(t, rec.Body.String(), `<label for="joinBot123456">Botの入退出を通知する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="joinBot123456" name="joinBot123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="everyoneMention123456">通知に@everyoneメンションをつける</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="everyoneMention123456" name="everyoneMention123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="defaultChannel123456">送信先チャンネル</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="defaultChannel123456" name="defaultChannelId123456" >`)
		assert.Contains(t, rec.Body.String(), `<option value="123">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="1234">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="12345">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="123456" selected>:🔊:test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)

		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgUserIds123456[]">NGユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users123456[]" name="vcSignalNgUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgRoleIds123456[]">NGロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles123456[]" name="vcSignalNgRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionUserIds123456[]">メンションユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_users123456[]" name="vcSignalMentionUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionRoleIds123456[]">メンションロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_roles123456[]" name="vcSignalMentionRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235" selected>test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<button type="submit" class="btn btn-primary">送信</button>`)
	})

	t.Run("正常に表示される(sendSignalをtureに)", func(t *testing.T) {
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
				{
					ID:       "123456",
					Name:     "test",
					Position: 4,
					Type:     discordgo.ChannelTypeGuildVoice,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID:       "123",
						Username: "test",
					},
				},
			},
			Roles: []*discordgo.Role{
				{
					ID:   "1235",
					Name: "test",
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
				if vcChannelID == "123456" {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     "123456",
						GuildID:         "123",
						SendSignal:      true,
						SendChannelID:   "123456",
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				}
				return &repository.VcSignalChannelAllColumn{}, nil
			},
			InsertVcSignalChannelFunc: func(ctx context.Context, vcChannelID string, guildID string, sendChannelID string) error {
				return nil
			},
			GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
		}

		handler := NewVcSignalViewHandler(indexService, repo)

		mux.HandleFunc("/guild/{guildId}/vc-signal", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123/vc-signal", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), "<h1>ボイスチャンネルの入退出通知</h1>")
		assert.Contains(t, rec.Body.String(), `<label for="sendSignal123456">通知を送信する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="sendSignal123456" name="sendSignal123456" checked />`)

		assert.Contains(t, rec.Body.String(), `<label for="joinBot123456">Botの入退出を通知する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="joinBot123456" name="joinBot123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="everyoneMention123456">通知に@everyoneメンションをつける</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="everyoneMention123456" name="everyoneMention123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="defaultChannel123456">送信先チャンネル</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="defaultChannel123456" name="defaultChannelId123456" >`)
		assert.Contains(t, rec.Body.String(), `<option value="123">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="1234">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="12345">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="123456" selected>:🔊:test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)

		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgUserIds123456[]">NGユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users123456[]" name="vcSignalNgUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgRoleIds123456[]">NGロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles123456[]" name="vcSignalNgRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionUserIds123456[]">メンションユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_users123456[]" name="vcSignalMentionUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionRoleIds123456[]">メンションロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_roles123456[]" name="vcSignalMentionRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<button type="submit" class="btn btn-primary">送信</button>`)
	})

	t.Run("正常に表示される(joinBotをtureに)", func(t *testing.T) {
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
				{
					ID:       "123456",
					Name:     "test",
					Position: 4,
					Type:     discordgo.ChannelTypeGuildVoice,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID:       "123",
						Username: "test",
					},
				},
			},
			Roles: []*discordgo.Role{
				{
					ID:   "1235",
					Name: "test",
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
				if vcChannelID == "123456" {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     "123456",
						GuildID:         "123",
						SendSignal:      false,
						SendChannelID:   "123456",
						JoinBot:         true,
						EveryoneMention: false,
					}, nil
				}
				return &repository.VcSignalChannelAllColumn{}, nil
			},
			InsertVcSignalChannelFunc: func(ctx context.Context, vcChannelID string, guildID string, sendChannelID string) error {
				return nil
			},
			GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
		}

		handler := NewVcSignalViewHandler(indexService, repo)

		mux.HandleFunc("/guild/{guildId}/vc-signal", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123/vc-signal", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), "<h1>ボイスチャンネルの入退出通知</h1>")
		assert.Contains(t, rec.Body.String(), `<label for="sendSignal123456">通知を送信する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="sendSignal123456" name="sendSignal123456"  />`)

		assert.Contains(t, rec.Body.String(), `<label for="joinBot123456">Botの入退出を通知する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="joinBot123456" name="joinBot123456"checked />`)

		assert.Contains(t, rec.Body.String(), `<label for="everyoneMention123456">通知に@everyoneメンションをつける</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="everyoneMention123456" name="everyoneMention123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="defaultChannel123456">送信先チャンネル</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="defaultChannel123456" name="defaultChannelId123456" >`)
		assert.Contains(t, rec.Body.String(), `<option value="123">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="1234">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="12345">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="123456" selected>:🔊:test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)

		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgUserIds123456[]">NGユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users123456[]" name="vcSignalNgUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgRoleIds123456[]">NGロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles123456[]" name="vcSignalNgRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionUserIds123456[]">メンションユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_users123456[]" name="vcSignalMentionUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionRoleIds123456[]">メンションロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_roles123456[]" name="vcSignalMentionRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<button type="submit" class="btn btn-primary">送信</button>`)
	})

	t.Run("正常に表示される(everyoneMentionをtureに)", func(t *testing.T) {
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
				{
					ID:       "123456",
					Name:     "test",
					Position: 4,
					Type:     discordgo.ChannelTypeGuildVoice,
				},
			},
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID:       "123",
						Username: "test",
					},
				},
			},
			Roles: []*discordgo.Role{
				{
					ID:   "1235",
					Name: "test",
				},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, indexService.DiscordBotState.Guilds, 1)

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
				if vcChannelID == "123456" {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     "123456",
						GuildID:         "123",
						SendSignal:      false,
						SendChannelID:   "123456",
						JoinBot:         false,
						EveryoneMention: true,
					}, nil
				}
				return &repository.VcSignalChannelAllColumn{}, nil
			},
			InsertVcSignalChannelFunc: func(ctx context.Context, vcChannelID string, guildID string, sendChannelID string) error {
				return nil
			},
			GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
			GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
				return []string{}, nil
			},
		}

		handler := NewVcSignalViewHandler(indexService, repo)

		mux.HandleFunc("/guild/{guildId}/vc-signal", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123/vc-signal", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), "<h1>ボイスチャンネルの入退出通知</h1>")
		assert.Contains(t, rec.Body.String(), `<label for="sendSignal123456">通知を送信する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="sendSignal123456" name="sendSignal123456"  />`)

		assert.Contains(t, rec.Body.String(), `<label for="joinBot123456">Botの入退出を通知する</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="joinBot123456" name="joinBot123456" />`)

		assert.Contains(t, rec.Body.String(), `<label for="everyoneMention123456">通知に@everyoneメンションをつける</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="checkbox" id="everyoneMention123456" name="everyoneMention123456"checked />`)

		assert.Contains(t, rec.Body.String(), `<label for="defaultChannel123456">送信先チャンネル</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="defaultChannel123456" name="defaultChannelId123456" >`)
		assert.Contains(t, rec.Body.String(), `<option value="123">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="1234">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="12345">:📝:test</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="123456" selected>:🔊:test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)

		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgUserIds123456[]">NGユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_users123456[]" name="vcSignalNgUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalNgRoleIds123456[]">NGロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="ng_roles123456[]" name="vcSignalNgRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionUserIds123456[]">メンションユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_users123456[]" name="vcSignalMentionUserIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="123">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<label for="vcSignalMentionRoleIds123456[]">メンションロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select id="mention_roles123456[]" name="vcSignalMentionRoleIds123456[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `</select>`)
		assert.Contains(t, rec.Body.String(), `<button type="submit" class="btn btn-primary">送信</button>`)
	})
}

func setCtxValue(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = ctxvalue.ContextWithDiscordPermission(ctx, &model.DiscordPermissionData{
		PermissionCode: 8,
		User: model.DiscordUser{
			ID:       "123",
			Username: "test",
			Avatar:   "test",
		},
		Permission: "all",
	})
	return r.WithContext(ctx)
}
