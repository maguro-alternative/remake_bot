# remake_bot 段階的リファクタリング計画

最終更新: 2026-06-13

## 0. 現状分析サマリー

コードベース全体(約33,000行)を調査した結果、以下の問題が確認された。

### 規模の把握

| 領域 | 規模 | 主な問題 |
|------|------|----------|
| `repository/` | 非テスト3,003行 / テスト5,514行 / 47ファイル | 15ペア以上のコピペ重複、81メソッドの単一巨大インターフェース |
| `bot/` | 2,358行 / 13ファイル | `on_message_create.go` が913行、メッセージ送信実装が3系統並存 |
| `web/` | ハンドラ2,388行 / 42ファイル + テンプレート/JS 2,000行超 | ハンドラごとのボイラープレート重複、暗号化フィールドマッピングの重複 |
| `pkg/` | 7サブパッケージ約2,500行 | クライアント設計の不統一(line vs lineworks)、デバッグ文の残骸 |
| `tasks/` | 193行 | YouTube/ニコニコの共通化は部分的、状態がティックごとに揮発 |
| 設定 | `core/config` `bot/config` `web/config` の3重実装(計344行) | ほぼ同一内容のコピー |

### 主要な問題点(優先度順)

1. **CI・lint が存在しない**: `.github/workflows/` なし、Makefile なし、golangci-lint 設定なし、docker-compose なし。テストは大量にある(repository層だけで5,514行)が自動実行されていない。
2. **セキュリティ上の懸念**:
   - `web/handler/api/linebot/internal/hmac.go`: HMAC署名比較が `==` による非定数時間比較で、不一致時に **`nil, nil` を黙って返す**(nilポインタ参照の温床)。
   - LINE系トークンを認証リクエストのたびに全フィールド復号している(キャッシュなし)。
3. **repository層の大量コピペ**: `line_ng_discord_user_id.go` と `line_ng_discord_role_id.go` のような「テーブル名とカラム名だけ違う」ファイルが15ペア以上。`InsertMany`ループ・`DeleteNotIn`パターンが10箇所以上で同型反復。typo `AllCoulmns` が10箇所以上に伝播。
4. **巨大ファイル/巨大関数**: `on_message_create.go`(913行)に message handler が3世代分(`onMessageCreateFunc` はコメントアウト済みの死蔵コード、`Func2` は265行、`Func3` が現行)。`onVoiceStateUpdateFunc` 295行。`lineworks_token.go` 403行。
5. **横断的な不統一**:
   - エラー処理: `cockroachdb/errors` / `fmt.Errorf` / ログだけして握りつぶし、の3様式が混在。
   - ロギング: `slog` と `fmt.Printf`/`fmt.Println` が混在(`pkg/db/db.go:159` にデバッグ用 `fmt.Println` が残存)。
   - context: イベントハンドラが毎回 `context.Background()` を生成し、キャンセル・タイムアウトが伝播しない。
6. **設定の3重管理**: `core/config` / `bot/config` / `web/config` がほぼ同じ内容。環境変数の追加時に3箇所修正が必要。
7. **マイグレーション基盤なし**: `schema.sql`(352行)を起動時に `go:embed` で流し込むだけ。スキーマ変更の履歴管理ができない。

---

## 方針

- **挙動を変えないことを各フェーズの原則とする**。機能追加・仕様変更はこの計画に含めない。
- **フェーズ間に依存関係がある**(後述)。Phase 1 完了前に Phase 2 以降に着手しない。
- **1フェーズ = 複数の小さなPR**。1PRの差分は概ね±500行以内を目安にし、各PRでCIがグリーンであることを必須とする。
- 既存テスト(repository 5,514行、web 多数、JS Jest)が回帰検知の主力。**テストを先に動かせる状態にする(Phase 1)のが全体の前提**。

---

## Phase 1: 安全網の構築(リファクタの前提整備)

> 目的: 「変更して壊れたらすぐ分かる」状態を作る。コード本体には一切手を入れない。

### 作業項目

1. **docker-compose.yml の追加**
   - PostgreSQL(repositoryテストが実DBを要求するため)+ アプリ起動用の構成。
   - テスト用DBの起動手順を README に追記。
2. **GitHub Actions ワークフロー追加**(`.github/workflows/ci.yml`)
   - `go vet` / `go build ./...` / `go test ./...`(servicesでPostgresを起動)
   - JSテスト: `npm test`(Jest)
   - Dockerビルド検証(`GITHUB_TOKEN` が必要な private 依存 `line-works-sdk-go` に注意。secrets 設定が必要)
3. **golangci-lint 導入**(`.golangci.yml`)
   - 初期は `errcheck, govet, staticcheck, unused, misspell` 程度の緩い構成で開始し、既存違反は `--new-from-rev` で増分のみチェック。フェーズが進むごとに厳格化。
4. **Makefile 追加**: `make test` / `make lint` / `make up`(docker-compose)/ `make build` を定義し、CIとローカルのコマンドを一致させる。
5. **カバレッジ計測の開始**: 現状値を記録し、以降のフェーズで低下させないことをルール化。

### 完了条件
- PRを出すと自動でビルド・lint・全テストが走り、mainブランチで全てグリーン。

### リスク
- ほぼなし(コード非変更)。repositoryテストがCI上のPostgresで通らない場合はスキーマ投入手順の整備が必要(`core/schema.sql` を流す)。

---

## Phase 2: 死蔵コード削除と低リスクな清掃

> 目的: 挙動に影響しないゴミを先に消し、以降のフェーズの差分を読みやすくする。

### 作業項目

1. **コメントアウトされたコードの削除**
   - `bot/cogs/on_message_create.go:39-42` の呼び出しと、それが指す `onMessageCreateFunc`(LINE Notify 実装、~113行)。LINE Notify はサービス終了済みのため復活の見込みなし。Git履歴に残るので物理削除でよい。
   - 付随して使われなくなるヘルパー・テストも削除。
2. **デバッグ残骸の除去**
   - `pkg/db/db.go:159` の `fmt.Println(query, err)`。
   - `bot/commands/command_handler.go` / `voicevox.go` 内の `fmt.Printf` 計16箇所前後を `slog` に置換。
3. **typo の一括修正**: `AllCoulmns` → `AllColumns`(10箇所以上、機械的rename。gopls/gorename で実施)。
4. **未使用コードの削除**: `voicevox.go` で渡されるが未使用の `repository.RepositoryFunc` 引数など、staticcheck の `unused` 検出結果に従って削除。
5. **lint設定の厳格化(第1段)**: misspell / unused を全コード対象に昇格。

### 完了条件
- `on_message_create.go` が約800行以下になり、コメントアウトブロックがゼロ。
- `fmt.Printf`/`fmt.Println` によるログ出力がゼロ(slogに統一)。
- 全テストグリーン。

### リスク
- 低。削除対象が本当に未参照であることを `go build ./...` + grep で確認してから消す。

---

## Phase 3: repository 層の重複排除とインターフェース分割

> 目的: 最大のコピペ温床(15ペア以上)を解消し、新テーブル追加コストを下げる。

### 作業項目

1. **共通クエリヘルパーの抽出**(Go 1.23 のジェネリクスを活用)
   - `insertMany[T any](ctx, db, query string, rows []T) error` — 10箇所以上の `NamedExecContext` ループを置換。
   - `deleteNotInProvidedList` — 8ファイル以上に同型で存在する「リスト外を削除、空なら全削除」パターン(`db.In` + `db.Rebind`)を1実装に集約。
   - `selectIDsBy(ctx, db, table, selectCol, whereCol, value)` 相当の薄いヘルパー。
   - **注意**: テーブル名・カラム名は呼び出し側の定数で渡す(動的SQL組み立てによるインジェクションを防ぐため、ユーザー入力は絶対に識別子に使わない)。
2. **対象ファイルの書き換え順序**(1PR=1〜2エンティティ群)
   1. `line_ng_discord_user_id` / `line_ng_discord_role_id`(最も単純なペア)
   2. `vc_signal_ng_*` / `vc_signal_mention_*`(4ファイル)
   3. `permissions_user_id` / `permissions_role_id`
   4. `webhook_word` / `webhook_user_mention` / `webhook_role_mention` / `webhook_thread`
   5. `line_bot` / `line_bot_iv` / `lineworks_bot*`(IV分離構造のため最後。Phase 5 の暗号化共通化と調整)
3. **`RepositoryFunc` インターフェース(81メソッド)の分割**
   - 利用側(cogs / web handlers / tasks)が実際に使うメソッド群を調査し、機能単位の小インターフェースに分割: `LineBotRepository` / `VcSignalRepository` / `WebhookRepository` / `PermissionRepository` / `LineWorksRepository`。
   - 既存の `RepositoryFunc` は分割後インターフェースの合成として残し、利用側の引数型を段階的に小インターフェースへ差し替え(後方互換を維持しながら移行)。
   - `RepositoryFuncMock`(手書き81フィールド)は分割に追従。将来的には moq / gomock などの生成に置き換える検討余地あり(本計画ではスコープ外、手書き維持でも可)。
4. **repositoryテスト(5,514行)はインターフェース経由の振る舞いテストなので原則無変更**。これが本フェーズの回帰検知装置。

### 完了条件
- `InsertMany` / `DeleteNotIn` の同型実装が各1箇所のみ。
- repository非テストコードが3,003行 → 2,000行程度に減少(目安)。
- 利用側が81メソッド全部入りインターフェースに依存している箇所がゼロ。
- 全テストグリーン、SQLの発行内容が変わっていないこと(テストで担保)。

### リスク
- 中。SQL文字列の組み立てを共通化する際の挙動差(空リスト時の分岐など)。→ 既存テストが空リスト系のエッジケースをカバーしているか先に確認し、不足分はテスト追加を先行させる。

---

## Phase 4: セキュリティ修正と暗号化処理の共通化

> 目的: 調査で見つかったセキュリティ上の懸念を解消し、暗号化フィールドの取り回しを一本化する。

### 作業項目

1. **HMAC署名検証の修正**(優先・小PRで即出し可能)
   - `web/handler/api/linebot/internal/hmac.go`: `==` 比較を `hmac.Equal()`(定数時間比較)に変更。署名不一致時の `nil, nil` リターンを明示的なエラー返却に変更し、呼び出し側で401を返す。
   - `lineworks_bot/internal/hmac.go` と実装様式を統一。
2. **暗号化フィールドマッピングの共通化**
   - 現状: `linetoken.go`(6フィールド)と `lineworks_token.go`(8フィールド)で「平文があれば Encrypt して値とIVを別構造体に詰める」を手書き反復。復号側も `line_oauth_check.go` 等で6連続の同型 Decrypt。
   - `pkg/crypto` に `EncryptedField`(値+IVのペア)型を導入し、`EncryptFields(map[string]string) (...)` / 構造体タグベースの一括暗号化・復号ヘルパーを実装。ハンドラ側はフィールド列挙のみにする。
   - DBスキーマ(値とIVが別テーブル)は**変更しない**(マイグレーション基盤がPhase 6まで無いため)。アプリ層のマッピングだけ共通化。
3. **復号結果の取り扱い改善**
   - 認証ミドルウェアが毎リクエスト全トークンを復号している点について、必要フィールドのみ復号するよう絞り込み(キャッシュ導入は挙動変更リスクがあるため本計画ではスコープ外とし、課題として記録)。
4. **セッションCookieの `Secure` フラグ**を本番環境で有効化(環境変数で切替)。

### 完了条件
- 署名検証が定数時間比較になり、不一致時にエラーが明示的に伝播する(テスト追加)。
- Encrypt/Decrypt の同型コードブロックが各ハンドラから消え、`pkg/crypto` のヘルパー1箇所に集約。
- 全テストグリーン。

### リスク
- 中。暗号化まわりはデータ破壊が許されない。**既存DBに保存済みの暗号文との互換性テスト**(既存fixtureの暗号文を新ヘルパーで復号できること)を先に書く。

---

## Phase 5: bot 層の分解 — メッセージ送信パスの一本化

> 目的: 913行の `on_message_create.go` を解体し、3系統あった送信実装を現行系に一本化する。

### 作業項目

1. **送信パスの整理**(要・事前確認)
   - 現状アクティブなのは `onMessageCreateFunc2`(LINE WORKS 直接送信+トークンリフレッシュ265行)と `onMessageCreateFunc3`(`lineworks_service` 経由のキュー送信48行)の両方。
   - **両方が同時に動く必要があるのか、Func3 への移行途中なのかをオーナーに確認**するのが最初のタスク。移行途中なら Func2 を削除して `lineworks_service` に一本化、両方必要なら共通前処理を抽出して並列維持。
   - トークンリフレッシュ処理(Func2 内の~120行)は `pkg/lineworks_service` 側の責務に移す。
2. **ファイル分割**
   - `on_message_create.go` → `message_handler.go`(エントリ+権限チェック)/ `message_builder.go`(`buildMessageText` / `processStickerItems` / `buildFinalMessage`)/ `credentials.go`(復号系ヘルパー。Phase 4 の共通ヘルパーを利用)/ `attachments.go`(ffmpeg変換)。
   - `on_voice_state_update.go` の295行関数を join/leave/camera/stream のイベント種別ごとの関数に分割。
3. **context の是正**
   - `bot/main.go` で生成した親 context(シグナルで cancel される)をハンドラ構造体に保持させ、イベントごとに `context.WithTimeout(parent, ...)` を切る形へ。`context.Background()` の直接生成をハンドラから排除。
4. **エラー処理の統一**
   - bot/ 内を `cockroachdb/errors` の `Wrap`/`WithStack` に統一。「ログだけして継続」する箇所は、継続が意図的であることをコメントで明示するか、エラーを返す形に修正。
5. **tasks/ の整理(小)**
   - YouTube/ニコニコの `run()` 共有は維持しつつ、フィード種別をパラメータ化して `tasks/main.go` の switch を除去。

### 完了条件
- `bot/cogs/` の1ファイル最大行数が300行以下、1関数最大100行以下。
- メッセージ送信の実装系統が(オーナー確認の結果に応じて)1系統、または共通前処理+2出口の構造。
- ハンドラ内の `context.Background()` 直接生成ゼロ。
- 既存の cogs テスト(SessionMock / RepositoryFuncMock 利用)グリーン。

### リスク
- 高(本計画中もっとも挙動変更リスクが高いフェーズ)。Func2/Func3 の並存理由の確認を必ず先行させる。分割→一本化の順で、分割だけのPRと送信パス変更のPRを必ず分ける。

---

## Phase 6: web 層のボイラープレート除去とテンプレート整理

> 目的: 全ハンドラで反復される定型処理をラッパーに集約し、ハンドラを「入力検証+業務処理」だけにする。

### 作業項目

1. **共通ハンドララッパーの導入**
   - 全 `ServeHTTP` で反復されている (a) HTTPメソッド検証 (b) `ctx == nil` チェック (c) `json.NewDecoder(r.Body).Decode` + 400返却 (d) slogへのエラー記録、を `handleJSON[Req any](method string, fn func(ctx, Req) (resp, error))` 形式のジェネリックラッパーに集約。
   - エラー → HTTPステータスのマッピングを1箇所に定義(sentinel error or エラー型)。
   - 1PR=2〜3ハンドラのペースで段階移行(対象は15ハンドラ前後)。
2. **巨大ハンドラの分割**: `lineworks_token.go`(403行)/ `linebot.go`(309行)/ `webhook.go`(278行)を「リクエスト解析」「業務処理(service関数)」「レスポンス構築」に分離。Phase 4 の暗号化ヘルパー適用でさらに縮む見込み。
3. **テンプレートの partial 化**
   - 各テンプレートに重複しているギルドヘッダ/ナビゲーションを `{{template "guild_header" .}}` partial に抽出。
   - チャンネルセレクタの共通 select パターンを partial 化。
   - `testform.html`(125行のテスト用フォーム)が本番ルーティングから参照されていないなら削除。
4. **セッションゲッターの整理**: 約24行×14ファイルの同型セッションゲッターをジェネリクスで1〜2ファイルに集約。
5. **JS の重複確認**: webhook.js(360行)等のフォームシリアライズ処理に共通化余地があれば共通モジュール化(Jestテストが既にあるため安全に実施可能)。

### 完了条件
- ハンドラからメソッド検証・JSONデコード・ctxチェックの直書きが消える。
- web/handler の合計行数が2,388行 → 1,600行程度に減少(目安)。
- Go テスト・Jest テストともグリーン。

### リスク
- 中。エラーレスポンスの文言・ステータスコードが変わるとフロントJSが壊れる可能性があるため、既存レスポンスのスナップショット(テスト)を先に固定する。

---

## Phase 7: 設定統一・マイグレーション導入・起動構成の整理

> 目的: 横断インフラの3重管理を解消し、スキーマ変更を安全にする。

### 作業項目

1. **config パッケージの統一**
   - `core/config`(101行)/ `bot/config`(117行)/ `web/config`(126行)を `internal/config` 1パッケージに統合。`caarlos0/env` + `sync.Once` の現方式は維持。
   - 3パッケージ間の差分(web専用のセッション系、bot専用のトークン系)はフィールドとして統合し、利用側のimportを一括置換。
2. **マイグレーションツールの導入**(golang-migrate を推奨)
   - 現行 `schema.sql`(352行)を `migrations/000001_init.up.sql` としてベースライン化。
   - 起動時 `ExecContext(schema)` を migrate 実行に置換(`IF NOT EXISTS` 前提の現挙動と互換にする)。
   - CI にマイグレーション適用→テストの順序を組み込む。
   - 以降のスキーマ変更(例: Phase 4 で見送ったIVテーブル統合)はマイグレーションとして実施可能になる。
3. **エントリポイントの整理**: `core/main.go`(197行)の起動シーケンス(DB → migrate → bot → web → tasks)を関数分割し、graceful shutdown(Phase 5 で導入した親context)と接続。
4. **Dockerfile の見直し**: ビルドキャッシュ効率(go.mod/go.sum を先にCOPY)、不要レイヤーの削減。CIでのビルド時間短縮。

### 完了条件
- config パッケージが1つ。環境変数の追加が1箇所の修正で済む。
- `migrations/` ディレクトリが存在し、CIで新規DBに対して migrate → 全テストが通る。

### リスク
- 中。マイグレーションのベースライン化は既存本番DBとの整合確認が必要(`schema_migrations` テーブルへの初期レコード投入手順を運用手順書として残す)。

---

## Phase 8: 仕上げ — 観測性・ドキュメント・lint最終化

1. **lint の最終厳格化**: `gocyclo`(複雑度上限)、`funlen` 等を有効化し、Phase 5/6 で達成した「1関数100行以下」を機械的に強制。
2. **ドキュメント整備**: アーキテクチャ図(bot / web / tasks / repository / pkg の依存関係)、開発手順(make コマンド)、運用手順(マイグレーション)を `docs/` に追加。`CLAUDE.md` の作成も検討。
3. **観測性**: slog のフィールド名規約統一(現状 `"Error:"` のようなコロン付きキーが混在)、接続ヘルスモニタ(`on_connection.go`)のメトリクス露出検討。
4. **積み残し課題の棚卸し**: Phase 4 で見送った復号キャッシュ、IVテーブル統合、`RepositoryFuncMock` の自動生成化、discordgo Session のインターフェース化、などを Issue 化。

---

## フェーズ依存関係と進め方

```
Phase 1 (CI/安全網)
  └─ Phase 2 (死蔵コード削除)
       ├─ Phase 3 (repository重複排除) ──┐
       ├─ Phase 4 (セキュリティ/暗号化) ──┼─ Phase 7 (config/migration)
       ├─ Phase 5 (bot層分解)※Phase 4後半に依存 │      └─ Phase 8 (仕上げ)
       └─ Phase 6 (web層)※Phase 4に依存 ──┘
```

- Phase 3 と Phase 4 は並行可能。Phase 5 / 6 は Phase 4 の暗号化ヘルパー完成後に着手。
- **Phase 4-1(HMAC修正)だけは依存なしで即時実施可能**であり、セキュリティ上の理由から最優先で小PRとして出すことを推奨。
- 各フェーズの規模目安: Phase 1-2 は小(数PR)、Phase 3・5・6 は大(各5〜10PR)、Phase 4・7 は中。

## 全体の完了条件

- CI(ビルド・lint・Goテスト・Jestテスト・Dockerビルド)が全PRで実行されグリーン。
- コピペ由来の同型実装(InsertMany / DeleteNotIn / 暗号化マッピング / セッションゲッター / ハンドラ定型処理)が各1箇所。
- 最大ファイル913行 → 300行以下、81メソッドインターフェース → 機能別5分割。
- 設定1箇所・マイグレーション管理あり・死蔵コードゼロ。
- テストカバレッジが開始時点を下回らない。
