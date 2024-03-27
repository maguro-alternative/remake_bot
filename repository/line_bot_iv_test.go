package repository


import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestInsertLineBotIv(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM line_bot_iv")

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
	t.Run("LineBotIvが正しく挿入されること", func(t *testing.T) {
		lineBotIv := &LineBotIv{
			GuildID:          "987654321",
		}
		err := repo.InsertLineBotIv(ctx, lineBotIv)
		assert.NoError(t, err)

		var lineBotIvResult LineBotIv
		err = tx.GetContext(ctx, &lineBotIvResult, "SELECT * FROM line_bot_iv WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBotIvResult.GuildID)
	})
}

func TestGetLineBotIvNotClient(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM line_bot_iv")

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
		lineBotIv, err := repo.GetLineBotIvNotClient(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, []byte("123456789"), lineBotIv.LineNotifyTokenIv[0])
		assert.Equal(t, []byte("123456789"), lineBotIv.LineBotTokenIv[0])
		assert.Equal(t, []byte("123456789"), lineBotIv.LineBotSecretIv[0])
		assert.Equal(t, []byte("987654321"), lineBotIv.LineGroupIDIv[0])
	})
}

func TestRepository_GetAllColumnsLineBotIv(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM line_bot_iv")

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineBotIv(ctx, func(lbi *fixtures.LineBotIv) {
			lbi.GuildID = "987654321"
			lbi.LineNotifyTokenIv = pq.ByteaArray{[]byte("123456789")}
			lbi.LineBotTokenIv = pq.ByteaArray{[]byte("123456789")}
			lbi.LineBotSecretIv = pq.ByteaArray{[]byte("123456789")}
			lbi.LineGroupIDIv = pq.ByteaArray{[]byte("987654321")}
			lbi.LineClientIDIv = pq.ByteaArray{[]byte("123456789")}
			lbi.LineClientSecretIv = pq.ByteaArray{[]byte("123456789")}
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからLineBotIvを取得できること", func(t *testing.T) {
		lineBotIv, err := repo.GetAllColumnsLineBotIv(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBotIv.LineNotifyTokenIv)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBotIv.LineBotTokenIv)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBotIv.LineBotSecretIv)
		assert.Equal(t, pq.ByteaArray{[]byte("987654321")}, lineBotIv.LineGroupIDIv)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBotIv.LineClientIDIv)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBotIv.LineClientSecretIv)
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

		tx.ExecContext(ctx, "DELETE FROM line_bot_iv")

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

		tx.ExecContext(ctx, "DELETE FROM line_bot_iv")

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

