# tasks
## ディレクトリ構成
```
├── tasks
│   ├── internal
│   │   ├── youtube_test.go
│   │   └── youtube.go
│   └── main.go
```

## internal
Webhookに投稿する内容の取得とメッセージを生成する関数を定義します。  
ここでは主にYouTubeのAPIを操作する関数を定義します。

## main.go
一定単位の時間でRSSを監視し、新しい投稿がされた場合にWebhookに通知します。
