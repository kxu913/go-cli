package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

)

const (
	SECRETKEY = "243223ffslsfsldfl412fdsfsdf"
	MAXAGE    = 60 * 60 * 24
	ISSUER    = "{{.Prefix}}"
)

type EmailClaims struct {
	Email    string
	Password string
	jwt.StandardClaims
}

func ParseToken(tokenString string) (*EmailClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &EmailClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(SECRETKEY), nil
	})
	if claims, ok := token.Claims.(*EmailClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}


func CreateToken(email string, password string) (string, int64, error) {
	expiredAt := time.Now().Add(time.Duration(MAXAGE) * time.Second).Unix()
	emailClaims := EmailClaims{
		Email:    email,
		Password: password,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(MAXAGE) * time.Second).Unix(),
			Issuer:    ISSUER,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, emailClaims)
	tokenString, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		return "", 0, fmt.Errorf("invalid token")
	}
	return tokenString, expiredAt, nil
}

func JwtAuth(fn CustomAuthValidator) echo.MiddlewareFunc {
	c := DefaultCustomAuthConfig
	c.Validator = fn
	return JwtAuthWithConfig(c)
}

func JwtAuthWithConfig(config CustomAuthConfig) echo.MiddlewareFunc {

	if config.Realm == "" {
		config.Realm = defaultRealm
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			l := len(bearer)

			if len(auth) < l {
				return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "authorization failed"}
			}
			if strings.ToLower(auth[:l]) == bearer {
				claims, err := ParseToken(auth[l+1:])
				if err != nil {
					return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "authorization failed, check the JWT token."}
				}
				isOK, err := config.Validator(claims.Email, claims.Password, c)
				if err != nil || !isOK {
					return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "authorization failed, use correct email and password to create token."}
				}
				c.Request().Header.Add("email", claims.Email)

			}

			return next(c)
		}
	}
}

