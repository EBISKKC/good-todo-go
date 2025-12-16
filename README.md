# Good Todo Go

PostgreSQL の **Row Level Security (RLS)** を用いたマルチテナント分離を学習・実践するために作成した Todo アプリケーションです。

Go + Echo によるバックエンド API と Next.js + shadcn/ui によるフロントエンドで構成されています。

## プロジェクトの目的

このプロジェクトは、以下の技術を学習・実践することを目的としています：

- **RLS (Row Level Security)** によるテナント間のデータ分離
- **マルチテナントアーキテクチャ** の実装パターン
- **クリーンアーキテクチャ** によるGoバックエンドの設計
- **OpenAPI (oapi-codegen)** を用いた型安全なAPI開発

## 開発状況

| 機能 | 状態 |
|------|------|
| Public API | ✅ 完成 |
| Admin API | 🚧 未実装 |
| Frontend (shadcn/ui) | ✅ 完成 |

## 技術スタック

### バックエンド
| カテゴリ | 技術 |
|----------|------|
| 言語 | Go 1.24 |
| フレームワーク | Echo v4 |
| ORM | Ent (コード生成型) |
| データベース | PostgreSQL 17 |
| マイグレーション | Atlas |
| 認証 | JWT (アクセストークン + リフレッシュトークン) |
| DI | Uber Dig |
| API仕様 | OpenAPI 3.0 + oapi-codegen |

### フロントエンド
| カテゴリ | 技術 |
|----------|------|
| フレームワーク | Next.js 16 (App Router) |
| 言語 | TypeScript 5 |
| UIライブラリ | React 19 |
| スタイリング | Tailwind CSS v4 |
| 状態管理 | TanStack React Query v5 |
| フォーム | React Hook Form + Zod |
| UIコンポーネント | Radix UI + shadcn/ui |
| APIクライアント | Orval (自動生成) |

### インフラ
- Docker & Docker Compose
- MailHog (開発用メールサーバー)
- Atlas (マイグレーション管理)

## アーキテクチャ

### クリーンアーキテクチャ

```
backend/
├── cmd/                    # エントリーポイント
│   ├── api/               # APIサーバー
│   └── test_rls/          # RLSテスト用
├── internal/
│   ├── domain/            # ドメイン層 (モデル、リポジトリインターフェース)
│   ├── infrastructure/    # インフラ層 (DB接続、リポジトリ実装)
│   ├── usecase/           # ユースケース層 (ビジネスロジック)
│   └── presentation/      # プレゼンテーション層
│       ├── public/        # 公開API (v1) ✅
│       └── admin/         # 管理API 🚧
├── openapi/               # OpenAPI仕様
└── Makefile               # ビルド・開発コマンド
```

### RLS によるテナント分離

このプロジェクトの核心となる機能です。PostgreSQL の Row Level Security を使用して、テナント間のデータを完全に分離しています。

```sql
-- RLS の有効化
ALTER TABLE "users" ENABLE ROW LEVEL SECURITY;
ALTER TABLE "todos" ENABLE ROW LEVEL SECURITY;

-- テナント分離ポリシー
CREATE POLICY "users_tenant_isolation" ON "users"
    FOR ALL
    USING ("tenant_id" = current_setting('app.current_tenant_id', true))
    WITH CHECK ("tenant_id" = current_setting('app.current_tenant_id', true));
```

**動作の仕組み:**
1. JWTトークンに `tenant_id` を含める
2. リクエスト時にミドルウェアで `SET app.current_tenant_id = '<tenant_id>'` を実行
3. RLSポリシーにより、自動的に該当テナントのデータのみにアクセス可能

## 主な機能

### 認証・認可
- メールアドレスによるユーザー登録
- メール認証 (トークン方式)
- JWT認証 (アクセストークン + リフレッシュトークン)
- 自動トークンリフレッシュ

### マルチテナント
- ユーザー登録時にテナント (ワークスペース) を自動作成
- RLS によるデータ分離 (アプリケーションコードでの明示的なフィルタリング不要)
- テナント間のデータ完全分離

### Todo管理
- Todo作成・編集・削除
- 完了状態の管理
- 公開/非公開設定 (テナント内での共有)
- 期日設定

## セットアップ

### 必要条件
- Go 1.24+
- Node.js 20+
- Docker & Docker Compose

### バックエンド

```bash
cd backend

# 環境変数ファイルを作成
cp .env.example .env

# Dockerサービスを起動 (PostgreSQL, MailHog, 自動マイグレーション)
docker compose up -d

# または、手動でマイグレーションを実行する場合
make run              # DBのみ起動
make migrate_apply    # マイグレーション適用

# 開発サーバーを起動
make dev
```

### フロントエンド

```bash
cd frontend

# 依存関係をインストール
npm install

# 環境変数ファイルを作成
cp .env.example .env

# APIクライアントを生成
npm run generate:api

# 開発サーバーを起動
npm run dev
```

### アクセス先
| サービス | URL |
|----------|-----|
| Frontend | http://localhost:3000 |
| Public API | http://localhost:8000 |
| MailHog (メール確認) | http://localhost:8025 |

## 開発コマンド

### バックエンド (Makefile)

```bash
# Docker
make run                 # Dockerサービス起動
make stop                # Dockerサービス停止

# 開発
make dev                 # 開発サーバー起動 (ホットリロードなし)

# コード生成
make generate_ent        # Ent ORMコード生成
make oapi-gen            # OpenAPIコード生成

# マイグレーション
make migrate_diff        # マイグレーション作成
make migrate_apply       # マイグレーション適用
make migrate_status      # マイグレーション状態確認
make migrate_down n=1    # マイグレーションロールバック

# テスト
make test_unit           # ユニットテスト
make test_integration    # 統合テスト

# コード品質
make fmt                 # フォーマット
make lint                # リント
make vet                 # vet
```

### フロントエンド (npm scripts)

```bash
npm run dev              # 開発サーバー起動
npm run build            # ビルド
npm run generate:api     # APIクライアント生成
npm run lint             # リント
```

## API エンドポイント

### Public API (実装済み)

#### 認証
| メソッド | パス | 説明 |
|---------|------|------|
| POST | `/api/v1/auth/register` | ユーザー登録 |
| POST | `/api/v1/auth/login` | ログイン |
| POST | `/api/v1/auth/verify-email` | メール認証 |
| POST | `/api/v1/auth/refresh` | トークンリフレッシュ |

#### ユーザー
| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/api/v1/me` | 現在のユーザー情報取得 |
| PUT | `/api/v1/me` | プロフィール更新 |

#### Todo
| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/api/v1/todos` | 自分のTodo一覧取得 |
| GET | `/api/v1/todos/public` | 公開Todo一覧取得 (テナント内) |
| POST | `/api/v1/todos` | Todo作成 |
| PUT | `/api/v1/todos/:id` | Todo更新 |
| DELETE | `/api/v1/todos/:id` | Todo削除 |

### Admin API (未実装)

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/api/admin/tenants` | テナント一覧 |
| POST | `/api/admin/tenants` | テナント作成 |
| ... | ... | ... |

## データベース設計

### テーブル構成

**tenants** - テナント (ワークスペース)
- `id` (UUID), `name`, `slug`, `created_at`, `updated_at`

**users** - ユーザー (RLS適用)
- `id` (UUID), `tenant_id`, `email`, `password_hash`, `name`, `role`
- `email_verified`, `verification_token`, `verification_token_expires_at`
- `created_at`, `updated_at`

**todos** - Todo (RLS適用)
- `id` (UUID), `tenant_id`, `user_id`, `title`, `description`
- `completed`, `is_public`, `due_date`, `completed_at`
- `created_at`, `updated_at`

### セキュリティ設計

```
┌─────────────────────────────────────────────────────────────┐
│                      Application                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              JWT Middleware                          │   │
│  │  - tenant_id をJWTから取得                           │   │
│  │  - SET app.current_tenant_id を実行                  │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     PostgreSQL                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Row Level Security                      │   │
│  │  - app.current_tenant_id でフィルタリング            │   │
│  │  - 他テナントのデータには一切アクセス不可            │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## ディレクトリ構成

```
good-todo-go/
├── backend/
│   ├── cmd/
│   │   ├── api/                    # メインエントリーポイント
│   │   └── test_rls/               # RLSテスト用
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── model/              # ドメインモデル
│   │   │   └── repository/         # リポジトリインターフェース
│   │   ├── ent/                    # Ent ORM (自動生成)
│   │   │   ├── schema/             # スキーマ定義
│   │   │   └── migrate/migrations/ # マイグレーションSQL
│   │   ├── infrastructure/         # リポジトリ実装、DB接続
│   │   ├── usecase/                # ユースケース (ビジネスロジック)
│   │   └── presentation/
│   │       ├── public/             # Public API ✅
│   │       └── admin/              # Admin API 🚧
│   ├── openapi/                    # OpenAPI仕様
│   ├── docker-compose.yml
│   ├── Dockerfile.api.local
│   ├── Dockerfile.migrate
│   └── Makefile
├── frontend/
│   ├── src/
│   │   ├── app/                    # Next.js App Router
│   │   │   ├── (auth)/             # 認証ページ
│   │   │   └── (main)/             # メインページ
│   │   ├── api/                    # 生成されたAPIクライアント
│   │   ├── components/
│   │   │   ├── ui/                 # shadcn/ui コンポーネント
│   │   │   ├── todo/               # Todoコンポーネント
│   │   │   └── user/               # ユーザーコンポーネント
│   │   ├── contexts/               # Reactコンテキスト
│   │   └── providers/              # プロバイダー
│   ├── orval.config.ts             # APIクライアント生成設定
│   └── package.json
└── README.md
```

## 環境変数

### バックエンド (.env)
```env
# PostgreSQL (管理者)
POSTGRES_DB_USER=postgres
POSTGRES_DB_PASSWORD=postgres
POSTGRES_DB_NAME=good_todo_go
POSTGRES_DB_PORT=5432
POSTGRES_DB_HOST=localhost

# PostgreSQL (アプリケーション用 - RLS適用)
POSTGRES_APP_USER=app
POSTGRES_APP_PASSWORD=app

# JWT
JWT_SECRET=your-super-secret-jwt-key

# Server
PUBLIC_API_PORT=8000
ADMIN_API_PORT=8001
```

### フロントエンド (.env)
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8000/api/v1
```

