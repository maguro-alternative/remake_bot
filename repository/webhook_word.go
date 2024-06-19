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
