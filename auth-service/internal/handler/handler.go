package handler

import (
	"auth-service/api/internal/handler/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"auth-service/api/pkg/model"
)

type Handler struct {
	UserService  model.UserService
	TokenService model.TokenService
}

type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	TokenService    model.TokenService
	BaseURL         string
	TimeoutDuration time.Duration
}

func NewHandler(c *Config) {
	h := &Handler{
		UserService:  c.UserService,
		TokenService: c.TokenService,
	}

	g := c.R.Group(c.BaseURL)

	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(c.TimeoutDuration, model.NewServiceUnavailable()))
		g.GET("/me", middleware.AuthUser(h.TokenService), h.Me)
		g.POST("/signout", middleware.AuthUser(h.TokenService), h.Signout)
	} else {
		g.GET("/me", h.Me)
		g.POST("/signout", h.Signout)
	}

	g.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"pong": "test"})
	})
	g.POST("/signup", h.Signup)
	g.POST("/signin", h.Signin)
	g.POST("/tokens", h.Tokens)
}
