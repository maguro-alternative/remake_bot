package fixtures

import (
	"context"
	"testing"
)

type PermissionsUserID struct {
	GuildID    string `db:"guild_id"`
	Type       string `db:"type"`
	UserID     string `db:"user_id"`
	Permission string `db:"permission"`
}

func NewPermissionsUserID(ctx context.Context, setter ...func(b *PermissionsUserID)) *ModelConnector {
	permissionsUserID := &PermissionsUserID{
		GuildID:    "1111111111111",
		Type:       "line",
		UserID:     "1111111111111",
		Permission: "read",
	}

	return &ModelConnector{
		Model: permissionsUserID,
		setter: func() {
			for _, s := range setter {
				s(permissionsUserID)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.PermissionsUserIDs = append(f.PermissionsUserIDs, permissionsUserID)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, permissionsUserID)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO permissions_user_id (
					guild_id,
					type,
					user_id,
					permission
				) VALUES (
					:guild_id,
					:type,
					:user_id,
					:permission
				)
			`, permissionsUserID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
