package fixtures

import (
	"context"
	"testing"

	"github.com/lib/pq"
)

type LineWorksBotIV struct {
	GuildID               string        `db:"guild_id"`
	LineWorksBotTokenIV     pq.ByteaArray `db:"line_works_bot_token_iv"`
	LineWorksRefreshTokenIV pq.ByteaArray `db:"line_works_refresh_token_iv"`
	LineWorksGroupIDIV      pq.ByteaArray `db:"line_works_group_id_iv"`
	LineWorksBotIDIV        pq.ByteaArray `db:"line_works_bot_id_iv"`
	LineWorksBotSecretIV    pq.ByteaArray `db:"line_works_bot_secret_iv"`
}

func NewLineWorksBotIv(ctx context.Context, setter ...func(b *LineWorksBotIV)) *ModelConnector {
	bytes := make([]byte, 13)
	copy(bytes, []byte("1111111111111"))

	lineWorksBotIv := &LineWorksBotIV{
		GuildID:               "1111111111111",
		LineWorksBotTokenIV:     pq.ByteaArray{bytes},
		LineWorksRefreshTokenIV: pq.ByteaArray{bytes},
		LineWorksGroupIDIV:      pq.ByteaArray{bytes},
		LineWorksBotIDIV:        pq.ByteaArray{bytes},
		LineWorksBotSecretIV:    pq.ByteaArray{bytes},
	}
	return &ModelConnector{
		Model: lineWorksBotIv,
		setter: func() {
			for _, s := range setter {
				s(lineWorksBotIv)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.LineWorksBotIvs = append(f.LineWorksBotIvs, lineWorksBotIv)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, lineWorksBotIv)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO line_works_bot_iv (
					guild_id,
					line_works_bot_token_iv,
					line_works_refresh_token_iv,
					line_works_group_id_iv,
					line_works_bot_id_iv,
					line_works_bot_secret_iv
				) VALUES (
					:guild_id,
					:line_works_bot_token_iv,
					:line_works_refresh_token_iv,
					:line_works_group_id_iv,
					:line_works_bot_id_iv,
					:line_works_bot_secret_iv
				)
			`, lineWorksBotIv)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}

