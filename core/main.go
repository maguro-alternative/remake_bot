package core

import (
	"context"
	"os"
	"os/signal"

	"github.com/maguro-alternative/remake_bot/bot"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/web"

	"github.com/gorilla/sessions"
)

func main() {
	ctx := context.Background()
	// データベースの接続を開始
	dbV1, cleanup, err := db.NewDB(ctx, "","")
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// ボットの起動
	discord, cleanupCommandHandlers, err := bot.BotOnReady(dbV1)
	if err != nil {
		panic(err)
	}
	defer discord.Close()
	defer cleanupCommandHandlers()

	// セッションストアを作成します。
	store := sessions.NewCookieStore([]byte("secret-key"))

	// サーバーの待ち受けを開始(ゴルーチンで非同期処理)
	// ここでサーバーを起動すると、Ctrl+Cで終了するまでサーバーが起動し続ける
	go func() {
		web.NewWebRouter(
			dbV1,
			store,
			discord,
		)
	}()
	// Ctrl+Cを受け取るためのチャンネル
	sc := make(chan os.Signal, 1)
	// Ctrl+Cを受け取る
	signal.Notify(sc, os.Interrupt)
	<-sc //プログラムが終了しないようロック

}