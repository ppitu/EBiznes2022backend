package handlers

import (
	"backend/database/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetAllCart(c echo.Context) error {
	var cart []models.Cart

	db, _ := c.Get("db").(*gorm.DB)

	if err := db.Preload("Product").Find(&cart).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, cart)
}

func CreateCart(c echo.Context) error {
	type RequestBody struct {
		ProductID uint64 `json:"product_id"`
	}
	var body RequestBody

	if err := c.Bind(&body); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var product models.Product

	if err := db.Find(&product, body.ProductID).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	cart := models.Cart{
		ProductID: body.ProductID,
	}

	db.Create(&cart)
	return c.JSON(http.StatusOK, cart)
}

func DeleteAllCart(c echo.Context) error {
	db, _ := c.Get("db").(*gorm.DB)

	db.Exec("DELETE FROM carts")

	return c.NoContent(http.StatusOK)
}
