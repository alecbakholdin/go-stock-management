package main

import (
	"cmp"
	"context"
	"database/sql"
	"os"
	"stock-management/internal/models"
	"stock-management/internal/task"
	"stock-management/internal/web/login"
	"stock-management/internal/web/root"

	"github.com/caarlos0/env/v6"
	_ "github.com/go-sql-driver/mysql"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type EnvConfig struct {
	SigningSecret         string `env:"SIGNING_SECRET,required"`
	AdminUsername         string `env:"ADMIN_USERNAME,required"`
	AdminPassword         string `env:"ADMIN_PASSWORD,required"`
	MySqlConnectionString string `env:"MYSQL_CONNECTION_STRING,required"`
}

func main() {
	var ec EnvConfig
	if err := env.Parse(&ec); err != nil {
		panic(err.Error())
	}

	db, err := sql.Open("mysql", ec.MySqlConnectionString)
	if err != nil {
		panic("Error connecting to mysql " + err.Error())
	}
	models.New(db)

	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetHeader("${time_rfc3339} [id=${id}] ${level}")
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: "${time_rfc3339} [id=${id}] ${method} ${path} ${status} ${error}\n"}))
	e.Use(echojwt.WithConfig(echojwt.Config{
		TokenLookup: "cookie:Authorization", SigningKey: []byte(ec.SigningSecret),
		Skipper: func(c echo.Context) bool {
			_, ok := c.Cookie("Authorization")
			return ok != nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			c.Logger().Error(err)
			return nil
		},
		SuccessHandler: func(c echo.Context) {
			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "user", c.Get("user"))))
		},
		ContinueOnIgnoredError: true,
	}))
	tasks := []task.TaskStatus{

	}

	e.GET("/", root.Handler(tasks))
	e.POST("/login", login.Handler(ec.SigningSecret, ec.AdminUsername, ec.AdminPassword))

	e.Logger.Fatal(e.Start(cmp.Or(os.Getenv("PORT"), ":1323")))
}
