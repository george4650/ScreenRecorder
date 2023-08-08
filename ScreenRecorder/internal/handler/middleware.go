package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Аутентификация пользователя
func (h *AuthHandler) userIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerValue := c.GetHeader("Authorization")

		if headerValue == "" {
			log.Printf("AuthHandler - userIdentity - c.GetHeader: %s", "empty auth header")
			c.AbortWithStatusJSON(http.StatusForbidden, fmt.Errorf("Доступ запрещён"))
			return
		}

		headerParts := strings.Split(headerValue, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			log.Printf("AuthHandler - userIdentity - c.GetHeader: %s", "invalid auth header")
			c.AbortWithStatusJSON(http.StatusForbidden, fmt.Errorf("Доступ запрещён"))
			return
		}

		if len(headerParts[1]) == 0 {
			log.Printf("AuthHandler - userIdentity - c.GetHeader: %s", "token is empty")
			c.AbortWithStatusJSON(http.StatusForbidden, fmt.Errorf("Доступ запрещён"))
			return
		}

		_, err := h.us.ParseToken(headerParts[1])
		if err != nil {
			log.Printf("AuthHandler - userIdentity - h.us.ParseToken: %v \n", err)
			c.AbortWithStatusJSON(http.StatusForbidden, fmt.Errorf("Доступ запрещён"))
			return
		}

		c.Next()
	}
}
