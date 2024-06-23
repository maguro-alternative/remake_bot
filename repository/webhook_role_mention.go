package repository

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type WebhookRoleMention struct {
	WebhookSerialID int64  `db:"webhook_serial_id"`
	RoleID          string `db:"role_id"`
}

func (r *Repository) InsertWebhookRoleMention(
	ctx context.Context,
	webhookSerialID int64,
	roleID string,
) error {
	query := `
		INSERT INTO webhook_role_mention (
			webhook_serial_id,
			role_id
		) VALUES (
			$1,
			$2
		) ON CONFLICT (webhook_serial_id, role_id) DO NOTHING
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		webhookSerialID,
		roleID,
	)
	return err
}

func (r *Repository) GetWebhookRoleMentionWithWebhookSerialID(
	ctx context.Context,
	webhookSerialID int64,
) ([]*WebhookRoleMention, error) {
	query := `
		SELECT
			*
		FROM
			webhook_role_mention
		WHERE
			webhook_serial_id = $1
	`
	var webhookRoleMention []*WebhookRoleMention
	err := r.db.SelectContext(ctx, &webhookRoleMention, query, webhookSerialID)
	return webhookRoleMention, err
}

func (r *Repository) GetWebhookRoleMentionWithWebhookSerialIDs(
	ctx context.Context,
	webhookSerialIDs []int64,
) ([]*WebhookRoleMention, error) {
	query := `
		SELECT
			*
		FROM
			webhook_role_mention
		WHERE
			webhook_serial_id IN (?)
	`
	var webhookRoleMention []*WebhookRoleMention
	query, args, err := db.In(query, webhookSerialIDs)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(2, query)
	err = r.db.SelectContext(ctx, &webhookRoleMention, query, args...)
	return webhookRoleMention, err
}

func (r *Repository) DeleteWebhookRoleMentionsNotInProvidedList(
	ctx context.Context,
	webhookSerialID int64,
	roleIDs []string,
) error {
	query := `
		DELETE FROM
			webhook_role_mention
		WHERE
			webhook_serial_id = ? AND
			role_id NOT IN (?)
	`
	if len(roleIDs) == 0 {
		query = `
			DELETE FROM
				webhook_role_mention
			WHERE
				webhook_serial_id = $1
		`
		_, err := r.db.ExecContext(ctx, query, webhookSerialID)
		return err
	}
	query, args, err := db.In(query, webhookSerialID, roleIDs)
	if err != nil {
		return err
	}
	query = db.Rebind(2, query)
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
