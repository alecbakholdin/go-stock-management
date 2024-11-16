package main

import (
	"cmp"
	"context"
	"database/sql"
	"stock-management/internal/models"
	"stock-management/internal/task"
	"stock-management/internal/task/yahoo"
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
	Port                  string `env:"PORT"`
	SigningSecret         string `env:"SIGNING_SECRET,required"`
	AdminUsername         string `env:"ADMIN_USERNAME,required"`
	AdminPassword         string `env:"ADMIN_PASSWORD,required"`
	MySqlConnectionString string `env:"MYSQL_CONNECTION_STRING,required"`

	ZacksUrl             string `env:"ZACKS_URL,required"`
	ZacksDailyFormValue  string `env:"ZACKS_DAILY_FORM_VALUE,required"`
	ZacksGrowthFormValue string `env:"ZACKS_GROWTH_FORM_VALUE,required"`

	YahooUrl string `env:"YAHOO_URL_PREFIX,required"`
}

func main() {
	log.SetHeader("${time_rfc3339} ${level}")
	log.SetLevel(log.INFO)

	ec := parseEnv()
	q := initDb(ec)

	tasks := initAndScheduleTasks(ec, q)

	e := initEchoClient(ec, tasks)
	e.Logger.Fatal(e.Start(":" + cmp.Or(ec.Port, "1323")))
}

func parseEnv() EnvConfig {
	var ec EnvConfig
	if err := env.Parse(&ec); err != nil {
		panic(err.Error())
	}
	return ec
}

func initDb(ec EnvConfig) *models.Queries {
	db, err := sql.Open("mysql", ec.MySqlConnectionString)
	if err != nil {
		panic("Error connecting to mysql " + err.Error())
	}
	return models.New(db)
}

func initAndScheduleTasks(ec EnvConfig, q *models.Queries) []task.TaskStatus {
	zacksDailyTask := task.New("Zacks Daily", "/zacksdaily", zacks.NewDaily(q, ec.ZacksUrl, ec.ZacksDailyFormValue))
	zacksGrowthTask := task.New("Zacks Growth", "/zacksgrowth", zacks.NewGrowth(q, ec.ZacksUrl, ec.ZacksGrowthFormValue))
	yahooInsightsTask := task.New("Yahoo Insights", "/yahooinsights", yahoo.NewInsights(q, ec.YahooUrl))

	return []task.TaskStatus{zacksDailyTask, zacksGrowthTask, yahooInsightsTask}
}

func initEchoClient(ec EnvConfig, tasks []task.TaskStatus) *echo.Echo {
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

	return e
}
