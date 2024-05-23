package cookies

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ajugalushkin/gofer-mart/config"
	"github.com/ajugalushkin/gofer-mart/internal/auth"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
)

type Claims struct {
	jwt.RegisteredClaims
	Login dto.Login
}

const TokenExp = time.Hour * 3

func buildJWTString(ctx context.Context, login dto.Login) (string, error) {
	cfg := config.FlagsFromContext(ctx)

	var err error
	login.Password, err = auth.HashPassword(login.Password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		Login: login,
	})

	tokenString, err := token.SignedString([]byte(cfg.TokenKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetLogin(ctx context.Context, tokenString string) *dto.Login {
	cfg := config.FlagsFromContext(ctx)
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.TokenKey), nil
	})
	if err != nil {
		return &dto.Login{}
	}

	return &claims.Login
}

func Create(ctx context.Context, nameCookie string, login dto.Login) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = nameCookie
	cookie.Value, _ = buildJWTString(ctx, login)
	cookie.Expires = time.Now().Add(TokenExp)
	return cookie
}

//func ReadCookieData(ctx context.Context, echoCtx echo.Context, name string) (*dto.Login, error) {
//	cookie, err := echoCtx.Cookie(name)
//	if err != nil {
//		return &dto.Login{}, err
//	}
//	return getLogin(ctx, cookie.Value), nil
//}
