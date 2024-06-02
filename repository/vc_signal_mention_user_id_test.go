package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestInsertVcSignalMentionUser(t *testing.T) {
	ctx := context.Background()
	t.Run("MentionUserを追加できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")
		repo := NewRepository(tx)

		err = repo.InsertVcSignalMentionUser(ctx, "111", "1111", "11111")
		assert.NoError(t, err)

		var mentionUserIDs []*VcSignalMentionUser
		err = tx.SelectContext(ctx, &mentionUserIDs, "SELECT * FROM vc_signal_mention_user_id")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 1)
		assert.Equal(t, "111", mentionUserIDs[0].VcChannelID)
		assert.Equal(t, "1111", mentionUserIDs[0].GuildID)
		assert.Equal(t, "11111", mentionUserIDs[0].UserID)
	})

	t.Run("既にある場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
		)

		err = repo.InsertVcSignalMentionUser(ctx, "111", "1111", "11111")
		assert.NoError(t, err)

		var mentionUserIDs []*VcSignalMentionUser
		err = tx.SelectContext(ctx, &mentionUserIDs, "SELECT * FROM vc_signal_mention_user_id")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 1)
		assert.Equal(t, "111", mentionUserIDs[0].VcChannelID)
		assert.Equal(t, "1111", mentionUserIDs[0].GuildID)
		assert.Equal(t, "11111", mentionUserIDs[0].UserID)
	})
}

func TestGetVcSignalMentionUsersByChannelID(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したchannelIDのMentionUserを取得できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.UserID = "22222"
			}),
		)

		mentionUserIDs, err := repo.GetVcSignalMentionUserIDsByVcChannelID(ctx, "111")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 1)
		assert.Equal(t, "11111", mentionUserIDs[0])
	})

	t.Run("指定したchannelIDのMentionUserがない場合は空の配列を返すこと", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")
		repo := NewRepository(tx)

		mentionUserIDs, err := repo.GetVcSignalMentionUserIDsByVcChannelID(ctx, "111")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 0)
	})
}

func TestDeleteVcSignalMentionUser(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したchannelID, guildID, userIDのMentionUserを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.UserID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionUser(ctx, "111", "1111", "11111")
		assert.NoError(t, err)

		var mentionUserIDs []*VcSignalMentionUser
		err = tx.SelectContext(ctx, &mentionUserIDs, "SELECT * FROM vc_signal_mention_user_id")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 1)
		assert.Equal(t, "222", mentionUserIDs[0].VcChannelID)
		assert.Equal(t, "2222", mentionUserIDs[0].GuildID)
		assert.Equal(t, "22222", mentionUserIDs[0].UserID)
	})

	t.Run("指定したchannelID, guildID, userIDのMentionUserがない場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.UserID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionUser(ctx, "333", "3333", "33333")
		assert.NoError(t, err)

		var mentionUserIDs []*VcSignalMentionUser
		err = tx.SelectContext(ctx, &mentionUserIDs, "SELECT * FROM vc_signal_mention_user_id")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 2)
	})
}

func TestDeleteVcSignalMentionUserByChannelID(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したchannelIDのMentionUserを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.UserID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionUsersByVcChannelID(ctx, "111")
		assert.NoError(t, err)

		var mentionUserIDs []*VcSignalMentionUser
		err = tx.SelectContext(ctx, &mentionUserIDs, "SELECT * FROM vc_signal_mention_user_id")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 1)
		assert.Equal(t, "222", mentionUserIDs[0].VcChannelID)
		assert.Equal(t, "2222", mentionUserIDs[0].GuildID)
		assert.Equal(t, "22222", mentionUserIDs[0].UserID)
	})

	t.Run("指定したchannelIDのMentionUserがない場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")

		repo := NewRepository(tx)

		err = repo.DeleteVcSignalMentionUsersByVcChannelID(ctx, "111")
		assert.NoError(t, err)

		var mentionUserIDs []*VcSignalMentionUser
		err = tx.SelectContext(ctx, &mentionUserIDs, "SELECT * FROM vc_signal_mention_user_id")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 0)
	})
}

func TestDeleteVcSignalMentionUserByGuildID(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したguildIDのMentionUserを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.UserID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionUsersByGuildID(ctx, "1111")
		assert.NoError(t, err)

		var mentionUserIDs []*VcSignalMentionUser
		err = tx.SelectContext(ctx, &mentionUserIDs, "SELECT * FROM vc_signal_mention_user_id")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 1)
		assert.Equal(t, "222", mentionUserIDs[0].VcChannelID)
		assert.Equal(t, "2222", mentionUserIDs[0].GuildID)
		assert.Equal(t, "22222", mentionUserIDs[0].UserID)
	})

	t.Run("指定したguildIDのMentionUserがない場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")

		repo := NewRepository(tx)

		err = repo.DeleteVcSignalMentionUsersByGuildID(ctx, "1111")
		assert.NoError(t, err)

		var mentionUserIDs []*VcSignalMentionUser
		err = tx.SelectContext(ctx, &mentionUserIDs, "SELECT * FROM vc_signal_mention_user_id")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 0)
	})
}

func TestDeleteVcSignalMentionUserByUserID(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したuserIDのMentionUserを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "222"
				v.GuildID = "2222"
				v.UserID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionUsersByUserID(ctx, "11111")
		assert.NoError(t, err)

		var mentionUserIDs []*VcSignalMentionUser
		err = tx.SelectContext(ctx, &mentionUserIDs, "SELECT * FROM vc_signal_mention_user_id")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 1)
		assert.Equal(t, "222", mentionUserIDs[0].VcChannelID)
		assert.Equal(t, "2222", mentionUserIDs[0].GuildID)
		assert.Equal(t, "22222", mentionUserIDs[0].UserID)
	})

	t.Run("指定したuserIDのMentionUserがない場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")

		repo := NewRepository(tx)

		err = repo.DeleteVcSignalMentionUsersByUserID(ctx, "11111")
		assert.NoError(t, err)

		var mentionUserIDs []*VcSignalMentionUser
		err = tx.SelectContext(ctx, &mentionUserIDs, "SELECT * FROM vc_signal_mention_user_id")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 0)
	})
}

func TestDeleteVcSignalMentionUsersNotInProvidedList(t *testing.T) {
	ctx := context.Background()
	t.Run("指定したchannelIDのMentionUserのうち、指定したguildID, userIDのMentionUser以外を削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionUsersNotInProvidedList(ctx, "111", []string{"11111"})
		assert.NoError(t, err)

		var mentionUserIDs []*VcSignalMentionUser
		err = tx.SelectContext(ctx, &mentionUserIDs, "SELECT * FROM vc_signal_mention_user_id")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 1)
		assert.Equal(t, "111", mentionUserIDs[0].VcChannelID)
		assert.Equal(t, "1111", mentionUserIDs[0].GuildID)
		assert.Equal(t, "11111", mentionUserIDs[0].UserID)
	})

	t.Run("指定したchannelIDのMentionUserのうち、指定したguildID, userIDのMentionUser以外がない場合は全て削除されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_mention_user_id")
		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "11111"
			}),
			fixtures.NewVcSignalMentionUserID(ctx, func(v *fixtures.VcSignalMentionUserID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.UserID = "22222"
			}),
		)

		err = repo.DeleteVcSignalMentionUsersNotInProvidedList(ctx, "111", []string{})
		assert.NoError(t, err)

		var mentionUserIDs []*VcSignalMentionUser
		err = tx.SelectContext(ctx, &mentionUserIDs, "SELECT * FROM vc_signal_mention_user_id")
		assert.NoError(t, err)

		assert.Len(t, mentionUserIDs, 0)
	})
}
