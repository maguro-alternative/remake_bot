package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	//"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestInsertVcSignalChannel(t *testing.T) {
	ctx := context.Background()
	t.Run("ChannelIDを追加できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_channel")

		repo := NewRepository(tx)
		err = repo.InsertVcSignalChannel(ctx, "123456789", "987654321", "1234567890")
		assert.NoError(t, err)

		var channels []VcSignalChannelAllColumns
		err = tx.SelectContext(ctx, &channels, "SELECT * FROM vc_signal_channel")
		assert.NoError(t, err)
		assert.Len(t, channels, 1)
		assert.Equal(t, "123456789", channels[0].VcChannelID)
		assert.Equal(t, "987654321", channels[0].GuildID)
		assert.Equal(t, "1234567890", channels[0].SendChannelID)
	})

	t.Run("ChannelIDが重複している場合はエラーは返さず挿入しないこと", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_channel")

		repo := NewRepository(tx)
		err = repo.InsertVcSignalChannel(ctx, "123456789", "987654321", "1234567890")
		assert.NoError(t, err)

		err = repo.InsertVcSignalChannel(ctx, "123456789", "987654321", "1234567890")
		assert.NoError(t, err)

		var channels []VcSignalChannelAllColumns
		err = tx.SelectContext(ctx, &channels, "SELECT * FROM vc_signal_channel")
		assert.NoError(t, err)
		assert.Len(t, channels, 1)
		assert.Equal(t, "123456789", channels[0].VcChannelID)
		assert.Equal(t, "987654321", channels[0].GuildID)
		assert.Equal(t, "1234567890", channels[0].SendChannelID)
	})
}
