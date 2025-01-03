package fixtures

import (
	"context"
	"testing"

	"github.com/lib/pq"
)

type LineWorksBotInfoIV struct {
	GuildID string `db:"guild_id"`
	LineWorksClientIDIV pq.ByteaArray `db:"line_works_client_id_iv"`
	LineWorksClientSecretIV pq.ByteaArray `db:"line_works_client_secret_iv"`
	LineWorksServiceAccountIV pq.ByteaArray `db:"line_works_service_account_iv"`
	LineWorksPrivateKeyIV pq.ByteaArray `db:"line_works_private_key_iv"`
	LineWorksDomainIDIV pq.ByteaArray `db:"line_works_domain_id_iv"`
	LineWorksAdminIDIV pq.ByteaArray `db:"line_works_admin_id_iv"`
}

func NewLineWorksBotInfoIv(ctx context.Context, setter ...func(b *LineWorksBotInfoIV)) *ModelConnector {
	bytes := make([]byte, 13)
	copy(bytes, []byte("1111111111111"))

	lineWorksBotInfoIv := &LineWorksBotInfoIV{
		GuildID:               "1111111111111",
		LineWorksClientIDIV:     pq.ByteaArray{bytes},
		LineWorksClientSecretIV: pq.ByteaArray{bytes},
		LineWorksServiceAccountIV:      pq.ByteaArray{bytes},
		LineWorksPrivateKeyIV:        pq.ByteaArray{bytes},
		LineWorksDomainIDIV:    pq.ByteaArray{bytes},
		LineWorksAdminIDIV: pq.ByteaArray{bytes},
	}
	return &ModelConnector{
		Model: lineWorksBotInfoIv,
		setter: func() {
			for _, s := range setter {
				s(lineWorksBotInfoIv)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.LineWorksBotInfoIvs = append(f.LineWorksBotInfoIvs, lineWorksBotInfoIv)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, lineWorksBotInfoIv)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO line_works_bot_info_iv (
					guild_id,
					line_works_client_id_iv,
					line_works_client_secret_iv,
					line_works_service_account_iv,
					line_works_private_key_iv,
					line_works_domain_id_iv,
					line_works_admin_id_iv
				) VALUES (
					:guild_id,
					:line_works_client_id_iv,
					:line_works_client_secret_iv,
					:line_works_service_account_iv,
					:line_works_private_key_iv,
					:line_works_domain_id_iv,
					:line_works_admin_id_iv
				)
			`, lineWorksBotInfoIv)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
