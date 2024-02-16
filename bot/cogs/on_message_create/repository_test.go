package on_message_create

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/crypto"
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

	keyString := "645E739A7F9F162725C1533DC2C5E827"
	key, err := hex.DecodeString(keyString)
	assert.NoError(t, err)

	notifyToken := "testnotifytoken"
	botToken := "testbottoken"
	botSecret := "testbotsecret"
	groupID := "testgroupid"

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineBot(ctx, func(lb *fixtures.LineBot) {
			lb.GuildID = "987654321"
			lb.LineNotifyToken = []byte("123456789")
			lb.LineBotToken = []byte("123456789")
			lb.LineBotSecret = []byte("123456789")
			lb.LineGroupID = []byte("987654321")
			lb.DefaultChannelID = "987654321"
			lb.DebugMode = false
		}),
		fixtures.NewLineBot(ctx, func(lb *fixtures.LineBot) {
			lb.GuildID = "123456789"
			lb.LineNotifyToken = []byte("nTOK7MAo4X69eZu/0rg0Gw==")
			lb.LineBotToken = []byte("uy2qtvYTnSoB5qIntwUdVQ==")
			lb.LineBotSecret = []byte("i2uHQCyn58wRR/b03fRw6w==")
			lb.LineGroupID = []byte("YgexFQQlLcaXmsw9mFN35Q==")
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
		assert.Equal(t, "987654321", lineBot.DefaultChannelID)
		assert.Equal(t, false, lineBot.DebugMode)
	})

	t.Run("GuildIDからLineBotを取得できること", func(t *testing.T) {
		lineBot, err := repo.GetLineBot(ctx, "123456789")
		assert.NoError(t, err)

		notifyTokenDecrypted, err := crypto.Decrypt(lineBot.LineNotifyToken, key, []byte("fc18bf369f91199353d81b4530beb521"))
		assert.NoError(t, err)
		botTokenDecrypted, err := crypto.Decrypt(lineBot.LineBotToken, key, []byte("baeff317cb83ef55b193b6d3de194124"))
		assert.NoError(t, err)
		botSecretDecrypted, err := crypto.Decrypt(lineBot.LineBotSecret, key, []byte("0ffa8ed72efcb5f1d834e4ce8463a62c"))
		assert.NoError(t, err)
		groupIDDecrypted, err := crypto.Decrypt(lineBot.LineGroupID, key, []byte("e14db710b23520766fd652c0f19d437a"))
		assert.NoError(t, err)

		assert.Equal(t, notifyToken, string(notifyTokenDecrypted))
		assert.Equal(t, botToken, string(botTokenDecrypted))
		assert.Equal(t, botSecret, string(botSecretDecrypted))
		assert.Equal(t, groupID, string(groupIDDecrypted))
		assert.Equal(t, "987654321", lineBot.DefaultChannelID)
		assert.Equal(t, false, lineBot.DebugMode)
	})
}

func TestGetLineBotIv(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineBotIv(ctx, func(lbi *fixtures.LineBotIv) {
			lbi.GuildID = "987654321"
			lbi.LineNotifyTokenIv = []byte("123456789")
			lbi.LineBotTokenIv = []byte("123456789")
			lbi.LineBotSecretIv = []byte("123456789")
			lbi.LineGroupIDIv = []byte("987654321")
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからLineBotIvを取得できること", func(t *testing.T) {
		lineBotIv, err := repo.GetLineBotIv(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, []byte("123456789"), lineBotIv.LineNotifyTokenIv)
		assert.Equal(t, []byte("123456789"), lineBotIv.LineBotTokenIv)
		assert.Equal(t, []byte("123456789"), lineBotIv.LineBotSecretIv)
		assert.Equal(t, []byte("987654321"), lineBotIv.LineGroupIDIv)
	})
}
