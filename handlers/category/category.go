package category

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

	var category models.Category

	if err := db.First(&category, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, category)
}

func GetAll(c echo.Context) error {
	var categories []models.Category

	db, _ := c.Get("db").(*gorm.DB)

	if err := db.Model(&models.Category{}).Find(&categories).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, categories)
}

func Create(c echo.Context) error {
	type RequestBody struct {
		Name string `json:"name"`
	}

	var body RequestBody

	if err := c.Bind(&body); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db, _ := c.Get("db").(*gorm.DB)

	category := models.Category{
		Name: body.Name,
	}

	db.Create(&category)

	return c.JSON(http.StatusOK, category)
}

func Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}

	type RequestBody struct {
		Name string `json:"name"`
	}

	var body RequestBody

	if err := c.Bind(&body); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var category models.Category
	if err := db.First(&category, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	category.Name = body.Name

	db.Save(&category)

	return c.JSON(http.StatusOK, category)
}

func Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var category models.Category

	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db.Delete(&category)

	return c.NoContent(http.StatusOK)
}
