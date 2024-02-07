package api

import (
	"net/http"
	"context"
	//"io"
	//"encoding/hex"

	"github.com/maguro-alternative/remake_bot/web/service"
	//"github.com/maguro-alternative/remake_bot/pkg/crypto"
)

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
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
}