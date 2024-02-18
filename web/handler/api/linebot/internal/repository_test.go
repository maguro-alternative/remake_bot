package internal

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/stretchr/testify/assert"
)

func TestRepository_UpdateLineBot(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineBot(ctx, func(lb *fixtures.LineBot) {
			lb.GuildID = "987654321"
			lb.LineNotifyToken = []byte("123456789")
			lb.LineBotToken = []byte("123456789")
			lb.LineBotSecret = []byte("123456789")
			lb.LineGroupID = []byte("987654321")
			lb.DefaultChannelID = "987654321"
			lb.DebugMode = false
		}),
	)
}