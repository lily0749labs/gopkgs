// Package jwtUtil 提供服务端 JWT 创建、解析和刷新工具。
package jwtUtil

import (
	"errors"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/sync/singleflight"
)

var (
	SigningKey  string
	ExpiresTime string
	BufferTime  string
	Issuer      string
)

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

var refreshGroup singleflight.Group

type JWT struct {
	SigningKey []byte
	now        func() time.Time
}

type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	ID       uint64
	Account  string
	Password string
	Roles    string
	Enable   bool
	IP       string
}

func NewJWT() *JWT {
	return &JWT{SigningKey: []byte(SigningKey)}
}

func (j *JWT) currentTime() time.Time {
	if j.now != nil {
		return j.now()
	}
	return time.Now()
}

func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	bufferDuration, _ := ParseDuration(BufferTime)
	expiresDuration, _ := ParseDuration(ExpiresTime)
	now := j.currentTime()

	return CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bufferDuration / time.Second),
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"GVA"},
			NotBefore: jwt.NewNumericDate(now.Add(-1000)),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresDuration)),
			Issuer:    Issuer,
		},
	}
}

// CreateToken 创建一个使用 HS256 签名的 Token。
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken 使用 singleflight 合并同一个旧 Token 的并发刷新。
func (j *JWT) CreateTokenByOldToken(oldToken string, claims CustomClaims) (string, error) {
	value, err, _ := refreshGroup.Do("JWT:"+oldToken, func() (any, error) {
		return j.CreateToken(claims)
	})
	if err != nil {
		return "", err
	}
	return value.(string), nil
}

// ParseToken 解析并验证 Token。
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(_ *jwt.Token) (any, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		default:
			return nil, err
		}
	}

	if token == nil {
		return nil, TokenInvalid
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, TokenInvalid
	}
	return claims, nil
}

// ParseDuration 在标准 time.Duration 基础上支持 7d、1d12h 等天数写法。
func ParseDuration(value string) (time.Duration, error) {
	value = strings.TrimSpace(value)
	duration, err := time.ParseDuration(value)
	if err == nil {
		return duration, nil
	}
	if strings.Contains(value, "d") {
		index := strings.Index(value, "d")
		days, _ := strconv.Atoi(value[:index])
		duration = 24 * time.Hour * time.Duration(days)
		remainder, remainderErr := time.ParseDuration(value[index+1:])
		if remainderErr != nil {
			return duration, nil
		}
		return duration + remainder, nil
	}

	nanoseconds, err := strconv.ParseInt(value, 10, 64)
	return time.Duration(nanoseconds), err
}
