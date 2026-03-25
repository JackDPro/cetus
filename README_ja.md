# Cetus

[English](README.md) | [中文](README_zh.md) | **日本語** | [한국어](README_ko.md)

[Gin](https://github.com/gin-gonic/gin) ベースの Go HTTP サービス開発ツールキット。設定管理、データベースアクセス（MySQL/PostgreSQL）、JWT 認証、構造化ログ、ミドルウェア、標準化された API レスポンスをすぐに利用できます。

## 機能

- **設定管理** - `.env` ファイルによる環境変数ベースの設定、シングルトンパターン
- **データベース** - GORM 経由で MySQL と PostgreSQL をサポート、ドライバ自動選択
- **JWT 認証** - RSA 署名トークン、Redis によるセッション管理
- **ログ** - Zap による構造化ログ（JSON/コンソール形式、ファイル/標準出力対応）
- **ミドルウェア** - リクエスト ID トラッキング、レート制限
- **API レスポンス** - Gin 用の標準化された JSON レスポンスヘルパー
- **ユーティリティ** - Bcrypt ハッシュ、ID 難読化（Optimus）、RSA 暗号、i18n、ファイル操作

## インストール

```bash
go get github.com/JackDPro/cetus
```

## クイックスタート

### 1. `.env` ファイルの作成

サンプルファイルをコピーし、環境に合わせて設定を変更します：

```bash
cp .env.example .env
```

主要な環境変数：

```env
APP_NAME=my-app
APP_ENV=dev

LOG_CONSOLE_OUT=true
LOG_FILE_OUT=false
LOG_LEVEL=debug
LOG_FORMAT=json

DB_TYPE=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=mydb
DB_USERNAME=root
DB_PASSWORD=password

REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_DATABASE=0
REDIS_PASSWORD=password

SERVER_HTTP_PORT=80
SERVER_GRPC_PORT=50051
```

すべてのオプションは [.env.example](.env.example) を参照してください。

### 2. サーバーの作成

```go
package main

import (
    "fmt"
    "github.com/JackDPro/cetus/config"
    "github.com/JackDPro/cetus/controller"
    "github.com/JackDPro/cetus/middleware"
    "github.com/JackDPro/cetus/provider"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.New()
    router.Use(gin.Recovery())

    // CORS 設定
    corsConf := cors.DefaultConfig()
    corsConf.AllowAllOrigins = true
    corsConf.AllowHeaders = []string{"Authorization", "Accept-Language"}
    router.Use(cors.New(corsConf))

    // リクエスト ID ミドルウェア
    router.Use(middleware.RequestId())

    // ヘルスチェックエンドポイント
    probeCtr := controller.NewProbeController()
    router.GET("/probe", probeCtr.Show)

    // ルートを追加
    // router.GET("/users", userController.Index)

    conf := config.GetApiConfig()
    addr := fmt.Sprintf("0.0.0.0:%d", conf.HttpPort)
    provider.GetLogger().Info("server starting", "address", addr)
    if err := router.Run(addr); err != nil {
        panic(err)
    }
}
```

## パッケージガイド

### config - 設定管理

すべての設定は環境変数（`.env` ファイル経由）から読み込まれます。各設定構造体はシングルトンパターンを使用します。

```go
import "github.com/JackDPro/cetus/config"

appConf := config.GetAppConfig()      // アプリ名、環境、データルート、公開URL
dbConf  := config.GetDatabaseConfig() // DB タイプ、ホスト、ポート、認証情報
apiConf := config.GetApiConfig()      // HTTP と gRPC ポート
authConf := config.GetAuthConf()      // JWT 証明書/鍵パス、有効期限
redisConf := config.GetRedisConfig()  // Redis 接続情報
logConf := config.GetLogConfig()      // ログレベル、形式、出力先
hashConf := config.GetHashIdConfig()  // Optimus ID ハッシュパラメータ
natsConf := config.GetNatsConf()      // NATS メッセージング設定
```

**`APP_ENV` で動作を制御：**
```go
if config.GetAppConfig().Env == "prod" {
    gin.SetMode(gin.ReleaseMode)
}
```

### provider - コアサービス

スレッドセーフなシングルトンプロバイダー。

#### ロガー

[Zap](https://github.com/uber-go/zap) による構造化ログ。コンソール/JSON 形式、ファイル/標準出力をサポート。

```go
import "github.com/JackDPro/cetus/provider"

logger := provider.GetLogger()
logger.Info("サーバー起動")
logger.Infow("ユーザー作成", "userId", 123)
logger.Errorw("リクエスト失敗", "error", err)
```

#### データベース (ORM)

[GORM](https://gorm.io/) ラッパー。`DB_TYPE` に基づいて MySQL/PostgreSQL ドライバを自動選択。

```go
import "github.com/JackDPro/cetus/provider"

db := provider.GetOrm().Db

// 標準的な GORM 操作
var users []User
db.Where("active = ?", true).Find(&users)

// オートマイグレーション
db.AutoMigrate(&User{}, &Post{})
```

#### Redis

```go
import "github.com/JackDPro/cetus/provider"

rdb := provider.GetRedisClient()
rdb.Set(ctx, "key", "value", time.Hour)
val, err := rdb.Get(ctx, "key").Result()
```

#### パスワードハッシュ (Bcrypt)

```go
import "github.com/JackDPro/cetus/provider"

// ハッシュ生成
hashed, err := provider.HashMake("my-password")

// 検証
err = provider.HashCheck("my-password", hashed)

// リハッシュが必要か確認
if provider.HashNeedRefresh(hashed) {
    // 現在のコストでリハッシュ
}
```

#### ID 難読化 (Optimus)

可逆的な整数 ID エンコーディング。API でデータベースの連番 ID を隠蔽します。

```go
import "github.com/JackDPro/cetus/provider"

encoded := provider.Hash().Encode(42)         // 例: 1580030173
decoded := provider.Hash().Decode(1580030173) // 42
```

`.env` に `OPTIMUS_PRIME`、`OPTIMUS_INVERSE`、`OPTIMUS_RANDOM` の設定が必要です。

#### RSA 暗号

```go
import "github.com/JackDPro/cetus/provider"

// ファイルから読み込み
rsaKey, err := provider.NewRsaByKeyFile("path/to/private.pem")

// 署名
signature, err := rsaKey.Signature(data)

// 署名検証
err = rsaKey.VerifySignature(data, signature)
```

#### 文字列・配列ユーティリティ

```go
import "github.com/JackDPro/cetus/provider"

provider.RandomString(16)            // "aB3kF9mNpQ2xR7wL"
provider.RandomInt(6)                // "483921"
provider.StringInSlice("a", []string{"a", "b"}) // true
```

#### ファイル操作

```go
import "github.com/JackDPro/cetus/provider"

provider.CopyFile("src.txt", "dst.txt")
provider.CopyDir("src_dir", "dst_dir") // 再帰コピー
```

#### 国際化 (i18n)

```go
import "github.com/JackDPro/cetus/provider"

t := provider.GetTranslate()
t.SetLanguage("ja")
msg := t.Tr("hello_world")

// グローバルショートカット
msg := provider.Tr("hello_world")
```

### model - データモデルとレスポンス構造体

#### BaseModel

モデルに `BaseModel` を埋め込むことで、シリアライズヘルパーを利用できます：

```go
import "github.com/JackDPro/cetus/model"

type User struct {
    model.BaseModel
    Id       uint64 `json:"id" gorm:"primaryKey"`
    Nickname string `json:"nickname"`
    Email    string `json:"email"`
}

// IModel インターフェースの実装
func (u *User) ToMap() (map[string]interface{}, error) {
    return u.BaseModel.ToMap(u)
}
```

`BaseModel` が提供するメソッド：
- `ToMap(model)` - `json` タグを使用して構造体を map に変換
- `ToJson(model)` - 構造体を JSON 文字列に変換
- `Include(model, relations, db)` - カンマ区切り文字列から GORM リレーションをプリロード

#### IModel インターフェース

`IModel` を実装したモデルはレスポンスヘルパーで使用できます：

```go
type IModel interface {
    ToMap() (map[string]interface{}, error)
}
```

#### レスポンス構造体

```go
// API レスポンスは DataWrapper でラップされます
type DataWrapper struct {
    Data interface{} `json:"data"`
    Meta interface{} `json:"meta,omitempty"`
}

// エラーレスポンス
type Error struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// ページネーションメタデータ
type Pagination struct {
    Count       int `json:"count"`
    CurrentPage int `json:"current_page"`
    PerPage     int `json:"per_page"`
    Total       int `json:"total"`
    TotalPages  int `json:"total_pages"`
}
```

### controller - レスポンスヘルパー

Gin ハンドラー用の標準化された HTTP レスポンス関数：

```go
import "github.com/JackDPro/cetus/controller"

func (ctr *UserController) Store(c *gin.Context) {
    // 入力検証
    if err := c.ShouldBindJSON(&req); err != nil {
        controller.ResponseUnprocessable(c, 1, "パラメータが無効です", err)
        return
    }

    // リソース作成
    user, err := createUser(req)
    if err != nil {
        controller.ResponseInternalError(c, 2, "作成に失敗しました", err)
        return
    }

    // 作成されたリソースを返す
    controller.ResponseCreated(c, user.Id)
}

func (ctr *UserController) Show(c *gin.Context) {
    user, err := findUser(id)
    if err != nil {
        controller.ResponseNotFound(c, "ユーザーが見つかりません")
        return
    }
    controller.ResponseItem(c, user) // user は IModel を実装する必要があります
}

func (ctr *UserController) Index(c *gin.Context) {
    users, meta := listUsers(page, perPage)
    controller.ResponseCollection(c, users, meta)
}
```

**利用可能なレスポンスヘルパー：**

| 関数 | HTTP ステータス | 用途 |
|------|---------------|------|
| `ResponseSuccess(c)` | 200 | 操作成功 |
| `ResponseItem(c, item)` | 200 | 単一リソースを返す |
| `ResponseCollection(c, items, meta)` | 200 | リストを返す（ページネーション付き） |
| `ResponseCreated(c, id)` | 201 | リソース作成済み（`Location` ヘッダー設定） |
| `ResponseAccepted(c)` | 202 | 非同期操作を受付 |
| `ResponseBadRequest(c, code, msg)` | 400 | 不正なリクエスト |
| `ResponseUnauthorized(c)` | 401 | 認証が必要 |
| `ResponseForbidden(c)` | 403 | アクセス権限なし |
| `ResponseNotFound(c, msg)` | 404 | リソースが見つからない |
| `ResponseUnprocessable(c, code, msg, err)` | 422 | バリデーション失敗 |
| `ResponseInternalError(c, code, msg, err)` | 500 | サーバーエラー（自動ログ記録） |

### jwt - JWT 認証

RSA 署名の JWT トークン、Redis によるセッションストレージ。アクセストークンとリフレッシュトークンをサポート。

**前提条件：**
- RSA 鍵ペア（PKCS#8 秘密鍵 + PEM 公開鍵証明書）
- Redis サーバー
- `.env` の JWT 設定

```env
JWT_EXPIRES_IN=72
JWT_REDIS_PREFIX=auth
JWT_CERT_PATH=storage/jwt.crt
JWT_KEY_PATH=storage/jwt.key
JWT_ISSUE=example.com
```

**RSA 鍵の生成：**
```bash
openssl genpkey -algorithm RSA -out storage/jwt.key -pkeyopt rsa_keygen_bits:2048
openssl rsa -in storage/jwt.key -pubout -out storage/jwt.crt
```

**使い方：**
```go
import "github.com/JackDPro/cetus/jwt"

guard, err := jwt.GetJwtGuard()

// アクセストークン + リフレッシュトークンの作成
accessToken, err := guard.CreateToken(userId, true) // true = 古いトークンを無効化
// accessToken.AccessToken  - JWT 文字列
// accessToken.RefreshToken - リフレッシュトークン JWT 文字列
// accessToken.Type         - "bearer"
// accessToken.ExpiresIn    - 有効期限（秒）

// アクセストークンのみ作成（例：API キー用）
accessToken, err := guard.CreateAccessToken(userId)

// トークン検証
validToken, err := guard.Attempt(tokenString)
// validToken.UserId - 認証済みユーザー ID
// validToken.Type   - "token"、"refresh"、または "access_key"

// トークン無効化（ログアウト）
err = guard.DeleteCredential(tokenString)
```

### middleware - HTTP ミドルウェア

#### リクエスト ID

すべてのリクエストに一意のリクエスト ID を付与します。まずリクエストヘッダー（`X-Request-Id`、`HTTP_X_REQUEST_ID`、`HTTP_REQUEST_ID`）を確認し、存在しない場合は UUID を自動生成します。

```go
import "github.com/JackDPro/cetus/middleware"

router.Use(middleware.RequestId())
```

#### レート制限

[tollbooth](https://github.com/didip/tollbooth) を使用。

```go
import (
    "github.com/JackDPro/cetus/middleware"
    "github.com/didip/tollbooth/v7/limiter"
)

lmt := limiter.New(nil).SetMax(100) // 毎秒 100 リクエスト
router.Use(middleware.LimitRate(lmt))
```

## GORM モデルのリレーションプリロード

`BaseModel.Include()` でクエリパラメータからリレーションを動的にプリロード：

```go
// GET /users?include=posts,comments
func (ctr *UserController) Index(c *gin.Context) {
    includeStr := c.Query("include")
    var users []User

    db := provider.GetOrm().Db
    db = (&model.BaseModel{}).Include(&User{}, includeStr, db)
    db.Find(&users)

    controller.ResponseCollection(c, users, nil)
}
```

## 完全な例

REST API の完全なセットアップ：

```go
package main

import (
    "fmt"
    "github.com/JackDPro/cetus/config"
    "github.com/JackDPro/cetus/controller"
    "github.com/JackDPro/cetus/middleware"
    "github.com/JackDPro/cetus/provider"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    appConf := config.GetAppConfig()
    if appConf.Env == "prod" {
        gin.SetMode(gin.ReleaseMode)
    }

    router := gin.New()
    router.Use(gin.Recovery())

    // CORS 設定
    corsConf := cors.DefaultConfig()
    corsConf.AllowAllOrigins = true
    corsConf.AllowHeaders = []string{"Authorization", "Accept-Language"}
    router.Use(cors.New(corsConf))

    // ミドルウェア
    router.Use(middleware.RequestId())

    // 公開ルート
    probeCtr := controller.NewProbeController()
    router.GET("/probe", probeCtr.Show)
    router.POST("/auth/login", authLogin)
    router.POST("/users", createUser)

    // 保護されたルート（認証ミドルウェアを追加）
    authorized := router.Group("/")
    // authorized.Use(yourAuthMiddleware())
    authorized.GET("/users/me", getMe)
    authorized.POST("/auth/logout", logout)

    // サーバー起動
    apiConf := config.GetApiConfig()
    addr := fmt.Sprintf("0.0.0.0:%d", apiConf.HttpPort)
    provider.GetLogger().Info("server starting", "address", addr)
    if err := router.Run(addr); err != nil {
        panic(err)
    }
}
```

## GORM Hook による自動パスワードハッシュ

GORM モデルフックで bcrypt による自動パスワード暗号化：

```go
import (
    "github.com/JackDPro/cetus/model"
    "github.com/JackDPro/cetus/provider"
    "gorm.io/gorm"
)

type User struct {
    model.BaseModel
    Id       uint64 `json:"id" gorm:"primaryKey"`
    Email    string `json:"email"`
    Password string `json:"password,omitempty"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
    if u.Password != "" {
        u.Password, err = provider.HashMake(u.Password)
    }
    return
}

func (u *User) ToMap() (map[string]interface{}, error) {
    return u.BaseModel.ToMap(u)
}
```

## 環境変数リファレンス

| 変数 | 説明 | デフォルト/例 |
|------|------|-------------|
| `APP_NAME` | アプリケーション名 | `my-app` |
| `APP_ENV` | 実行環境（`dev`/`prod`） | `dev` |
| `APP_DATA_ROOT` | データ保存ルートパス | `/usr/app` |
| `APP_PUBLIC_RES_URL` | 公開静的リソース URL | `http://example.com/static` |
| `LOG_CONSOLE_OUT` | コンソールにログ出力 | `true` |
| `LOG_FILE_OUT` | ファイルにログ出力 | `false` |
| `LOG_FILE_PATH` | ログファイルパス | `/var/log/app.log` |
| `LOG_LEVEL` | ログレベル（`debug`/`info`/`warn`/`error`） | `debug` |
| `LOG_FORMAT` | ログ形式（`json`/`console`） | `json` |
| `DB_TYPE` | データベースタイプ（`mysql`/`postgres`） | `mysql` |
| `DB_HOST` | データベースホスト | `127.0.0.1` |
| `DB_PORT` | データベースポート | `3306` |
| `DB_DATABASE` | データベース名 | `mydb` |
| `DB_USERNAME` | データベースユーザー名 | `root` |
| `DB_PASSWORD` | データベースパスワード | - |
| `DB_SSLMODE` | PostgreSQL SSL モード | `disable` |
| `DB_MIGRATE_SELF_ONLY` | マイグレーション範囲を制限 | `false` |
| `JWT_EXPIRES_IN` | トークン有効期限（時間） | `72` |
| `JWT_REDIS_PREFIX` | Redis トークンプレフィックス | `auth` |
| `JWT_CERT_PATH` | RSA 公開鍵パス | `storage/jwt.crt` |
| `JWT_KEY_PATH` | RSA 秘密鍵パス | `storage/jwt.key` |
| `JWT_ISSUE` | JWT 発行者 | `example.com` |
| `OPTIMUS_PRIME` | Optimus 素数 | - |
| `OPTIMUS_INVERSE` | Optimus 逆数 | - |
| `OPTIMUS_RANDOM` | Optimus ランダムシード | - |
| `REDIS_HOST` | Redis ホスト | `127.0.0.1` |
| `REDIS_PORT` | Redis ポート | `6379` |
| `REDIS_DATABASE` | Redis データベースインデックス | `0` |
| `REDIS_PASSWORD` | Redis パスワード | - |
| `NATS_HOST` | NATS サーバーホスト | `127.0.0.1` |
| `NATS_USERNAME` | NATS ユーザー名 | - |
| `NATS_PASSWORD` | NATS パスワード | - |
| `SERVER_HTTP_PORT` | HTTP サーバーポート | `80` |
| `SERVER_GRPC_PORT` | gRPC サーバーポート | `50051` |

## ライセンス

MIT
