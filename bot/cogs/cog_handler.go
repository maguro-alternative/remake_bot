package cogs

import (
	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type CogHandler struct {
	DB db.Driver
}

func NewCogHandler(db db.Driver) *CogHandler {
	return &CogHandler{
		DB: db,
	}
}
