package fixtures

import (
	"context"
	"testing"
)

type PermissionsRoleID struct {
	GuildID    string `db:"guild_id"`
	Type       string `db:"type"`
	TargetID   string `db:"target_id"`
	Permission string `db:"permission"`
}

func NewPermissionsRoleID(ctx context.Context, setter ...func(b *PermissionsRoleID)) *ModelConnector {
	permissionsRoleID := &PermissionsRoleID{
		GuildID:    "1111111111111",
		Type:       "line",
		TargetID:   "1111111111111",
		Permission: "read",
	}

	return &ModelConnector{
		Model: permissionsRoleID,
		setter: func() {
			for _, s := range setter {
				s(permissionsRoleID)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.PermissionsRoleIDs = append(f.PermissionsRoleIDs, permissionsRoleID)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, permissionsRoleID)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO permissions_role_id (
					guild_id,
					type,
					target_id,
					permission
				) VALUES (
					:guild_id,
					:type,
					:target_id,
					:permission
				)
			`, permissionsRoleID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
