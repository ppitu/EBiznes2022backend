package user

import (
	userHandler "backend/handlers/user"

	"github.com/labstack/echo/v4"
)

type UserRouter struct{}

func (ps UserRouter) Init(g *echo.Group) {
	g.GET("", userHandler.GetAll)
	g.GET("/:id", userHandler.Get)
	g.POST("", userHandler.Create)
	g.PUT("/:id", userHandler.Update)
	g.DELETE("/:id", userHandler.Delete)
	g.GET("/auth/google/login", userHandler.OauthGoogleLogin)
	g.GET("/auth/google/callback", userHandler.OauthGoogleCallback)
	g.GET("/auth/github/login", userHandler.OAuthGithubLogin)
	g.GET("/auth/github/callback", userHandler.OauthGithubCallback)
	g.POST("/create-payment-intent", userHandler.HandleCreatePaymentIntent)
}
