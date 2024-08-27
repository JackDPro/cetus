package jwt

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JackDPro/cetus/config"
	"github.com/JackDPro/cetus/provider"
	"github.com/go-kit/log/level"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"os"
	"sync"
	"time"
)

type Guard struct {
	KeyPath    string
	PubPath    string
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	redis      *redis.Client
}

var guardInstance *Guard
var guardOnce sync.Once

func GetJwtGuard() *Guard {
	guardOnce.Do(func() {
		var conf = config.GetAuthConf()

		keyBuffer, err := os.ReadFile(conf.KeyPath)
		if err != nil {
			_ = level.Error(provider.GetLogger()).Log("id", "jwt", "method", "GetJwtGuard", "message", "load private key content failed", "error", err)
			return
		}
		privateKey, err := x509.ParsePKCS8PrivateKey(keyBuffer)
		if err != nil {
			_ = level.Error(provider.GetLogger()).Log("id", "jwt", "method", "GetJwtGuard", "message", "parse private key failed", "error", err)
			return
		}
		publicKeyByte, err := os.ReadFile(conf.CertPath)
		if err != nil {
			_ = level.Error(provider.GetLogger()).Log("id", "jwt", "method", "GetJwtGuard", "message", "load public key failed", "error", err)
			return
		}
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
		rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
		if !ok {
			_ = level.Error(provider.GetLogger()).Log("id", "jwt", "method", "GetJwtGuard", "message", "parse public key failed", "error", err)
			return
		}
		guardInstance = &Guard{
			PubPath:    conf.CertPath,
			KeyPath:    conf.KeyPath,
			privateKey: rsaPrivateKey,
			publicKey:  publicKey,
			redis:      provider.GetRedisClient(),
		}
	})
	return guardInstance
}

func (guard *Guard) CreateAccessToken(userId uint64) (*AccessToken, error) {
	ctx := context.Background()
	var conf = config.GetAuthConf()
	now := time.Now()
	// 创建持久化 jti
	tokenJti := &ValidToken{
		Id:        provider.RandomString(12),
		UserId:    userId,
		Token:     provider.RandomString(32),
		Type:      "access_key",
		ExpiredAt: now.Add(time.Duration(conf.ExpiresIn) * time.Hour),
		Now:       now,
		Audience:  conf.Audience,
	}
	tokenKey := fmt.Sprintf("%s:%d:%s", conf.RedisPrefix, userId, tokenJti.Id)
	jsonStr, err := json.Marshal(tokenJti.ToJsonMap())
	if err != nil {
		return nil, err
	}
	_, err = guard.redis.SetNX(ctx, tokenKey, jsonStr, time.Duration(conf.ExpiresIn)*time.Hour).Result()
	if err != nil {
		return nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": conf.Issue,                //issuer 谁创建的颁发的 token
		"aud": tokenJti.Audience,         //audience 颁发给谁的
		"jti": tokenJti.Token,            //JWT ID token 唯一标识
		"iat": tokenJti.Now.Unix(),       //issued at 颁发时间
		"exp": tokenJti.ExpiredAt.Unix(), //expiration time 过期时间
		"uid": userId,                    // 登录人 id
		"nbf": tokenJti.Now.Unix(),       //Not valid before 在什么时间之前不可用
		"typ": tokenJti.Type,             // 类型 token / refresh
		"tid": tokenJti.Id,               // 该数据 id
	})
	accessTokenStr, err := token.SignedString(guard.privateKey)
	if err != nil {
		return nil, err
	}
	accessToken := &AccessToken{
		AccessToken: accessTokenStr,
		Type:        "bearer",
		ExpiresIn:   int64(time.Duration(conf.ExpiresIn) * time.Hour),
	}
	return accessToken, nil
}

func (guard *Guard) CreateToken(userId uint64, clean bool) (*AccessToken, error) {
	ctx := context.Background()
	var conf = config.GetAuthConf()
	now := time.Now()
	// 创建持久化 jti
	tokenJti := &ValidToken{
		Id:        provider.RandomString(12),
		UserId:    userId,
		Token:     provider.RandomString(32),
		Type:      "token",
		ExpiredAt: now.Add(time.Duration(conf.ExpiresIn) * time.Hour),
		Now:       now,
		Audience:  conf.Audience,
	}
	tokenKey := fmt.Sprintf("%s:%d:%s", conf.RedisPrefix, userId, tokenJti.Id)
	jsonStr, err := json.Marshal(tokenJti.ToJsonMap())
	if err != nil {
		return nil, err
	}
	_, err = guard.redis.SetNX(ctx, tokenKey, jsonStr, time.Duration(conf.ExpiresIn)*time.Hour).Result()
	if err != nil {
		return nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": conf.Issue,                //issuer 谁创建的颁发的 token
		"aud": tokenJti.Audience,         //audience 颁发给谁的
		"jti": tokenJti.Token,            //JWT ID token 唯一标识
		"iat": tokenJti.Now.Unix(),       //issued at 颁发时间
		"exp": tokenJti.ExpiredAt.Unix(), //expiration time 过期时间
		"uid": userId,                    // 登录人 id
		"nbf": tokenJti.Now.Unix(),       //Not valid before 在什么时间之前不可用
		"typ": tokenJti.Type,             // 类型 token / refresh
		"tid": tokenJti.Id,               // 该数据 id
	})
	accessTokenStr, err := token.SignedString(guard.privateKey)
	if err != nil {
		return nil, err
	}

	// 刷新令牌
	refreshJti := &ValidToken{
		Id:        provider.RandomString(12),
		UserId:    userId,
		Token:     provider.RandomString(32),
		Type:      "refresh",
		ExpiredAt: now.Add(time.Duration(conf.ExpiresIn) * 3 * time.Hour),
		Now:       now,
		Audience:  conf.Audience,
	}
	refreshKey := fmt.Sprintf("%s:%d:%s", conf.RedisPrefix, userId, refreshJti.Id)
	refresh := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": conf.Issue,                  //issuer 谁创建的颁发的 token
		"aud": refreshJti.Audience,         //audience 颁发给谁的
		"jti": refreshJti.Token,            //JWT ID token 唯一标识
		"iat": refreshJti.Now.Unix(),       //issued at 颁发时间
		"exp": refreshJti.ExpiredAt.Unix(), //expiration time 过期时间
		"uid": userId,                      // 登录人 id
		"nbf": now.Unix(),                  //Not valid before 在什么时间之前不可用
		"typ": refreshJti.Type,
		"tid": refreshJti.Id, // 该数据 id
	})
	refreshJsonStr, err := json.Marshal(refreshJti.ToJsonMap())
	if err != nil {
		return nil, err
	}
	_, err = guard.redis.SetNX(ctx, refreshKey, refreshJsonStr, time.Duration(conf.ExpiresIn)*3*time.Hour).Result()
	if err != nil {
		return nil, err
	}
	refreshTokenStr, err := refresh.SignedString(guard.privateKey)
	if err != nil {
		return nil, err
	}
	accessToken := &AccessToken{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		Type:         "bearer",
		ExpiresIn:    int64(time.Duration(conf.ExpiresIn) * time.Hour),
	}

	// 清空 token
	if clean {
		iter := guard.redis.Scan(ctx, 0, fmt.Sprintf("%s:%d:*", conf.RedisPrefix, userId), 0).Iterator()
		for iter.Next(ctx) {
			if iter.Val() != tokenKey && iter.Val() != refreshKey {
				err := guard.redis.Del(ctx, iter.Val()).Err()
				if err != nil {
					return nil, err
				}
			}
		}
		if err := iter.Err(); err != nil {
			return nil, err
		}
	}
	return accessToken, nil
}

func (guard *Guard) DeleteCredential(credential string) error {
	conf := config.GetAuthConf()
	token, err := guard.Attempt(credential)
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, err = guard.redis.Del(ctx, fmt.Sprintf("%s:%d:%s", conf.RedisPrefix, token.UserId, token.Id)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (guard *Guard) Attempt(credentials string) (*ValidToken, error) {
	ctx := context.Background()
	// 基于公钥验证Token合法性
	token, err := jwt.Parse(credentials, func(token *jwt.Token) (interface{}, error) {
		// 基于JWT的第一部分中的alg字段值进行一次验证
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("验证 Token 的加密类型错误")
		}
		return guard.publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	conf := config.GetAuthConf()
	// 签名验证
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 验证 jti
		jti := claims["jti"]
		uid, ok := claims["uid"].(float64)
		tokenId := claims["tid"]
		tokenType := claims["typ"]
		if !ok {
			return nil, errors.New("uid is not a number")
		}
		userId := uint64(uid)
		validTokenStr, err := guard.redis.Get(ctx, fmt.Sprintf("%s:%d:%s", conf.RedisPrefix, userId, tokenId)).Result()
		if err != nil {
			return nil, err
		}
		validToken := &ValidToken{}
		err = json.Unmarshal([]byte(validTokenStr), &validToken)
		if err != nil {
			return nil, err
		}
		if validToken.Token != jti || validToken.Type != tokenType {
			return nil, errors.New("invalid token")
		}

		return validToken, nil
	}
	return nil, errors.New("token 签名不合法")
}
