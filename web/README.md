# web
## ディレクトリ構成
```
├── web
│   ├── components
│   │   ├── channel_select.go
│   │   ├── discord_account_pop.go
│   │   ├── entity.go
│   │   ├── line_account_pop.go
│   │   ├── line_post_discord_chennel.go
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
│   │   │   │   │   ├── hmac.go
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
│   │   └── views
│   │       ├── group
│   │       │   ├── internal
│   │       │   │   ├── entity.go
│   │       │   │   ├── repository_test.go
│   │       │   │   └── repository.go
│   │       │   └── group.go
│   │       ├── guildid
│   │       │   ├── line_post_discord_chennel
│   │       │   │   ├── internal
│   │       │   │   │   ├── entity.go
│   │       │   │   │   ├── repository_test.go
│   │       │   │   │   └── repository.go
│   │       │   │   └── line_post_discord_chennel.go
│   │       │   ├── line_token
│   │       │   │   ├── internal
│   │       │   │   │   ├── entity.go
│   │       │   │   │   ├── repository_test.go
│   │       │   │   │   └── repository.go
│   │       │   │   └── linetoken.go
│   │       │   └── guildid.go
│   │       ├── guilds
│   │       │   └── guilds.go
│   │       └── index.go
│   ├── middleware
│   │   ├── discord_oauth_check.go
│   │   └── log.go
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
│   ├── templates
│   │   ├── static
│   │   │   ├── img
│   │   │   │   └── logo.png
│   │   │   └── js
│   │   │       ├── group.js
│   │   │       ├── line_post_discord_chennel.js
│   │   │       ├── linetoken.js
│   │   │       └── popover.js
│   │   ├── views
│   │   │   ├── group
│   │   │   │   └── group.html
│   │   │   ├── guildid
│   │   │   │   ├── line_post_discord_chennel.html
│   │   │   │   └── linetoken.html
│   │   │   ├── guilds
│   │   │   │   └── guilds.html
|   │   │   ├── login
│   │   │   │   └── line_login.html
│   │   │   └── guildid.html
│   │   ├── index.html
│   │   └── layout.html
│   └── main.go
```

## components
HTMLのコンポーネントを管理するディレクトリです。

## config
環境変数を読み込みます。
```web```内ではここから読み込んだ環境変数を参照します。

## handler
APIのエンドポイントを管理するディレクトリです。
```views```ディレクトリ内ではHTMLファイルを表示するための処理を行います。

## middleware
ミドルウェアを管理するディレクトリです。

## service
プロパティを管理するディレクトリです。

## shared
```web```共通の処理を管理するディレクトリです。

## templates
HTMLテンプレートを管理するディレクトリです。

## main.go
webのメインとなるコードです。
