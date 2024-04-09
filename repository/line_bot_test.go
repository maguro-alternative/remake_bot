package repository

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"
	"github.com/maguro-alternative/remake_bot/web/config"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestRepository_InsertLineBots(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM line_bot")

	repo := NewRepository(tx)
	t.Run("LineBotが正しく登録されること", func(t *testing.T) {
		lineBot := LineBot{
			GuildID:          "987654321",
			LineNotifyToken:  pq.ByteaArray{[]byte("123456789")},
			LineBotToken:     pq.ByteaArray{[]byte("123456789")},
			LineBotSecret:    pq.ByteaArray{[]byte("123456789")},
			LineGroupID:      pq.ByteaArray{[]byte("987654321")},
			LineClientID:     pq.ByteaArray{[]byte("123456789")},
			LineClientSecret: pq.ByteaArray{[]byte("123456789")},
			DefaultChannelID: "987654321",
			DebugMode:        false,
		}
		err := repo.InsertLineBot(ctx, &lineBot)
		assert.NoError(t, err)

		var lineBotRes LineBot
		err = tx.GetContext(ctx, &lineBotRes, "SELECT * FROM line_bot WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBotRes.GuildID)
		assert.Equal(t, "987654321", lineBotRes.DefaultChannelID)
		assert.Equal(t, false, lineBotRes.DebugMode)
	})
}

func TestRepository_GetLineBots(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM line_bot")

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineBot(ctx, func(lb *fixtures.LineBot) {
			lb.GuildID = "987654321"
			lb.LineNotifyToken = pq.ByteaArray{[]byte("123456789")}
			lb.LineBotToken = pq.ByteaArray{[]byte("123456789")}
			lb.LineBotSecret = pq.ByteaArray{[]byte("123456789")}
			lb.LineGroupID = pq.ByteaArray{[]byte("987654321")}
			lb.LineClientID = pq.ByteaArray{[]byte("123456789")}
			lb.LineClientSecret = pq.ByteaArray{[]byte("123456789")}
			lb.DefaultChannelID = "987654321"
			lb.DebugMode = false
		}),
	)

	repo := NewRepository(tx)
	t.Run("GuildIDからLineBotを取得できること", func(t *testing.T) {
		lineBots, err := repo.GetAllColumnsLineBots(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(lineBots))
		assert.Equal(t, "987654321", lineBots[0].GuildID)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBots[0].LineNotifyToken)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBots[0].LineBotToken)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBots[0].LineBotSecret)
		assert.Equal(t, pq.ByteaArray{[]byte("987654321")}, lineBots[0].LineGroupID)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBots[0].LineClientID)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBots[0].LineClientSecret)
		assert.Equal(t, "987654321", lineBots[0].DefaultChannelID)
		assert.Equal(t, false, lineBots[0].DebugMode)
	})
}

func TestRepository_GetLineBot(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM line_bot")

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineBot(ctx, func(lb *fixtures.LineBot) {
			lb.GuildID = "987654321"
			lb.LineNotifyToken = pq.ByteaArray{[]byte("123456789")}
			lb.LineBotToken = pq.ByteaArray{[]byte("123456789")}
			lb.LineBotSecret = pq.ByteaArray{[]byte("123456789")}
			lb.LineGroupID = pq.ByteaArray{[]byte("987654321")}
			lb.LineClientID = pq.ByteaArray{[]byte("123456789")}
			lb.LineClientSecret = pq.ByteaArray{[]byte("123456789")}
			lb.DefaultChannelID = "987654321"
			lb.DebugMode = false
		}),
	)

	repo := NewRepository(tx)
	t.Run("GuildIDからLineBotを取得できること", func(t *testing.T) {
		lineBot, err := repo.GetAllColumnsLineBot(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBot.GuildID)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBot.LineNotifyToken)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBot.LineBotToken)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBot.LineBotSecret)
		assert.Equal(t, pq.ByteaArray{[]byte("987654321")}, lineBot.LineGroupID)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBot.LineClientID)
		assert.Equal(t, pq.ByteaArray{[]byte("123456789")}, lineBot.LineClientSecret)
		assert.Equal(t, "987654321", lineBot.DefaultChannelID)
		assert.Equal(t, false, lineBot.DebugMode)
	})
}

func TestGetLineBotNotClient(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM line_bot")

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
			lb.LineNotifyToken = pq.ByteaArray{[]byte("123456789")}
			lb.LineBotToken = pq.ByteaArray{[]byte("123456789")}
			lb.LineBotSecret = pq.ByteaArray{[]byte("123456789")}
			lb.LineGroupID = pq.ByteaArray{[]byte("987654321")}
			lb.DefaultChannelID = "987654321"
			lb.DebugMode = false
		}),
		fixtures.NewLineBot(ctx, func(lb *fixtures.LineBot) {
			lb.GuildID = "123456789"
			lb.LineNotifyToken = pq.ByteaArray{[]byte("X+P6kmO6DnEjM3TVqXkwNA==")} //95 227 250 146 99 186 14 113 35 51 116 213 169 121 48 52
			lb.LineBotToken = pq.ByteaArray{[]byte("uy2qtvYTnSoB5qIntwUdVQ==")}
			lb.LineBotSecret = pq.ByteaArray{[]byte("i2uHQCyn58wRR/b03fRw6w==")}
			lb.LineGroupID = pq.ByteaArray{[]byte("YgexFQQlLcaXmsw9mFN35Q==")}
			lb.DefaultChannelID = "987654321"
			lb.DebugMode = false
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからClient以外のLineBotを取得できること", func(t *testing.T) {
		lineBot, err := repo.GetLineBotNotClient(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, []byte("123456789"), lineBot.LineNotifyToken[0])
		assert.Equal(t, []byte("123456789"), lineBot.LineBotToken[0])
		assert.Equal(t, []byte("123456789"), lineBot.LineBotSecret[0])
		assert.Equal(t, []byte("987654321"), lineBot.LineGroupID[0])
		assert.Equal(t, "987654321", lineBot.DefaultChannelID)
		assert.Equal(t, false, lineBot.DebugMode)
	})

	t.Run("GuildIDからClient以外のLineBotを取得し、復号できること", func(t *testing.T) {
		lineBot, err := repo.GetLineBotNotClient(ctx, "123456789")
		assert.NoError(t, err)

		decodeNotifyToken, err := hex.DecodeString("aa7c5fe80002633327f0fefe67a565de")
		assert.NoError(t, err)
		lineNotifyStr, err := base64.StdEncoding.DecodeString(string(lineBot.LineNotifyToken[0]))
		assert.NoError(t, err)
		decodeBotToken, err := hex.DecodeString("baeff317cb83ef55b193b6d3de194124")
		assert.NoError(t, err)
		lineBotStr, err := base64.StdEncoding.DecodeString(string(lineBot.LineBotToken[0]))
		assert.NoError(t, err)
		decodeBotSecret, err := hex.DecodeString("0ffa8ed72efcb5f1d834e4ce8463a62c")
		assert.NoError(t, err)
		lineBotSecretStr, err := base64.StdEncoding.DecodeString(string(lineBot.LineBotSecret[0]))
		assert.NoError(t, err)
		decodeGroupID, err := hex.DecodeString("e14db710b23520766fd652c0f19d437a")
		assert.NoError(t, err)
		lineGroupStr, err := base64.StdEncoding.DecodeString(string(lineBot.LineGroupID[0]))
		assert.NoError(t, err)

		notifyTokenDecrypted, err := crypto.Decrypt(lineNotifyStr, key, decodeNotifyToken)
		assert.NoError(t, err)
		botTokenDecrypted, err := crypto.Decrypt(lineBotStr, key, decodeBotToken)
		assert.NoError(t, err)
		botSecretDecrypted, err := crypto.Decrypt(lineBotSecretStr, key, decodeBotSecret)
		assert.NoError(t, err)
		groupIDDecrypted, err := crypto.Decrypt(lineGroupStr, key, decodeGroupID)
		assert.NoError(t, err)

		assert.Equal(t, notifyToken, string(notifyTokenDecrypted))
		assert.Equal(t, botToken, string(botTokenDecrypted))
		assert.Equal(t, botSecret, string(botSecretDecrypted))
		assert.Equal(t, groupID, string(groupIDDecrypted))
		assert.Equal(t, "987654321", lineBot.DefaultChannelID)
		assert.Equal(t, false, lineBot.DebugMode)
	})
}

func TestRepository_GetLineBotDefaultChannelID(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM line_bot")

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineBot(ctx, func(lb *fixtures.LineBot) {
			lb.GuildID = "987654321"
			lb.LineNotifyToken = pq.ByteaArray{[]byte("123456789")}
			lb.LineBotToken = pq.ByteaArray{[]byte("123456789")}
			lb.LineBotSecret = pq.ByteaArray{[]byte("123456789")}
			lb.LineGroupID = pq.ByteaArray{[]byte("987654321")}
			lb.LineClientID = pq.ByteaArray{[]byte("123456789")}
			lb.LineClientSecret = pq.ByteaArray{[]byte("123456789")}
			lb.DefaultChannelID = "987654321"
			lb.DebugMode = false
		}),
	)

	repo := NewRepository(tx)
	t.Run("GuildIDからLineBotのDefaultChannelUDを取得できること", func(t *testing.T) {
		lineBot, err := repo.GetLineBotDefaultChannelID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBot.DefaultChannelID)
	})
}

func TestRepository_UpdateLineBot(t *testing.T) {
	ctx := context.Background()
	t.Run("LineBotが正しく更新されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM line_bot")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewLineBot(ctx, func(lb *fixtures.LineBot) {
				lb.GuildID = "987654321"
				lb.LineNotifyToken = pq.ByteaArray{[]byte("123456789")}
				lb.LineBotToken = pq.ByteaArray{[]byte("123456789")}
				lb.LineBotSecret = pq.ByteaArray{[]byte("123456789")}
				lb.LineGroupID = pq.ByteaArray{[]byte("987654321")}
				lb.LineClientID = pq.ByteaArray{[]byte("987654321")}
				lb.LineClientSecret = pq.ByteaArray{[]byte("987654321")}
				lb.DefaultChannelID = "987654321"
				lb.DebugMode = false
			}),
		)

		repo := NewRepository(tx)
		updateLineBot := &LineBot{
			GuildID:          "987654321",
			DefaultChannelID: "123456789",
		}
		err = repo.UpdateLineBot(ctx, updateLineBot)
		assert.NoError(t, err)

		var lineBot LineBot
		err = tx.GetContext(ctx, &lineBot, "SELECT * FROM line_bot WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBot.GuildID)
		assert.Equal(t, "123456789", lineBot.DefaultChannelID)
	})

	t.Run("LineBotの1部分(notifyとbottoken)が正しく更新されること", func(t *testing.T) {
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
				lb.LineNotifyToken = pq.ByteaArray{[]byte("123456789")}
				lb.LineBotToken = pq.ByteaArray{[]byte("123456789")}
				lb.LineBotSecret = pq.ByteaArray{[]byte("123456789")}
				lb.LineGroupID = pq.ByteaArray{[]byte("987654321")}
				lb.LineClientID = pq.ByteaArray{[]byte("987654321")}
				lb.LineClientSecret = pq.ByteaArray{[]byte("987654321")}
				lb.DefaultChannelID = "987654321"
				lb.DebugMode = false
			}),
		)

		repo := NewRepository(tx)
		updateLineBot := &LineBot{
			GuildID:          "987654321",
			DefaultChannelID: "123456789",
		}
		err = repo.UpdateLineBot(ctx, updateLineBot)
		assert.NoError(t, err)

		var lineBot LineBot
		err = tx.GetContext(ctx, &lineBot, "SELECT * FROM line_bot WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", lineBot.GuildID)
		assert.Equal(t, "123456789", lineBot.DefaultChannelID)
	})
}
