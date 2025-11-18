package middlewares

import (
	"net/http"
	"strings"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/services"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT(role models.JwtServiceRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			return
		}

		if len(strings.Split(authHeader, " ")) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		if strings.ToLower(strings.Split(authHeader, " ")[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
		}

		tokenString := strings.Split(authHeader, " ")[1]
		jwtService, _ := services.NewJWTService()
		token, err := jwtService.ValidateToken(tokenString, role)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if token.Valid {
			claims := token.Claims.(*services.AuthClaims)
			c.Set("identifier", claims.Identifier)
			c.Set("role", role)
			c.Set("id", claims.ID)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}
	}
}

func AuthorizeProductionLineJWT() gin.HandlerFunc {
	return AuthorizeJWT(models.JwtServiceRoleProductionLine)
}
