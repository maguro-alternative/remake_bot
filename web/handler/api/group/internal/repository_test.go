package internal

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/web/config"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestRepository_UpdateLineBot(t *testing.T) {
	ctx := context.Background()
	t.Run("LineBotが正しく更新されること", func(t *testing.T) {
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
				lb.LineClientID = pq.ByteaArray{[]byte("987654321")}
				lb.LineClientSecret = pq.ByteaArray{[]byte("987654321")}
				lb.DefaultChannelID = "987654321"
				lb.DebugMode = false
			}),
		)

		repo := NewRepository(tx)
		updateLineBot := &LineBot{
			GuildID:          "987654321",
			DefaultChannelID: "123456789",
		}
		err = repo.UpdateLineBot(ctx, updateLineBot)
		assert.NoError(t, err)

		var lineBot LineBot
		err = tx.GetContext(ctx, &lineBot, "SELECT * FROM line_bot WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBot.GuildID)
		assert.Equal(t, "123456789", lineBot.DefaultChannelID)
	})

	t.Run("LineBotの1部分(notifyとbottoken)が正しく更新されること", func(t *testing.T) {
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
				lb.LineClientID = pq.ByteaArray{[]byte("987654321")}
				lb.LineClientSecret = pq.ByteaArray{[]byte("987654321")}
				lb.DefaultChannelID = "987654321"
				lb.DebugMode = false
			}),
		)

		repo := NewRepository(tx)
		updateLineBot := &LineBot{
			GuildID:          "987654321",
			DefaultChannelID: "123456789",
		}
		err = repo.UpdateLineBot(ctx, updateLineBot)
		assert.NoError(t, err)

		var lineBot LineBot
		err = tx.GetContext(ctx, &lineBot, "SELECT * FROM line_bot WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBot.GuildID)
		assert.Equal(t, "123456789", lineBot.DefaultChannelID)
	})
}
