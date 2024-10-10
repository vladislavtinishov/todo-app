package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vladislavtinishov/todo-app/utils"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) == 1 {
		NewErrorResponse(c, http.StatusUnauthorized, "empty token")
		return
	}

	claims, err := utils.ParseJWT(headerParts[1])

	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	userId := claims.UserID

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)

	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)

	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "user is of invalid type")
		return 0, errors.New("user is of invalid type")
	}

	return idInt, nil
}
