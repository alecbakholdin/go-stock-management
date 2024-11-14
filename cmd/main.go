package main

import (
	"cmp"
	"database/sql"
	"os"
	"stock-management/internal/models"

	_ "github.com/go-sql-driver/mysql"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_CONNECTION_STRING"))
	if err != nil {
		panic("Error connecting to mysql " + err.Error())
	}
	secret := os.Getenv("SIGNING_SECRET");
	if secret == "" {
		panic("Must include SIGNING_SECRET env variable")
	}
	models.New(db)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: "${time_rfc3339} [id=${id}] ${method} ${path} ${status} ${error}"}))
	e.Use(echojwt.WithConfig(echojwt.Config{TokenLookup: "cookie:Authorization", SigningKey: []byte(secret)}))

	e.Logger.Fatal(e.Start(cmp.Or(os.Getenv("PORT"), "1323")))
}
