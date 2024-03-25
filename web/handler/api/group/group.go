package group

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/pkg/line"

	"github.com/maguro-alternative/remake_bot/web/handler/api/group/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/permission"
	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

//go:generate go run github.com/matryer/moq -out mock_test.go . Repository
type Repository interface {
	UpdateLineBot(ctx context.Context, lineBot *internal.LineBot) error
}

//go:generate go run github.com/matryer/moq -out permission_mock_test.go . OAuthPermission
type OAuthPermission interface {
	CheckLinePermission(ctx context.Context, r *http.Request, guildId string) (lineProfile line.LineProfile, lineLoginUser *model.LineOAuthSession, err error)
}

type LineGroupHandler struct {
	IndexService    *service.IndexService
	repo            Repository
	oauthPermission OAuthPermission
}

func NewLineGroupHandler(indexService *service.IndexService) *LineGroupHandler {
	return &LineGroupHandler{
		IndexService: indexService,
	}
}

func (g *LineGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		slog.ErrorContext(ctx, "/api/group Method Not Allowed")
		return
	}
	var lineGroupJson internal.LineBotJson
	var repo Repository
	var oauthPermission OAuthPermission
	if err := json.NewDecoder(r.Body).Decode(&lineGroupJson); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "jsonの読み取りに失敗しました:"+err.Error())
		return
	}
	if err := lineGroupJson.Validate(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "jsonのバリデーションに失敗しました:"+err.Error())
		return
	}
	guildId := r.PathValue("guildId")
	oauthPermission = permission.NewPermissionHandler(r, g.IndexService)
	if g.oauthPermission != nil {
		oauthPermission = g.oauthPermission
	}
	_, _, err := oauthPermission.CheckLinePermission(
		ctx,
		r,
		guildId,
	)
	if err != nil {
		http.Redirect(w, r, "/login/line", http.StatusFound)
		slog.InfoContext(ctx, "Redirect to /login/line")
		return
	}
	repo = internal.NewRepository(g.IndexService.DB)
	// mockの場合はmockを使用
	if g.repo != nil {
		repo = g.repo
	}
	err = repo.UpdateLineBot(ctx, &internal.LineBot{
		GuildID:          guildId,
		DefaultChannelID: lineGroupJson.DefaultChannelID,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "DBの更新に失敗しました:"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}
