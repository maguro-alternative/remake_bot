package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB // DBは*sql.DB型の変数、グローバル変数

// RetryLimit はリトライする最大回数を指定します。
//const retryLimit = 3

// RetryInterval はリトライの間隔を指定します（ミリ秒単位）。
const retryInterval = 1 // 1秒


type DBHandler struct {
	Driver           *sqlx.DB
	DBPing           func(context.Context) error
	QueryxContext    func(context.Context, string, ...interface{}) (*sqlx.Rows, error)
	QueryRowxContent func(context.Context, string, ...interface{}) (*sqlx.Row, error)
	GetContent       func(context.Context, interface{}, string, ...interface{}) error
	SelectContent    func(context.Context, interface{}, string, ...interface{}) error
	ExecContext      func(ctx context.Context, query string, args ...any) (*sql.Result, error)
	NamedExecContext func(ctx context.Context, query string, arg interface{}) (*sql.Result, error)
}

func retryOperation(ctx context.Context, operation func() error) error {
	/*
		リトライする関数

		引数
			ctx: context.Context型の変数
			operation: func() error型の変数

		戻り値
			error型の変数
	*/
	retryBackoff := backoff.NewExponentialBackOff()
	retryBackoff.MaxElapsedTime = time.Second * retryInterval

	err := backoff.RetryNotify(func() error {
		err := operation()
		if err != nil {
			return err
		}
		err = backoff.Permanent(err)
		return errors.WithStack(err)
	}, retryBackoff, func(err error, duration time.Duration) {
		//slog.WarnContext(ctx, fmt.Sprintf("%v retrying in %v...", err, duration))
	})
	return errors.WithStack(err)
}

func In(query string, args ...any) (string, []any, error) {
	/*
		クエリを生成する関数

		引数
			query: sql文
			args: sql文に入れる値

		戻り値
			sql文
			sql文に入れる値
			error型の変数
	*/
	for _, arg := range args {
		if arg == nil {
			return "", nil, errors.New("nil arguments are not allowed")
		}
	}

	query, param, err := sqlx.In(query, args...)
	return query, param, errors.WithStack(err)
}

func Rebind(bindType int, query string) string {
	/*
		クエリを生成する関数
		postgresの場合はbindTypeにはlen(args)を入れる

		引数
			bindType: len(args)
			query: sql文

		戻り値
			sql文
	*/
	return sqlx.Rebind(bindType, query)
}

func NewDB(ctx context.Context, driverName string, path string) (*DB, func(), error) {
	/*
		データベースに接続する関数

		引数
			ctx: context.Context型の変数
			driverName: データベースの種類
			path: データベースのパス

		戻り値
			*DB型の変数
			データベースの接続を閉じる関数
			error型の変数
	*/
	db, err := sql.Open(driverName, path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to sql.Open(): ")
	}

	// pingが通らない場合はエラーを返して終了する
	ctx, cancel := context.WithTimeout(ctx, time.Second * retryInterval)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, errors.Wrap(err, "failed to (*sql.DB).PingContext(): ")
	}
	xDriver := sqlx.NewDb(db, driverName)

	return &DB{driver: xDriver}, func() { _ = db.Close() }, nil
}

type DB struct {
	driver *sqlx.DB
}

func (db *DB) PingDB(ctx context.Context) error {
	if err := db.driver.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (db *DB) PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	/*
		*sqlx.Stmtオブジェクトを生成する関数
		SQLステートメントを事前に解析（"prepare"）するために使用

		引数
			ctx: context.Context型の変数
			query: sql文

		戻り値
			*sqlx.Stmt型の変数
			error型の変数
	*/
	var err error
	var stmt *sqlx.Stmt
	operation := func() error {
		stmt, err = db.driver.PreparexContext(ctx, query)
		fmt.Println(query,err)
		return errors.WithStack(err)
	}
	if err := retryOperation(ctx, func() error { return operation() }); err != nil {
		return nil, err
	}

	return stmt, nil
}

func (db *DB) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	/*
		複数問い合わせの結果を*sqlx.Rowsで返す関数

		引数
			ctx: context.Context型の変数
			query: sql文
			args: sql文に入れる値

		戻り値
			*sqlx.Rows型の変数
			error型の変数
	*/
	var err error
	var rows *sqlx.Rows
	operation := func() error {
		rows, err = db.driver.QueryxContext(ctx, query, args...)
		return errors.WithStack(err)
	}
	if err := retryOperation(ctx, func() error { return operation() }); err != nil {
		return nil, err
	}

	return rows, nil
}

func (db *DB) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	/*
		1行問い合わせの結果を*sqlx.Rowで返す関数

		引数
			ctx: context.Context型の変数
			query: sql文
			args: sql文に入れる値

		戻り値
			*sqlx.Row型の変数
	*/
	return db.driver.QueryRowxContext(ctx, query, args...)
}

func (db *DB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	/*
		1行問い合わせの結果をdestで返す関数

		引数
			ctx: context.Context型の変数
			dest: sql文の結果を格納する変数
			query: sql文
			args: sql文に入れる値

		戻り値
			error型の変数
	*/
	operation := func() error {
		err := db.driver.GetContext(ctx, dest, query, args...)
		return errors.WithStack(err)
	}
	return retryOperation(ctx, func() error { return operation() })
}

func (db *DB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	/*
		複数問い合わせの結果をdestで返す関数

		引数
			ctx: context.Context型の変数
			dest: sql文の結果を格納する変数
			query: sql文
			args: sql文に入れる値

		戻り値
			error型の変数
	*/
	operation := func() error {
		err := db.driver.SelectContext(ctx, dest, query, args...)
		return errors.WithStack(err)
	}
	return retryOperation(ctx, func() error { return operation() })
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	/*
		クエリ実行を行う関数

		引数
			ctx: context.Context型の変数
			query: sql文
			args: sql文に入れる値

		戻り値
			sql.Result型の変数
			error型の変数
	*/
	var err error
	var result sql.Result
	operation := func() error {
		result, err = db.driver.ExecContext(ctx, query, args...)
		return errors.WithStack(err)
	}
	if err := retryOperation(ctx, func() error { return operation() }); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *DB) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	/*
		クエリ実行を行う関数
		argには構造体を入れ、sql文には構造体のフィールド名を入れる

		引数
			ctx: context.Context型の変数
			query: sql文
			arg: 構造体

		戻り値
			sql.Result型の変数
			error型の変数
	*/
	var err error
	var result sql.Result
	operation := func() error {
		result, err = db.driver.NamedExecContext(ctx, query, arg)
		return errors.WithStack(err)
	}
	if err := retryOperation(ctx, func() error { return operation() }); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *DB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	/*
		トランザクションを開始する関数

		引数
			ctx: context.Context型の変数
			opts: sql.TxOptions型の変数

		戻り値
			*Tx型の変数
			error型の変数
	*/
	var err error
	var tx *sqlx.Tx
	operation := func() error {
		tx, err = db.driver.BeginTxx(ctx, opts)
		return errors.WithStack(err)
	}
	if err := retryOperation(ctx, func() error { return operation() }); err != nil {
		return nil, err
	}

	return &Tx{driver: tx}, nil
}

type Tx struct {
	driver *sqlx.Tx
}

func (tx *Tx) PingDB(ctx context.Context) error {
	/*
		データベースの接続を確認する関数

		引数
			ctx: context.Context型の変数

		戻り値
			error型の変数
	*/
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}


func (tx *Tx) PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	/*
		*sqlx.Stmtオブジェクトを生成する関数
		SQLステートメントを事前に解析（"prepare"）するために使用

		引数
			ctx: context.Context型の変数
			query: sql文

		戻り値
			*sqlx.Stmt型の変数
			error型の変数
	*/
	var err error
	var stmt *sqlx.Stmt
	operation := func() error {
		stmt, err = tx.driver.PreparexContext(ctx, query)
		return errors.WithStack(err)
	}
	if err := retryOperation(ctx, func() error { return operation() }); err != nil {
		return nil, err
	}

	return stmt, nil
}

func (tx *Tx) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	/*
		複数問い合わせの結果を*sqlx.Rowsで返す関数

		引数
			ctx: context.Context型の変数
			query: sql文
			args: sql文に入れる値

		戻り値
			*sqlx.Rows型の変数
			error型の変数
	*/
	var err error
	var rows *sqlx.Rows
	operation := func() error {
		rows, err = tx.driver.QueryxContext(ctx, query, args...)
		return errors.WithStack(err)
	}
	if err := retryOperation(ctx, func() error { return operation() }); err != nil {
		return nil, err
	}

	return rows, nil
}

func (tx *Tx) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	/*
		1行問い合わせの結果を*sqlx.Rowで返す関数

		引数
			ctx: context.Context型の変数
			query: sql文
			args: sql文に入れる値

		戻り値
			*sqlx.Row型の変数
	*/
	return tx.driver.QueryRowxContext(ctx, query, args...)
}

func (tx *Tx) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	/*
		1行問い合わせの結果をdestで返す関数

		引数
			ctx: context.Context型の変数
			dest: sql文の結果を格納する変数
			query: sql文
			args: sql文に入れる値

		戻り値
			error型の変数
	*/
	operation := func() error {
		err := tx.driver.GetContext(ctx, dest, query, args...)
		return errors.WithStack(err)
	}
	return retryOperation(ctx, func() error { return operation() })
}

func (tx *Tx) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	/*
		複数問い合わせの結果をdestで返す関数

		引数
			ctx: context.Context型の変数
			dest: sql文の結果を格納する変数
			query: sql文
			args: sql文に入れる値

		戻り値
			error型の変数
	*/
	operation := func() error {
		err := tx.driver.SelectContext(ctx, dest, query, args...)
		return errors.WithStack(err)
	}
	return retryOperation(ctx, func() error { return operation() })
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	/*
		クエリ実行を行う関数

		引数
			ctx: context.Context型の変数
			query: sql文
			args: sql文に入れる値

		戻り値
			sql.Result型の変数
			error型の変数
	*/
	var err error
	var result sql.Result
	operation := func() error {
		result, err = tx.driver.ExecContext(ctx, query, args...)
		return errors.WithStack(err)
	}
	if err := retryOperation(ctx, func() error { return operation() }); err != nil {
		return nil, err
	}

	return result, nil
}

func (tx *Tx) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	/*
		クエリ実行を行う関数
		argには構造体を入れ、sql文には構造体のフィールド名を入れる

		引数
			ctx: context.Context型の変数
			query: sql文
			arg: 構造体

		戻り値
			sql.Result型の変数
			error型の変数
	*/
	var err error
	var result sql.Result
	operation := func() error {
		result, err = tx.driver.NamedExecContext(ctx, query, arg)
		return errors.WithStack(err)
	}
	if err := retryOperation(ctx, func() error { return operation() }); err != nil {
		return nil, err
	}

	return result, nil
}

func (tx *Tx) CommitCtx(ctx context.Context) error {
	/*
		トランザクションをコミットする関数

		引数
			ctx: context.Context型の変数

		戻り値
			error型の変数
	*/
	err := retryOperation(ctx, func() error {
		err := tx.driver.Commit()
		return errors.WithStack(err)
	})
	if err != nil {
		return err
	}

	return nil
}

func (tx *Tx) RollbackCtx(ctx context.Context) error {
	/*
		トランザクションをロールバックする関数

		引数
			ctx: context.Context型の変数

		戻り値
			error型の変数
	*/
	err := retryOperation(ctx, func() error {
		err := tx.driver.Rollback()
		return errors.WithStack(err)
	})
	if err != nil {
		return err
	}

	return nil
}


type Driver interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
}

// 上記のインターフェースを満たす構造体
// この構造体はDBとTxの両方で使用できる
// 関数に一つでも欠けがあるとコンパイルエラーになる
var (
	_ Driver = (*Tx)(nil)
	_ Driver = (*DB)(nil)
	_ Driver = (*sqlx.Tx)(nil)
	_ Driver = (*sqlx.DB)(nil)
)