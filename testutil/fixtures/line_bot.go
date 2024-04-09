package fixtures

import (
	"context"
	"testing"

	"github.com/lib/pq"
)

type LineBot struct {
	GuildID          string `db:"guild_id"`
	LineNotifyToken  pq.ByteaArray `db:"line_notify_token"`
	LineBotToken     pq.ByteaArray `db:"line_bot_token"`
	LineBotSecret    pq.ByteaArray `db:"line_bot_secret"`
	LineGroupID      pq.ByteaArray `db:"line_group_id"`
	LineClientID     pq.ByteaArray `db:"line_client_id"`
	LineClientSecret pq.ByteaArray `db:"line_client_secret"`
	DefaultChannelID string `db:"default_channel_id"`
	DebugMode        bool   `db:"debug_mode"`
}

func NewLineBot(ctx context.Context, setter ...func(b *LineBot)) *ModelConnector {
	bytes := make([]byte, 13)
	copy(bytes, []byte("1111111111111"))

	lineBot := &LineBot{
		GuildID:          "1111111111111",
		LineNotifyToken:  pq.ByteaArray{bytes},
		LineBotToken:     pq.ByteaArray{bytes},
		LineBotSecret:    pq.ByteaArray{bytes},
		LineGroupID:      pq.ByteaArray{bytes},
		LineClientID:     pq.ByteaArray{bytes},
		LineClientSecret: pq.ByteaArray{bytes},
		DefaultChannelID: "1111111111111",
		DebugMode:        true,
	}

	return &ModelConnector{
		Model: lineBot,
		setter: func() {
			for _, s := range setter {
				s(lineBot)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.LineBots = append(f.LineBots, lineBot)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel := connectingModel.(type) {
			case *LineBotIv:
				lineBotIv := connectingModel
				lineBotIv.GuildID = lineBot.GuildID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, lineBot)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO line_bot (
					guild_id,
					line_notify_token,
					line_bot_token,
					line_bot_secret,
					line_group_id,
					line_client_id,
					line_client_secret,
					default_channel_id,
					debug_mode
				) VALUES (
					:guild_id,
					:line_notify_token,
					:line_bot_token,
					:line_bot_secret,
					:line_group_id,
					:line_client_id,
					:line_client_secret,
					:default_channel_id,
					:debug_mode
				)
			`, lineBot)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
