package linetoken

import (
	"net/http"
	"encoding/json"

	"github.com/maguro-alternative/remake_bot/web/handler/api/linetoken/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
)

type LineTokenHandler struct {
	IndexService *service.IndexService
}

func NewLineTokenHandler(indexService *service.IndexService) *LineTokenHandler {
	return &LineTokenHandler{
		IndexService: indexService,
	}
}

func (h *LineTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var lineTokenJson internal.LineBotJson
	if err := json.NewDecoder(r.Body).Decode(&lineTokenJson); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := lineTokenJson.Validate(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	if ctx == nil {
		ctx = r.Context()
	}
	repo := internal.NewRepository(h.IndexService.DB)
	if err := repo.UpdateLineBot(ctx, &lineTokenJson); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
