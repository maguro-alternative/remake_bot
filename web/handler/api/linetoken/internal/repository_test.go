package internal

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/web/config"

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
				lb.LineNotifyToken = []byte("123456789")
				lb.LineBotToken = []byte("123456789")
				lb.LineBotSecret = []byte("123456789")
				lb.LineGroupID = []byte("987654321")
				lb.LineClientID = []byte("987654321")
				lb.LineClientSecret = []byte("987654321")
				lb.DefaultChannelID = "987654321"
				lb.DebugMode = false
			}),
		)

		repo := NewRepository(tx)
		updateLineBot := &LineBot{
			GuildID:          "987654321",
			LineNotifyToken:  []byte("987654321"),
			LineBotToken:     []byte("987654321"),
			LineBotSecret:    []byte("987654321"),
			LineGroupID:      []byte("123456789"),
			LineClientID:     []byte("123456789"),
			LineClientSecret: []byte("123456789"),
			DefaultChannelID: "123456789",
			DebugMode:        true,
		}
		err = repo.UpdateLineBot(ctx, updateLineBot)
		assert.NoError(t, err)

		var lineBot LineBot
		err = tx.GetContext(ctx, &lineBot, "SELECT * FROM line_bot WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBot.GuildID)
		assert.Equal(t, []byte("987654321"), lineBot.LineNotifyToken)
		assert.Equal(t, []byte("987654321"), lineBot.LineBotToken)
		assert.Equal(t, []byte("987654321"), lineBot.LineBotSecret)
		assert.Equal(t, []byte("123456789"), lineBot.LineGroupID)
		assert.Equal(t, []byte("123456789"), lineBot.LineClientID)
		assert.Equal(t, []byte("123456789"), lineBot.LineClientSecret)
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
				lb.LineNotifyToken = []byte("123456789")
				lb.LineBotToken = []byte("123456789")
				lb.LineBotSecret = []byte("123456789")
				lb.LineGroupID = []byte("987654321")
				lb.LineClientID = []byte("987654321")
				lb.LineClientSecret = []byte("987654321")
				lb.DefaultChannelID = "987654321"
				lb.DebugMode = false
			}),
		)

		repo := NewRepository(tx)
		updateLineBot := &LineBot{
			GuildID:          "987654321",
			LineNotifyToken:  []byte("987654321"),
			LineBotToken:     []byte("987654321"),
			LineBotSecret:    []byte("987654321"),
			DefaultChannelID: "123456789",
			DebugMode:        true,
		}
		err = repo.UpdateLineBot(ctx, updateLineBot)
		assert.NoError(t, err)

		var lineBot LineBot
		err = tx.GetContext(ctx, &lineBot, "SELECT * FROM line_bot WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBot.GuildID)
		assert.Equal(t, []byte("987654321"), lineBot.LineNotifyToken)
		assert.Equal(t, []byte("987654321"), lineBot.LineBotToken)
		assert.Equal(t, []byte("987654321"), lineBot.LineBotSecret)
		assert.Equal(t, []byte("987654321"), lineBot.LineGroupID)
		assert.Equal(t, []byte("987654321"), lineBot.LineClientID)
		assert.Equal(t, []byte("987654321"), lineBot.LineClientSecret)
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
				lb.LineNotifyTokenIv = []byte("123456789")
				lb.LineBotTokenIv = []byte("123456789")
				lb.LineBotSecretIv = []byte("123456789")
				lb.LineGroupIDIv = []byte("123456789")
				lb.LineClientIDIv = []byte("123456789")
				lb.LineClientSecretIv = []byte("123456789")
			}),
		)

		updateLineBotIv := &LineBotIv{
			GuildID:            "987654321",
			LineNotifyTokenIv:  []byte("987654321"),
			LineBotTokenIv:     []byte("987654321"),
			LineBotSecretIv:    []byte("987654321"),
			LineGroupIDIv:      []byte("987654321"),
			LineClientIDIv:     []byte("987654321"),
			LineClientSecretIv: []byte("987654321"),
		}

		repo := NewRepository(tx)
		err = repo.UpdateLineBotIv(ctx, updateLineBotIv)
		assert.NoError(t, err)

		var lineBotIv LineBotIv
		err = tx.GetContext(ctx, &lineBotIv, "SELECT * FROM line_bot_iv WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBotIv.GuildID)
		assert.Equal(t, []byte("987654321"), lineBotIv.LineNotifyTokenIv)
		assert.Equal(t, []byte("987654321"), lineBotIv.LineBotTokenIv)
		assert.Equal(t, []byte("987654321"), lineBotIv.LineBotSecretIv)
		assert.Equal(t, []byte("987654321"), lineBotIv.LineGroupIDIv)
		assert.Equal(t, []byte("987654321"), lineBotIv.LineClientIDIv)
		assert.Equal(t, []byte("987654321"), lineBotIv.LineClientSecretIv)
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
				lb.LineNotifyTokenIv = []byte("123456789")
				lb.LineBotTokenIv = []byte("123456789")
				lb.LineBotSecretIv = []byte("123456789")
				lb.LineGroupIDIv = []byte("123456789")
				lb.LineClientIDIv = []byte("123456789")
				lb.LineClientSecretIv = []byte("123456789")
			}),
		)

		updateLineBotIv := &LineBotIv{
			GuildID:           "987654321",
			LineNotifyTokenIv: []byte("987654321"),
			LineBotTokenIv:    []byte("987654321"),
		}

		repo := NewRepository(tx)
		err = repo.UpdateLineBotIv(ctx, updateLineBotIv)
		assert.NoError(t, err)

		var lineBotIv LineBotIv
		err = tx.GetContext(ctx, &lineBotIv, "SELECT * FROM line_bot_iv WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBotIv.GuildID)
		assert.Equal(t, []byte("987654321"), lineBotIv.LineNotifyTokenIv)
		assert.Equal(t, []byte("987654321"), lineBotIv.LineBotTokenIv)
		assert.Equal(t, []byte("123456789"), lineBotIv.LineBotSecretIv)
		assert.Equal(t, []byte("123456789"), lineBotIv.LineGroupIDIv)
		assert.Equal(t, []byte("123456789"), lineBotIv.LineClientIDIv)
		assert.Equal(t, []byte("123456789"), lineBotIv.LineClientSecretIv)
	})
}
