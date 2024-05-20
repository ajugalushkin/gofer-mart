package cookies

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ajugalushkin/gofer-mart/config"
)

type Claims struct {
	jwt.RegisteredClaims
	UserData string
}

const TokenExp = time.Hour * 3

func buildJWTString(ctx context.Context, data string) (string, error) {
	cfg := config.FlagsFromContext(ctx)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserData: data,
	})

	tokenString, err := token.SignedString([]byte(cfg.TokenKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//func GetUserID(ctx context.Context, tokenString string) int {
//	flags := config.FlagsFromContext(ctx)
//	claims := &Claims{}
//	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
//		return []byte(flags.SecretKey), nil
//	})
//	if err != nil {
//		return 0
//	}
//
//	return claims.UserID
//}

func Create(ctx context.Context, nameCookie string, data string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = nameCookie
	cookie.Value, _ = buildJWTString(ctx, data)
	cookie.Expires = time.Now().Add(TokenExp)
	return cookie
}

/*func Read(echoCtx echo.Context, name string) (string, error) {
	cookie, err := echoCtx.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}*/
