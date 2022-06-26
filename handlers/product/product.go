package product

import (
	"backend/database/models"
	"fmt"
	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Get(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	db, _ := c.Get("db").(*gorm.DB)

	var product models.Product
	if err := db.Preload("Category").First(&product, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, product)
}

func GetAll(c echo.Context) error {

	var products []models.Product

	db, _ := c.Get("db").(*gorm.DB)

	err := db.Preload("Category").Find(&products)

	if err.Error != nil {
		c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, products)
}

func Create(c echo.Context) error {
	type RequestBody struct {
		Name       string `json:"name"`
		CategoryID uint64 `json:"category_id"`
	}

	var body RequestBody
	fmt.Println("Test1")
	if err := c.Bind(&body); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var category models.Category
	fmt.Println("Test2")
	if err := db.First(&category, body.CategoryID).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	product := models.Product{
		Name:       body.Name,
		CategoryID: body.CategoryID,
	}
	fmt.Println("Test3")
	db.Create(&product)
	return c.JSON(http.StatusOK, product)
}

func Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	type RequestBody struct {
		Name string `json:"name"`
	}

	var body RequestBody

	if err := c.Bind(&body); err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var product models.Product
	if err := db.First(&product, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	product.Name = body.Name

	db.Save(&product)

	return c.JSON(http.StatusOK, product)
}

func Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusPreconditionFailed)
	}

	db, _ := c.Get("db").(*gorm.DB)

	var product models.Product

	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	db.Delete(&product)

	return c.NoContent(http.StatusOK)
}
