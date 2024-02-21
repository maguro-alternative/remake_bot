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
			LineNotifyToken:  pq.ByteaArray{[]byte("987654321")},
			LineBotToken:     pq.ByteaArray{[]byte("987654321")},
			LineBotSecret:    pq.ByteaArray{[]byte("987654321")},
			LineGroupID:      pq.ByteaArray{[]byte("123456789")},
			LineClientID:     pq.ByteaArray{[]byte("123456789")},
			LineClientSecret: pq.ByteaArray{[]byte("123456789")},
			DefaultChannelID: "123456789",
			DebugMode:        true,
		}
		err = repo.UpdateLineBot(ctx, updateLineBot)
		assert.NoError(t, err)

		var lineBot LineBot
		err = tx.GetContext(ctx, &lineBot, "SELECT * FROM line_bot WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBot.GuildID)
		assert.Equal(t, []byte("987654321"), lineBot.LineNotifyToken[0])
		assert.Equal(t, []byte("987654321"), lineBot.LineBotToken[0])
		assert.Equal(t, []byte("987654321"), lineBot.LineBotSecret[0])
		assert.Equal(t, []byte("123456789"), lineBot.LineGroupID[0])
		assert.Equal(t, []byte("123456789"), lineBot.LineClientID[0])
		assert.Equal(t, []byte("123456789"), lineBot.LineClientSecret[0])
		assert.Equal(t, "123456789", lineBot.DefaultChannelID)
		assert.Equal(t, true, lineBot.DebugMode)
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
			LineNotifyToken:  pq.ByteaArray{[]byte("987654321")},
			LineBotToken:     pq.ByteaArray{[]byte("987654321")},
			LineBotSecret:    pq.ByteaArray{[]byte("987654321")},
			DefaultChannelID: "123456789",
			DebugMode:        true,
		}
		err = repo.UpdateLineBot(ctx, updateLineBot)
		assert.NoError(t, err)

		var lineBot LineBot
		err = tx.GetContext(ctx, &lineBot, "SELECT * FROM line_bot WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBot.GuildID)
		assert.Equal(t, []byte("987654321"), lineBot.LineNotifyToken[0])
		assert.Equal(t, []byte("987654321"), lineBot.LineBotToken[0])
		assert.Equal(t, []byte("987654321"), lineBot.LineBotSecret[0])
		assert.Equal(t, []byte("987654321"), lineBot.LineGroupID[0])
		assert.Equal(t, []byte("987654321"), lineBot.LineClientID[0])
		assert.Equal(t, []byte("987654321"), lineBot.LineClientSecret[0])
		assert.Equal(t, "123456789", lineBot.DefaultChannelID)
		assert.Equal(t, true, lineBot.DebugMode)
	})
}

func TestRepository_UpdateLineBotIv(t *testing.T) {
	ctx := context.Background()
	t.Run("LineBotのIVが正しく更新されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewLineBotIv(ctx, func(lb *fixtures.LineBotIv) {
				lb.GuildID = "987654321"
				lb.LineNotifyTokenIv = pq.ByteaArray{[]byte("123456789")}
				lb.LineBotTokenIv = pq.ByteaArray{[]byte("123456789")}
				lb.LineBotSecretIv = pq.ByteaArray{[]byte("123456789")}
				lb.LineGroupIDIv = pq.ByteaArray{[]byte("123456789")}
				lb.LineClientIDIv = pq.ByteaArray{[]byte("123456789")}
				lb.LineClientSecretIv = pq.ByteaArray{[]byte("123456789")}
			}),
		)

		updateLineBotIv := &LineBotIv{
			GuildID:            "987654321",
			LineNotifyTokenIv:  pq.ByteaArray{[]byte("987654321")},
			LineBotTokenIv:     pq.ByteaArray{[]byte("987654321")},
			LineBotSecretIv:    pq.ByteaArray{[]byte("987654321")},
			LineGroupIDIv:      pq.ByteaArray{[]byte("987654321")},
			LineClientIDIv:     pq.ByteaArray{[]byte("987654321")},
			LineClientSecretIv: pq.ByteaArray{[]byte("987654321")},
		}

		repo := NewRepository(tx)
		err = repo.UpdateLineBotIv(ctx, updateLineBotIv)
		assert.NoError(t, err)

		var lineBotIv LineBotIv
		err = tx.GetContext(ctx, &lineBotIv, "SELECT * FROM line_bot_iv WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBotIv.GuildID)
		assert.Equal(t, []byte("987654321"), lineBotIv.LineNotifyTokenIv[0])
		assert.Equal(t, []byte("987654321"), lineBotIv.LineBotTokenIv[0])
		assert.Equal(t, []byte("987654321"), lineBotIv.LineBotSecretIv[0])
		assert.Equal(t, []byte("987654321"), lineBotIv.LineGroupIDIv[0])
		assert.Equal(t, []byte("987654321"), lineBotIv.LineClientIDIv[0])
		assert.Equal(t, []byte("987654321"), lineBotIv.LineClientSecretIv[0])
	})

	t.Run("LineBotのIVの1部分(notifyとbottoken)が正しく更新されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewLineBotIv(ctx, func(lb *fixtures.LineBotIv) {
				lb.GuildID = "987654321"
				lb.LineNotifyTokenIv = pq.ByteaArray{[]byte("123456789")}
				lb.LineBotTokenIv = pq.ByteaArray{[]byte("123456789")}
				lb.LineBotSecretIv = pq.ByteaArray{[]byte("123456789")}
				lb.LineGroupIDIv = pq.ByteaArray{[]byte("123456789")}
				lb.LineClientIDIv = pq.ByteaArray{[]byte("123456789")}
				lb.LineClientSecretIv = pq.ByteaArray{[]byte("123456789")}
			}),
		)

		updateLineBotIv := &LineBotIv{
			GuildID:           "987654321",
			LineNotifyTokenIv: pq.ByteaArray{[]byte("987654321")},
			LineBotTokenIv:    pq.ByteaArray{[]byte("987654321")},
		}

		repo := NewRepository(tx)
		err = repo.UpdateLineBotIv(ctx, updateLineBotIv)
		assert.NoError(t, err)

		var lineBotIv LineBotIv
		err = tx.GetContext(ctx, &lineBotIv, "SELECT * FROM line_bot_iv WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBotIv.GuildID)
		assert.Equal(t, []byte("987654321"), lineBotIv.LineNotifyTokenIv[0])
		assert.Equal(t, []byte("987654321"), lineBotIv.LineBotTokenIv[0])
		assert.Equal(t, []byte("123456789"), lineBotIv.LineBotSecretIv[0])
		assert.Equal(t, []byte("123456789"), lineBotIv.LineGroupIDIv[0])
		assert.Equal(t, []byte("123456789"), lineBotIv.LineClientIDIv[0])
		assert.Equal(t, []byte("123456789"), lineBotIv.LineClientSecretIv[0])
	})
}
