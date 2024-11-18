package main

import (
	"cmp"
	"context"
	"database/sql"
	"stock-management/internal/models"
	"stock-management/internal/task"
	"stock-management/internal/task/tipranks"
	"stock-management/internal/task/yahoo"
	"stock-management/internal/task/zacks"
	"stock-management/internal/util/must"
	"stock-management/internal/web/login"
	"stock-management/internal/web/root"
	"strings"

	"github.com/caarlos0/env/v6"
	_ "github.com/go-sql-driver/mysql"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pressly/goose/v3"
	"github.com/robfig/cron/v3"
)

type EnvConfig struct {
	Port string `env:"PORT"`

	SigningSecret string `env:"SIGNING_SECRET,required"`
	AdminUsername string `env:"ADMIN_USERNAME,required"`
	AdminPassword string `env:"ADMIN_PASSWORD,required"`

	MySqlConnectionString string `env:"MYSQL_CONNECTION_STRING,required"`
	GooseMigrationsDir    string `env:"GOOSE_MIGRATIONS_DIR,required"`

	ZacksUrl             string `env:"ZACKS_URL,required"`
	ZacksDailyFormValue  string `env:"ZACKS_DAILY_FORM_VALUE,required"`
	ZacksGrowthFormValue string `env:"ZACKS_GROWTH_FORM_VALUE,required"`

	YahooInsightsUrl    string `env:"YAHOO_INSIGHTS_URL_PREFIX,required"`
	YahooQuotesCrumbUrl string `env:"YAHOO_QUOTES_GET_CRUMB_URL,required"`
	YahooQuotesUrl      string `env:"YAHOO_QUOTES_URL_PREFIX,required"`

	TipRanksUrl string `env:"TIPRANKS_URL_PREFIX,required"`
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
	if atLoc := strings.Index(ec.MySqlConnectionString, "@"); atLoc < 0 {
		log.Warn("MYSQL_CONNECTION_STRING missing @?")
	} else if strings.Contains(ec.MySqlConnectionString[atLoc:], "?") {
		panic("MYSQL_CONNECTION_STRING should not contain any query params")
	}
	db, err := sql.Open("mysql", ec.MySqlConnectionString+"?parseTime=true&multiStatements=true")
	if err != nil {
		panic("Error connecting to mysql " + err.Error())
	}
	if err := goose.SetDialect(string(goose.DialectMySQL)); err != nil {
		log.Fatalf("Error setting MYSQL dialect: %s", err.Error())
	} else if err := goose.Up(db, ec.GooseMigrationsDir); err != nil {
		log.Fatalf("Error migrating DB: %s", err.Error())
	}
	return models.New(db)
}

func initAndScheduleTasks(ec EnvConfig, q *models.Queries) []task.Task {
	zacksDailyTask := task.New(q, "Zacks Daily", "/zacksdaily", zacks.NewDaily(q, ec.ZacksUrl, ec.ZacksDailyFormValue))
	zacksGrowthTask := task.New(q, "Zacks Growth", "/zacksgrowth", zacks.NewGrowth(q, ec.ZacksUrl, ec.ZacksGrowthFormValue))
	yahooInsightsTask := task.New(q, "Yahoo Insights", "/yahooinsights", yahoo.NewInsights(q, ec.YahooInsightsUrl))
	tipranksTask := task.New(q, "TipRanks", "/tipranks", tipranks.New(q, ec.TipRanksUrl))
	yahooQuotesTask := task.New(q, "Yahoo Quotes", "/yahooquotes", yahoo.NewQuotes(q, ec.YahooQuotesCrumbUrl, ec.YahooQuotesUrl))

	tasks := []task.Task{zacksDailyTask, zacksGrowthTask, yahooInsightsTask, yahooQuotesTask, tipranksTask}

	allTasks := func() {
		for _, task := range tasks {
			task.Execute()
		}
	}

	c := cron.New(cron.WithLocation(must.MustLoadLocation("America/New_York")))
	if _, err := c.AddFunc("0 10,14 * * MON-FRI", allTasks); err != nil {
		log.Fatalf("Error adding cron job 1: %s", err.Error())
	} else if _, err := c.AddFunc("30 16 * * MON-FRI", allTasks); err != nil {

		log.Fatalf("Error adding cron job 2: %s", err.Error())
	}
	c.Start()

	return tasks
}

func initEchoClient(ec EnvConfig, tasks []task.Task) *echo.Echo {
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
	for _, task := range tasks {
		e.GET(task.UrlPath(), task.GetHandler)
		e.POST(task.UrlPath(), task.PostHandler)
	}

	return e
}
