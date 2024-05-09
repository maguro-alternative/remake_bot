package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertVcSignalNgUserID(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	require.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	require.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
	t.Run("NgUserIDを追加できること", func(t *testing.T) {
		repo := NewRepository(tx)
		err = repo.InsertVcSignalNgUser(ctx, "111","1111","11111")
		assert.NoError(t, err)

		var ngUser VcSignalNgUserAllColumn
		err = tx.GetContext(ctx, &ngUser, "SELECT * FROM vc_signal_ng_user_id WHERE user_id = '11111'")
		assert.NoError(t, err)

		assert.Equal(t, "111", ngUser.VcChannelID)
		assert.Equal(t, "1111", ngUser.GuildID)
		assert.Equal(t, "11111", ngUser.UserID)
	})

	t.Run("既にある場合はエラーを返さずそのまま", func(t *testing.T) {
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
		)

		repo := NewRepository(tx)
		err = repo.InsertVcSignalNgUser(ctx, "111","1111","11111")
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)
		assert.Len(t, ngUsers, 1)
	})
}

func TestGetVcSignalNgUsersByChannelIDAllColumn(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	require.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	require.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
	t.Run("NgUserIDを取得できること", func(t *testing.T) {
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
		)

		repo := NewRepository(tx)
		ngUsers, err := repo.GetVcSignalNgUsersByChannelIDAllColumn(ctx, "111")
		assert.NoError(t, err)

		assert.Len(t, ngUsers, 1)
		assert.Equal(t, "111", ngUsers[0].VcChannelID)
		assert.Equal(t, "1111", ngUsers[0].GuildID)
		assert.Equal(t, "11111", ngUsers[0].UserID)
	})

	t.Run("存在しない場合は空のスライスを返す", func(t *testing.T) {
		repo := NewRepository(tx)
		ngUsers, err := repo.GetVcSignalNgUsersByChannelIDAllColumn(ctx, "111")
		assert.NoError(t, err)

		assert.Len(t, ngUsers, 0)
	})
}

func TestDeleteVcNgUserByChannelID(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	require.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	require.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
	t.Run("NgUserIDを削除できること", func(t *testing.T) {
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteVcNgUserByChannelID(ctx, "111")
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)

		assert.Len(t, ngUsers, 0)
	})

	t.Run("存在しない場合はエラーを返さずに終了すること", func(t *testing.T) {
		repo := NewRepository(tx)
		err = repo.DeleteVcNgUserByChannelID(ctx, "111")
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)
		assert.Len(t, ngUsers, 0)
	})
}

func TestDeleteVcNgUserByUserID(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	require.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	require.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
	t.Run("NgUserIDを削除できること", func(t *testing.T) {
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteVcNgUserByUserID(ctx, "11111")
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)

		assert.Len(t, ngUsers, 0)
	})

	t.Run("存在しない場合はエラーを返さずに終了すること", func(t *testing.T) {
		repo := NewRepository(tx)
		err = repo.DeleteVcNgUserByUserID(ctx, "11111")
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)
		assert.Len(t, ngUsers, 0)
	})
}
