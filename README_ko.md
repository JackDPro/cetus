# Cetus

[English](README.md) | [中文](README_zh.md) | [日本語](README_ja.md) | **한국어**

[Gin](https://github.com/gin-gonic/gin) 기반의 Go HTTP 서비스 개발 툴킷. 설정 관리, 데이터베이스 액세스(MySQL/PostgreSQL), JWT 인증, 구조화된 로깅, 미들웨어, 표준화된 API 응답을 즉시 사용할 수 있습니다.

## 기능

- **설정 관리** - `.env` 파일 기반 환경 변수 설정, 싱글턴 패턴
- **데이터베이스** - GORM을 통한 MySQL 및 PostgreSQL 지원, 드라이버 자동 선택
- **JWT 인증** - RSA 서명 토큰, Redis 기반 세션 관리
- **로깅** - Zap 기반 구조화된 로깅 (JSON/콘솔 형식, 파일/표준출력)
- **미들웨어** - 요청 ID 추적, 속도 제한
- **API 응답** - Gin용 표준화된 JSON 응답 헬퍼
- **유틸리티** - Bcrypt 해싱, ID 난독화(Optimus), RSA 암호화, i18n, 파일 작업

## 설치

```bash
go get github.com/JackDPro/cetus
```

## 빠른 시작

### 1. `.env` 파일 생성

예제 파일을 복사하고 환경에 맞게 설정을 수정합니다:

```bash
cp .env.example .env
```

주요 환경 변수:

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

전체 옵션은 [.env.example](.env.example)을 참조하세요.

### 2. 서버 생성

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

    // CORS 설정
    corsConf := cors.DefaultConfig()
    corsConf.AllowAllOrigins = true
    corsConf.AllowHeaders = []string{"Authorization", "Accept-Language"}
    router.Use(cors.New(corsConf))

    // 요청 ID 미들웨어
    router.Use(middleware.RequestId())

    // 헬스 체크 엔드포인트
    probeCtr := controller.NewProbeController()
    router.GET("/probe", probeCtr.Show)

    // 라우트 추가
    // router.GET("/users", userController.Index)

    conf := config.GetApiConfig()
    addr := fmt.Sprintf("0.0.0.0:%d", conf.HttpPort)
    provider.GetLogger().Info("server starting", "address", addr)
    if err := router.Run(addr); err != nil {
        panic(err)
    }
}
```

## 패키지 가이드

### config - 설정 관리

모든 설정은 환경 변수(`.env` 파일)에서 로드됩니다. 각 설정 구조체는 싱글턴 패턴을 사용합니다.

```go
import "github.com/JackDPro/cetus/config"

appConf := config.GetAppConfig()      // 앱 이름, 환경, 데이터 루트, 공개 URL
dbConf  := config.GetDatabaseConfig() // DB 타입, 호스트, 포트, 인증 정보
apiConf := config.GetApiConfig()      // HTTP 및 gRPC 포트
authConf := config.GetAuthConf()      // JWT 인증서/키 경로, 만료 시간
redisConf := config.GetRedisConfig()  // Redis 연결 정보
logConf := config.GetLogConfig()      // 로그 레벨, 형식, 출력 대상
hashConf := config.GetHashIdConfig()  // Optimus ID 해싱 파라미터
natsConf := config.GetNatsConf()      // NATS 메시징 설정
```

**`APP_ENV`로 동작 제어:**
```go
if config.GetAppConfig().Env == "prod" {
    gin.SetMode(gin.ReleaseMode)
}
```

### provider - 코어 서비스

스레드 안전한 싱글턴 프로바이더.

#### 로거

[Zap](https://github.com/uber-go/zap) 기반 구조화된 로깅. 콘솔/JSON 형식, 파일/표준출력 지원.

```go
import "github.com/JackDPro/cetus/provider"

logger := provider.GetLogger()
logger.Info("서버 시작")
logger.Infow("사용자 생성", "userId", 123)
logger.Errorw("요청 실패", "error", err)
```

#### 데이터베이스 (ORM)

[GORM](https://gorm.io/) 래퍼. `DB_TYPE`에 따라 MySQL/PostgreSQL 드라이버 자동 선택.

```go
import "github.com/JackDPro/cetus/provider"

db := provider.GetOrm().Db

// 표준 GORM 작업
var users []User
db.Where("active = ?", true).Find(&users)

// 자동 마이그레이션
db.AutoMigrate(&User{}, &Post{})
```

#### Redis

```go
import "github.com/JackDPro/cetus/provider"

rdb := provider.GetRedisClient()
rdb.Set(ctx, "key", "value", time.Hour)
val, err := rdb.Get(ctx, "key").Result()
```

#### 비밀번호 해싱 (Bcrypt)

```go
import "github.com/JackDPro/cetus/provider"

// 해시 생성
hashed, err := provider.HashMake("my-password")

// 검증
err = provider.HashCheck("my-password", hashed)

// 리해시 필요 여부 확인
if provider.HashNeedRefresh(hashed) {
    // 현재 cost로 리해시
}
```

#### ID 난독화 (Optimus)

가역적 정수 ID 인코딩. API에서 데이터베이스 자동 증가 ID를 숨깁니다.

```go
import "github.com/JackDPro/cetus/provider"

encoded := provider.Hash().Encode(42)         // 예: 1580030173
decoded := provider.Hash().Decode(1580030173) // 42
```

`.env`에 `OPTIMUS_PRIME`, `OPTIMUS_INVERSE`, `OPTIMUS_RANDOM` 설정이 필요합니다.

#### RSA 암호화

```go
import "github.com/JackDPro/cetus/provider"

// 파일에서 로드
rsaKey, err := provider.NewRsaByKeyFile("path/to/private.pem")

// 서명
signature, err := rsaKey.Signature(data)

// 서명 검증
err = rsaKey.VerifySignature(data, signature)
```

#### 문자열 및 배열 유틸리티

```go
import "github.com/JackDPro/cetus/provider"

provider.RandomString(16)            // "aB3kF9mNpQ2xR7wL"
provider.RandomInt(6)                // "483921"
provider.StringInSlice("a", []string{"a", "b"}) // true
```

#### 파일 작업

```go
import "github.com/JackDPro/cetus/provider"

provider.CopyFile("src.txt", "dst.txt")
provider.CopyDir("src_dir", "dst_dir") // 재귀 복사
```

#### 국제화 (i18n)

```go
import "github.com/JackDPro/cetus/provider"

t := provider.GetTranslate()
t.SetLanguage("ko")
msg := t.Tr("hello_world")

// 글로벌 단축 함수
msg := provider.Tr("hello_world")
```

### model - 데이터 모델 및 응답 구조체

#### BaseModel

모델에 `BaseModel`을 임베드하여 직렬화 헬퍼를 사용할 수 있습니다:

```go
import "github.com/JackDPro/cetus/model"

type User struct {
    model.BaseModel
    Id       uint64 `json:"id" gorm:"primaryKey"`
    Nickname string `json:"nickname"`
    Email    string `json:"email"`
}

// IModel 인터페이스 구현
func (u *User) ToMap() (map[string]interface{}, error) {
    return u.BaseModel.ToMap(u)
}
```

`BaseModel` 제공 메서드:
- `ToMap(model)` - `json` 태그를 사용하여 구조체를 map으로 변환
- `ToJson(model)` - 구조체를 JSON 문자열로 변환
- `Include(model, relations, db)` - 쉼표로 구분된 문자열에서 GORM 관계 프리로드

#### IModel 인터페이스

`IModel`을 구현한 모델은 응답 헬퍼와 함께 사용할 수 있습니다:

```go
type IModel interface {
    ToMap() (map[string]interface{}, error)
}
```

#### 응답 구조체

```go
// API 응답은 DataWrapper로 래핑됩니다
type DataWrapper struct {
    Data interface{} `json:"data"`
    Meta interface{} `json:"meta,omitempty"`
}

// 에러 응답
type Error struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// 페이지네이션 메타데이터
type Pagination struct {
    Count       int `json:"count"`
    CurrentPage int `json:"current_page"`
    PerPage     int `json:"per_page"`
    Total       int `json:"total"`
    TotalPages  int `json:"total_pages"`
}
```

### controller - 응답 헬퍼

Gin 핸들러용 표준화된 HTTP 응답 함수:

```go
import "github.com/JackDPro/cetus/controller"

func (ctr *UserController) Store(c *gin.Context) {
    // 입력 검증
    if err := c.ShouldBindJSON(&req); err != nil {
        controller.ResponseUnprocessable(c, 1, "파라미터가 유효하지 않습니다", err)
        return
    }

    // 리소스 생성
    user, err := createUser(req)
    if err != nil {
        controller.ResponseInternalError(c, 2, "생성 실패", err)
        return
    }

    // 생성된 리소스 반환
    controller.ResponseCreated(c, user.Id)
}

func (ctr *UserController) Show(c *gin.Context) {
    user, err := findUser(id)
    if err != nil {
        controller.ResponseNotFound(c, "사용자를 찾을 수 없습니다")
        return
    }
    controller.ResponseItem(c, user) // user는 IModel을 구현해야 합니다
}

func (ctr *UserController) Index(c *gin.Context) {
    users, meta := listUsers(page, perPage)
    controller.ResponseCollection(c, users, meta)
}
```

**사용 가능한 응답 헬퍼:**

| 함수 | HTTP 상태 코드 | 용도 |
|------|---------------|------|
| `ResponseSuccess(c)` | 200 | 작업 성공 |
| `ResponseItem(c, item)` | 200 | 단일 리소스 반환 |
| `ResponseCollection(c, items, meta)` | 200 | 목록 반환 (페이지네이션 포함) |
| `ResponseCreated(c, id)` | 201 | 리소스 생성됨 (`Location` 헤더 설정) |
| `ResponseAccepted(c)` | 202 | 비동기 작업 수락 |
| `ResponseBadRequest(c, code, msg)` | 400 | 잘못된 요청 |
| `ResponseUnauthorized(c)` | 401 | 인증 필요 |
| `ResponseForbidden(c)` | 403 | 접근 권한 없음 |
| `ResponseNotFound(c, msg)` | 404 | 리소스를 찾을 수 없음 |
| `ResponseUnprocessable(c, code, msg, err)` | 422 | 유효성 검사 실패 |
| `ResponseInternalError(c, code, msg, err)` | 500 | 서버 오류 (자동 로그 기록) |

### jwt - JWT 인증

RSA 서명 JWT 토큰, Redis 기반 세션 저장소. 액세스 토큰과 리프레시 토큰을 지원합니다.

**사전 요구 사항:**
- RSA 키 쌍 (PKCS#8 개인 키 + PEM 공개 키 인증서)
- Redis 서버
- `.env`의 JWT 설정

```env
JWT_EXPIRES_IN=72
JWT_REDIS_PREFIX=auth
JWT_CERT_PATH=storage/jwt.crt
JWT_KEY_PATH=storage/jwt.key
JWT_ISSUE=example.com
```

**RSA 키 생성:**
```bash
openssl genpkey -algorithm RSA -out storage/jwt.key -pkeyopt rsa_keygen_bits:2048
openssl rsa -in storage/jwt.key -pubout -out storage/jwt.crt
```

**사용법:**
```go
import "github.com/JackDPro/cetus/jwt"

guard, err := jwt.GetJwtGuard()

// 액세스 토큰 + 리프레시 토큰 생성
accessToken, err := guard.CreateToken(userId, true) // true = 이전 토큰 폐기
// accessToken.AccessToken  - JWT 문자열
// accessToken.RefreshToken - 리프레시 토큰 JWT 문자열
// accessToken.Type         - "bearer"
// accessToken.ExpiresIn    - 만료 시간(초)

// 액세스 토큰만 생성 (예: API 키용)
accessToken, err := guard.CreateAccessToken(userId)

// 토큰 검증
validToken, err := guard.Attempt(tokenString)
// validToken.UserId - 인증된 사용자 ID
// validToken.Type   - "token", "refresh", 또는 "access_key"

// 토큰 폐기 (로그아웃)
err = guard.DeleteCredential(tokenString)
```

### middleware - HTTP 미들웨어

#### 요청 ID

모든 요청에 고유한 요청 ID를 부여합니다. 먼저 요청 헤더(`X-Request-Id`, `HTTP_X_REQUEST_ID`, `HTTP_REQUEST_ID`)를 확인하고, 없으면 UUID를 자동 생성합니다.

```go
import "github.com/JackDPro/cetus/middleware"

router.Use(middleware.RequestId())
```

#### 속도 제한

[tollbooth](https://github.com/didip/tollbooth) 기반.

```go
import (
    "github.com/JackDPro/cetus/middleware"
    "github.com/didip/tollbooth/v7/limiter"
)

lmt := limiter.New(nil).SetMax(100) // 초당 100 요청
router.Use(middleware.LimitRate(lmt))
```

## GORM 모델 관계 프리로딩

`BaseModel.Include()`로 쿼리 파라미터에서 관계를 동적으로 프리로드:

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

## 전체 예제

완전한 REST API 설정:

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

    // CORS 설정
    corsConf := cors.DefaultConfig()
    corsConf.AllowAllOrigins = true
    corsConf.AllowHeaders = []string{"Authorization", "Accept-Language"}
    router.Use(cors.New(corsConf))

    // 미들웨어
    router.Use(middleware.RequestId())

    // 공개 라우트
    probeCtr := controller.NewProbeController()
    router.GET("/probe", probeCtr.Show)
    router.POST("/auth/login", authLogin)
    router.POST("/users", createUser)

    // 보호된 라우트 (인증 미들웨어 추가)
    authorized := router.Group("/")
    // authorized.Use(yourAuthMiddleware())
    authorized.GET("/users/me", getMe)
    authorized.POST("/auth/logout", logout)

    // 서버 시작
    apiConf := config.GetApiConfig()
    addr := fmt.Sprintf("0.0.0.0:%d", apiConf.HttpPort)
    provider.GetLogger().Info("server starting", "address", addr)
    if err := router.Run(addr); err != nil {
        panic(err)
    }
}
```

## GORM Hook을 이용한 자동 비밀번호 해싱

GORM 모델 Hook에서 bcrypt를 이용한 자동 비밀번호 암호화:

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

## 환경 변수 참조

| 변수 | 설명 | 기본값/예시 |
|------|------|-----------|
| `APP_NAME` | 애플리케이션 이름 | `my-app` |
| `APP_ENV` | 실행 환경 (`dev`/`prod`) | `dev` |
| `APP_DATA_ROOT` | 데이터 저장 루트 경로 | `/usr/app` |
| `APP_PUBLIC_RES_URL` | 공개 정적 리소스 URL | `http://example.com/static` |
| `LOG_CONSOLE_OUT` | 콘솔에 로그 출력 | `true` |
| `LOG_FILE_OUT` | 파일에 로그 출력 | `false` |
| `LOG_FILE_PATH` | 로그 파일 경로 | `/var/log/app.log` |
| `LOG_LEVEL` | 로그 레벨 (`debug`/`info`/`warn`/`error`) | `debug` |
| `LOG_FORMAT` | 로그 형식 (`json`/`console`) | `json` |
| `DB_TYPE` | 데이터베이스 타입 (`mysql`/`postgres`) | `mysql` |
| `DB_HOST` | 데이터베이스 호스트 | `127.0.0.1` |
| `DB_PORT` | 데이터베이스 포트 | `3306` |
| `DB_DATABASE` | 데이터베이스 이름 | `mydb` |
| `DB_USERNAME` | 데이터베이스 사용자 이름 | `root` |
| `DB_PASSWORD` | 데이터베이스 비밀번호 | - |
| `DB_SSLMODE` | PostgreSQL SSL 모드 | `disable` |
| `DB_MIGRATE_SELF_ONLY` | 마이그레이션 범위 제한 | `false` |
| `JWT_EXPIRES_IN` | 토큰 만료 시간(시간) | `72` |
| `JWT_REDIS_PREFIX` | Redis 토큰 접두사 | `auth` |
| `JWT_CERT_PATH` | RSA 공개 키 경로 | `storage/jwt.crt` |
| `JWT_KEY_PATH` | RSA 개인 키 경로 | `storage/jwt.key` |
| `JWT_ISSUE` | JWT 발급자 | `example.com` |
| `OPTIMUS_PRIME` | Optimus 소수 | - |
| `OPTIMUS_INVERSE` | Optimus 역수 | - |
| `OPTIMUS_RANDOM` | Optimus 랜덤 시드 | - |
| `REDIS_HOST` | Redis 호스트 | `127.0.0.1` |
| `REDIS_PORT` | Redis 포트 | `6379` |
| `REDIS_DATABASE` | Redis 데이터베이스 인덱스 | `0` |
| `REDIS_PASSWORD` | Redis 비밀번호 | - |
| `NATS_HOST` | NATS 서버 호스트 | `127.0.0.1` |
| `NATS_USERNAME` | NATS 사용자 이름 | - |
| `NATS_PASSWORD` | NATS 비밀번호 | - |
| `SERVER_HTTP_PORT` | HTTP 서버 포트 | `80` |
| `SERVER_GRPC_PORT` | gRPC 서버 포트 | `50051` |

## 데모 프로젝트

사용자 등록, JWT 인증, CRUD 작업을 포함한 전체 예제는 [cetus-demo](https://github.com/JackDPro/cetus-demo)를 참조하세요.

## 라이선스

MIT
