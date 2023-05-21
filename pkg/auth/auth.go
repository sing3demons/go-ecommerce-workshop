package auth

import (
	"fmt"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sing3demons/shop/config"
	"github.com/sing3demons/shop/modules/users"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "api_key"
)

type IAuth interface{
	SignToken() string
}

type auth struct {
	mapClaims *MapClaims
	cfg       config.IJwtConfig
}

type MapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

func NewAuth(tt TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IAuth, error) {
	switch tt {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefresh(cfg, claims), nil
	default:
		return nil, fmt.Errorf("invalid token type: %s", tt)
	}

}

func (a *auth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString([]byte(a.cfg.SecretKey()))
	return ss
}

func jwtTimeDuration(expire int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(expire) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IAuth {
	return &auth{
		mapClaims: &MapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   claims.Id,
				ExpiresAt: jwtTimeDuration(cfg.AccessExpire()),
				Issuer:    "sing3demons",
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
		cfg: cfg,
	}
}
func newRefresh(cfg config.IJwtConfig, claims *users.UserClaims) IAuth {
	return &auth{
		mapClaims: &MapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   claims.Id,
				ExpiresAt: jwtTimeDuration(cfg.RefreshExpire()),
				Issuer:    "sing3demons",
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
		cfg: cfg,
	}
}
