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

func TestGetLineWorksBotsInfo(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot_info")
	assert.NoError(t, err)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineWorksBotInfo(ctx, func(l *fixtures.LineWorksBotInfo) {
			l.GuildID = "987654321"
			l.LineWorksClientID = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksClientSecret = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksPrivateKey = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksServiceAccount = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksDomainID = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksAdminID = pq.ByteaArray{[]byte("123456789")}
		}),
	)

	repo := NewRepository(tx)
	t.Run("LineWorksBotを取得できること", func(t *testing.T) {
		info, err := repo.GetLineWorksBotInfoByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", info.GuildID)
		assert.Equal(t, "123456789", string(info.LineWorksClientID[0]))
		assert.Equal(t, "123456789", string(info.LineWorksClientSecret[0]))
		assert.Equal(t, "123456789", string(info.LineWorksServiceAccount[0]))
		assert.Equal(t, "123456789", string(info.LineWorksPrivateKey[0]))
		assert.Equal(t, "123456789", string(info.LineWorksDomainID[0]))
		assert.Equal(t, "123456789", string(info.LineWorksAdminID[0]))
	})

	t.Run("存在しないGuildIDの場合はエラーが返ること", func(t *testing.T) {
		_, err := repo.GetLineWorksBotInfoByGuildID(ctx, "123456789")
		assert.Error(t, err)
	})
}

func TestInsertLineWorksBotInfo(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot_info")
	assert.NoError(t, err)

	repo := NewRepository(tx)
	t.Run("LineWorksBotを追加できること", func(t *testing.T) {
		info := NewLineWorksBotInfo(
			"987654321",
			pq.ByteaArray{[]byte("123456789")},
			pq.ByteaArray{[]byte("123456789")},
			pq.ByteaArray{[]byte("123456789")},
			pq.ByteaArray{[]byte("123456789")},
			pq.ByteaArray{[]byte("123456789")},
			pq.ByteaArray{[]byte("123456789")},
		)
		err := repo.InsertLineWorksBotInfo(ctx, info)
		assert.NoError(t, err)

		info, err = repo.GetLineWorksBotInfoByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", info.GuildID)
		assert.Equal(t, "123456789", string(info.LineWorksClientID[0]))
		assert.Equal(t, "123456789", string(info.LineWorksClientSecret[0]))
		assert.Equal(t, "123456789", string(info.LineWorksServiceAccount[0]))
		assert.Equal(t, "123456789", string(info.LineWorksPrivateKey[0]))
		assert.Equal(t, "123456789", string(info.LineWorksDomainID[0]))
		assert.Equal(t, "123456789", string(info.LineWorksAdminID[0]))
	})
}

func TestUpdateLineWorksBotInfo(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot_info")
	assert.NoError(t, err)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineWorksBotInfo(ctx, func(l *fixtures.LineWorksBotInfo) {
			l.GuildID = "987654321"
			l.LineWorksClientID = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksClientSecret = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksPrivateKey = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksServiceAccount = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksDomainID = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksAdminID = pq.ByteaArray{[]byte("123456789")}
		}),
	)

	repo := NewRepository(tx)
	t.Run("LineWorksBotを更新できること", func(t *testing.T) {
		err := repo.UpdateLineWorksBotInfo(ctx, &LineWorksBotInfo{
			GuildID:               "987654321",
			LineWorksClientID:     pq.ByteaArray{[]byte("987654321")},
			LineWorksClientSecret: pq.ByteaArray{[]byte("987654321")},
			LineWorksPrivateKey:   pq.ByteaArray{[]byte("987654321")},
			LineWorksServiceAccount: pq.ByteaArray{[]byte("987654321")},
			LineWorksDomainID:     pq.ByteaArray{[]byte("987654321")},
			LineWorksAdminID:      pq.ByteaArray{[]byte("987654321")},
		})
		assert.NoError(t, err)

		info, err := repo.GetLineWorksBotInfoByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", info.GuildID)
		assert.Equal(t, "987654321", string(info.LineWorksClientID[0]))
		assert.Equal(t, "987654321", string(info.LineWorksClientSecret[0]))
		assert.Equal(t, "987654321", string(info.LineWorksPrivateKey[0]))
		assert.Equal(t, "987654321", string(info.LineWorksServiceAccount[0]))
		assert.Equal(t, "987654321", string(info.LineWorksDomainID[0]))
		assert.Equal(t, "987654321", string(info.LineWorksAdminID[0]))
	})
}
