# Cetus

[English](README.md) | **中文** | [日本語](README_ja.md) | [한국어](README_ko.md)

一个基于 [Gin](https://github.com/gin-gonic/gin) 的 Go HTTP 服务开发工具包。开箱即用，提供配置管理、数据库访问（MySQL/PostgreSQL）、JWT 认证、结构化日志、中间件及标准化 API 响应等功能。

## 特性

- **配置管理** - 基于 `.env` 文件的环境变量配置，单例模式
- **数据库** - 通过 GORM 支持 MySQL 和 PostgreSQL，自动选择驱动
- **JWT 认证** - RSA 签名的令牌，Redis 存储会话管理
- **日志系统** - 基于 Zap 的结构化日志（支持 JSON/控制台格式，文件/标准输出）
- **中间件** - 请求 ID 追踪、限流
- **API 响应** - 标准化的 Gin JSON 响应辅助函数
- **工具集** - Bcrypt 密码哈希、ID 混淆（Optimus）、RSA 加密、国际化、文件操作

## 安装

```bash
go get github.com/JackDPro/cetus
```

## 快速开始

### 1. 创建 `.env` 文件

复制示例文件并根据你的环境修改配置：

```bash
cp .env.example .env
```

主要环境变量：

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

完整配置项请参考 [.env.example](.env.example)。

### 2. 创建服务

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

    // 跨域配置
    corsConf := cors.DefaultConfig()
    corsConf.AllowAllOrigins = true
    corsConf.AllowHeaders = []string{"Authorization", "Accept-Language"}
    router.Use(cors.New(corsConf))

    // 请求 ID 中间件
    router.Use(middleware.RequestId())

    // 健康检查端点
    probeCtr := controller.NewProbeController()
    router.GET("/probe", probeCtr.Show)

    // 在此添加你的路由
    // router.GET("/users", userController.Index)

    conf := config.GetApiConfig()
    addr := fmt.Sprintf("0.0.0.0:%d", conf.HttpPort)
    provider.GetLogger().Info("server starting", "address", addr)
    if err := router.Run(addr); err != nil {
        panic(err)
    }
}
```

## 包使用指南

### config - 配置管理

所有配置通过环境变量（`.env` 文件）加载，每个配置结构体使用单例模式。

```go
import "github.com/JackDPro/cetus/config"

appConf := config.GetAppConfig()      // 应用名称、环境、数据根目录、公共资源 URL
dbConf  := config.GetDatabaseConfig() // 数据库类型、主机、端口、凭证
apiConf := config.GetApiConfig()      // HTTP 和 gRPC 端口
authConf := config.GetAuthConf()      // JWT 证书/密钥路径、过期时间
redisConf := config.GetRedisConfig()  // Redis 连接信息
logConf := config.GetLogConfig()      // 日志级别、格式、输出目标
hashConf := config.GetHashIdConfig()  // Optimus ID 哈希参数
natsConf := config.GetNatsConf()      // NATS 消息配置
```

**通过 `APP_ENV` 控制运行模式：**
```go
if config.GetAppConfig().Env == "prod" {
    gin.SetMode(gin.ReleaseMode)
}
```

### provider - 核心服务

线程安全的单例服务提供者。

#### 日志

基于 [Zap](https://github.com/uber-go/zap) 的结构化日志，支持控制台/JSON 格式，文件/标准输出。

```go
import "github.com/JackDPro/cetus/provider"

logger := provider.GetLogger()
logger.Info("服务已启动")
logger.Infow("用户已创建", "userId", 123)
logger.Errorw("请求失败", "error", err)
```

#### 数据库 (ORM)

[GORM](https://gorm.io/) 封装，根据 `DB_TYPE` 自动选择 MySQL/PostgreSQL 驱动。

```go
import "github.com/JackDPro/cetus/provider"

db := provider.GetOrm().Db

// 使用标准 GORM 操作
var users []User
db.Where("active = ?", true).Find(&users)

// 自动迁移
db.AutoMigrate(&User{}, &Post{})
```

#### Redis

```go
import "github.com/JackDPro/cetus/provider"

rdb := provider.GetRedisClient()
rdb.Set(ctx, "key", "value", time.Hour)
val, err := rdb.Get(ctx, "key").Result()
```

#### 密码哈希 (Bcrypt)

```go
import "github.com/JackDPro/cetus/provider"

// 生成哈希
hashed, err := provider.HashMake("my-password")

// 验证密码
err = provider.HashCheck("my-password", hashed)

// 检查是否需要重新哈希
if provider.HashNeedRefresh(hashed) {
    // 使用当前 cost 重新哈希
}
```

#### ID 混淆 (Optimus)

可逆的整数 ID 编码，用于在 API 中隐藏数据库自增 ID。

```go
import "github.com/JackDPro/cetus/provider"

encoded := provider.Hash().Encode(42)         // 例如 1580030173
decoded := provider.Hash().Decode(1580030173) // 42
```

需要在 `.env` 中配置 `OPTIMUS_PRIME`、`OPTIMUS_INVERSE`、`OPTIMUS_RANDOM`。

**获取方式：**

1. 访问 http://primes.utm.edu/lists/small/millions/，下载任意一个 `.txt` 文件，打开后**随机选取一个**小于 `2,147,483,647` 的质数，作为 `OPTIMUS_PRIME`。
2. 使用 [cetus-demo](https://github.com/JackDPro/cetus-demo) 中的工具计算另外两个值：

```bash
go run storage/optimus_gen.go 104393867

# 输出：
# OPTIMUS_PRIME=104393867
# OPTIMUS_INVERSE=1990279033
# OPTIMUS_RANDOM=1333095938
```

> **重要提示：** 一旦部署到生产环境，切勿更改这些值——否则所有已编码的 ID 将失效。

#### RSA 加密

```go
import "github.com/JackDPro/cetus/provider"

// 从文件加载
rsaKey, err := provider.NewRsaByKeyFile("path/to/private.pem")

// 签名
signature, err := rsaKey.Signature(data)

// 验证签名
err = rsaKey.VerifySignature(data, signature)
```

#### 字符串和数组工具

```go
import "github.com/JackDPro/cetus/provider"

provider.RandomString(16)            // "aB3kF9mNpQ2xR7wL"
provider.RandomInt(6)                // "483921"
provider.StringInSlice("a", []string{"a", "b"}) // true
```

#### 文件操作

```go
import "github.com/JackDPro/cetus/provider"

provider.CopyFile("src.txt", "dst.txt")
provider.CopyDir("src_dir", "dst_dir") // 递归复制
```

#### 国际化翻译

```go
import "github.com/JackDPro/cetus/provider"

t := provider.GetTranslate()
t.SetLanguage("zh")
msg := t.Tr("hello_world")

// 或使用全局快捷方式
msg := provider.Tr("hello_world")
```

### model - 数据模型与响应结构

#### BaseModel

在你的模型中嵌入 `BaseModel` 即可获得序列化辅助方法：

```go
import "github.com/JackDPro/cetus/model"

type User struct {
    model.BaseModel
    Id       uint64 `json:"id" gorm:"primaryKey"`
    Nickname string `json:"nickname"`
    Email    string `json:"email"`
}

// 实现 IModel 接口
func (u *User) ToMap() (map[string]interface{}, error) {
    return u.BaseModel.ToMap(u)
}
```

`BaseModel` 提供以下方法：
- `ToMap(model)` - 根据 `json` 标签将结构体转为 map
- `ToJson(model)` - 将结构体转为 JSON 字符串
- `Include(model, relations, db)` - 根据逗号分隔的字符串预加载 GORM 关联

#### IModel 接口

任何实现 `IModel` 的模型都可以配合响应辅助函数使用：

```go
type IModel interface {
    ToMap() (map[string]interface{}, error)
}
```

#### 响应结构体

```go
// API 响应使用 DataWrapper 包装
type DataWrapper struct {
    Data interface{} `json:"data"`
    Meta interface{} `json:"meta,omitempty"`
}

// 错误响应
type Error struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// 分页元数据
type Pagination struct {
    Count       int `json:"count"`
    CurrentPage int `json:"current_page"`
    PerPage     int `json:"per_page"`
    Total       int `json:"total"`
    TotalPages  int `json:"total_pages"`
}
```

### controller - 响应辅助函数

为 Gin 处理器提供标准化的 HTTP 响应函数：

```go
import "github.com/JackDPro/cetus/controller"

func (ctr *UserController) Store(c *gin.Context) {
    // 验证输入
    if err := c.ShouldBindJSON(&req); err != nil {
        controller.ResponseUnprocessable(c, 1, "参数无效", err)
        return
    }

    // 创建资源
    user, err := createUser(req)
    if err != nil {
        controller.ResponseInternalError(c, 2, "创建失败", err)
        return
    }

    // 返回已创建的资源
    controller.ResponseCreated(c, user.Id)
}

func (ctr *UserController) Show(c *gin.Context) {
    user, err := findUser(id)
    if err != nil {
        controller.ResponseNotFound(c, "用户不存在")
        return
    }
    controller.ResponseItem(c, user) // user 必须实现 IModel 接口
}

func (ctr *UserController) Index(c *gin.Context) {
    users, meta := listUsers(page, perPage)
    controller.ResponseCollection(c, users, meta)
}
```

**可用的响应辅助函数：**

| 函数 | HTTP 状态码 | 用途 |
|------|-----------|------|
| `ResponseSuccess(c)` | 200 | 操作成功 |
| `ResponseItem(c, item)` | 200 | 返回单个资源 |
| `ResponseCollection(c, items, meta)` | 200 | 返回列表（可带分页） |
| `ResponseCreated(c, id)` | 201 | 资源已创建（设置 `Location` 头） |
| `ResponseAccepted(c)` | 202 | 异步操作已接受 |
| `ResponseBadRequest(c, code, msg)` | 400 | 无效请求 |
| `ResponseUnauthorized(c)` | 401 | 需要认证 |
| `ResponseForbidden(c)` | 403 | 权限不足 |
| `ResponseNotFound(c, msg)` | 404 | 资源不存在 |
| `ResponseUnprocessable(c, code, msg, err)` | 422 | 验证失败 |
| `ResponseInternalError(c, code, msg, err)` | 500 | 服务器错误（自动记录日志） |

### jwt - JWT 认证

基于 RSA 签名的 JWT 令牌，使用 Redis 存储会话。支持访问令牌和刷新令牌。

**前置要求：**
- RSA 密钥对（PKCS#8 私钥 + PEM 公钥证书）
- Redis 服务
- `.env` 中的 JWT 配置

```env
JWT_EXPIRES_IN=72
JWT_REDIS_PREFIX=auth
JWT_CERT_PATH=storage/jwt.crt
JWT_KEY_PATH=storage/jwt.key
JWT_ISSUE=example.com
```

**生成 RSA 密钥：**

私钥必须为 **PKCS#8 DER** 格式，公钥为 **PEM** 格式：

```bash
# 1. 生成 RSA 私钥（PKCS#1 PEM 格式）
openssl genrsa -out jwt1.pem 2048

# 2. 转换为 PKCS#8 DER 格式（cetus 要求的格式）
openssl pkcs8 -topk8 -inform PEM -outform DER \
  -in jwt1.pem -out storage/jwt.key -nocrypt

# 3. 导出公钥（PEM 格式）
openssl rsa -in jwt1.pem -pubout -out storage/jwt.crt

# 4. 清理临时文件
rm jwt1.pem
```

**使用方式：**
```go
import "github.com/JackDPro/cetus/jwt"

guard, err := jwt.GetJwtGuard()

// 创建访问令牌 + 刷新令牌
accessToken, err := guard.CreateToken(userId, true) // true = 撤销旧令牌
// accessToken.AccessToken  - JWT 字符串
// accessToken.RefreshToken - 刷新令牌 JWT 字符串
// accessToken.Type         - "bearer"
// accessToken.ExpiresIn    - 过期时间（秒）

// 仅创建访问令牌（例如用于 API Key）
accessToken, err := guard.CreateAccessToken(userId)

// 验证令牌
validToken, err := guard.Attempt(tokenString)
// validToken.UserId - 已认证的用户 ID
// validToken.Type   - "token"、"refresh" 或 "access_key"

// 撤销令牌（退出登录）
err = guard.DeleteCredential(tokenString)
```

### middleware - HTTP 中间件

#### 请求 ID

为每个请求添加唯一的请求 ID。优先检查请求头中的 `X-Request-Id`、`HTTP_X_REQUEST_ID`、`HTTP_REQUEST_ID`，如果都没有则自动生成 UUID。

```go
import "github.com/JackDPro/cetus/middleware"

router.Use(middleware.RequestId())
```

#### 限流

基于 [tollbooth](https://github.com/didip/tollbooth) 实现。

```go
import (
    "github.com/JackDPro/cetus/middleware"
    "github.com/didip/tollbooth/v7/limiter"
)

lmt := limiter.New(nil).SetMax(100) // 每秒 100 次请求
router.Use(middleware.LimitRate(lmt))
```

## GORM 模型关联预加载

使用 `BaseModel.Include()` 根据查询参数动态预加载关联：

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

## 完整示例

一个完整的 REST API 配置：

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

    // 跨域配置
    corsConf := cors.DefaultConfig()
    corsConf.AllowAllOrigins = true
    corsConf.AllowHeaders = []string{"Authorization", "Accept-Language"}
    router.Use(cors.New(corsConf))

    // 中间件
    router.Use(middleware.RequestId())

    // 公开路由
    probeCtr := controller.NewProbeController()
    router.GET("/probe", probeCtr.Show)
    router.POST("/auth/login", authLogin)
    router.POST("/users", createUser)

    // 受保护路由（添加你的认证中间件）
    authorized := router.Group("/")
    // authorized.Use(yourAuthMiddleware())
    authorized.GET("/users/me", getMe)
    authorized.POST("/auth/logout", logout)

    // 启动服务
    apiConf := config.GetApiConfig()
    addr := fmt.Sprintf("0.0.0.0:%d", apiConf.HttpPort)
    provider.GetLogger().Info("server starting", "address", addr)
    if err := router.Run(addr); err != nil {
        panic(err)
    }
}
```

## GORM Hook 自动密码哈希

在 GORM 模型 Hook 中使用 bcrypt 自动加密密码：

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

## 环境变量参考

| 变量 | 说明 | 默认值/示例 |
|------|------|------------|
| `APP_NAME` | 应用名称 | `my-app` |
| `APP_ENV` | 运行环境（`dev`/`prod`） | `dev` |
| `APP_DATA_ROOT` | 数据存储根路径 | `/usr/app` |
| `APP_PUBLIC_RES_URL` | 公共静态资源 URL | `http://example.com/static` |
| `LOG_CONSOLE_OUT` | 输出日志到控制台 | `true` |
| `LOG_FILE_OUT` | 输出日志到文件 | `false` |
| `LOG_FILE_PATH` | 日志文件路径 | `/var/log/app.log` |
| `LOG_LEVEL` | 日志级别（`debug`/`info`/`warn`/`error`） | `debug` |
| `LOG_FORMAT` | 日志格式（`json`/`console`） | `json` |
| `DB_TYPE` | 数据库类型（`mysql`/`postgres`） | `mysql` |
| `DB_HOST` | 数据库主机 | `127.0.0.1` |
| `DB_PORT` | 数据库端口 | `3306` |
| `DB_DATABASE` | 数据库名称 | `mydb` |
| `DB_USERNAME` | 数据库用户名 | `root` |
| `DB_PASSWORD` | 数据库密码 | - |
| `DB_SSLMODE` | PostgreSQL SSL 模式 | `disable` |
| `DB_MIGRATE_SELF_ONLY` | 限制迁移范围 | `false` |
| `JWT_EXPIRES_IN` | 令牌过期时间（小时） | `72` |
| `JWT_REDIS_PREFIX` | Redis 令牌前缀 | `auth` |
| `JWT_CERT_PATH` | RSA 公钥路径 | `storage/jwt.crt` |
| `JWT_KEY_PATH` | RSA 私钥路径 | `storage/jwt.key` |
| `JWT_ISSUE` | JWT 签发者 | `example.com` |
| `OPTIMUS_PRIME` | Optimus 质数 | - |
| `OPTIMUS_INVERSE` | Optimus 逆元 | - |
| `OPTIMUS_RANDOM` | Optimus 随机种子 | - |
| `REDIS_HOST` | Redis 主机 | `127.0.0.1` |
| `REDIS_PORT` | Redis 端口 | `6379` |
| `REDIS_DATABASE` | Redis 数据库索引 | `0` |
| `REDIS_PASSWORD` | Redis 密码 | - |
| `NATS_HOST` | NATS 服务主机 | `127.0.0.1` |
| `NATS_USERNAME` | NATS 用户名 | - |
| `NATS_PASSWORD` | NATS 密码 | - |
| `SERVER_HTTP_PORT` | HTTP 服务端口 | `80` |
| `SERVER_GRPC_PORT` | gRPC 服务端口 | `50051` |

## 示例项目

查看 [cetus-demo](https://github.com/JackDPro/cetus-demo)，一个包含用户注册、JWT 认证和 CRUD 操作的完整示例。

## 开源许可

MIT
