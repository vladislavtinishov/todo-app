package handler

import (
	"github.com/gin-gonic/gin"
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

	userId, err := h.services.Authorization.ParseToken(headerParts[1])

	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}
