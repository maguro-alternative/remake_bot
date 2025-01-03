package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetLineWorksBotsIv(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot_iv")
	assert.NoError(t, err)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineWorksBotIv(ctx, func(l *fixtures.LineWorksBotIV) {
			l.GuildID = "987654321"
			l.LineWorksBotTokenIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksRefreshTokenIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksGroupIDIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksBotIDIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksBotSecretIV = pq.ByteaArray{[]byte("123456789")}
		}),
	)

	repo := NewRepository(tx)
	t.Run("全てのLineWorksBotを取得できること", func(t *testing.T) {
		iv, err := repo.GetLineWorksBotIVByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", iv.GuildID)
		assert.Equal(t, "123456789", string(iv.LineWorksBotTokenIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksRefreshTokenIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksGroupIDIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksBotIDIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksBotSecretIV[0]))
	})

	t.Run("存在しないGuildIDの場合はエラーが返ること", func(t *testing.T) {
		_, err := repo.GetLineWorksBotIVByGuildID(ctx, "123456789")
		assert.Error(t, err)
	})
}

func TestInsertLineWorksBotIv(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot_iv")
	assert.NoError(t, err)

	repo := NewRepository(tx)
	t.Run("LineWorksBotIVを追加できること", func(t *testing.T) {
		iv := &LineWorksBotIV{
			GuildID:               "987654321",
			LineWorksBotTokenIV:     pq.ByteaArray{[]byte("123456789")},
			LineWorksRefreshTokenIV: pq.ByteaArray{[]byte("123456789")},
			LineWorksGroupIDIV:      pq.ByteaArray{[]byte("123456789")},
			LineWorksBotIDIV:        pq.ByteaArray{[]byte("123456789")},
			LineWorksBotSecretIV:    pq.ByteaArray{[]byte("123456789")},
		}
		err := repo.InsertLineWorksBotIV(ctx, iv)
		assert.NoError(t, err)

		iv2, err := repo.GetLineWorksBotIVByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", iv2.GuildID)
		assert.Equal(t, "123456789", string(iv2.LineWorksBotTokenIV[0]))
		assert.Equal(t, "123456789", string(iv2.LineWorksRefreshTokenIV[0]))
		assert.Equal(t, "123456789", string(iv2.LineWorksGroupIDIV[0]))
		assert.Equal(t, "123456789", string(iv2.LineWorksBotIDIV[0]))
		assert.Equal(t, "123456789", string(iv2.LineWorksBotSecretIV[0]))
	})
}

func TestUpdateLineWorksBotIv(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot_iv")
	assert.NoError(t, err)

	repo := NewRepository(tx)
	t.Run("LineWorksBotIVを更新できること", func(t *testing.T) {
		iv := &LineWorksBotIV{
			GuildID:               "987654321",
			LineWorksBotTokenIV:     pq.ByteaArray{[]byte("123456789")},
			LineWorksRefreshTokenIV: pq.ByteaArray{[]byte("123456789")},
			LineWorksGroupIDIV:      pq.ByteaArray{[]byte("123456789")},
			LineWorksBotIDIV:        pq.ByteaArray{[]byte("123456789")},
			LineWorksBotSecretIV:    pq.ByteaArray{[]byte("123456789")},
		}
		err := repo.InsertLineWorksBotIV(ctx, iv)
		assert.NoError(t, err)

		iv.LineWorksBotTokenIV = pq.ByteaArray{[]byte("987654321")}
		iv.LineWorksRefreshTokenIV = pq.ByteaArray{[]byte("987654321")}
		iv.LineWorksGroupIDIV = pq.ByteaArray{[]byte("987654321")}
		iv.LineWorksBotIDIV = pq.ByteaArray{[]byte("987654321")}
		iv.LineWorksBotSecretIV = pq.ByteaArray{[]byte("987654321")}
		err = repo.UpdateLineWorksBotIV(ctx, iv)
		assert.NoError(t, err)

		iv2, err := repo.GetLineWorksBotIVByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", iv2.GuildID)
		assert.Equal(t, "987654321", string(iv2.LineWorksBotTokenIV[0]))
		assert.Equal(t, "987654321", string(iv2.LineWorksRefreshTokenIV[0]))
		assert.Equal(t, "987654321", string(iv2.LineWorksGroupIDIV[0]))
		assert.Equal(t, "987654321", string(iv2.LineWorksBotIDIV[0]))
		assert.Equal(t, "987654321", string(iv2.LineWorksBotSecretIV[0]))
	})
}
