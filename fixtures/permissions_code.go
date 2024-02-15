package fixtures

import (
	"context"
	"testing"
)

type PermissionsCode struct {
	GuildID    string `db:"guild_id"`
	Type       string `db:"type"`
	Code       string `db:"code"`
}

func NewPermissionsCode(ctx context.Context, setter ...func(b *PermissionsCode)) *ModelConnector {
	permissionsCode := &PermissionsCode{
		GuildID:    "1111111111111",
		Type:       "line",
		Code:       "8",
	}

	return &ModelConnector{
		Model: permissionsCode,
		setter: func() {
			for _, s := range setter {
				s(permissionsCode)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.PermissionsCodes = append(f.PermissionsCodes, permissionsCode)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, permissionsCode)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO permissions_code (
					guild_id,
					type,
					code
				) VALUES (
					:guild_id,
					:type,
					:code
				)
			`, permissionsCode)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
