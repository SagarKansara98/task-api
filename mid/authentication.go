package mid

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Mid struct {
	TokenSeceret string
}

func (m Mid) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the JWT token from the Authorization header
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return echo.ErrUnauthorized
		}
		tokenString = strings.Split(tokenString, " ")[1]
		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method:sa %v", token.Header["alg"])
			}
			return []byte(m.TokenSeceret), nil
		})

		if err != nil || !token.Valid {
			return echo.ErrUnauthorized
		}

		// Get the 'sub' claim from the token and set it in the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			sub := claims["sub"].(float64)
			c.Set("user_id", int(sub)) // Set 'sub' claim in the context
			return next(c)
		}

		return echo.ErrUnauthorized
	}
}
