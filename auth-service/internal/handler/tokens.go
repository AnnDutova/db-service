package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"auth-service/api/pkg/model"
)

type tokensReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// Tokens handler
func (h *Handler) Tokens(c *gin.Context) {
	var req tokensReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	refreshToken, err := h.TokenService.ValidateRefreshToken(req.RefreshToken)

	if err != nil {
		c.JSON(model.Status(err), gin.H{
			"error": err,
		})
		return
	}

	u, err := h.UserService.Get(ctx, refreshToken.UID)

	if err != nil {
		c.JSON(model.Status(err), gin.H{
			"error": err,
		})
		return
	}
	
	tokens, err := h.TokenService.NewToken(ctx, u, refreshToken.ID.String())

	if err != nil {
		log.Printf("Failed to create tokens for user: %+v. Error: %v\n", u, err.Error())

		c.JSON(model.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
