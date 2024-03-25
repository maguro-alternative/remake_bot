package group

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maguro-alternative/remake_bot/pkg/line"

	"github.com/maguro-alternative/remake_bot/web/handler/api/group/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"

	"github.com/stretchr/testify/assert"
)

func TestLineGroupHandler_ServeHTTP(t *testing.T) {
	bodyJson, err := json.Marshal(internal.LineBotJson{
		DefaultChannelID: "123456789",
		DebugMode: 	  true,
	})
	assert.NoError(t, err)

	t.Run("MethodがPOST以外の場合、Method Not Allowedが返ること", func(t *testing.T) {
		h := &LineGroupHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/group", nil)
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("jsonの読み取りに失敗すると、Bad Requestが返ること", func(t *testing.T) {
		h := &LineGroupHandler{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/group", bytes.NewReader([]byte("")))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("jsonのバリデーションに失敗すると、Unprocessable Entityが返ること", func(t *testing.T) {
		h := &LineGroupHandler{
			IndexService: &service.IndexService{},
			repo: &RepositoryMock{
				UpdateLineBotFunc: func(ctx context.Context, lineBot *internal.LineBot) error {
					return nil
				},
			},
			oauthPermission: &OAuthPermissionMock{
				CheckLinePermissionFunc: func(ctx context.Context, r *http.Request, guildId string) (lineProfile line.LineProfile, lineLoginUser *model.LineOAuthSession, err error) {
					return line.LineProfile{}, &model.LineOAuthSession{}, nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/group", bytes.NewReader([]byte(`{"channel_id":"123456789"}`)))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("LineBotの更新が成功すること", func(t *testing.T) {
		h := &LineGroupHandler{
			IndexService: &service.IndexService{},
			repo: &RepositoryMock{
				UpdateLineBotFunc: func(ctx context.Context, lineBot *internal.LineBot) error {
					return nil
				},
			},
			oauthPermission: &OAuthPermissionMock{
				CheckLinePermissionFunc: func(ctx context.Context, r *http.Request, guildId string) (lineProfile line.LineProfile, lineLoginUser *model.LineOAuthSession, err error) {
					return line.LineProfile{}, &model.LineOAuthSession{}, nil
				},
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/987654321/group", bytes.NewReader(bodyJson))
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
