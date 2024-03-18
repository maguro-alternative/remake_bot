# discordfastリプレイス

https://github.com/maguro-alternative/discordfast

上記のリポジトリのリプレイスです。

# 進捗

- [x] LINE→Discordのメッセージの送信
- [x] Discord→LINEのメッセージの送信
- [x] LINE→Discordの画像の送信
- [x] Discord→LINEの画像の送信
- [x] LINE→Discordの動画の送信
- [x] Discord→LINEの動画の送信
- [x] LINE→Discordの音声の送信
- [x] Discord→LINEの音声の送信
- [x] LINE→Discordのスタンプの送信
- [x] Discord→LINEのスタンプの送信
- [] LINE→Discordのファイルの送信
- [] Discord→LINEのファイルの送信

# ディレクトリ構造
<datails><summary>長すぎるので折り畳み</summary>
```plaintext
.
├── bot
│   ├── cogs
│   │   ├── on_message_create       // メッセージが送信されたときのデータベースの操作
│   │   │   ├── entity.go
│   │   │   ├── repository_test.go
│   │   │   └── repository.go
│   │   ├── cog_handler.go          // ここでcogを登録
│   │   ├── on_message_create.go    // discord内でメッセージが送信されたときの処理
│   │   └── vc_signal.go            // ボイスチャンネルのステータス変化時の処理
|   ├── commands
|   |   ├── command_handler.go      // ここでコマンドを登録
|   |   └── ping.go                 // pingコマンド
│   ├── config                      // 環境変数設定ファイル
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
│   └── main.go                     // Botメイン関数
├── core
│   ├── config                      // 環境変数設定ファイル
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
|   ├── main.go
│   └── schema.sql                  // データベースのスキーマ
├── fixtures
│   ├── fixtures.go
│   ├── line_bot_iv.go
│   ├── line_bot.go
|   ├── line_ng_discord_id.go
|   ├── line_ng_discord_message_type.go
│   ├── line_post_discord_chennel.go
│   ├── permissions_code.go
│   ├── permissions_id.go
│   ├── vc_signal_channel.go
│   ├── vc_signal_mention_id.go
│   ├── vc_signal_ng_id.go
│   ├── webhook_mention.go
│   ├── webhook_word.go
│   └── webhook.go
├── pkg
│   ├── crypto
│   │   ├── aes.go
│   │   └── aes_test.go
│   ├── db
│   │   ├── db.go
│   │   └── db_test.go
│   ├── line
│   │   ├── get_bot_info.go
│   │   ├── get_friend_count.go
│   │   ├── get_group_count.go
│   │   ├── get_message_content.go
│   │   ├── get_profile.go
│   │   ├── get_pushlimit.go
│   │   ├── get_totalpush_count.go
│   │   ├── line_message.go
│   │   ├── line_notify.go
│   │   └── line.go
│   └── youtube
│       ├── create_client_secret.go
│       ├── create_oauth2.go
│       └── youtube.go
├── web
│   ├── components
│   │   ├── discord_account_pop.go
│   │   ├── entity.go
│   │   ├── line_account_pop.go
│   │   ├── line_post_discord_chennel.go
│   │   ├── linetoken.go
│   │   └── submittag.go
│   ├── config
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
│   ├── handler
│   │   ├── api
│   │   │   ├── group
│   │   │   │   ├── internal
│   │   │   │   │   ├── entity.go
│   │   │   │   │   ├── repository_test.go
│   │   │   │   │   └── repository.go
│   │   │   │   └── group.go
│   │   |   ├── line_post_discord_chennel
│   │   │   │   ├── internal
│   │   │   │   │   ├── entity.go
│   │   │   │   │   ├── repository_test.go
│   │   │   │   │   └── repository.go
│   │   │   │   └── line_post_discord_chennel.go
│   │   │   ├── line_bot
│   │   │   │   ├── internal
│   │   │   │   │   ├── entity.go
│   │   │   │   │   ├── repository_test.go
│   │   │   │   │   └── repository.go
│   │   │   │   ├── entity.go
│   │   │   │   └── line_bot.go
│   │   │   └── linetoken
│   │   │       ├── internal
│   │   │       │   ├── entity.go
│   │   │       │   ├── repository_test.go
│   │   │       │   └── repository.go
│   │   │       └── linetoken.go
│   │   ├── callback
│   │   │   ├── discord_callback
│   │   │   │   └── callback.go
│   │   │   └── line_callback
│   │   │       └── callback.go
│   │   ├── login
│   │   │   ├── discord_login
│   │   │   │   └── discord_login.go
│   │   │   └── line_login
│   │   │       ├── internal
│   │   │       │   ├── entity.go
│   │   │       │   ├── repository_test.go
│   │   │       │   └── repository.go
│   │   │       └── line_login.go
│   │   ├── logout
│   │   │   ├── discord_logout
│   │   │   │   └── discord_logout.go
│   │   │   └── line_logout
│   │   │       └── line_logout.go
│   │   ├── views
│   │   │   ├── group
│   │   │   │   ├── internal
│   │   │   │   │   ├── entity.go
│   │   │   │   │   ├── repository_test.go
│   │   │   │   │   └── repository.go
│   │   │   │   └── group.go
│   │   │   ├── guildid
│   │   │   │   ├── line_post_discord_chennel
│   │   │   │   │   ├── internal
│   │   │   │   │   │   ├── entity.go
│   │   │   │   │   │   ├── repository_test.go
│   │   │   │   │   │   └── repository.go
│   │   │   │   │   └── line_post_discord_chennel.go
│   │   │   │   ├── line_token
│   │   │   │   │   ├── internal
│   │   │   │   │   │   ├── entity.go
│   │   │   │   │   │   ├── repository_test.go
│   │   │   │   │   │   └── repository.go
│   │   │   │   │   └── linetoken.go
│   │   │   │   └── guildid.go
│   │   │   ├── guilds
│   │   │   │   └── guilds.go
│   │   │   └── index.go
│   ├── middleware
│   │   └── middleware.go
│   ├── service
│   │   ├── discord_oauth2.go
│   │   └── index.go
│   ├── shared
│   │   ├── permission
│   │   │   ├── internal
│   │   │   │   ├── entity.go
│   │   │   │   ├── repository_test.go
│   │   │   │   └── repository.go
│   │   │   ├── check_discord_permission.go
│   │   │   └── check_line_permission.go
│   │   └── session
│   │       ├── getoauth
│   │       │   ├── get_discord_oauth.go
│   │       │   └── get_line_oauth.go
│   │       └── model
│   │           └── entity.go
│   └── templates
│       ├── static
│       │   ├── img
│       │   │   └── logo.png
│       │   └── js
│       │       ├── group.js
│       │       ├── line_post_discord_chennel.js
│       │       ├── linetoken.js
│       │       └── popover.js
│       ├── views
│       │   ├── group
│       │   │   └── group.html
│       │   ├── guildid
│       │   │   ├── line_post_discord_chennel.html
│       │   │   └── linetoken.html
│       │   ├── guilds
│       │   │   └── guilds.html
|       │   ├── login
│       │   │   └── line_login.html
│       │   └── guildid.html
│       ├── index.html
│       └── layout.html
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```
</details>

