package category

import (
	categoryHandler "backend/handlers/category"

	"github.com/labstack/echo/v4"
)

type CategoryRouter struct{}

func (ps CategoryRouter) Init(g *echo.Group) {
	g.GET("", categoryHandler.GetAll)
	g.GET("/:id", categoryHandler.Get)
	g.POST("", categoryHandler.Create)
	g.PUT("/:id", categoryHandler.Update)
	g.DELETE("/:id", categoryHandler.Delete)
}
