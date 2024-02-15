package on_message_create

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/stretchr/testify/assert"
)

func TestGetLineChannel(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineChannel(ctx, func(lc *fixtures.LineChannel) {
			lc.ChannelID = "123456789"
			lc.GuildID = "987654321"
			lc.Ng = false
			lc.BotMessage = false
		}),
	)
	repo := NewRepository(tx)
	t.Run("ChannelIDから送信しないかどうか取得できること", func(t *testing.T) {
		channel, err := repo.GetLineChannel(ctx, "123456789")
		assert.NoError(t, err)
		assert.Equal(t, false, channel.Ng)
		assert.Equal(t, false, channel.BotMessage)
	})
}

func TestInsertLineChannel(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	repo := NewRepository(tx)

	var channel TestLineChannel
	t.Run("ChannelIDを追加できること", func(t *testing.T) {
		err := repo.InsertLineChannel(ctx, "123456789", "987654321")
		assert.NoError(t, err)
		query := `
			SELECT
				*
			FROM
				line_channel
			WHERE
				channel_id = $1
		`
		err = tx.SelectContext(ctx, &channel, query, "123456789")
		assert.NoError(t, err)

		assert.Equal(t, "123456789", channel.ChannelID)
		assert.Equal(t, "987654321", channel.GuildID)
		assert.Equal(t, false, channel.Ng)
		assert.Equal(t, false, channel.BotMessage)
	})
}

func TestGetLineNgType(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineNgType(ctx, func(lnt *fixtures.LineNgType) {
			lnt.ChannelID = "987654321"
			lnt.Type = 6
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからNGタイプを取得できること", func(t *testing.T) {
		ngTypes, err := repo.GetLineNgType(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, []int{6}, ngTypes)
	})
}

func TestGetLineBot(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineBot(ctx, func(lb *fixtures.LineBot) {
			lb.GuildID = "987654321"
			lb.LineNotifyToken = []byte("123456789")
			lb.LineBotToken = []byte("123456789")
			lb.LineBotSecret = []byte("123456789")
			lb.LineGroupID = []byte("987654321")
			lb.Iv = []byte("123456789")
			lb.DefaultChannelID = "987654321"
			lb.DebugMode = false
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからLineBotを取得できること", func(t *testing.T) {
		lineBot, err := repo.GetLineBot(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, []byte("123456789"), lineBot.LineNotifyToken)
		assert.Equal(t, []byte("123456789"), lineBot.LineBotToken)
		assert.Equal(t, []byte("123456789"), lineBot.LineBotSecret)
		assert.Equal(t, []byte("987654321"), lineBot.LineGroupID)
		assert.Equal(t, []byte("123456789"), lineBot.Iv)
		assert.Equal(t, "987654321", lineBot.DefaultChannelID)
		assert.Equal(t, false, lineBot.DebugMode)
	})
}
