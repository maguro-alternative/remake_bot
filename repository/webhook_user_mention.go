package repository

import (
	"context"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type WebhookUserMention struct {
	WebhookSerialID int64  `db:"webhook_serial_id"`
	UserID          string `db:"user_id"`
}

func (r *Repository) InsertWebhookUserMention(
	ctx context.Context,
	webhookSerialID int64,
	userID string,
) error {
	query := `
		INSERT INTO webhook_user_mention (
			webhook_serial_id,
			user_id
		) VALUES (
			$1,
			$2
		) ON CONFLICT (webhook_serial_id, user_id) DO NOTHING
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		webhookSerialID,
		userID,
	)
	return err
}

func (r *Repository) GetWebhookUserMentionWithWebhookSerialID(
	ctx context.Context,
	webhookSerialID int64,
) ([]*WebhookUserMention, error) {
	query := `
		SELECT
			*
		FROM
			webhook_user_mention
		WHERE
			webhook_serial_id = $1
	`
	var webhookUserMention []*WebhookUserMention
	err := r.db.SelectContext(ctx, &webhookUserMention, query, webhookSerialID)
	return webhookUserMention, err
}

func (r *Repository) DeleteWebhookUserMentionsNotInProvidedList(
	ctx context.Context,
	webhookSerialID int64,
	userIDs []string,
) error {
	query := `
		DELETE FROM
			webhook_user_mention
		WHERE
			webhook_serial_id = ? AND
			user_id NOT IN (?)
	`
	if len(userIDs) == 0 {
		query = `
			DELETE FROM
				webhook_user_mention
			WHERE
				webhook_serial_id = $1
		`
		_, err := r.db.ExecContext(ctx, query, webhookSerialID)
		return err
	}
	query, args, err := db.In(query, webhookSerialID, userIDs)
	if err != nil {
		return err
	}
	query = db.Rebind(2, query)
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
