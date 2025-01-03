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

func TestGetLineWorksBotInfoIv(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot_info_iv")
	assert.NoError(t, err)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineWorksBotInfoIv(ctx, func(l *fixtures.LineWorksBotInfoIV) {
			l.GuildID = "987654321"
			l.LineWorksClientIDIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksClientSecretIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksPrivateKeyIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksServiceAccountIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksDomainIDIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksAdminIDIV = pq.ByteaArray{[]byte("123456789")}
		}),
	)

	repo := NewRepository(tx)
	t.Run("LineWorksBotInfoIvを取得できること", func(t *testing.T) {
		iv, err := repo.GetLineWorksBotInfoIVByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", iv.GuildID)
		assert.Equal(t, "123456789", string(iv.LineWorksClientIDIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksClientSecretIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksServiceAccountIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksPrivateKeyIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksDomainIDIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksAdminIDIV[0]))
	})

	t.Run("存在しないGuildIDの場合はエラーが返ること", func(t *testing.T) {
		_, err := repo.GetLineWorksBotInfoIVByGuildID(ctx, "123456789")
		assert.Error(t, err)
	})
}

func TestInsertLineWorksBotInfoIv(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot_info_iv")
	assert.NoError(t, err)

	repo := NewRepository(tx)
	t.Run("LineWorksBotInfoIvを登録できること", func(t *testing.T) {
		iv := NewLineWorksBotInfoIV(
			"987654321",
			pq.ByteaArray{[]byte("123456789")},
			pq.ByteaArray{[]byte("123456789")},
			pq.ByteaArray{[]byte("123456789")},
			pq.ByteaArray{[]byte("123456789")},
			pq.ByteaArray{[]byte("123456789")},
			pq.ByteaArray{[]byte("123456789")},
		)
		err := repo.InsertLineWorksBotInfoIV(ctx, iv)
		assert.NoError(t, err)

		iv, err = repo.GetLineWorksBotInfoIVByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", iv.GuildID)
		assert.Equal(t, "123456789", string(iv.LineWorksClientIDIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksClientSecretIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksServiceAccountIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksPrivateKeyIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksDomainIDIV[0]))
		assert.Equal(t, "123456789", string(iv.LineWorksAdminIDIV[0]))
	})
}

func TestUpdateLineWorksBotInfoIv(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot_info_iv")
	assert.NoError(t, err)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineWorksBotInfoIv(ctx, func(l *fixtures.LineWorksBotInfoIV) {
			l.GuildID = "987654321"
			l.LineWorksClientIDIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksClientSecretIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksPrivateKeyIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksServiceAccountIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksDomainIDIV = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksAdminIDIV = pq.ByteaArray{[]byte("123456789")}
		}),
	)

	defer tx.RollbackCtx(ctx)

	repo := NewRepository(tx)

	t.Run("LineWorksBotInfoIvを更新できること", func(t *testing.T) {
		err := repo.UpdateLineWorksBotInfoIV(ctx, &LineWorksBotInfoIV{
			GuildID:               "987654321",
			LineWorksClientIDIV:     pq.ByteaArray{[]byte("987654321")},
			LineWorksClientSecretIV: pq.ByteaArray{[]byte("987654321")},
			LineWorksPrivateKeyIV:   pq.ByteaArray{[]byte("987654321")},
			LineWorksServiceAccountIV: pq.ByteaArray{[]byte("987654321")},
			LineWorksDomainIDIV:     pq.ByteaArray{[]byte("987654321")},
			LineWorksAdminIDIV:      pq.ByteaArray{[]byte("987654321")},
		})
		assert.NoError(t, err)

		iv, err := repo.GetLineWorksBotInfoIVByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", iv.GuildID)
		assert.Equal(t, "987654321", string(iv.LineWorksClientIDIV[0]))
		assert.Equal(t, "987654321", string(iv.LineWorksClientSecretIV[0]))
		assert.Equal(t, "987654321", string(iv.LineWorksPrivateKeyIV[0]))
		assert.Equal(t, "987654321", string(iv.LineWorksServiceAccountIV[0]))
		assert.Equal(t, "987654321", string(iv.LineWorksDomainIDIV[0]))
		assert.Equal(t, "987654321", string(iv.LineWorksAdminIDIV[0]))
	})
}
