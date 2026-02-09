package middleware

import (
	"commmunity/app/internal/response"
	"commmunity/app/internal/service/login"
	"commmunity/app/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.FailWithMessage(c, "未登录")
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.FailWithMessage(c, "token格式不对")
			c.Abort()
			return
		}
		tokenString := parts[1]
		if !login.IsTokenValid(tokenString) {
			response.FailWithMessage(c, "请重新登录")
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			response.FailWithMessage(c, "token过期或无效")
			c.Abort()
			return
		}
		c.Set("account", (*claims).Account)
		c.Set("role", (*claims).Role)
		c.Next()
	}
}

func CorsMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// "*" 代表允许所有域名。生产环境建议换成具体的前端域名，比如 "https://your-domain.com"
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, session, X_Requested_With, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Allow-Credentials", "true")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}
		c.Next()
	}
}
