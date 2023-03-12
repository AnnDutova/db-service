package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"auth-service/api/pkg/model"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func AuthUser(s model.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}

		if err := c.ShouldBindHeader(&h); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				err := model.BadRequestError("Invalid request parameters. See invalidArgs")

				c.JSON(err.Status, gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				c.Abort()
				return
			}

			err := model.InternalError()
			c.JSON(err.Status, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		idTokenHeader := strings.Split(h.IDToken, "Bearer ")

		if len(idTokenHeader) < 2 {
			err := model.UnauthorizedError("Must provide Authorization header with format `Bearer {token}`")

			c.JSON(err.Status, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		user, err := s.ValidateIDToken(idTokenHeader[1])

		if err != nil {
			err := model.UnauthorizedError("Provided token is invalid")
			c.JSON(err.Status, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
