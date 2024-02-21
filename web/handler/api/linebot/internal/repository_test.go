package internal

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/web/config"
	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/stretchr/testify/assert"
	"github.com/lib/pq"
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
			lb.LineNotifyToken = pq.ByteaArray{[]byte("123456789")}
			lb.LineBotToken = pq.ByteaArray{[]byte("123456789")}
			lb.LineBotSecret = pq.ByteaArray{[]byte("123456789")}
			lb.LineGroupID = pq.ByteaArray{[]byte("987654321")}
			lb.DefaultChannelID = "987654321"
			lb.DebugMode = false
		}),
	)

	repo := NewRepository(tx)
	t.Run("GuildIDからLineBotを取得できること", func(t *testing.T) {
		lineBots, err := repo.GetLineBots(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(lineBots))
		assert.Equal(t, "987654321", lineBots[0].GuildID)
		assert.Equal(t, []byte("123456789"), lineBots[0].LineNotifyToken)
		assert.Equal(t, []byte("123456789"), lineBots[0].LineBotToken)
		assert.Equal(t, []byte("123456789"), lineBots[0].LineBotSecret)
		assert.Equal(t, []byte("987654321"), lineBots[0].LineGroupID)
		assert.Equal(t, "987654321", lineBots[0].DefaultChannelID)
		assert.Equal(t, false, lineBots[0].DebugMode)
	})
}

func TestRepository_GetLineBotIv(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineBotIv(ctx, func(lbi *fixtures.LineBotIv) {
			lbi.GuildID = "987654321"
			lbi.LineNotifyTokenIv = pq.ByteaArray{[]byte("123456789")}
			lbi.LineBotTokenIv = pq.ByteaArray{[]byte("123456789")}
			lbi.LineBotSecretIv = pq.ByteaArray{[]byte("123456789")}
			lbi.LineGroupIDIv = pq.ByteaArray{[]byte("987654321")}
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからLineBotIvを取得できること", func(t *testing.T) {
		lineBotIv, err := repo.GetLineBotIv(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, []byte("123456789"), lineBotIv.LineNotifyTokenIv)
		assert.Equal(t, []byte("123456789"), lineBotIv.LineBotTokenIv)
		assert.Equal(t, []byte("123456789"), lineBotIv.LineBotSecretIv)
		assert.Equal(t, []byte("987654321"), lineBotIv.LineGroupIDIv)
	})
}
