package routes

import (
	"backend/routes/address"
	"backend/routes/category"
	creditcard "backend/routes/credit_card"
	"backend/routes/product"
	"backend/routes/user"

	"github.com/labstack/echo/v4"
)

func Routes(g *echo.Group) {
	product.ProductRouter{}.Init(g.Group("/products"))
	category.CategoryRouter{}.Init(g.Group("/categories"))
	user.UserRouter{}.Init(g.Group("/users"))
	address.AddressRouter{}.Init(g.Group("/addresses"))
	creditcard.CreditCardRouter{}.Init(g.Group("/credit_cards"))
	CartRouter{}.Init(g.Group("/cart"))
}
