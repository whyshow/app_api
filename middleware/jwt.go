package middleware

import (
	"app_api/config"
	"app_api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

// GenerateJWT 生成JWT令牌
func GenerateJWT(uid string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   uid,
		ExpiresAt: expirationTime.Unix(),
		Issuer:    "app_api",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Get().JWTSecret))
}

// JWTAuth 鉴权中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "请先登录"})
			return
		}

		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Get().JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "无效凭证"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if _, err := models.GetUserByUID(claims["sub"].(string)); err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "用户不存在"})
			return
		}

		c.Set("uid", claims["sub"])
		c.Next()
	}
}
