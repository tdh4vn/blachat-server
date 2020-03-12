package middlewares

import (
	"blachat-server/config"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

func PartnerAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := config.GetConfig()
		reqKey := c.Request.Header.Get("X-Auth-Key")
		reqSecret := c.Request.Header.Get("X-Auth-Secret")

		var key string
		var secret string

		if key = config.GetString("http_key"); len(strings.TrimSpace(key)) == 0 {
			c.AbortWithStatus(500)
		}

		if secret = config.GetString("http_secret"); len(strings.TrimSpace(secret)) == 0 {
			c.AbortWithStatus(401)
		}

		if key != reqKey || secret != reqSecret {
			c.AbortWithStatus(401)
			return
		}
		c.Next()

	}
}


func UserAuthMiddleware() gin.HandlerFunc {
	return func (c *gin.Context) {
		env := config.GetConfig()
		tokenString := c.Request.Header.Get("Authorization")

		tokenString = tokenString[7:]

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(env.GetString("service_sceret")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userID", claims["userId"])
			c.Next()
		} else {
			c.AbortWithStatus(401)
		}

	}
}
