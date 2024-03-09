package internal


import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/crypto"
	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

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
			lb.LineNotifyToken = pq.ByteaArray{[]byte("123456789")}
			lb.LineBotToken = pq.ByteaArray{[]byte("123456789")}
			lb.LineBotSecret = pq.ByteaArray{[]byte("123456789")}
			lb.LineGroupID = pq.ByteaArray{[]byte("987654321")}
			lb.DefaultChannelID = "987654321"
			lb.DebugMode = false
		}),
		fixtures.NewLineBot(ctx, func(lb *fixtures.LineBot) {
			lb.GuildID = "123456789"
			lb.LineNotifyToken = pq.ByteaArray{[]byte("X+P6kmO6DnEjM3TVqXkwNA==")}//95 227 250 146 99 186 14 113 35 51 116 213 169 121 48 52
			lb.LineBotToken = pq.ByteaArray{[]byte("uy2qtvYTnSoB5qIntwUdVQ==")}
			lb.LineBotSecret = pq.ByteaArray{[]byte("i2uHQCyn58wRR/b03fRw6w==")}
			lb.LineGroupID = pq.ByteaArray{[]byte("YgexFQQlLcaXmsw9mFN35Q==")}
			lb.DefaultChannelID = "987654321"
			lb.DebugMode = false
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからLineBotを取得できること", func(t *testing.T) {
		lineBot, err := repo.GetLineBot(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, []byte("123456789"), lineBot.LineNotifyToken[0])
		assert.Equal(t, []byte("123456789"), lineBot.LineBotToken[0])
		assert.Equal(t, []byte("123456789"), lineBot.LineBotSecret[0])
		assert.Equal(t, []byte("987654321"), lineBot.LineGroupID[0])
		assert.Equal(t, "987654321", lineBot.DefaultChannelID)
		assert.Equal(t, false, lineBot.DebugMode)
	})

	t.Run("GuildIDからLineBotを取得し、復号できること", func(t *testing.T) {
		lineBot, err := repo.GetLineBot(ctx, "123456789")
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
			lbi.LineNotifyTokenIv = pq.ByteaArray{[]byte("123456789")}
			lbi.LineBotTokenIv = pq.ByteaArray{[]byte("123456789")}
			lbi.LineBotSecretIv = pq.ByteaArray{[]byte("123456789")}
			lbi.LineGroupIDIv = pq.ByteaArray{[]byte("987654321")}
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからLineBotIvを取得できること", func(t *testing.T) {
		lineBotIv, err := repo.GetLineBotIv(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, []byte("123456789"), lineBotIv.LineNotifyTokenIv[0])
		assert.Equal(t, []byte("123456789"), lineBotIv.LineBotTokenIv[0])
		assert.Equal(t, []byte("123456789"), lineBotIv.LineBotSecretIv[0])
		assert.Equal(t, []byte("987654321"), lineBotIv.LineGroupIDIv[0])
	})
}

func TestGetPermissionsCode(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewPermissionsCode(ctx, func(p *fixtures.PermissionsCode) {
			p.GuildID = "987654321"
			p.Type = "line_bot"
			p.Code = 8
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからPermissionsCodeを取得できること", func(t *testing.T) {
		permissionCode, err := repo.GetPermissionCode(ctx, "987654321", "line_bot")
		assert.NoError(t, err)
		assert.Equal(t, int64(8), permissionCode)
	})
}

func TestGetPermissionIDs(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewPermissionsID(ctx, func(p *fixtures.PermissionsID) {
			p.GuildID = "987654321"
			p.TargetID = "123456789"
			p.TargetType = "user"
			p.Type = "line_bot"
			p.Permission = "read"
		}),
		fixtures.NewPermissionsID(ctx, func(p *fixtures.PermissionsID) {
			p.GuildID = "987654321"
			p.TargetID = "345678912"
			p.TargetType = "user"
			p.Type = "line_bot"
			p.Permission = "write"
		}),
		fixtures.NewPermissionsID(ctx, func(p *fixtures.PermissionsID) {
			p.GuildID = "987654321"
			p.TargetID = "567891234"
			p.TargetType = "user"
			p.Type = "line_bot"
			p.Permission = "all"
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからPermissionIDを取得できること", func(t *testing.T) {
		permissionIDs, err := repo.GetPermissionIDs(ctx, "987654321", "line_bot")
		assert.NoError(t, err)
		assert.Equal(t, "123456789", permissionIDs[0].TargetID)
		assert.Equal(t, "user", permissionIDs[0].TargetType)
		assert.Equal(t, "read", permissionIDs[0].Permission)
		assert.Equal(t, "345678912", permissionIDs[1].TargetID)
		assert.Equal(t, "user", permissionIDs[1].TargetType)
		assert.Equal(t, "write", permissionIDs[1].Permission)
		assert.Equal(t, "567891234", permissionIDs[2].TargetID)
		assert.Equal(t, "user", permissionIDs[2].TargetType)
		assert.Equal(t, "all", permissionIDs[2].Permission)
	})
}
