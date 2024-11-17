// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package models

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type TaskHistoryTaskStatus string

const (
	TaskHistoryTaskStatusSucceeded TaskHistoryTaskStatus = "Succeeded"
	TaskHistoryTaskStatusFailed    TaskHistoryTaskStatus = "Failed"
)

func (e *TaskHistoryTaskStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TaskHistoryTaskStatus(s)
	case string:
		*e = TaskHistoryTaskStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for TaskHistoryTaskStatus: %T", src)
	}
	return nil
}

type NullTaskHistoryTaskStatus struct {
	TaskHistoryTaskStatus TaskHistoryTaskStatus
	Valid                 bool // Valid is true if TaskHistoryTaskStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTaskHistoryTaskStatus) Scan(value interface{}) error {
	if value == nil {
		ns.TaskHistoryTaskStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TaskHistoryTaskStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTaskHistoryTaskStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TaskHistoryTaskStatus), nil
}

type Company struct {
	Symbol string
}

type TaskHistory struct {
	ID         string
	TaskName   string
	TaskStatus TaskHistoryTaskStatus
	StartTime  time.Time
	EndTime    time.Time
	Details    string
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
