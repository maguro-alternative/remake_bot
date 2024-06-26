package webhook

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/ctxvalue"
	"github.com/maguro-alternative/remake_bot/web/shared/model"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWebhookViewHandler(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(cwd))
	})
	require.NoError(t, os.Chdir("../../../../../"))
	t.Run("正常に表示される(データベースは空)", func(t *testing.T) {
		indexService := &service.IndexService{
			DiscordSession: &mock.SessionMock{
				GuildWebhooksFunc: func(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Webhook, err error) {
					st = []*discordgo.Webhook{
						{
							ID:        "11",
							Type:      1,
							GuildID:   "123",
							ChannelID: "234",
							Name:      "test",
						},
					}
					return st, err
				},
			},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "123",
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID:       "1234",
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

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetAllColumnsWebhooksByGuildIDFunc: func(ctx context.Context, guildID string) ([]*repository.Webhook, error) {
				return []*repository.Webhook{}, nil
			},
			GetWebhookWordWithWebhookSerialIDAndConditionFunc: func(ctx context.Context, webhookSerialID int64, condition string) ([]*repository.WebhookWord, error) {
				return []*repository.WebhookWord{}, nil
			},
			GetWebhookUserMentionWithWebhookSerialIDFunc: func(ctx context.Context, webhookSerialID int64) ([]*repository.WebhookUserMention, error) {
				return []*repository.WebhookUserMention{}, nil
			},
			GetWebhookRoleMentionWithWebhookSerialIDFunc: func(ctx context.Context, webhookSerialID int64) ([]*repository.WebhookRoleMention, error) {
				return []*repository.WebhookRoleMention{}, nil
			},
		}

		handler := NewWebhookViewHandler(indexService, repo)

		mux.HandleFunc("/guild/{guildId}/webhook", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123/webhook", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWebhook()">追加</button>`)
		assert.Contains(t, rec.Body.String(), `<label for="newWebhookType1">Webhook</label>`)
		assert.Contains(t, rec.Body.String(), `<select name="newWebhookType1" id="newWebhookType1">`)
		assert.Contains(t, rec.Body.String(), `<option value="11">test</option>`)
		assert.Contains(t, rec.Body.String(), `<label for="newSubscriptionName1">サービス名</label>`)
		assert.Contains(t, rec.Body.String(), `<select name="newSubscriptionName1" id="newSubscriptionName1" />`)
		assert.Contains(t, rec.Body.String(), `<option value="youtube">youtube</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="niconico">niconico</option>`)
		assert.Contains(t, rec.Body.String(), `<label for="newSubscriptionId1">サービスID</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="text" name="newSubscriptionId1" id="newSubscriptionId1" value="" />`)
		assert.Contains(t, rec.Body.String(), `<label for="newMemberMention1[]">メンションするユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select name="newMemberMention1[]" id="newMemberMention1[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1234">test</option>`)
		assert.Contains(t, rec.Body.String(), `<label for="newRoleMention1[]">メンションするロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select name="newRoleMention1[]" id="newRoleMention1[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235">test</option>`)
		assert.Contains(t, rec.Body.String(), `<div id="newNgOrWords1">`)
		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWord('newNgOr', 1)">NGワードOR検索追加</button>`)
		assert.Contains(t, rec.Body.String(), `<div id="newNgAndWords1">`)
		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWord('newNgAnd', 1)">NGワードAND検索追加</button>`)
		assert.Contains(t, rec.Body.String(), `<div id="newSearchOrWords1">`)
		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWord('newSearchOr', 1)">キーワードOR検索追加</button>`)
		assert.Contains(t, rec.Body.String(), `<div id="newSearchAndWords1">`)
		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWord('newSearchAnd', 1)">キーワードAND検索追加</button>`)
		assert.Contains(t, rec.Body.String(), `<div id="newMentionOrWords1">`)
		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWord('newMentionOr', 1)">メンションOR検索追加</button>`)
		assert.Contains(t, rec.Body.String(), `<div id="newMentionAndWords1">`)
		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWord('newMentionAnd', 1)">メンションAND検索追加</button>`)
		assert.Contains(t, rec.Body.String(), `<button type="submit" class="btn btn-primary">送信</button>`)
	})

	t.Run("正常に表示される(データベースにデータがある)", func(t *testing.T) {
		serialId := int64(1)
		indexService := &service.IndexService{
			DiscordSession: &mock.SessionMock{
				GuildWebhooksFunc: func(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Webhook, err error) {
					st = []*discordgo.Webhook{
						{
							ID:        "11",
							Type:      1,
							GuildID:   "123",
							ChannelID: "234",
							Name:      "test",
						},
					}
					return st, err
				},
			},
		}
		indexService.DiscordBotState = discordgo.NewState()
		err := indexService.DiscordBotState.GuildAdd(&discordgo.Guild{
			ID: "123",
			Members: []*discordgo.Member{
				{
					User: &discordgo.User{
						ID:       "1234",
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

		mux := http.NewServeMux()

		repo := &repository.RepositoryFuncMock{
			GetAllColumnsWebhooksByGuildIDFunc: func(ctx context.Context, guildID string) ([]*repository.Webhook, error) {
				return []*repository.Webhook{
					{
						WebhookSerialID: &serialId,
						GuildID:         "123",
						WebhookID:       "11",
						SubscriptionType: "youtube",
						SubscriptionID:   "test",
					},
				}, nil
			},
			GetWebhookWordWithWebhookSerialIDAndConditionFunc: func(ctx context.Context, webhookSerialID int64, condition string) ([]*repository.WebhookWord, error) {
				return []*repository.WebhookWord{
					{
						WebhookSerialID: 1,
						Condition:       "NgOr",
						Word:            "test",
					},
				}, nil
			},
			GetWebhookUserMentionWithWebhookSerialIDFunc: func(ctx context.Context, webhookSerialID int64) ([]*repository.WebhookUserMention, error) {
				return []*repository.WebhookUserMention{
					{
						WebhookSerialID: 1,
						UserID:          "1234",
					},
				}, nil
			},
			GetWebhookRoleMentionWithWebhookSerialIDFunc: func(ctx context.Context, webhookSerialID int64) ([]*repository.WebhookRoleMention, error) {
				return []*repository.WebhookRoleMention{
					{
						WebhookSerialID: 1,
						RoleID:          "1235",
					},
				}, nil
			},
		}

		handler := NewWebhookViewHandler(indexService, repo)

		mux.HandleFunc("/guild/{guildId}/webhook", handler.Index)

		req := httptest.NewRequest(http.MethodGet, "/guild/123/webhook", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, setCtxValue(req))

		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWebhook()">追加</button>`)
		assert.Contains(t, rec.Body.String(), `<label for="updateWebhookType1">Webhook</label>`)
		assert.Contains(t, rec.Body.String(), `<select name="updateWebhookType1" id="updateWebhookType1">`)
		assert.Contains(t, rec.Body.String(), `<option value="11" selected>test</option>`)
		assert.Contains(t, rec.Body.String(), `<label for="updateSubscriptionName1">サービス名</label>`)
		assert.Contains(t, rec.Body.String(), `<select name="updateSubscriptionName1" id="updateSubscriptionName1" />`)
		assert.Contains(t, rec.Body.String(), `<option value="youtube" selected>youtube</option>`)
		assert.Contains(t, rec.Body.String(), `<option value="niconico">niconico</option>`)
		assert.Contains(t, rec.Body.String(), `<label for="updateSubscriptionId1">サービスID</label>`)
		assert.Contains(t, rec.Body.String(), `<input type="text" name="updateSubscriptionId1" id="updateSubscriptionId1" value="test" />`)
		assert.Contains(t, rec.Body.String(), `<label for="updateMemberMention1[]">メンションするユーザー</label>`)
		assert.Contains(t, rec.Body.String(), `<select name="updateMemberMention1[]" id="updateMemberMention1[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1234" selected>test</option>`)
		assert.Contains(t, rec.Body.String(), `<label for="updateRoleMention1[]">メンションするロール</label>`)
		assert.Contains(t, rec.Body.String(), `<select name="updateRoleMention1[]" id="updateRoleMention1[]" multiple>`)
		assert.Contains(t, rec.Body.String(), `<option value="1235" selected>test</option>`)
		assert.Contains(t, rec.Body.String(), `<div id="updateNgOrWords1">`)
		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWord('updateNgOr', 1)">NGワードOR検索追加</button>`)
		assert.Contains(t, rec.Body.String(), `<input type="text" id="updateNgOrWord1[]" name="updateNgOrWord1[]" value="test">`)
		assert.Contains(t, rec.Body.String(), `<div id="updateNgAndWords1">`)
		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWord('updateNgAnd', 1)">NGワードAND検索追加</button>`)
		assert.Contains(t, rec.Body.String(), `<div id="updateSearchOrWords1">`)
		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWord('updateSearchOr', 1)">キーワードOR検索追加</button>`)
		assert.Contains(t, rec.Body.String(), `<div id="updateSearchAndWords1">`)
		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWord('updateSearchAnd', 1)">キーワードAND検索追加</button>`)
		assert.Contains(t, rec.Body.String(), `<div id="updateMentionOrWords1">`)
		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWord('updateMentionOr', 1)">メンションOR検索追加</button>`)
		assert.Contains(t, rec.Body.String(), `<div id="updateMentionAndWords1">`)
		assert.Contains(t, rec.Body.String(), `<button type="button" onclick="addWord('updateMentionAnd', 1)">メンションAND検索追加</button>`)
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
