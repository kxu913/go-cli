package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

)

type (
	Skipper func(echo.Context) bool

	CustomAuthConfig struct {
		// // Skipper defines a function to skip middleware.
		Skipper Skipper

		// Validator is a function to validate BasicAuth credentials.
		// Do Nothing.
		Validator CustomAuthValidator

		// Default value "Restricted".
		Realm string
	}

	// BasicAuthValidator defines a function to validate BasicAuth credentials.
	CustomAuthValidator func(string, string, echo.Context) (bool, error)
)

// DefaultBasicAuthConfig is the default BasicAuth middleware config.
var DefaultCustomAuthConfig = CustomAuthConfig{
	Skipper: func(ctx echo.Context) bool {
		fmt.Println(ctx.Path())

		if ctx.Path() == "/api/user" && ctx.Request().Method == "POST" {
			return true

		}
		return false
	},
	Realm: defaultRealm,
}

const (
	basic        = "basic"
	defaultRealm = "Restricted"
	bearer       = "bearer"
)

func CustomAuth() echo.MiddlewareFunc {
	c := DefaultCustomAuthConfig
	// c.Validator = fn

	return CustomAuthWithConfig(c)
}

func CustomAuthWithSkipper(fn CustomAuthValidator, skipper Skipper) echo.MiddlewareFunc {
	c := DefaultCustomAuthConfig
	c.Validator = fn
	c.Skipper = skipper
	return CustomAuthWithConfig(c)
}

func CustomAuthWithConfig(config CustomAuthConfig) echo.MiddlewareFunc {

	if config.Realm == "" {
		config.Realm = defaultRealm
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			l := len(basic)

			if len(auth) < l {
				return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "authorization failed"}
			}
			if strings.ToLower(auth[:l]) == basic {

				b, err := base64.StdEncoding.DecodeString(auth[l+1:])
				if err != nil {
					return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "authorization failed"}
				}
				cred := string(b)
				for i := 0; i < len(cred); i++ {
					if cred[i] == ':' {
						email := cred[:i]
						password := cred[i+1:]
						success := basicAuth(email, password)
						if !success {
							return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "authorization failed, check your email and password."}
						}
						// you can set some fileds into header
						//SetAccountFromHeader(c.Request().Header, acct)
						break
					}
				}

			}

			return next(c)
		}
	}
}

func basicAuth(email string, password string) bool {

	//Implement using your basic auth
	if(email=="kevin" && password=="888888"){
		return true
	}
	return false
}

