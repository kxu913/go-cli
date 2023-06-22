package middleware

import (
	"encoding/base64"
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
		Validator CustomAuthValidator
		// Default value "Restricted".
		Realm string
	}

	// BasicAuthValidator defines a function to validate BasicAuth credentials.
	CustomAuthValidator func(string, string, echo.Context) (bool, error)
)
var DefaultCustomAuthConfig = CustomAuthConfig{
	Skipper: func(ctx echo.Context) bool {
		// allow some api to skip validator, don't recommand do as this
		// if ctx.Path() == "/api/user" && ctx.Request().Method == "POST" {
		//	return true

		// }
		return false
	},
	Realm: defaultRealm,
}

const (
	basic        = "basic"
	defaultRealm = "Restricted"
	bearer       = "bearer"
)

func BasicAuth(fn CustomAuthValidator) echo.MiddlewareFunc {
	c := DefaultCustomAuthConfig
	c.Validator = fn

	return BasicAuthWithConfig(c)
}

func BasicAuthWithConfig(config CustomAuthConfig) echo.MiddlewareFunc {

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
						isOK, err := config.Validator(email, password, c)
						if err != nil || !isOK {
							return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "authorization failed, check your email and password."}
						}
						c.Request().Header.Add("email", email)
						c.Request().Header.Add("email", password)
						break
					}
				}

			}

			return next(c)
		}
	}
}
