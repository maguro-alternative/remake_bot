package internal

import (
	"context"
	"fmt"
	"reflect"

	"github.com/maguro-alternative/remake_bot/pkg/db"
)

type Repository struct {
	db db.Driver
}

func NewRepository(db db.Driver) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) UpdateLineBot(ctx context.Context, lineBotJson *LineBotJson) error {
	var setNameQuery string
	structTypeOf := reflect.TypeOf(lineBotJson)

	// 受け取った構造体のフィールドのみを更新する
	for i := 0; i < structTypeOf.NumField(); i++ {
		field := structTypeOf.Field(i).Tag.Get("db")
		if field == "" || field == "guild_id"{
			continue
		}
		if i == structTypeOf.NumField()-1 {
			setNameQuery += field + " = :" + field
			continue
		}
		setNameQuery += field + " = :" + field + ","
	}
	query := fmt.Sprintf(`
		UPDATE
			line_bot
		SET
			%s
		WHERE
			guild_id = :guild_id
	`, setNameQuery)
	_, err := r.db.NamedExecContext(ctx, query, lineBotJson)
	return err
}
