package fixtures

import (
	"context"
	"testing"

	"github.com/lib/pq"
)

type LineBotIv struct {
	GuildID            string `db:"guild_id"`
	LineNotifyTokenIv  pq.ByteaArray `db:"line_notify_token_iv"`
	LineBotTokenIv     pq.ByteaArray `db:"line_bot_token_iv"`
	LineBotSecretIv    pq.ByteaArray `db:"line_bot_secret_iv"`
	LineGroupIDIv      pq.ByteaArray `db:"line_group_id_iv"`
	LineClientIDIv     pq.ByteaArray `db:"line_client_id_iv"`
	LineClientSecretIv pq.ByteaArray `db:"line_client_secret_iv"`
}

func NewLineBotIv(ctx context.Context, setter ...func(b *LineBotIv)) *ModelConnector {
	lineBotIv := &LineBotIv{
		GuildID:            "1111111111111",
		LineNotifyTokenIv:  pq.ByteaArray{[]byte("1111111111111")},
		LineBotTokenIv:     pq.ByteaArray{[]byte("1111111111111")},
		LineBotSecretIv:    pq.ByteaArray{[]byte("1111111111111")},
		LineGroupIDIv:      pq.ByteaArray{[]byte("1111111111111")},
		LineClientIDIv:     pq.ByteaArray{[]byte("1111111111111")},
		LineClientSecretIv: pq.ByteaArray{[]byte("1111111111111")},
	}

	return &ModelConnector{
		Model: lineBotIv,
		setter: func() {
			for _, s := range setter {
				s(lineBotIv)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.LineBotIvs = append(f.LineBotIvs, lineBotIv)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, lineBotIv)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO line_bot_iv (
					guild_id,
					line_notify_token_iv,
					line_bot_token_iv,
					line_bot_secret_iv,
					line_group_id_iv,
					line_client_id_iv,
					line_client_secret_iv
				) VALUES (
					:guild_id,
					:line_notify_token_iv,
					:line_bot_token_iv,
					:line_bot_secret_iv,
					:line_group_id_iv,
					:line_client_id_iv,
					:line_client_secret_iv
				)
			`, lineBotIv)
			if err != nil {
				t.Fatal(err)
			}
		},
	}
}
