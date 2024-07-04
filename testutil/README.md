# testutil
## ディレクトリ構成
```
├── testutil
│   ├── files
│   │  └── video.xml
│   ├── fixtures
│   │   ├── fixtures.go
│   │   ├── line_bot_iv.go
│   │   ├── line_bot.go
│   |   ├── line_ng_discord_message_type.go
│   |   ├── line_ng_discord_role_id.go
│   |   ├── line_ng_discord_user_id.go
│   │   ├── line_post_discord_chennel.go
│   │   ├── permissions_code.go
│   │   ├── permissions_role_id.go
│   │   ├── permissions_user_id.go
│   │   ├── vc_signal_channel.go
│   │   ├── vc_signal_mention_role_id.go
│   │   ├── vc_signal_mention_user_id.go
│   │   ├── vc_signal_ng_role_id.go
│   │   ├── vc_signal_ng_user_id.go
│   │   ├── webhook_role_mention.go
│   │   ├── webhook_user_mention.go
│   │   ├── webhook_word.go
│   │   └── webhook.go
│   └── mock
│       ├── client.go
│       ├── repository.go
│       └── session.go
```

## fixtures
テストで使用するデータベースのデータを定義します。

## mock
テストで使用するモックを定義します。
### client.go
http.Clientのスタブを定義します。
### repository.go
データベース操作のリポジトリのモックを定義します。
### session.go
discordgoを操作するセッションのモックを定義します。
