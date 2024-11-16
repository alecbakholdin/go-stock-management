// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package models

import (
	"database/sql"
	"time"
)

type Company struct {
	Symbol string
}

type YahooInsight struct {
	Created         time.Time
	Symbol          string
	ShortTerm       sql.NullString
	MidTerm         sql.NullString
	LongTerm        sql.NullString
	FairValue       sql.NullString
	EstimatedReturn sql.NullInt32
}

type ZacksDaily struct {
	Created       time.Time
	Symbol        string
	Company       sql.NullString
	Price         float64
	DollarChange  float64
	PercentChange float64
	IndustryRank  sql.NullInt32
	ZacksRank     sql.NullInt32
	ValueScore    sql.NullString
	GrowthScore   sql.NullString
	MomentumScore sql.NullString
	VgmScore      sql.NullString
}

type ZacksGrowth struct {
	Created                            time.Time
	Symbol                             string
	Company                            sql.NullString
	Price                              float64
	GrowthScore                        sql.NullString
	YearOverYearQ0Growth               float64
	LongTermGrowthPercent              float64
	LastFinancialYearActual            float64
	ThisFinancialYearEst               float64
	NextFinanicalYearEst               float64
	Q1Est                              float64
	EarningsExpectedSurprisePrediction float64
	NextReportDate                     sql.NullTime
}
