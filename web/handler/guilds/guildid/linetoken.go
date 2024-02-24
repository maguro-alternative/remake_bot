package guildid

import (
	"html/template"
	"os"
	"net/http"

	"github.com/maguro-alternative/remake_bot/web/service"
)

type GuildIdHandler struct {
	IndexService *service.IndexService
}

func NewGuildIdHandler(indexService *service.IndexService) *GuildIdHandler {
	return &GuildIdHandler{
		IndexService: indexService,
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Hello, World",
	}
	t := template.Must(template.New("a.html").ParseFiles("a.html"))
    t.Execute(os.Stdout, data)
}
