package linelogin

import (
	"context"
	"encoding/hex"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/line"
	"github.com/maguro-alternative/remake_bot/web/shared/ctxvalue"

	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
	"github.com/maguro-alternative/remake_bot/web/shared/session"
)

type LineLoginHandler struct {
	IndexService *service.IndexService
	Repo         repository.RepositoryFunc
}

func NewLineLoginHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
) *LineLoginHandler {
	return &LineLoginHandler{
		IndexService: indexService,
		Repo:         repo,
	}
}

func (h *LineLoginHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var lineBotIv repository.LineBotIv
	var lineLoginHtmlBuilder strings.Builder
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	// 暗号化キーの取得
	privateKey := config.PrivateKey()
	// 暗号化キーのバイトへの変換
	keyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		slog.ErrorContext(ctx, "暗号化キーのバイトへの変換に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	lineBots, err := h.Repo.GetAllColumnsLineBots(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "line_botの取得に失敗しました。")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	for _, lineBot := range lineBots {
		lineBotIv, err = h.Repo.GetAllColumnsLineBotIv(ctx, lineBot.GuildID)
		if err != nil {
			slog.ErrorContext(ctx, "line_bot_ivの取得に失敗しました。")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		lineNotifyTokenByte, err := crypto.Decrypt(lineBot.LineNotifyToken[0], keyBytes, lineBotIv.LineNotifyTokenIv[0])
		if err != nil {
			slog.ErrorContext(ctx, "line_notify_tokenの復号に失敗しました。")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		lineBotTokenByte, err := crypto.Decrypt(lineBot.LineBotToken[0], keyBytes, lineBotIv.LineBotTokenIv[0])
		if err != nil {
			slog.ErrorContext(ctx, "line_bot_tokenの復号に失敗しました。")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		lineGroupByte, err := crypto.Decrypt(lineBot.LineGroupID[0], keyBytes, lineBotIv.LineGroupIDIv[0])
		if err != nil {
			slog.ErrorContext(ctx, "line_group_idの復号に失敗しました。")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		lineRequ := line.NewLineRequest(
			*h.IndexService.Client,
			string(lineNotifyTokenByte),
			string(lineBotTokenByte),
			string(lineGroupByte),
		)
		lineBotProfile, err := lineRequ.GetBotInfo(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "LineBotの情報取得に失敗しました。"+err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		lineLoginHtmlBuilder.WriteString(fmt.Sprintf(`
			<a href="/login/line/%s">
				<img src="%s"/>
				<li>%s</li>
			</a>
			<br><br>
		`, lineBot.GuildID, lineBotProfile.PictureURL, lineBotProfile.DisplayName))
	}
	// Discordの認証情報なしでもアクセス可能なためエラーレスポンスは出さない
	discordLoginUser, err := ctxvalue.DiscordUserFromContext(ctx)
	if err != nil {
		discordLoginUser = &model.DiscordOAuthSession{}
	}
	// Lineの認証情報なしでもアクセス可能なためエラーレスポンスは出さない
	lineSession, err := ctxvalue.LineUserFromContext(ctx)
	if err != nil {
		lineSession = &model.LineOAuthSession{}
	}
	accountVer := strings.Builder{}
	accountVer.WriteString(components.CreateDiscordAccountVer(discordLoginUser.User))
	accountVer.WriteString(components.CreateLineAccountVer(lineSession.User))
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/views/login/line_login.html"))
	err = tmpl.Execute(w, struct {
		Title         string
		AccountVer    template.HTML
		JsScriptTag   template.HTML
		LineLoginList template.HTML
	}{
		Title:         "LINEログイン選択",
		LineLoginList: template.HTML(lineLoginHtmlBuilder.String()),
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "template error: "+err.Error())
		return
	}
}

func (h *LineLoginHandler) LineLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	guildID := r.PathValue("guildId")
	state := uuid.New().String()
	nonce := uuid.New().String()
	sessionStore, err := session.NewSessionStore(r, h.IndexService.CookieStore, config.SessionSecret())
	if err != nil {
		slog.ErrorContext(r.Context(), "sessionの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	// 暗号化キーの取得
	privateKey := config.PrivateKey()
	// 暗号化キーのバイトへの変換
	keyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		slog.ErrorContext(ctx, "暗号化キーのバイトへの変換に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	lineBot, err := h.Repo.GetAllColumnsLineBot(ctx, guildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_botの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	lineBotIv, err := h.Repo.GetAllColumnsLineBotIv(ctx, lineBot.GuildID)
	if err != nil {
		slog.ErrorContext(ctx, "line_bot_ivの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	lineClientIDByte, err := crypto.Decrypt(lineBot.LineClientID[0], keyBytes, lineBotIv.LineClientIDIv[0])
	if err != nil {
		slog.ErrorContext(ctx, "line_client_idの復号に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sessionStore.SetLineState(state)
	sessionStore.SetLineNonce(nonce)
	sessionStore.SetGuildID(guildID)
	err = sessionStore.SessionSave(r, w)
	if err != nil {
		slog.ErrorContext(ctx, "セッションの初期化に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.IndexService.CookieStore.Save(r, w, sessionStore.GetSession())
	if err != nil {
		slog.ErrorContext(ctx, "セッションの保存に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	redirectUriEncode := url.QueryEscape(config.ServerUrl() + "/callback/line-callback/")
	lineOAuthUrl := fmt.Sprintf("https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=profile%%20openid%%20email&nonce=%s", string(lineClientIDByte), redirectUriEncode, state, nonce)
	http.Redirect(w, r, lineOAuthUrl, http.StatusSeeOther)
}
