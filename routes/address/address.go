package address

import (
	addressHandler "backend/handlers/address"

	"github.com/labstack/echo/v4"
)

type AddressRouter struct{}

func (ps AddressRouter) Init(g *echo.Group) {
	g.GET("", addressHandler.GetAll)
	g.GET("/:id", addressHandler.Get)
	g.POST("", addressHandler.Create)
	g.PUT("/:id", addressHandler.Update)
	g.DELETE("/:id", addressHandler.Delete)
}
