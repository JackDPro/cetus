# Cetus

**English** | [中文](README_zh.md) | [日本語](README_ja.md) | [한국어](README_ko.md)

A Go toolkit for building HTTP services with [Gin](https://github.com/gin-gonic/gin). Provides configuration management, database access (MySQL/PostgreSQL), JWT authentication, structured logging, middleware, and standardized API responses out of the box.

## Features

- **Configuration** - Environment-based config with `.env` files, singleton pattern
- **Database** - MySQL & PostgreSQL via GORM, auto driver selection
- **JWT Authentication** - RSA-signed tokens with Redis-backed session management
- **Logging** - Structured logging with Zap (JSON/console, file/stdout)
- **Middleware** - Request ID tracking, rate limiting
- **API Responses** - Standardized JSON response helpers for Gin
- **Utilities** - Bcrypt hashing, ID obfuscation (Optimus), RSA crypto, i18n, file operations

## Installation

```bash
go get github.com/JackDPro/cetus
```

## Quick Start

### 1. Create `.env` file

Copy the example and modify values for your environment:

```bash
cp .env.example .env
```

Key environment variables:

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

See [.env.example](.env.example) for all available options.

### 2. Create your server

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

    // CORS
    corsConf := cors.DefaultConfig()
    corsConf.AllowAllOrigins = true
    corsConf.AllowHeaders = []string{"Authorization", "Accept-Language"}
    router.Use(cors.New(corsConf))

    // Request ID middleware
    router.Use(middleware.RequestId())

    // Health check endpoint
    probeCtr := controller.NewProbeController()
    router.GET("/probe", probeCtr.Show)

    // Your routes here
    // router.GET("/users", userController.Index)

    conf := config.GetApiConfig()
    addr := fmt.Sprintf("0.0.0.0:%d", conf.HttpPort)
    provider.GetLogger().Info("server starting", "address", addr)
    if err := router.Run(addr); err != nil {
        panic(err)
    }
}
```

## Package Guide

### config - Configuration Management

All configuration is loaded from environment variables (via `.env` file). Each config struct uses the singleton pattern.

```go
import "github.com/JackDPro/cetus/config"

appConf := config.GetAppConfig()     // App name, env, data root, public URL
dbConf  := config.GetDatabaseConfig() // DB type, host, port, credentials
apiConf := config.GetApiConfig()     // HTTP and gRPC ports
authConf := config.GetAuthConf()     // JWT cert/key paths, expiration
redisConf := config.GetRedisConfig() // Redis connection info
logConf := config.GetLogConfig()     // Log level, format, output targets
hashConf := config.GetHashIdConfig() // Optimus ID hashing params
natsConf := config.GetNatsConf()     // NATS messaging config
```

**Use `APP_ENV` to control behavior:**
```go
if config.GetAppConfig().Env == "prod" {
    gin.SetMode(gin.ReleaseMode)
}
```

### provider - Core Services

Thread-safe singleton providers for common infrastructure.

#### Logger

Structured logging powered by [Zap](https://github.com/uber-go/zap). Supports console/JSON format, file/stdout output.

```go
import "github.com/JackDPro/cetus/provider"

logger := provider.GetLogger()
logger.Info("server started")
logger.Infow("user created", "userId", 123)
logger.Errorw("request failed", "error", err)
```

#### Database (ORM)

[GORM](https://gorm.io/) wrapper with automatic MySQL/PostgreSQL driver selection based on `DB_TYPE`.

```go
import "github.com/JackDPro/cetus/provider"

db := provider.GetOrm().Db

// Use standard GORM operations
var users []User
db.Where("active = ?", true).Find(&users)

// Auto-migrate
db.AutoMigrate(&User{}, &Post{})
```

#### Redis

```go
import "github.com/JackDPro/cetus/provider"

rdb := provider.GetRedisClient()
rdb.Set(ctx, "key", "value", time.Hour)
val, err := rdb.Get(ctx, "key").Result()
```

#### Password Hashing (Bcrypt)

```go
import "github.com/JackDPro/cetus/provider"

// Hash a password
hashed, err := provider.HashMake("my-password")

// Verify
err = provider.HashCheck("my-password", hashed)

// Check if re-hash is needed
if provider.HashNeedRefresh(hashed) {
    // re-hash with current cost
}
```

#### ID Obfuscation (Optimus)

Reversible integer ID encoding to hide sequential database IDs in APIs.

```go
import "github.com/JackDPro/cetus/provider"

encoded := provider.Hash().Encode(42)    // e.g. 1580030173
decoded := provider.Hash().Decode(1580030173) // 42
```

Requires `OPTIMUS_PRIME`, `OPTIMUS_INVERSE`, `OPTIMUS_RANDOM` in `.env`.

#### RSA Cryptography

```go
import "github.com/JackDPro/cetus/provider"

// Load from file
rsaKey, err := provider.NewRsaByKeyFile("path/to/private.pem")

// Sign data
signature, err := rsaKey.Signature(data)

// Verify signature
err = rsaKey.VerifySignature(data, signature)
```

#### String & Array Utilities

```go
import "github.com/JackDPro/cetus/provider"

provider.RandomString(16)            // "aB3kF9mNpQ2xR7wL"
provider.RandomInt(6)                // "483921"
provider.StringInSlice("a", []string{"a", "b"}) // true
```

#### File Operations

```go
import "github.com/JackDPro/cetus/provider"

provider.CopyFile("src.txt", "dst.txt")
provider.CopyDir("src_dir", "dst_dir") // recursive copy
```

#### i18n Translation

```go
import "github.com/JackDPro/cetus/provider"

t := provider.GetTranslate()
t.SetLanguage("zh")
msg := t.Tr("hello_world")

// Or use the global shorthand
msg := provider.Tr("hello_world")
```

### model - Data Models & Response Structures

#### BaseModel

Embed `BaseModel` in your models to get serialization helpers:

```go
import "github.com/JackDPro/cetus/model"

type User struct {
    model.BaseModel
    Id       uint64 `json:"id" gorm:"primaryKey"`
    Nickname string `json:"nickname"`
    Email    string `json:"email"`
}

// Implement IModel interface
func (u *User) ToMap() (map[string]interface{}, error) {
    return u.BaseModel.ToMap(u)
}
```

`BaseModel` provides:
- `ToMap(model)` - Convert struct to map using `json` tags
- `ToJson(model)` - Convert struct to JSON string
- `Include(model, relations, db)` - Preload GORM relations from a comma-separated string

#### IModel Interface

Any model implementing `IModel` can be used with the response helpers:

```go
type IModel interface {
    ToMap() (map[string]interface{}, error)
}
```

#### Response Structures

```go
// API responses are wrapped in DataWrapper
type DataWrapper struct {
    Data interface{} `json:"data"`
    Meta interface{} `json:"meta,omitempty"`
}

// Error response
type Error struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// Pagination metadata
type Pagination struct {
    Count      int `json:"count"`
    CurrentPage int `json:"current_page"`
    PerPage    int `json:"per_page"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}
```

### controller - Response Helpers

Standardized HTTP response functions for Gin handlers:

```go
import "github.com/JackDPro/cetus/controller"

func (ctr *UserController) Store(c *gin.Context) {
    // Validate input
    if err := c.ShouldBindJSON(&req); err != nil {
        controller.ResponseUnprocessable(c, 1, "invalid params", err)
        return
    }

    // Create resource
    user, err := createUser(req)
    if err != nil {
        controller.ResponseInternalError(c, 2, "create failed", err)
        return
    }

    // Return created resource
    controller.ResponseCreated(c, user.Id)
}

func (ctr *UserController) Show(c *gin.Context) {
    user, err := findUser(id)
    if err != nil {
        controller.ResponseNotFound(c, "user not found")
        return
    }
    controller.ResponseItem(c, user) // user must implement IModel
}

func (ctr *UserController) Index(c *gin.Context) {
    users, meta := listUsers(page, perPage)
    controller.ResponseCollection(c, users, meta)
}
```

**Available response helpers:**

| Function | HTTP Status | Use Case |
|----------|------------|----------|
| `ResponseSuccess(c)` | 200 | Action succeeded |
| `ResponseItem(c, item)` | 200 | Return single resource |
| `ResponseCollection(c, items, meta)` | 200 | Return list with optional pagination |
| `ResponseCreated(c, id)` | 201 | Resource created (sets `Location` header) |
| `ResponseAccepted(c)` | 202 | Async action accepted |
| `ResponseBadRequest(c, code, msg)` | 400 | Invalid request |
| `ResponseUnauthorized(c)` | 401 | Authentication required |
| `ResponseForbidden(c)` | 403 | Permission denied |
| `ResponseNotFound(c, msg)` | 404 | Resource not found |
| `ResponseUnprocessable(c, code, msg, err)` | 422 | Validation failed |
| `ResponseInternalError(c, code, msg, err)` | 500 | Server error (auto-logged) |

### jwt - JWT Authentication

RSA-signed JWT tokens with Redis-backed session storage. Supports access tokens and refresh tokens.

**Prerequisites:**
- RSA key pair (PKCS#8 private key + PEM public certificate)
- Redis server
- JWT config in `.env`

```env
JWT_EXPIRES_IN=72
JWT_REDIS_PREFIX=auth
JWT_CERT_PATH=storage/jwt.crt
JWT_KEY_PATH=storage/jwt.key
JWT_ISSUE=example.com
```

**Generate RSA keys:**
```bash
openssl genpkey -algorithm RSA -out storage/jwt.key -pkeyopt rsa_keygen_bits:2048
openssl rsa -in storage/jwt.key -pubout -out storage/jwt.crt
```

**Usage:**
```go
import "github.com/JackDPro/cetus/jwt"

guard, err := jwt.GetJwtGuard()

// Create access token + refresh token
accessToken, err := guard.CreateToken(userId, true) // true = revoke old tokens
// accessToken.AccessToken  - JWT string
// accessToken.RefreshToken - Refresh JWT string
// accessToken.Type         - "bearer"
// accessToken.ExpiresIn    - Expiration in seconds

// Create access token only (e.g. for API keys)
accessToken, err := guard.CreateAccessToken(userId)

// Validate token
validToken, err := guard.Attempt(tokenString)
// validToken.UserId - Authenticated user ID
// validToken.Type   - "token", "refresh", or "access_key"

// Revoke token (logout)
err = guard.DeleteCredential(tokenString)
```

### middleware - HTTP Middleware

#### Request ID

Adds a unique request ID to every request. Checks incoming headers first (`X-Request-Id`, `HTTP_X_REQUEST_ID`, `HTTP_REQUEST_ID`), generates a UUID if none found.

```go
import "github.com/JackDPro/cetus/middleware"

router.Use(middleware.RequestId())
```

#### Rate Limiting

Powered by [tollbooth](https://github.com/didip/tollbooth).

```go
import (
    "github.com/JackDPro/cetus/middleware"
    "github.com/didip/tollbooth/v7/limiter"
)

lmt := limiter.New(nil).SetMax(100) // 100 requests per second
router.Use(middleware.LimitRate(lmt))
```

## GORM Model with Preloading

Use `BaseModel.Include()` to dynamically preload relations from query parameters:

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

## Full Example

A complete REST API setup:

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

    // CORS
    corsConf := cors.DefaultConfig()
    corsConf.AllowAllOrigins = true
    corsConf.AllowHeaders = []string{"Authorization", "Accept-Language"}
    router.Use(cors.New(corsConf))

    // Middleware
    router.Use(middleware.RequestId())

    // Public routes
    probeCtr := controller.NewProbeController()
    router.GET("/probe", probeCtr.Show)
    router.POST("/auth/login", authLogin)
    router.POST("/users", createUser)

    // Protected routes (add your auth middleware)
    authorized := router.Group("/")
    // authorized.Use(yourAuthMiddleware())
    authorized.GET("/users/me", getMe)
    authorized.POST("/auth/logout", logout)

    // Start
    apiConf := config.GetApiConfig()
    addr := fmt.Sprintf("0.0.0.0:%d", apiConf.HttpPort)
    provider.GetLogger().Info("server starting", "address", addr)
    if err := router.Run(addr); err != nil {
        panic(err)
    }
}
```

## Password Hashing with GORM Hooks

Use bcrypt hashing in GORM model hooks for automatic password encryption:

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

## Environment Variables Reference

| Variable | Description | Default/Example |
|----------|-------------|-----------------|
| `APP_NAME` | Application name | `my-app` |
| `APP_ENV` | Environment (`dev`/`prod`) | `dev` |
| `APP_DATA_ROOT` | Data storage root path | `/usr/app` |
| `APP_PUBLIC_RES_URL` | Public static resource URL | `http://example.com/static` |
| `LOG_CONSOLE_OUT` | Output logs to console | `true` |
| `LOG_FILE_OUT` | Output logs to file | `false` |
| `LOG_FILE_PATH` | Log file path | `/var/log/app.log` |
| `LOG_LEVEL` | Log level (`debug`/`info`/`warn`/`error`) | `debug` |
| `LOG_FORMAT` | Log format (`json`/`console`) | `json` |
| `DB_TYPE` | Database type (`mysql`/`postgres`) | `mysql` |
| `DB_HOST` | Database host | `127.0.0.1` |
| `DB_PORT` | Database port | `3306` |
| `DB_DATABASE` | Database name | `mydb` |
| `DB_USERNAME` | Database username | `root` |
| `DB_PASSWORD` | Database password | - |
| `DB_SSLMODE` | PostgreSQL SSL mode | `disable` |
| `DB_MIGRATE_SELF_ONLY` | Limit migration scope | `false` |
| `JWT_EXPIRES_IN` | Token expiration (hours) | `72` |
| `JWT_REDIS_PREFIX` | Redis key prefix for tokens | `auth` |
| `JWT_CERT_PATH` | RSA public key path | `storage/jwt.crt` |
| `JWT_KEY_PATH` | RSA private key path | `storage/jwt.key` |
| `JWT_ISSUE` | JWT issuer claim | `example.com` |
| `OPTIMUS_PRIME` | Optimus prime number | - |
| `OPTIMUS_INVERSE` | Optimus inverse | - |
| `OPTIMUS_RANDOM` | Optimus random seed | - |
| `REDIS_HOST` | Redis host | `127.0.0.1` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_DATABASE` | Redis database index | `0` |
| `REDIS_PASSWORD` | Redis password | - |
| `NATS_HOST` | NATS server host | `127.0.0.1` |
| `NATS_USERNAME` | NATS username | - |
| `NATS_PASSWORD` | NATS password | - |
| `SERVER_HTTP_PORT` | HTTP server port | `80` |
| `SERVER_GRPC_PORT` | gRPC server port | `50051` |

## Demo Project

See [cetus-demo](https://github.com/JackDPro/cetus-demo) for a complete working example with user registration, JWT authentication, and CRUD operations.

## License

MIT
