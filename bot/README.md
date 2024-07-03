# bot
DiscordBotのソースコードです。

# ディレクトリ構成
```
├── bot
│   ├── cogs
│   │   ├── internal
│   │   │   └── entity.go
│   │   ├── cog_handler.go
│   │   ├── on_message_create.go
│   │   ├── on_message_create_test.go
│   │   ├── on_voice_state_update.go
│   │   └── on_voice_state_update_test.go
|   ├── commands
|   |   ├── command_handler.go
|   |   ├── command_handler_test.go
|   |   ├── ping_test.go
|   |   ├── ping.go
|   |   ├── voicevox_test.go
|   |   └── voicevox.go
│   ├── config
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
│   ├── ffmpeg
│   │   ├── ffmpeg_test.go
│   │   └── ffmpeg.go
│   └── main.go
```

## cogs
Discordのイベントを管理するディレクトリです。
ファイル名と同じディレクトリ名のものは、そのイベントで使用するデータベースの処理を行います。
|ファイル名|イベント|説明|
|---|---|---|
|on_message_create|Discordにメッセージが投稿されたとき|メッセージをLINEに送信します。|
|on_voice_state_update|ユーザーがボイスチャンネルに参加したり退出したりしたとき|ボイスチャンネルに参加した場合、指定されたメッセージに通知します。|

## commands
スラッシュコマンドを管理するディレクトリです。
|ファイル名|コマンド名|説明|
|---|---|---|
|ping|ping|Botの応答速度を測定します。|
|voicevox|voicevox|テキストを音声に変換し、ボイスチャンネルで読み上げます。|

## config
環境変数を読み込みます。
```bot```内ではここから読み込んだ環境変数を参照します。

## main.go
discordのセッションを```core```に引き渡します。
引き渡し前に各イベントとスラッシュコマンドを登録させておきます。
```cleanupCommandHandlers```は登録したスラッシュコマンドを削除する関数で、シャットダウン時に実行させます。
