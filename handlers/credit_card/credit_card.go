package creditcard

import (
	"backend/database/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var creditcard models.CreditCard
	if err := db.First(&creditcard, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, creditcard)
}

func GetAll(c echo.Context) error {
	var creditcards []models.CreditCard

	db, _ := c.Get("db").(*gorm.DB)

	if err := db.Model(&models.CreditCard{}).Find(&creditcards).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, creditcards)
}

func Create(c echo.Context) error {
	type RequestBody struct {
		Number string `json:"number"`
	}

	var body RequestBody

	if err := c.Bind(&body); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db, _ := c.Get("db").(*gorm.DB)

	creditcard := models.CreditCard{
		Number: body.Number,
	}

	db.Create(&creditcard)

	return c.JSON(http.StatusOK, creditcard)
}

func Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}

	type RequestBody struct {
		Number string `json:"number"`
	}

	var body RequestBody

	if err := c.Bind(&body); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var creditcard models.CreditCard
	if err := db.First(&creditcard, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	creditcard.Number = body.Number

	db.Save(&creditcard)

	return c.JSON(http.StatusOK, creditcard)
}

func Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var creditcard models.CreditCard

	if err := db.Where("id = ?", id).First(&creditcard).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db.Delete(&creditcard)

	return c.NoContent(http.StatusOK)
}
