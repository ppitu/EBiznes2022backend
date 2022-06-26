package main

import (
	"backend/database/models"
	"backend/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseContext struct {
	echo.Context
	db *gorm.DB
}

func main() {
	e := echo.New()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Category{})
	db.Migrator().CreateConstraint(&models.Product{}, "Category")
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Address{})
	db.AutoMigrate(&models.CreditCard{})
	db.AutoMigrate(&models.Cart{})
	db.Migrator().CreateConstraint(&models.Cart{}, "Product")

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	routes.Routes(e.Group(""))

	e.Logger.Fatal(e.Start(":1323"))
}
