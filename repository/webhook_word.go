package repository

import (
	"context"
)

type WebhookWord struct {
	ID        int64  `db:"id"`
	Condition string `db:"condition"`
	Word      string `db:"word"`
}

func (r *Repository) InsertWebhookWord(
	ctx context.Context,
	condition string,
	word string,
) error {
	query := `
		INSERT INTO webhook_word (
			condition,
			word
		) VALUES (
			$1,
			$2
		)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		condition,
		word,
	)
	return err
}

func (r *Repository) GetWebhookWordWithWebhookSerialIDAndCondition(
	ctx context.Context,
	webhookSerialID int64,
	condition string,
) ([]*WebhookWord, error) {
	query := `
		SELECT
			*
		FROM
			webhook_word
		WHERE
			webhook_serial_id = $1
			condition = $2
	`
	var webhookWord []*WebhookWord
	err := r.db.SelectContext(ctx, &webhookWord, query, webhookSerialID, condition)
	return webhookWord, err
}
