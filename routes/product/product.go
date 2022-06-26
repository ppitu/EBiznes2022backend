package product

import (
	productHandler "backend/handlers/product"

	"github.com/labstack/echo/v4"
)

type ProductRouter struct{}

func (ps ProductRouter) Init(g *echo.Group) {
	g.GET("", productHandler.GetAll)
	g.GET("/:id", productHandler.Get)
	g.POST("", productHandler.Create)
	g.PUT("/:id", productHandler.Update)
	g.DELETE("/:id", productHandler.Delete)
}
