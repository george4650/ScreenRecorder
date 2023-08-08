package handler

import (
	"myapp/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	us usecase.AuthUseCases
}


// Authorization Авторизация
//
//	@Summary		Авторизация
//	@Description	Авторизация.
//	@Tags			Auth User
//	@Produce		json
//	@Param			as	query	handler.AuthUser.AuthUserRequest	true	"Модификатор вида списка"
//	@Success		200	string	token
//	@Router			/sign-in/ [post]
func (h *AuthHandler) AuthUser(c *gin.Context) {

	type AuthUserRequest struct {
		Login    string `form:"login" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	user := AuthUserRequest{}

	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.us.Authentication(user.Login, user.Password)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "invalid credentials"):
			c.JSON(http.StatusForbidden, "Недействительные аутентификационные данные.")
			return
		default:
			c.JSON(http.StatusInternalServerError, "Внутренняя ошибка сервера")
			return
		}
	}

	token, err := h.us.GenerateToken(user.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
