# bot
DiscordBotのソースコードです。

# ディレクトリ構成
```
├── bot
│   ├── cogs
│   │   ├── on_message_create
│   │   │   ├── entity.go
│   │   │   ├── repository_test.go
│   │   │   └── repository.go
│   │   ├── cog_handler.go
│   │   ├── on_message_create.go
│   │   └── vc_signal.go
|   ├── commands
|   |   ├── command_handler.go
|   |   └── ping.go
│   ├── config
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
│   └── main.go
```

## cogs
Discordのイベントを管理するディレクトリです。
ファイル名と同じディレクトリ名のものは、そのイベントで使用するデータベースの処理を行います。
|ファイル名|イベント|説明|
|---|---|---|
|on_message_create|Discordにメッセージが投稿されたとき|メッセージをLINEに送信します。|

## commands
スラッシュコマンドを管理するディレクトリです。
|ファイル名|コマンド名|説明|
|---|---|---|
|ping|ping|Botの応答速度を測定します。|

## config
環境変数を読み込みます。
```bot```内ではここから読み込んだ環境変数を参照します。

## main.go
discordのセッションを```core```に引き渡します。
引き渡し前に各イベントとスラッシュコマンドを登録させておきます。
```cleanupCommandHandlers```は登録したスラッシュコマンドを削除する関数で、シャットダウン時に実行させます。
