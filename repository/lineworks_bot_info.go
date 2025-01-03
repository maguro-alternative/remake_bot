package repository

import (
	"context"
	"fmt"
	"strings"

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

func NewLineWorksBotInfo(
	guildID string,
	lineWorksClientID pq.ByteaArray,
	lineWorksClientSecret pq.ByteaArray,
	lineWorksServiceAccount pq.ByteaArray,
	lineWorksPrivateKey pq.ByteaArray,
	lineWorksDomainID pq.ByteaArray,
	lineWorksAdminID pq.ByteaArray,
) *LineWorksBotInfo {
	return &LineWorksBotInfo{
		GuildID: guildID,
		LineWorksClientID: lineWorksClientID,
		LineWorksClientSecret: lineWorksClientSecret,
		LineWorksServiceAccount: lineWorksServiceAccount,
		LineWorksPrivateKey: lineWorksPrivateKey,
		LineWorksDomainID: lineWorksDomainID,
		LineWorksAdminID: lineWorksAdminID,
	}
}

func (r *Repository) InsertLineWorksBotInfo(ctx context.Context, lineWorksBotInfo *LineWorksBotInfo) error {
	query := `
		INSERT INTO
			line_works_bot_info (
				guild_id,
				line_works_client_id,
				line_works_client_secret,
				line_works_service_account,
				line_works_private_key,
				line_works_domain_id,
				line_works_admin_id
			)
		VALUES (
			:guild_id,
			:line_works_client_id,
			:line_works_client_secret,
			:line_works_service_account,
			:line_works_private_key,
			:line_works_domain_id,
			:line_works_admin_id
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, lineWorksBotInfo)
	return err
}

func (r *Repository) GetLineWorksBotInfoByGuildID(ctx context.Context, guildID string) (*LineWorksBotInfo, error) {
	var lineWorksBotInfo LineWorksBotInfo
	query := `
		SELECT
			*
		FROM
			line_works_bot_info
		WHERE
			guild_id = $1
	`
	err := r.db.GetContext(ctx, &lineWorksBotInfo, query, guildID)
	return &lineWorksBotInfo, err
}

func (r *Repository) UpdateLineWorksBotInfo(ctx context.Context, lineWorksBotInfo *LineWorksBotInfo) error {
	var setNameQuery string
	var setQueryArray []string

	if len(lineWorksBotInfo.LineWorksClientID) > 0 {
		setQueryArray = append(setQueryArray, "line_works_client_id = :line_works_client_id")
	}
	if len(lineWorksBotInfo.LineWorksClientSecret) > 0 {
		setQueryArray = append(setQueryArray, "line_works_client_secret = :line_works_client_secret")
	}
	if len(lineWorksBotInfo.LineWorksServiceAccount) > 0 {
		setQueryArray = append(setQueryArray, "line_works_service_account = :line_works_service_account")
	}
	if len(lineWorksBotInfo.LineWorksPrivateKey) > 0 {
		setQueryArray = append(setQueryArray, "line_works_private_key = :line_works_private_key")
	}
	if len(lineWorksBotInfo.LineWorksDomainID) > 0 {
		setQueryArray = append(setQueryArray, "line_works_domain_id = :line_works_domain_id")
	}
	if len(lineWorksBotInfo.LineWorksAdminID) > 0 {
		setQueryArray = append(setQueryArray, "line_works_admin_id = :line_works_admin_id")
	}
	setNameQuery = strings.Join(setQueryArray, ",")
	if setNameQuery == "" {
		return nil
	}
	query := fmt.Sprintf(`
		UPDATE
			line_works_bot_info
		SET
			%s
		WHERE
			guild_id = :guild_id
	`, setNameQuery)
	_, err := r.db.NamedExecContext(ctx, query, lineWorksBotInfo)
	return err
}
