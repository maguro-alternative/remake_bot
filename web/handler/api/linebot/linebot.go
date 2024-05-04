package linebot

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/pkg/line"
	"github.com/maguro-alternative/remake_bot/pkg/youtube"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linebot/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
)

// A LineBotHandler handles requests for the line bot.
type LineBotHandler struct {
	indexService *service.IndexService
	repo         repository.RepositoryFunc
}

// NewLineBotHandler returns new LineBotHandler.
func NewLineBotHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
) *LineBotHandler {
	return &LineBotHandler{
		indexService: indexService,
		repo:         repo,
	}
}

// ServeHTTP handles HTTP requests.
func (h *LineBotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var lineResponses internal.LineResponses
	var lineBotDecrypt *internal.LineBotDecrypt
	var lineBotIv repository.LineBotIvNotClient
	// 暗号化キーの取得
	privateKey := config.PrivateKey()
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	lineBots, err := h.repo.GetAllColumnsLineBots(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "line_botの取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if len(lineBots) == 0 {
		slog.ErrorContext(ctx, "line_botの情報が見つかりませんでした。")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// リクエストボディの読み込み
	requestBodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "リクエストの読み込みに失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for i, lineBot := range lineBots {
		lineBotIv, err = h.repo.GetLineBotIvNotClient(ctx, lineBot.GuildID)
		if err != nil {
			slog.ErrorContext(ctx, "line_bot_ivの取得に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// リクエストボディの検証
		lineBotDecrypt, err = internal.LineHmac(privateKey, requestBodyByte, lineBot, lineBotIv, r.Header.Get("X-Line-Signature"))
		if err != nil {
			slog.ErrorContext(ctx, "署名の検証に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		// 署名が一致した場合はループを抜ける
		if lineBotDecrypt != nil {
			break
		}
		// 署名が一致しなかった場合は最後のループでエラーを返す
		if i == len(lineBots)-1 {
			slog.ErrorContext(ctx, "line_botの情報が見つかりませんでした。")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	// リクエストボディのバイトから構造体への変換
	err = json.Unmarshal(requestBodyByte, &lineResponses)
	if err != nil {
		slog.ErrorContext(ctx, "jsonの読み込みに失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// バリデーションチェック
	err = lineResponses.Validate()
	if err != nil {
		slog.ErrorContext(ctx, "バリデーションチェックに失敗しました。", "エラー:", err.Error())
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	// 応答確認の場合は200を返す
	if len(lineResponses.Events) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	lineEvent := lineResponses.Events[0]

	lineRequ := line.NewLineRequest(
		*h.indexService.Client,
		lineBotDecrypt.LineNotifyToken,
		lineBotDecrypt.LineBotToken,
		lineBotDecrypt.LineGroupID,
	)
	// ユーザー情報の取得
	lineProfile, err := lineRequ.GetProfile(ctx, lineEvent.Source.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "LINEユーザー情報の取得に失敗しました。", "エラー:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(ctx, "LINEユーザー名: "+lineProfile.DisplayName)

	// メッセージの種類によって処理を分岐
	switch lineEvent.Message.Type {
	case "text":
		_, err = h.indexService.DiscordSession.ChannelMessageSend(
			lineBotDecrypt.DefaultChannelID,
			lineProfile.DisplayName+"\n「 "+lineResponses.Events[0].Message.Text+" 」",
		)
		if err != nil {
			slog.ErrorContext(ctx, "discordへのメッセージ送信に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case "sticker":
		_, err = h.indexService.DiscordSession.ChannelMessageSend(
			lineBotDecrypt.DefaultChannelID,
			lineProfile.DisplayName+"\nスタンプを送信しました\nhttps://stickershop.line-scdn.net/stickershop/v1/sticker/"+lineResponses.Events[0].Message.StickerID+"/iPhone/sticker.png",
		)
		if err != nil {
			slog.ErrorContext(ctx, "discordへのメッセージ送信に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case "image":
		imageContent, err := lineRequ.GetContent(ctx, lineResponses.Events[0].Message.ID)
		if err != nil {
			slog.ErrorContext(ctx, "LINE画像の取得に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// 画像の種類の取得
		imageType, err := magicNumberRead(bytes.NewReader(imageContent.Content))
		if err != nil {
			slog.ErrorContext(ctx, "マジックナンバーの取得に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, err = h.indexService.DiscordSession.ChannelFileSendWithMessage(
			lineBotDecrypt.DefaultChannelID,
			lineProfile.DisplayName+"\n ",
			"image."+imageType,
			bytes.NewReader(imageContent.Content),
		)
		if err != nil {
			slog.ErrorContext(ctx, "discordへのメッセージ送信に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case "video":
		videoContent, err := lineRequ.GetContent(ctx, lineResponses.Events[0].Message.ID)
		if err != nil {
			slog.ErrorContext(ctx, "LINE動画の取得に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// 25MB以下の動画はdiscordにアップロードさせる
		if videoContent.ContentLength <= 25_000_000 {
			// 動画の種類の取得
			videoType, err := magicNumberRead(bytes.NewReader(videoContent.Content))
			if err != nil {
				slog.ErrorContext(ctx, "マジックナンバーの取得に失敗しました。", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			_, err = h.indexService.DiscordSession.ChannelFileSendWithMessage(
				lineBotDecrypt.DefaultChannelID,
				lineProfile.DisplayName+"\n ",
				"video."+videoType,
				bytes.NewReader(videoContent.Content),
			)
			if err != nil {
				slog.ErrorContext(ctx, "discordへのメッセージ送信に失敗しました。", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			// 25MB以上の動画はYoutubeにアップロードさせる
			youtubeAPI := youtube.NewYoutubeAPI(
				config.YoutubeAccessToken(),
				config.YoutubeTokenExpiry(),
				config.YoutubeClientID(),
				config.YoutubeClientSecret(),
				config.YoutubeRefreshToken(),
				config.YoutubeProjectID(),
			)
			videoID, err := youtubeAPI.VideoUpload(
				ctx,
				bytes.NewReader(videoContent.Content),
				lineProfile.DisplayName+"の動画",
				"LINEからの動画投稿",
				"22",
				"unlisted",
				[]string{"LINE", "動画"},
			)
			if err != nil {
				slog.ErrorContext(ctx, "Youtubeへの動画アップロードに失敗しました。", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			_, err = h.indexService.DiscordSession.ChannelMessageSend(
				lineBotDecrypt.DefaultChannelID,
				lineProfile.DisplayName+"\nhttps://www.youtube.com/watch?v="+videoID,
			)
			if err != nil {
				slog.ErrorContext(ctx, "discordへのメッセージ送信に失敗しました。", "エラー:", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	case "audio":
		audioContent, err := lineRequ.GetContent(ctx, lineResponses.Events[0].Message.ID)
		if err != nil {
			slog.ErrorContext(ctx, "LINE音声の取得に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// 音声の種類の取得
		audioType, err := magicNumberRead(bytes.NewReader(audioContent.Content))
		if err != nil {
			slog.ErrorContext(ctx, "マジックナンバーの取得に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, err = h.indexService.DiscordSession.ChannelFileSendWithMessage(
			lineBotDecrypt.DefaultChannelID,
			lineProfile.DisplayName+"\n ",
			"audio."+audioType,
			bytes.NewReader(audioContent.Content),
		)
		if err != nil {
			slog.ErrorContext(ctx, "discordへのメッセージ送信に失敗しました。", "エラー:", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	// レスポンスの書き込み
	w.WriteHeader(http.StatusOK)
}

// マジックナンバーからファイルの種類を取得
func magicNumberRead(content io.Reader) (string, error) {
	var buf bytes.Buffer
	data := make([]byte, 12)
	if _, err := io.Copy(&buf, content); err != nil {
		return "", err
	}
	count, err := buf.Read(data)
	if err != nil {
		return "", err
	}
	// 最初の4バイトを読み取る
	if bytes.Equal(data[:count], []byte{0, 0, 0, 24, 102, 116, 121, 112, 109, 112, 52, 50}) {
		return "mp4", nil
	} else if bytes.Equal(data[:count], []byte{79, 103, 103, 83, 0, 2, 0, 0, 0, 0, 0, 0}) {
		return "ogg", nil
	} else if bytes.Equal(data[:4], []byte{82, 73, 70, 70}) && bytes.Equal(data[8:count], []byte{87, 65, 86, 69}) {
		return "wav", nil
	} else if bytes.Equal(data[8:count], []byte{77, 52, 65, 32}) {
		return "m4a", nil
	} else if bytes.Equal(data[:4], []byte{73, 68, 51, 4}) {
		return "mp3", nil
	} else if bytes.Equal(data[:6], []byte{71, 73, 70, 56, 57, 97}) {
		return "gif", nil
	} else if bytes.Equal(data[:3], []byte{255, 216, 255}) {
		return "jpg", nil
	}
	return "", nil
}
