package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestInsertVcSignalNgUserID(t *testing.T) {
	ctx := context.Background()
	t.Run("NgUserIDを追加できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
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
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
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
	t.Run("NgUserIDを取得できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
		)

		repo := NewRepository(tx)
		ngUsers, err := repo.GetVcSignalNgUserIDsByVcChannelID(ctx, "111")
		assert.NoError(t, err)

		assert.Len(t, ngUsers, 1)
		assert.Equal(t, "11111", ngUsers[0])
	})

	t.Run("存在しない場合は空のスライスを返す", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
		repo := NewRepository(tx)
		ngUsers, err := repo.GetVcSignalNgUserIDsByVcChannelID(ctx, "111")
		assert.NoError(t, err)

		assert.Len(t, ngUsers, 0)
	})
}

func TestDeleteVcNgUserByChannelID(t *testing.T) {
	ctx := context.Background()
	t.Run("NgUserIDを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteVcNgUserByVcChannelID(ctx, "111")
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)

		assert.Len(t, ngUsers, 0)
	})

	t.Run("存在しない場合はエラーを返さずに終了すること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
		repo := NewRepository(tx)
		err = repo.DeleteVcNgUserByVcChannelID(ctx, "111")
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)
		assert.Len(t, ngUsers, 0)
	})
}

func TestDeleteVcNgUserByGuildID(t *testing.T) {
	ctx := context.Background()
	t.Run("NgUserIDを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteVcNgUserByGuildID(ctx, "1111")
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)

		assert.Len(t, ngUsers, 0)
	})

	t.Run("存在しない場合はエラーを返さずに終了すること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
		repo := NewRepository(tx)
		err = repo.DeleteVcNgUserByGuildID(ctx, "1111")
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)
		assert.Len(t, ngUsers, 0)
	})
}

func TestDeleteVcNgUserByUserID(t *testing.T) {
	ctx := context.Background()
	t.Run("NgUserIDを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
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
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
		repo := NewRepository(tx)
		err = repo.DeleteVcNgUserByUserID(ctx, "11111")
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)
		assert.Len(t, ngUsers, 0)
	})
}

func TestDeleteNgUsersNotInProvidedList(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したリストに含まれないNgUserIDを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "22222"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteVcSignalNgUsersNotInProvidedList(ctx, "111", []string{"11111", "33333"})
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)

		assert.Len(t, ngUsers, 1)
		assert.Equal(t, "111", ngUsers[0].VcChannelID)
		assert.Equal(t, "1111", ngUsers[0].GuildID)
		assert.Equal(t, "11111", ngUsers[0].UserID)
	})

	t.Run("指定したリストに含まれないNgUserIDがない場合は何もしないこと", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.UserID = "22222"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteVcSignalNgUsersNotInProvidedList(ctx, "111", []string{"11111", "22222"})
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)

		assert.Len(t, ngUsers, 2)
	})

	t.Run("指定したリストが空の場合、チャンネルidのものをすべて削除すること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
			fixtures.NewVcSignalNgUserID(ctx, func(v *fixtures.VcSignalNgUserID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.UserID = "22222"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteVcSignalNgUsersNotInProvidedList(ctx, "111", []string{})
		assert.NoError(t, err)

		var ngUsers []VcSignalNgUserAllColumn
		err = tx.SelectContext(ctx, &ngUsers, "SELECT * FROM vc_signal_ng_user_id")
		assert.NoError(t, err)

		assert.Len(t, ngUsers, 1)
	})
}
