package cookies

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ajugalushkin/gofer-mart/config"
	"github.com/ajugalushkin/gofer-mart/internal/auth"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/logger"
)

type Claims struct {
	jwt.RegisteredClaims
	Login dto.User
}

const TokenExp = time.Hour * 3

func buildJWTString(ctx context.Context, login dto.User) string {
	cfg := config.FlagsFromContext(ctx)

	var err error
	login.Password, err = auth.HashPassword(login.Password)
	if err != nil {
		logger.LogFromContext(ctx).Debug("cookies.buildJWTString: unable to hash password")
		return ""
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		Login: login,
	})

	tokenString, err := token.SignedString([]byte(cfg.TokenKey))
	if err != nil {
		logger.LogFromContext(ctx).Debug("cookies.buildJWTString: unable to sign token")
		return ""
	}

	return tokenString
}

func GetLogin(ctx context.Context, tokenString string) *dto.User {
	cfg := config.FlagsFromContext(ctx)
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.TokenKey), nil
	})
	if err != nil {
		logger.LogFromContext(ctx).Debug("cookies.GetLogin: unable to parse token")
		return &dto.User{}
	}

	return &claims.Login
}

func Create(ctx context.Context, nameCookie string, login dto.User) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = nameCookie
	cookie.Value = buildJWTString(ctx, login)
	cookie.Expires = time.Now().Add(TokenExp)
	return cookie
}
