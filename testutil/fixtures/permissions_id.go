package fixtures

import (
	"context"
	"testing"
)

type PermissionsID struct {
	GuildID    string `db:"guild_id"`
	Type       string `db:"type"`
	TargetType string `db:"target_type"`
	TargetID   string `db:"target_id"`
	Permission string `db:"permission"`
}

func NewPermissionsID(ctx context.Context, setter ...func(b *PermissionsID)) *ModelConnector {
	permissionsID := &PermissionsID{
		GuildID:    "1111111111111",
		Type:       "line",
		TargetType: "user",
		TargetID:   "1111111111111",
		Permission: "read",
	}

	return &ModelConnector{
		Model: permissionsID,
		setter: func() {
			for _, s := range setter {
				s(permissionsID)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.PermissionsIDs = append(f.PermissionsIDs, permissionsID)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, permissionsID)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO permissions_id (
					guild_id,
					type,
					target_type,
					target_id,
					permission
				) VALUES (
					:guild_id,
					:type,
					:target_type,
					:target_id,
					:permission
				)
			`, permissionsID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
