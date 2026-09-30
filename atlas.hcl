# Atlas 設定(宣言的スキーマ管理 / ridgepole 相当)
# schema/schema.sql をあるべき姿とし、`atlas schema apply` で DB との差分をその場で適用する。
# 手順は docs/SCHEMA.md を参照。

diff {
  # schema.sql から消えたテーブル・スキーマを自動で DROP しない(ridgepole の --drop-table 無し相当)
  skip {
    drop_schema = true
    drop_table  = true
  }
}

# ローカル開発用: 差分計算に Docker の一時 Postgres を使う
env "local" {
  src = "file://schema/schema.sql"
  url = urlqueryset(getenv("DATABASE_URL"), "search_path", "public")
  dev = "docker://postgres/16/dev?search_path=public"
}

# 本番(Railway pre-deploy)用: Docker が使えないため、同じサーバー上の空DBを差分計算に使う
env "prod" {
  src = "file://schema/schema.sql"
  url = urlqueryset(getenv("DATABASE_URL"), "search_path", "public")
  dev = getenv("ATLAS_DEV_URL")
}
