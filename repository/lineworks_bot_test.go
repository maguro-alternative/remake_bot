package repository

import (
	"context"
	"testing"
	"time"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetAllLineWorksBots(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot")
	assert.NoError(t, err)

	now := time.Now().UTC()

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineWorksBot(ctx, func(l *fixtures.LineWorksBot) {
			l.GuildID = "987654321"
			l.LineWorksBotToken = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksRefreshToken = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksGroupID = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksBotID = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksBotSecret = pq.ByteaArray{[]byte("123456789")}
			l.RefreshTokenExpiresAt = pq.NullTime{Time: now, Valid: true}
			l.DefaultChannelID = "123456789"
			l.DebugMode = false
		}),
	)

	repo := NewRepository(tx)
	t.Run("全てのLineWorksBotを取得できること", func(t *testing.T) {
		bots, err := repo.GetAllLineWorksBots(ctx)
		assert.NoError(t, err)
		assert.Len(t, bots, 1)
		assert.Equal(t, "987654321", bots[0].GuildID)
		assert.Equal(t, "123456789", string(bots[0].LineWorksBotToken[0]))
		assert.Equal(t, "123456789", string(bots[0].LineWorksRefreshToken[0]))
		assert.Equal(t, "123456789", string(bots[0].LineWorksGroupID[0]))
		assert.Equal(t, "123456789", string(bots[0].LineWorksBotID[0]))
		assert.Equal(t, "123456789", string(bots[0].LineWorksBotSecret[0]))
		assert.Equal(t, now.Format(time.RFC3339), bots[0].RefreshTokenExpiresAt.Time.Format(time.RFC3339))
		assert.Equal(t, "123456789", bots[0].DefaultChannelID)
		assert.False(t, bots[0].DebugMode)
	})
}

func TestGetLineWorksBotByGuildID(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot")
	assert.NoError(t, err)

	now := time.Now().UTC()

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineWorksBot(ctx, func(l *fixtures.LineWorksBot) {
			l.GuildID = "987654321"
			l.LineWorksBotToken = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksRefreshToken = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksGroupID = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksBotID = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksBotSecret = pq.ByteaArray{[]byte("123456789")}
			l.RefreshTokenExpiresAt = pq.NullTime{Time: now, Valid: true}
			l.DefaultChannelID = "123456789"
			l.DebugMode = false
		}),
	)

	repo := NewRepository(tx)
	t.Run("指定したGuildIDのLineWorksBotを取得できること", func(t *testing.T) {
		bot, err := repo.GetLineWorksBotByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", bot.GuildID)
		assert.Equal(t, "123456789", string(bot.LineWorksBotToken[0]))
		assert.Equal(t, "123456789", string(bot.LineWorksRefreshToken[0]))
		assert.Equal(t, "123456789", string(bot.LineWorksGroupID[0]))
		assert.Equal(t, "123456789", string(bot.LineWorksBotID[0]))
		assert.Equal(t, "123456789", string(bot.LineWorksBotSecret[0]))
		assert.Equal(t, now.Format(time.RFC3339), bot.RefreshTokenExpiresAt.Time.Format(time.RFC3339))
		assert.Equal(t, "123456789", bot.DefaultChannelID)
		assert.False(t, bot.DebugMode)
	})

	t.Run("指定したGuildIDのLineWorksBotが存在しない場合はエラーが返ること", func(t *testing.T) {
		_, err := repo.GetLineWorksBotByGuildID(ctx, "123456789")
		assert.Error(t, err)
	})
}

func TestInsertLineWorksBot(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot")
	assert.NoError(t, err)

	repo := NewRepository(tx)
	t.Run("LineWorksBotが新規作成されること", func(t *testing.T) {
		err := repo.InsertLineWorksBot(ctx, &LineWorksBot{
			GuildID:               "987654321",
			LineWorksBotToken:     pq.ByteaArray{[]byte("123456789")},
			LineWorksRefreshToken: pq.ByteaArray{[]byte("123456789")},
			LineWorksGroupID:      pq.ByteaArray{[]byte("123456789")},
			LineWorksBotID:        pq.ByteaArray{[]byte("123456789")},
			LineWorksBotSecret:    pq.ByteaArray{[]byte("123456789")},
			RefreshTokenExpiresAt: pq.NullTime{Valid: false},
			DefaultChannelID:      "123456789",
			DebugMode:             false,
		})
		assert.NoError(t, err)

		bot, err := repo.GetLineWorksBotByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", bot.GuildID)
		assert.Equal(t, "123456789", string(bot.LineWorksBotToken[0]))
		assert.Equal(t, "123456789", string(bot.LineWorksRefreshToken[0]))
		assert.Equal(t, "123456789", string(bot.LineWorksGroupID[0]))
		assert.Equal(t, "123456789", string(bot.LineWorksBotID[0]))
		assert.Equal(t, "123456789", string(bot.LineWorksBotSecret[0]))
		assert.False(t, bot.RefreshTokenExpiresAt.Valid)
		assert.Equal(t, "123456789", bot.DefaultChannelID)
		assert.False(t, bot.DebugMode)
	})
}

func TestUpdateLineWorksBot(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	_, err = tx.ExecContext(ctx, "DELETE FROM line_works_bot")
	assert.NoError(t, err)

	now := time.Now().UTC()

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineWorksBot(ctx, func(l *fixtures.LineWorksBot) {
			l.GuildID = "987654321"
			l.LineWorksBotToken = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksRefreshToken = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksGroupID = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksBotID = pq.ByteaArray{[]byte("123456789")}
			l.LineWorksBotSecret = pq.ByteaArray{[]byte("123456789")}
			l.RefreshTokenExpiresAt = pq.NullTime{Time: now, Valid: true}
			l.DefaultChannelID = "123456789"
			l.DebugMode = false
		}),
	)

	repo := NewRepository(tx)
	t.Run("LineWorksBotが更新されること", func(t *testing.T) {
		err := repo.UpdateLineWorksBot(ctx, &LineWorksBot{
			GuildID:               "987654321",
			LineWorksBotToken:     pq.ByteaArray{[]byte("987654321")},
			LineWorksRefreshToken: pq.ByteaArray{[]byte("987654321")},
			LineWorksGroupID:      pq.ByteaArray{[]byte("987654321")},
			LineWorksBotID:        pq.ByteaArray{[]byte("987654321")},
			LineWorksBotSecret:    pq.ByteaArray{[]byte("987654321")},
			RefreshTokenExpiresAt: pq.NullTime{Time: now, Valid: true},
			DefaultChannelID:      "987654321",
			DebugMode:             true,
		})
		assert.NoError(t, err)

		bot, err := repo.GetLineWorksBotByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", bot.GuildID)
		assert.Equal(t, "987654321", string(bot.LineWorksBotToken[0]))
		assert.Equal(t, "987654321", string(bot.LineWorksRefreshToken[0]))
		assert.Equal(t, "987654321", string(bot.LineWorksGroupID[0]))
		assert.Equal(t, "987654321", string(bot.LineWorksBotID[0]))
		assert.Equal(t, "987654321", string(bot.LineWorksBotSecret[0]))
		assert.Equal(t, now.Format(time.RFC3339), bot.RefreshTokenExpiresAt.Time.Format(time.RFC3339))
		assert.Equal(t, "987654321", bot.DefaultChannelID)
		assert.True(t, bot.DebugMode)
	})
}
