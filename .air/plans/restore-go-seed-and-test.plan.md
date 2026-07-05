## 1. Goal

`npm run db:seed` でテスト用データを投入でき、`npm run dev` で `http://localhost:3000/api/v1` のGo APIを起動し、`test` ディレクトリの `npm run test` が実行できる状態にします。

## 2. Approach

既存の黒箱テストは `test/README.md:17-32` と `test/tests/config.ts:1-2` で `npm run db:seed`、`npm run dev`、`http://localhost:3000/api/v1` を前提にしているため、テスト側は変更せずGo側をその契約に合わせます。現在の `main.go:15-79` は `/health` のみなので、まずDB接続・モデル・ルーティングを小さく分割し、seedとHTTPハンドラから同じGORMモデルを使う構成にします。Express時代のnpm操作感を保つため、ルートに薄い `package.json` を追加してGoコマンドを呼び出します。

## 3. File Changes

- **Create** `package.json`
  - ルート用npmスクリプトを追加します。
  - `dev`: `go run .`
  - `db:seed`: `go run ./cmd/seed`
  - `test`: `cd test && npm run test`
  - 必要に応じて `test:install`: `cd test && npm install` を追加します。

- **Modify** `go.mod`
  - 現在は `github.com/gorilla/mux`、GORM、Postgresのみです。
  - UUID生成用に `github.com/google/uuid`、パスワードハッシュ/JWT用に既存間接依存の `golang.org/x/crypto` を直接依存化し、JWTライブラリとして `github.com/golang-jwt/jwt/v5` を追加します。
  - `go 1.26` はローカル環境の実際のGoバージョン次第でビルド不能になる可能性があるため、実装時に `go version` を確認し、必要なら安定版の `go 1.23` などへ調整します。

- **Modify** `go.sum`
  - `go.mod` 変更に伴うチェックサムを更新します。

- **Modify** `main.go:15-79`
  - DB接続とルータ登録を新しい内部パッケージへ委譲します。
  - `:8080` 固定の `main.go:64` を、テスト既定に合わせてデフォルト `3000`、`PORT` 環境変数で上書き可能にします。
  - `/health` は維持しつつ、APIは `/api/v1` 配下に登録します。

- **Create** `internal/config/config.go`
  - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `PORT`, `JWT_SECRET` を読み込みます。
  - `docker-compose.yml:11-15` のDB環境変数と、ローカル実行用のデフォルト値を合わせます。

- **Create** `internal/database/database.go`
  - GORMのPostgreSQL接続、`AutoMigrate`、SQLコネクション取得を提供します。
  - `main.go:16-34` にあるDSN生成と接続処理を移動します。

- **Create** `internal/models/models.go`
  - `User`, `Project`, `Task` モデルを定義します。
  - `test/tests/support/contracts.ts:48-60` が要求する認証レスポンス、`contracts.ts:63-73` のUser、`contracts.ts:76-98` のProject、`contracts.ts:136-153` のTaskに対応できるカラムを持たせます。
  - `User.Email` と `Project.Slug` はユニーク制約にします。

- **Create** `internal/http/response.go`
  - `{"data": ...}` と `{"data": [...], "pageInfo": ...}` のJSONレスポンスを共通化します。
  - `test/tests/support/assertions.ts:36-47` の `pageInfo` 形式に合わせて `totalCount`, `limit`, `page`, `hasNext`, `hasPrevious` を返します。

- **Create** `internal/http/auth.go`
  - JWT発行・検証、Bearerトークンのmiddleware、パスワードハッシュ/検証を実装します。
  - `test/README.md:69-80` と `test/tests/authRequired.test.ts` の全認証必須エンドポイントで、トークンなしを `401` にします。

- **Create** `internal/http/handlers.go`
  - 以下のルートを `/api/v1` 配下に登録します。
  - `POST /auth/login`: `test/tests/auth.test.ts` が要求するseedユーザー認証、失敗時 `401`。
  - `POST /auth/signup`: 一意なメールで作成、重複時 `409`。
  - `GET /users/me`: `test/tests/user.test.ts` のseedユーザー情報を返却。
  - `GET /users/projects`: `test/tests/projects.test.ts` の `limit` / `page` ページネーションとslug順を返却。
  - `GET /users/projects/{slug}`: slug検索、未存在時 `404`。
  - `GET /users/tasks`: `status` 絞り込みとページネーション。
  - `POST /users/tasks`: タスク作成、正常時 `201`、不正payload時 `400`。
  - `GET /users/tasks/{id}`: 未存在時 `404`。
  - `PATCH /users/tasks/{id}`: タスク更新、未存在時 `404`、不正payload時 `400`。
  - `DELETE /users/tasks/{id}`: 削除済みデータを返して `200`、未存在時 `404`。

- **Create** `internal/http/validation.go`
  - `test/tests/tasks.test.ts` の `invalidTaskPayloadCases` が期待する `status`, `kind`, `projectId`, `deadline` の検証を実装します。
  - 許可値は `status`: `scheduled`, `completed`, `archived`、`kind`: `task`, `milestone` とします。
  - `startingAt` / `deadline` は `YYYY-MM-DD` として受け取り、DB保存時は日付として扱います。

- **Create** `internal/seed/seed.go`
  - seed処理を関数化し、テスト用ユーザー・3プロジェクト・scheduledタスクを投入します。
  - `test/tests/testData.ts:1-12` と `test/README.md:35-52` のseed契約に合わせます。
  - 冪等性を持たせ、複数回 `npm run db:seed` しても重複せず、必要データが正しい値に更新されるようにします。

- **Create** `cmd/seed/main.go`
  - `internal/config`, `internal/database`, `internal/seed` を呼び出すCLIエントリポイントを追加します。
  - `package.json` の `db:seed` から実行します。

- **Modify** `Dockerfile:1-17`
  - アプリのデフォルトポートを `3000` に合わせるため `EXPOSE 3000` に変更します。
  - 必要なら `CMD ["./server"]` は維持します。

- **Modify** `docker-compose.yml:7-15`
  - `ports` を `3000:3000` に変更し、`PORT: 3000` と `JWT_SECRET` を追加します。
  - DB設定は既存の `DB_HOST=db`, `DB_PORT=5432`, `DB_USER=app_user`, `DB_PASSWORD=password`, `DB_NAME=express_db` を維持します。

- **Modify** `README.md:17-35` and `README.md:48-63`
  - Go化後の実行手順を `docker compose up -d db`、`npm run db:seed`、`npm run dev`、`cd test && npm install && npm run test` に更新します。
  - APIのデフォルトURLが `http://localhost:3000/api/v1` であることを明記します。

## 4. Implementation Steps

### Task 1: ルートnpmスクリプトを復旧する

1. `package.json` を作成し、`db:seed`, `dev`, `test` をGo/Vitestに委譲するスクリプトとして定義します。
2. `README.md:17-35` のテスト実行手順と `README.md:48-63` のチェック項目を、Go版のコマンドに合わせて更新します。

### Task 2: DB接続とモデルを分離する

1. `internal/config/config.go` を作成し、`main.go:16-23` の環境変数読み取りを移します。
2. `internal/database/database.go` を作成し、`main.go:25-34` のGORM接続と `AutoMigrate` を共通化します。
3. `internal/models/models.go` を作成し、User/Project/TaskのGORMモデルとJSONレスポンスDTOを定義します。
4. `main.go:15-79` を薄い起動処理に変更し、DB接続、migration、ルータ登録、サーバ起動だけにします。

### Task 3: seedコマンドを追加する

1. `internal/seed/seed.go` で `test/tests/testData.ts:1-12` のユーザーとプロジェクト、`test/README.md:49-50` のscheduledタスクを投入します。
2. seedユーザーのパスワードはbcryptで保存し、`POST /auth/login` で `password` が検証できるようにします。
3. プロジェクトは `programming`, `english`, `design` の順序を保持する `sort_order` カラムを持たせ、`test/tests/projects.test.ts` のページ順に合わせます。
4. `cmd/seed/main.go` を追加し、`npm run db:seed` から実行可能にします。

### Task 4: 認証APIとmiddlewareを実装する

1. `internal/http/auth.go` でJWTの発行・検証とBearer middlewareを実装します。
2. `internal/http/handlers.go` に `POST /api/v1/auth/login` と `POST /api/v1/auth/signup` を追加し、`test/tests/auth.test.ts` の `200`, `401`, `409` を返します。
3. 認証成功レスポンスは `test/tests/support/contracts.ts:48-60` に合わせ、`data.uuid`, `data.accessToken`, `data.refreshToken` を返します。
4. `/api/v1/users/*` 配下にmiddlewareを適用し、`test/README.md:69-80` の全エンドポイントでトークンなしを `401` にします。

### Task 5: ユーザー・プロジェクトAPIを実装する

1. `GET /api/v1/users/me` を追加し、`test/tests/support/contracts.ts:63-73` に合う `id`, `username`, `email`, `status` を返します。
2. `GET /api/v1/users/projects` を追加し、`limit`, `page` を解釈して `data` と `pageInfo` を返します。
3. `GET /api/v1/users/projects/{slug}` を追加し、見つからないslugでは `404` を返します。
4. Projectレスポンスは `contracts.ts:76-98` に合わせ、最低限 `id`, `name`, `slug`, `createdAt`, `updatedAt` を返し、必要なら `stats` も含めます。

### Task 6: タスクCRUD APIを実装する

1. `internal/http/validation.go` を作成し、Task payloadの必須項目、`status`, `kind`, `deadline`, `projectId` 存在チェックを実装します。
2. `GET /api/v1/users/tasks` で `status` 絞り込み、`limit` / `page`、`pageInfo` を返します。
3. `POST /api/v1/users/tasks` でタスクを作成し、正常時 `201`、不正payload時 `400` を返します。
4. `GET /api/v1/users/tasks/{id}` で作成済みタスクを返し、未存在UUIDは `404` を返します。
5. `PATCH /api/v1/users/tasks/{id}` でタスクを更新し、未存在は `404`、検証エラーは `400` を返します。
6. `DELETE /api/v1/users/tasks/{id}` で削除対象のタスク内容を返したうえで削除し、未存在は `404` を返します。

### Task 7: ポートとDocker設定をテスト契約に合わせる

1. `main.go:64` のデフォルトポートを `3000` にし、`PORT` で変更可能にします。
2. `Dockerfile:15` を `EXPOSE 3000` に変更します。
3. `docker-compose.yml:7-15` を `3000:3000` と `PORT=3000` に変更し、`JWT_SECRET` を追加します。

## 5. Acceptance Criteria

- ルートで `npm run db:seed` を実行すると、PostgreSQLに `test@example.com` / `password` / `Test User` / `active` のユーザーが作成または更新される。
- `npm run db:seed` を2回連続で実行しても、seedユーザー、`programming` / `english` / `design` プロジェクト、seedタスクが重複しない。
- `npm run dev` で `http://localhost:3000/api/v1` にGo APIが起動する。
- `POST /api/v1/auth/login` に `test@example.com` と `password` を送ると `200`、`data.uuid`, `data.accessToken`, `data.refreshToken` を返す。
- `POST /api/v1/auth/login` に誤ったpasswordを送ると `401` を返す。
- `POST /api/v1/auth/signup` は一意なemailで `200`、seedユーザーと同じemailで `409` を返す。
- `GET /api/v1/users/me` はBearer tokenなしで `401`、seedユーザーtokenありで `200` とseedユーザー情報を返す。
- `GET /api/v1/users/projects?limit=1&page=1..3` は順に `programming`, `english`, `design` を1件ずつ返す。
- `GET /api/v1/users/projects/missing-project-for-api-test` は `404` を返す。
- `GET /api/v1/users/tasks?limit=20&page=1&status=scheduled` は `scheduled` のタスクを1件以上返す。
- `POST /api/v1/users/tasks` は正しいpayloadで `201`、空title・不正status・存在しないprojectId・不正kind・不正deadlineで `400` を返す。
- `GET`, `PATCH`, `DELETE /api/v1/users/tasks/{id}` は作成済みタスクに対してそれぞれ `200` を返し、未存在UUIDに対して `404` を返す。
- `cd test && npm run test` が全件成功する。

## 6. Verification Steps

1. DBを起動します: `docker compose up -d db`
2. Go依存を整理します: `go mod tidy`
3. seedを実行します: `npm run db:seed`
4. 冪等性を確認します: `npm run db:seed` をもう一度実行し、エラーや重複がないことを確認します。
5. APIを起動します: `npm run dev`
6. 別ターミナルでテスト依存を入れます: `cd test && npm install`
7. 黒箱テストを実行します: `cd test && npm run test`
8. 必要に応じて直接確認します:
   - `curl -i http://localhost:3000/api/v1/users/me` が `401`。
   - `curl -i -X POST http://localhost:3000/api/v1/auth/login -H 'Content-Type: application/json' -d '{"email":"test@example.com","password":"password"}'` が `200`。

## 7. Risks & Mitigations

- `go.mod` の `go 1.26` がローカルGoで未対応の場合、`go run` 自体が失敗します。実装時に `go version` を確認し、必要なら `go.mod` のGoバージョンを実環境に合わせます。
- テストは `http://localhost:3000/api/v1` 固定が既定なので、現在の `main.go:64` の `:8080` と `docker-compose.yml:7-8` の `8080:8080` のままだと接続できません。デフォルトを `3000` に統一し、`PORT` で上書きできるようにします。
- seedとAPIが別々のモデル定義を持つとデータ不整合が起きやすいため、`internal/models` と `internal/database` をseed/APIで共有します。
- `DELETE /users/tasks/:id` 後に返すレスポンスは、削除後に再取得できないため、削除前のレコードを保持してから削除し、その内容を `data` として返します。
- テストは並列にタスクを作成する箇所があるため、IDはDB側のUUIDまたはGoのUUIDで必ず一意に生成し、一覧は `created_at` と `id` の安定順で返します。