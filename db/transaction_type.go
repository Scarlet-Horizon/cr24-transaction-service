package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"main/model"
)

func (receiver TransactionDB) GetTypes(ctx *gin.Context) ([]model.TransactionType, error) {
	stmt, err := receiver.DB.Prepare("SELECT * FROM transaction_type;")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			_ = ctx.Error(errors.New(fmt.Sprintf("stmt.Close() error: %v", err)))
		}
	}(stmt)

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			_ = ctx.Error(errors.New(fmt.Sprintf("rows.Close() error: %v", err)))
		}
	}(rows)

	var types []model.TransactionType

	for rows.Next() {
		var result model.TransactionType
		if err := rows.Scan(&result.ID, &result.Type); err != nil {
			_ = ctx.Error(errors.New(fmt.Sprintf("rows.Scan() error: %v", err)))
			continue
		}
		types = append(types, result)
	}
	if err := rows.Err(); err != nil {
		_ = ctx.Error(errors.New(fmt.Sprintf("rows.Err() error: %v", err)))
	}

	return types, nil
}
