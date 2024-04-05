package cogs

import (
	"net/http"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type CogHandler struct {
	DB db.Driver
	client *http.Client
}

func NewCogHandler(
	db db.Driver,
	client *http.Client,
) *CogHandler {
	return &CogHandler{
		DB: db,
		client: client,
	}
}
