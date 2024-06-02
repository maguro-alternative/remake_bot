package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestInsertVcSignalMentionRole(t *testing.T) {
	ctx := context.Background()
	t.Run("MentionRoleを追加できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")
		repo := NewRepository(tx)

		err = repo.InsertVcSignalMentionRole(ctx, "111", "1111", "11111")
		assert.NoError(t, err)

		var mentionRoleIDs []*VcSignalMentionRole
		err = tx.SelectContext(ctx, &mentionRoleIDs, "SELECT * FROM vc_signal_mention_role_id")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 1)
		assert.Equal(t, "111", mentionRoleIDs[0].VcChannelID)
		assert.Equal(t, "1111", mentionRoleIDs[0].GuildID)
		assert.Equal(t, "11111", mentionRoleIDs[0].RoleID)
	})

	t.Run("既にある場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
		)

		err = repo.InsertVcSignalMentionRole(ctx, "111", "1111", "11111")
		assert.NoError(t, err)

		var mentionRoleIDs []*VcSignalMentionRole
		err = tx.SelectContext(ctx, &mentionRoleIDs, "SELECT * FROM vc_signal_mention_role_id")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 1)
		assert.Equal(t, "111", mentionRoleIDs[0].VcChannelID)
		assert.Equal(t, "1111", mentionRoleIDs[0].GuildID)
		assert.Equal(t, "11111", mentionRoleIDs[0].RoleID)
	})
}

func TestGetVcSignalMentionRolesByChannelID(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したchannelIDのMentionRoleを取得できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.RoleID = "22222"
			}),
		)

		mentionRoleIDs, err := repo.GetVcSignalMentionRoleIDsByVcChannelID(ctx, "111")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 1)
		assert.Equal(t, "11111", mentionRoleIDs[0])
	})

	t.Run("指定したchannelIDのMentionRoleがない場合は空の配列を返すこと", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")
		repo := NewRepository(tx)

		mentionRoleIDs, err := repo.GetVcSignalMentionRoleIDsByVcChannelID(ctx, "111")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 0)
	})
}

func TestDeleteVcSignalMentionRole(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したchannelID, guildID, roleIDのMentionRoleを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.RoleID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionRole(ctx, "111", "1111", "11111")
		assert.NoError(t, err)

		var mentionRoleIDs []*VcSignalMentionRole
		err = tx.SelectContext(ctx, &mentionRoleIDs, "SELECT * FROM vc_signal_mention_role_id")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 1)
		assert.Equal(t, "222", mentionRoleIDs[0].VcChannelID)
		assert.Equal(t, "2222", mentionRoleIDs[0].GuildID)
		assert.Equal(t, "22222", mentionRoleIDs[0].RoleID)
	})

	t.Run("指定したchannelID, guildID, roleIDのMentionRoleがない場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.RoleID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionRole(ctx, "333", "3333", "33333")
		assert.NoError(t, err)

		var mentionRoleIDs []*VcSignalMentionRole
		err = tx.SelectContext(ctx, &mentionRoleIDs, "SELECT * FROM vc_signal_mention_role_id")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 2)
	})
}

func TestDeleteVcSignalMentionRoleByChannelID(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したchannelIDのMentionRoleを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.RoleID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionRolesByVcChannelID(ctx, "111")
		assert.NoError(t, err)

		var mentionRoleIDs []*VcSignalMentionRole
		err = tx.SelectContext(ctx, &mentionRoleIDs, "SELECT * FROM vc_signal_mention_role_id")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 1)
		assert.Equal(t, "222", mentionRoleIDs[0].VcChannelID)
		assert.Equal(t, "2222", mentionRoleIDs[0].GuildID)
		assert.Equal(t, "22222", mentionRoleIDs[0].RoleID)
	})

	t.Run("指定したchannelIDのMentionRoleがない場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")

		repo := NewRepository(tx)

		err = repo.DeleteVcSignalMentionRolesByVcChannelID(ctx, "111")
		assert.NoError(t, err)

		var mentionRoleIDs []*VcSignalMentionRole
		err = tx.SelectContext(ctx, &mentionRoleIDs, "SELECT * FROM vc_signal_mention_role_id")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 0)
	})
}

func TestDeleteVcSignalMentionRoleByGuildID(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したguildIDのMentionRoleを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.RoleID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionRolesByGuildID(ctx, "1111")
		assert.NoError(t, err)

		var mentionRoleIDs []*VcSignalMentionRole
		err = tx.SelectContext(ctx, &mentionRoleIDs, "SELECT * FROM vc_signal_mention_role_id")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 1)
		assert.Equal(t, "222", mentionRoleIDs[0].VcChannelID)
		assert.Equal(t, "2222", mentionRoleIDs[0].GuildID)
		assert.Equal(t, "22222", mentionRoleIDs[0].RoleID)
	})

	t.Run("指定したguildIDのMentionRoleがない場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")

		repo := NewRepository(tx)

		err = repo.DeleteVcSignalMentionRolesByGuildID(ctx, "1111")
		assert.NoError(t, err)

		var mentionRoleIDs []*VcSignalMentionRole
		err = tx.SelectContext(ctx, &mentionRoleIDs, "SELECT * FROM vc_signal_mention_role_id")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 0)
	})
}

func TestDeleteVcSignalMentionRoleByRoleID(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したroleIDのMentionRoleを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.RoleID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionRolesByRoleID(ctx, "11111")
		assert.NoError(t, err)

		var mentionRoleIDs []*VcSignalMentionRole
		err = tx.SelectContext(ctx, &mentionRoleIDs, "SELECT * FROM vc_signal_mention_role_id")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 1)
		assert.Equal(t, "222", mentionRoleIDs[0].VcChannelID)
		assert.Equal(t, "2222", mentionRoleIDs[0].GuildID)
		assert.Equal(t, "22222", mentionRoleIDs[0].RoleID)
	})

	t.Run("指定したroleIDのMentionRoleがない場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")

		repo := NewRepository(tx)

		err = repo.DeleteVcSignalMentionRolesByRoleID(ctx, "11111")
		assert.NoError(t, err)

		var mentionRoleIDs []*VcSignalMentionRole
		err = tx.SelectContext(ctx, &mentionRoleIDs, "SELECT * FROM vc_signal_mention_role_id")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 0)
	})
}

func TestDeleteVcSignalMentionRolesNotInProvidedList(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したchannelIDのMentionRoleのうち、指定したguildID, roleIDのMentionRole以外を削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionRolesNotInProvidedList(ctx, "111", []string{"11111"})
		assert.NoError(t, err)

		var mentionRoleIDs []*VcSignalMentionRole
		err = tx.SelectContext(ctx, &mentionRoleIDs, "SELECT * FROM vc_signal_mention_role_id")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 1)
		assert.Equal(t, "111", mentionRoleIDs[0].VcChannelID)
		assert.Equal(t, "1111", mentionRoleIDs[0].GuildID)
		assert.Equal(t, "11111", mentionRoleIDs[0].RoleID)
	})

	t.Run("指定したchannelIDのMentionRoleのうち、指定したguildID, roleIDのMentionRole以外がない場合は全て削除されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_role_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
			fixtures.NewVcSignalMentionRoleID(ctx, func(v *fixtures.VcSignalMentionRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionRolesNotInProvidedList(ctx, "111", []string{})
		assert.NoError(t, err)

		var mentionRoleIDs []*VcSignalMentionRole
		err = tx.SelectContext(ctx, &mentionRoleIDs, "SELECT * FROM vc_signal_mention_role_id")
		assert.NoError(t, err)

		assert.Len(t, mentionRoleIDs, 0)
	})
}
