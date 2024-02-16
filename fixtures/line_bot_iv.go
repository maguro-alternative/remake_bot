package fixtures

import (
	"context"
	"testing"
)

type LineBotIv struct {
	GuildID            string `db:"guild_id"`
	LineNotifyTokenIv  []byte `db:"line_notify_token_iv"`
	LineBotTokenIv     []byte `db:"line_bot_token_iv"`
	LineBotSecretIv    []byte `db:"line_bot_secret_iv"`
	LineGroupIDIv      []byte `db:"line_group_id_iv"`
	LineClientIDIv     []byte `db:"line_client_id_iv"`
	LineClientSecretIv []byte `db:"line_client_secret_iv"`
}

func NewLineBotIv(ctx context.Context, setter ...func(b *LineBotIv)) *ModelConnector {
	lineBotIv := &LineBotIv{
		GuildID:            "1111111111111",
		LineNotifyTokenIv:  []byte("1111111111111"),
		LineBotTokenIv:     []byte("1111111111111"),
		LineBotSecretIv:    []byte("1111111111111"),
		LineGroupIDIv:      []byte("1111111111111"),
		LineClientIDIv:     []byte("1111111111111"),
		LineClientSecretIv: []byte("1111111111111"),
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
