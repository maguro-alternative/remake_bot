# core
## ディレクトリ構成
```
├── core
│   ├── config
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
|   ├── main.go
│   └── schema.sql
```

## config
環境変数を読み込みます。
```core```内ではここから読み込んだ環境変数を参照します。

## main.go
botのメインとなるコードです。
ここを実行することで起動します。

## schema.sql
データベースのスキーマです。
