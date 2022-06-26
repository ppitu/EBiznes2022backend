package creditcard

import (
	creaditcardHandler "backend/handlers/credit_card"

	"github.com/labstack/echo/v4"
)

type CreditCardRouter struct{}

func (ps CreditCardRouter) Init(g *echo.Group) {
	g.GET("", creaditcardHandler.GetAll)
	g.GET("/:id", creaditcardHandler.Get)
	g.POST("", creaditcardHandler.Create)
	g.PUT("/:id", creaditcardHandler.Update)
	g.DELETE("/:id", creaditcardHandler.Delete)
}
