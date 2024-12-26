package fixtures

import (
	"context"
	"testing"

	"github.com/lib/pq"
)

type LineWorksBotInfo struct {
	GuildID string `db:"guild_id"`
	LineWorksClientID pq.ByteaArray `db:"line_works_client_id"`
	LineWorksClientSecret pq.ByteaArray `db:"line_works_client_secret"`
	LineWorksServiceAccount pq.ByteaArray `db:"line_works_service_account"`
	LineWorksPrivateKey pq.ByteaArray `db:"line_works_private_key"`
	LineWorksDomainID pq.ByteaArray `db:"line_works_domain_id"`
	LineWorksAdminID pq.ByteaArray `db:"line_works_admin_id"`
}

func NewLineWorksBotInfo(ctx context.Context, setter ...func(b *LineWorksBotInfo)) *ModelConnector {
	bytes := make([]byte, 13)
	copy(bytes, []byte("1111111111111"))

	lineWorksBotInfo := &LineWorksBotInfo{
		GuildID:               "1111111111111",
		LineWorksClientID:     pq.ByteaArray{bytes},
		LineWorksClientSecret: pq.ByteaArray{bytes},
		LineWorksServiceAccount:      pq.ByteaArray{bytes},
		LineWorksPrivateKey:        pq.ByteaArray{bytes},
		LineWorksDomainID:    pq.ByteaArray{bytes},
		LineWorksAdminID: pq.ByteaArray{bytes},
	}
	return &ModelConnector{
		Model: lineWorksBotInfo,
		setter: func() {
			for _, s := range setter {
				s(lineWorksBotInfo)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.LineWorksBotInfos = append(f.LineWorksBotInfos, lineWorksBotInfo)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, lineWorksBotInfo)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO line_works_bot_info (
					guild_id,
					line_works_client_id,
					line_works_client_secret,
					line_works_service_account,
					line_works_private_key,
					line_works_domain_id,
					line_works_admin_id
				) VALUES (
					:guild_id,
					:line_works_client_id,
					:line_works_client_secret,
					:line_works_service_account,
					:line_works_private_key,
					:line_works_domain_id,
					:line_works_admin_id
				)
			`, lineWorksBotInfo)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
