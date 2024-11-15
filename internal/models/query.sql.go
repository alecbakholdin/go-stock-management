// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package models

import (
	"context"
	"database/sql"
)

const listCompanies = `-- name: ListCompanies :many
SELECT symbol, name
FROM company
`

func (q *Queries) ListCompanies(ctx context.Context) ([]Company, error) {
	rows, err := q.db.QueryContext(ctx, listCompanies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Company
	for rows.Next() {
		var i Company
		if err := rows.Scan(&i.Symbol, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveZacksDailyRow = `-- name: SaveZacksDailyRow :exec
INSERT INTO zacks_daily (
        symbol,
        company,
        price,
        dollar_change,
        percent_change,
        industry_rank,
        zacks_rank,
        value_score,
        growth_score,
        momentum_score,
        vgm_score
    )
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type SaveZacksDailyRowParams struct {
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

func (q *Queries) SaveZacksDailyRow(ctx context.Context, arg SaveZacksDailyRowParams) error {
	_, err := q.db.ExecContext(ctx, saveZacksDailyRow,
		arg.Symbol,
		arg.Company,
		arg.Price,
		arg.DollarChange,
		arg.PercentChange,
		arg.IndustryRank,
		arg.ZacksRank,
		arg.ValueScore,
		arg.GrowthScore,
		arg.MomentumScore,
		arg.VgmScore,
	)
	return err
}
