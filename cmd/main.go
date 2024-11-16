package main

import (
	"cmp"
	"context"
	"database/sql"
	"os"
	"stock-management/internal/models"
	"stock-management/internal/task"
	"stock-management/internal/task/zacks"
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

	ZacksUrl string `env:"ZACKS_URL,required"`
	ZacksDailyFormValue string `env:"ZACKS_DAILY_FORM_VALUE,required"`
	ZacksGrowthFormValue string `env:"ZACKS_GROWTH_FORM_VALUE,required"`
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
	q := models.New(db)

	log.SetHeader("${time_rfc3339} ${level}")
	log.SetLevel(log.INFO)
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
		task.New(e, "Zacks Daily", "/zacksdaily", zacks.NewDaily(q, ec.ZacksUrl, ec.ZacksDailyFormValue)),
		task.New(e, "Zacks Growth", "/zacksgrowth", zacks.NewGrowth(q, ec.ZacksUrl, ec.ZacksGrowthFormValue)),
	}

	e.GET("/", root.Handler(tasks))
	e.POST("/login", login.Handler(ec.SigningSecret, ec.AdminUsername, ec.AdminPassword))
	type TaskUrlGetPost interface {
		UrlPath() string
		GetHandler(c echo.Context) error
		PostHandler(c echo.Context) error
	}
	for _, task := range tasks {
		if handlers, ok := task.(TaskUrlGetPost); ok {
			e.GET(handlers.UrlPath(), handlers.GetHandler)
			e.POST(handlers.UrlPath(), handlers.PostHandler)
		} else {
			panic("task does not have handlers")
		}
	}

	e.Logger.Fatal(e.Start(cmp.Or(os.Getenv("PORT"), ":1323")))
}
