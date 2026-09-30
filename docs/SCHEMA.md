# スキーマ運用手順

スキーマは Atlas の宣言的ワークフロー(`atlas schema apply`)で管理する。ridgepole と同様に、`schema/schema.sql` に「あるべき姿」だけを書き、DB との差分はその場で計算・適用する。マイグレーションファイルの履歴は持たない。

| ファイル | 役割 |
| --- | --- |
| `schema/schema.sql` | スキーマのあるべき姿。変更はここにだけ書く |
| `atlas.hcl` | Atlas の設定(`env "local"` / `env "prod"`、DROP の抑止) |
| `scripts/schema-apply.sh` | 本番適用スクリプト(Railway の pre-deploy から実行) |
| `railway.json` | Railway の pre-deploy コマンド設定 |

## 必要なもの(ローカル)

- Atlas CLI(`curl -sSf https://atlasgo.sh | sh` または `brew install ariga/tap/atlas`)
- Docker(`env "local"` は差分計算用の一時 Postgres を `docker://postgres/16/dev` で起動する)

## スキーマを変更する

1. `schema/schema.sql` を編集する(`IF NOT EXISTS` は付けない)。
2. ローカルDBで差分を確認して適用する。

   ```sh
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
   atlas schema apply --env local --dry-run      # 実行されるSQLを確認
   atlas schema apply --env local --auto-approve # 適用
   ```

3. 本番に対しても事前に dry-run で実行予定のSQLを確認しておく。

   ```sh
   DATABASE_URL="<本番>" ATLAS_DEV_URL="<本番サーバーの atlas_dev>" atlas schema apply --env prod --dry-run
   ```

4. `schema/schema.sql` をコミットしてデプロイすると、pre-deploy で本番に適用される。

## 注意点

- **テーブル・スキーマの DROP は自動実行されない**(`atlas.hcl` の `diff.skip`)。`schema.sql` からテーブル定義を消しても本番には残る。不要になったテーブルは手動で DROP する。
- **カラムの DROP は実行される**。`schema.sql` からカラムを消すと、デプロイ時にデータごと削除される。
- **リネームは DROP + ADD になる**(ridgepole の `renamed_from` 相当は無い)。カラム名の変更は、先に本番で手動 `ALTER TABLE ... RENAME COLUMN` してから `schema.sql` を合わせる。
- NOT NULL カラムの追加は、既存行があると失敗する。`DEFAULT` を付けるか、NULL 許容で追加→データ投入→NOT NULL 化の順に分けてデプロイする。

## 本番(Railway)の設定

デプロイ時、`railway.json` の `preDeployCommand` により `/app/scripts/schema-apply.sh` がアプリ起動前に実行される。失敗した場合はデプロイが中断され、旧バージョンが動き続ける。

初回のみ以下を設定する。

1. 本番の Postgres に差分計算用の空DBを作成する(中身は Atlas が毎回クリーンに使う。データを入れないこと)。

   ```sql
   CREATE DATABASE atlas_dev;
   ```

2. サービスの環境変数に `ATLAS_DEV_URL` を追加する。`DATABASE_URL` のDB名を `atlas_dev` に変え、`search_path=public` を付けたもの。

   ```
   ATLAS_DEV_URL=postgresql://<user>:<password>@<host>:<port>/atlas_dev?search_path=public
   ```

3. 初回デプロイ前に `atlas schema apply --env prod --dry-run` を本番に対して実行し、`Schema is synced` になることを確認する(旧方式の `schema.sql` で作成済みの本番と一致している前提)。差分が出た場合は内容を確認し、`schema/schema.sql` を本番に合わせるか、意図した変更であることを確認してからデプロイする。

## CI

`.github/workflows/schema.yml` で、空の Postgres に `schema.sql` を適用できること、再適用で差分が出ないことを検査する。
