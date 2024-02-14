package linebot

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/maguro-alternative/remake_bot/pkg/line"
	"github.com/maguro-alternative/remake_bot/pkg/youtube"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/web/handler/api/linebot/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
)

type Repository interface {
	GetLineBots(ctx context.Context) ([]*internal.LineBot, error)
}

// A LineBotHandler handles requests for the line bot.
type LineBotHandler struct {
	IndexService *service.IndexService
}

// NewLineBotHandler returns new LineBotHandler.
func NewLineBotHandler(indexService *service.IndexService) *LineBotHandler {
	return &LineBotHandler{
		IndexService: indexService,
	}
}

// ServeHTTP handles HTTP requests.
func (h *LineBotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var lineResponses LineResponses
	// 暗号化キーの取得
	privateKey := config.PrivateKey()
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	repo := internal.NewRepository(h.IndexService.DB)
	lineBots, err := repo.GetLineBots(ctx)
	if err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}

	// リクエストボディの読み込み
	requestBodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}

	// リクエストボディの検証
	lineBotDecrypt, err := internal.LineHmac(privateKey, requestBodyByte, lineBots, r.Header.Get("X-Line-Signature"))
	if err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}

	// リクエストボディのバイトから構造体への変換
	err = json.Unmarshal(requestBodyByte, &lineResponses)
	if err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}
	// バリデーションチェック
	if err := lineResponses.Validate(); err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}


	lineRequ := line.NewLineRequest(lineBotDecrypt.LineNotifyToken, lineBotDecrypt.LineBotToken, lineBotDecrypt.LineGroupID)
	// ユーザー情報の取得
	lineProfile, err := lineRequ.GetProfile(ctx, lineResponses.Events[0].Source.UserID)
	if err != nil {
		log.Println("Failed to Load Request")
		http.Error(w, "Failed to Load Request", http.StatusBadRequest)
		return
	}

	// メッセージの種類によって処理を分岐
	if lineResponses.Events[0].Type == "text" {
		_, err = h.IndexService.DiscordSession.ChannelMessageSend(
			lineBotDecrypt.DefaultChannelID,
			lineProfile.DisplayName+"\n「 "+lineResponses.Events[0].Message.Text+" 」",
		)
		if err != nil {
			log.Println("Failed to Load Request")
			http.Error(w, "Failed to Load Request", http.StatusBadRequest)
			return
		}
	}
	if lineResponses.Events[0].Type == "sticker" {
		_, err = h.IndexService.DiscordSession.ChannelMessageSend(
			lineBotDecrypt.DefaultChannelID,
			lineProfile.DisplayName+"\nスタンプを送信しました\nhttps://stickershop.line-scdn.net/stickershop/v1/sticker/"+lineResponses.Events[0].Message.StickerID+"/iPhone/sticker.png",
		)
		if err != nil {
			log.Println("Failed to Load Request")
			http.Error(w, "Failed to Load Request", http.StatusBadRequest)
			return
		}
	}
	if lineResponses.Events[0].Type == "image" {
		imageContent, err := lineRequ.GetContent(ctx, lineResponses.Events[0].Message.ID)
		if err != nil {
			log.Println("Failed to Load Request")
			http.Error(w, "Failed to Load Request", http.StatusBadRequest)
			return
		}
		// 画像の種類の取得
		imageType, err := magicNumberRead(imageContent.Content)
		_, err = h.IndexService.DiscordSession.ChannelFileSendWithMessage(
			lineBotDecrypt.DefaultChannelID,
			lineProfile.DisplayName+"\n ",
			"image."+imageType,
			imageContent.Content,
		)
		if err != nil {
			log.Println("Failed to Load Request")
			http.Error(w, "Failed to Load Request", http.StatusBadRequest)
			return
		}
	}
	if lineResponses.Events[0].Type == "video" {
		videoContent, err := lineRequ.GetContent(ctx, lineResponses.Events[0].Message.ID)
		if err != nil {
			log.Println("Failed to Load Request")
			http.Error(w, "Failed to Load Request", http.StatusBadRequest)
			return
		}
		if videoContent.ContentLength <= 25000000 {
			// 動画の種類の取得
			videoType, err := magicNumberRead(videoContent.Content)
			_, err = h.IndexService.DiscordSession.ChannelFileSendWithMessage(
				lineBotDecrypt.DefaultChannelID,
				lineProfile.DisplayName+"\n ",
				"video."+videoType,
				videoContent.Content,
			)
			if err != nil {
				log.Println("Failed to Load Request")
				http.Error(w, "Failed to Load Request", http.StatusBadRequest)
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
				videoContent.Content,
				lineProfile.DisplayName+"の動画",
				"LINEからの動画投稿",
				"22",
				"unlisted",
				[]string{"LINE", "動画"},
			)
			if err != nil {
				log.Println("Failed to Load Request")
				http.Error(w, "Failed to Load Request", http.StatusBadRequest)
				return
			}
			_, err = h.IndexService.DiscordSession.ChannelMessageSend(
				lineBotDecrypt.DefaultChannelID,
				lineProfile.DisplayName+"\nhttps://www.youtube.com/watch?v="+videoID,
			)
			if err != nil {
				log.Println("Failed to Load Request")
				http.Error(w, "Failed to Load Request", http.StatusBadRequest)
				return
			}
		}
	}
	if lineResponses.Events[0].Type == "audio" {
		audioContent, err := lineRequ.GetContent(ctx, lineResponses.Events[0].Message.ID)
		if err != nil {
			log.Println("Failed to Load Request")
			http.Error(w, "Failed to Load Request", http.StatusBadRequest)
			return
		}
		// 音声の種類の取得
		audioType, err := magicNumberRead(audioContent.Content)
		_, err = h.IndexService.DiscordSession.ChannelFileSendWithMessage(
			lineBotDecrypt.DefaultChannelID,
			lineProfile.DisplayName+"\n ",
			"audio."+audioType,
			audioContent.Content,
		)
		if err != nil {
			log.Println("Failed to Load Request")
			http.Error(w, "Failed to Load Request", http.StatusBadRequest)
			return
		}
	}
	// レスポンスの書き込み
	w.WriteHeader(http.StatusOK)
}

// マジックナンバーからファイルの種類を取得
func magicNumberRead(content io.ReadCloser) (string, error) {
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