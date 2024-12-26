package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type LineWorksBotIV struct {
	GuildID               string        `db:"guild_id"`
	LineWorksBotTokenIV     pq.ByteaArray `db:"line_works_bot_token_iv"`
	LineWorksRefreshTokenIV pq.ByteaArray `db:"line_works_refresh_token_iv"`
	LineWorksGroupIDIV      pq.ByteaArray `db:"line_works_group_id_iv"`
	LineWorksBotIDIV        pq.ByteaArray `db:"line_works_bot_id_iv"`
	LineWorksBotSecretIV    pq.ByteaArray `db:"line_works_bot_secret_iv"`
}

func NewLineWorksBotIV(
	guildID string,
	lineWorksBotTokenIV pq.ByteaArray,
	lineWorksRefreshTokenIV pq.ByteaArray,
	lineWorksGroupIDIV pq.ByteaArray,
	lineWorksBotIDIV pq.ByteaArray,
	lineWorksBotSecretIV pq.ByteaArray,
) *LineWorksBotIV {
	return &LineWorksBotIV{
		GuildID:               guildID,
		LineWorksBotTokenIV:     lineWorksBotTokenIV,
		LineWorksRefreshTokenIV: lineWorksRefreshTokenIV,
		LineWorksGroupIDIV:      lineWorksGroupIDIV,
		LineWorksBotIDIV:        lineWorksBotIDIV,
		LineWorksBotSecretIV:    lineWorksBotSecretIV,
	}
}

func (r *Repository) InsertLineWorksBotIV(ctx context.Context, lineWorksBotIV *LineWorksBotIV) error {
	query := `
		INSERT INTO
			line_works_bot_iv (
				guild_id,
				line_works_bot_token_iv,
				line_works_refresh_token_iv,
				line_works_group_id_iv,
				line_works_bot_id_iv,
				line_works_bot_secret_iv
			)
		VALUES (
			:guild_id,
			:line_works_bot_token_iv,
			:line_works_refresh_token_iv,
			:line_works_group_id_iv,
			:line_works_bot_id_iv,
			:line_works_bot_secret_iv
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, lineWorksBotIV)
	return err
}

func (r *Repository) GetLineWorksBotIVByGuildID(ctx context.Context, guildID string) (*LineWorksBotIV, error) {
	var lineWorksBotIV LineWorksBotIV
	query := `
		SELECT
			*
		FROM
			line_works_bot_iv
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &lineWorksBotIV, query, guildID)
	return &lineWorksBotIV, err
}

func (r *Repository) UpdateLineWorksBotIV(ctx context.Context, lineWorksBotIV *LineWorksBotIV) error {
	var setNameQuery string
	var setQueryArray []string

	if len(lineWorksBotIV.LineWorksBotTokenIV) > 0 {
		setQueryArray = append(setQueryArray, "line_works_bot_token_iv = :line_works_bot_token_iv")
	}
	if len(lineWorksBotIV.LineWorksRefreshTokenIV) > 0 {
		setQueryArray = append(setQueryArray, "line_works_refresh_token_iv = :line_works_refresh_token_iv")
	}
	if len(lineWorksBotIV.LineWorksGroupIDIV) > 0 {
		setQueryArray = append(setQueryArray, "line_works_group_id_iv = :line_works_group_id_iv")
	}
	if len(lineWorksBotIV.LineWorksBotIDIV) > 0 {
		setQueryArray = append(setQueryArray, "line_works_bot_id_iv = :line_works_bot_id_iv")
	}
	if len(lineWorksBotIV.LineWorksBotSecretIV) > 0 {
		setQueryArray = append(setQueryArray, "line_works_bot_secret_iv = :line_works_bot_secret_iv")
	}
	setNameQuery = strings.Join(setQueryArray, ",")
	if setNameQuery == "" {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE
			line_works_bot_iv
		SET
			%s
		WHERE
			guild_id = :guild_id
	`, setNameQuery)

	_, err := r.db.NamedExecContext(ctx, query, lineWorksBotIV)
	return err
}
