package address

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

	var address models.Address
	if err := db.First(&address, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, address)
}

func GetAll(c echo.Context) error {
	var addresses []models.Address

	db, _ := c.Get("db").(*gorm.DB)

	if err := db.Model(&models.Address{}).Find(&addresses).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, addresses)
}

func Create(c echo.Context) error {
	type RequestBody struct {
		City string `json:"city"`
	}

	var body RequestBody

	if err := c.Bind(&body); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db, _ := c.Get("db").(*gorm.DB)

	address := models.Address{
		City: body.City,
	}

	db.Create(&address)

	return c.JSON(http.StatusOK, address)
}

func Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}

	type RequestBody struct {
		City string `json:"city"`
	}

	var body RequestBody

	if err := c.Bind(&body); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var address models.Address
	if err := db.First(&address, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	address.City = body.City

	db.Save(&address)

	return c.JSON(http.StatusOK, address)
}

func Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var address models.Address

	if err := db.Where("id = ?", id).First(&address).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db.Delete(&address)

	return c.NoContent(http.StatusOK)
}
