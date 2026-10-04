# core
## ディレクトリ構成
```
├── core
│   ├── config
│   │   ├── internal
│   │   │   └── env.go
│   │   └── config.go
|   └── main.go
```

## config
環境変数を読み込みます。
```core```内ではここから読み込んだ環境変数を参照します。

## main.go
botのメインとなるコードです。
ここを実行することで起動します。

DBスキーマは起動前に `atlas schema apply` で適用します(`docs/SCHEMA.md` 参照)。
