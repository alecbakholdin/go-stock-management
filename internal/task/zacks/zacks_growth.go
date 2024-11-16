package zacks

import (
	"context"
	"stock-management/internal/models"
	"time"

	"github.com/labstack/gommon/log"
)

type ZacksGrowthRowSaver interface {
	SaveZacksGrowthRow(context.Context, models.SaveZacksGrowthRowParams) error
}

type zacksGrowth struct {
	q ZacksGrowthRowSaver
}

func NewGrowth(q ZacksGrowthRowSaver, url, formValue string) *zacksExecutor[zacksGrowthCsvRow] {
	return &zacksExecutor[zacksGrowthCsvRow]{
		ms:        &zacksGrowth{q: q},
		url:       url,
		formValue: formValue,
		tableName: "Growth",
	}
}

func (g *zacksGrowth) save(csvRow zacksGrowthCsvRow) error {
	nextReportDate, err := time.Parse("1/2/06", csvRow.NextReportDate)
	if err != nil {
		log.Warn("Error parsing time string ", csvRow.NextReportDate, " during Zacks Growth task: ", err)
	}

	sqlRow := models.SaveZacksGrowthRowParams{
		Symbol:                             csvRow.Symbol,
		Company:                            models.NullStringIfMatch(csvRow.Company, "NA"),
		Price:                              csvRow.Price,
		GrowthScore:                        models.NullStringIfMatch(csvRow.GrowthScore, "NA"),
		YearOverYearQ0Growth:               csvRow.YearOverYearQ0Growth,
		LongTermGrowthPercent:              csvRow.LongTermGrowthPercent,
		LastFinancialYearActual:            csvRow.LastFinancialYearActual,
		ThisFinancialYearEst:               csvRow.ThisFinancialYearEst,
		NextFinanicalYearEst:               csvRow.NextFinanicalYearEst,
		Q1Est:                              csvRow.Q1Est,
		EarningsExpectedSurprisePrediction: csvRow.EarningsExpectedSurprisePrediction,
		NextReportDate:                     models.NullTimeIfZero(nextReportDate),
	}
	return g.q.SaveZacksGrowthRow(context.Background(), sqlRow)
}

type zacksGrowthCsvRow struct {
	Symbol                             string
	Company                            string
	Price                              float64
	GrowthScore                        string  `csv:"Growth Score"`
	YearOverYearQ0Growth               float64 `csv:"YR/YR Q0 Growth"`
	LongTermGrowthPercent              float64 `csv:"LTG %"`
	LastFinancialYearActual            float64 `csv:"Last FY Actual"`
	ThisFinancialYearEst               float64 `csv:"This FY Est"`
	NextFinanicalYearEst               float64 `csv:"Next FY Est"`
	Q1Est                              float64 `csv:"Q1 Est"`
	EarningsExpectedSurprisePrediction float64 `csv:"Earnings ESP"`
	NextReportDate                     string  `csv:"Next Report Date"`
}
