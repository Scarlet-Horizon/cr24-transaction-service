package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"main/model"
	"time"
)

type TransactionDB struct {
	DB *sql.DB
}

func (receiver TransactionDB) Create(transaction model.Transaction, ctx *gin.Context) error {
	stmt, err := receiver.DB.Prepare("INSERT INTO account_transaction (id_transaction, sender_id, recipient_id, " +
		"amount, t_date, fk_t_type) VALUES (?,?,?,?,?,?);")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			_ = ctx.Error(errors.New(fmt.Sprintf("stmt.Close() error: %v", err)))
		}
	}(stmt)

	_, err = stmt.Exec(transaction.ID, transaction.SenderID, transaction.RecipientID, transaction.Amount,
		transaction.GetDate(), transaction.Type.ID)
	return err
}

func (receiver TransactionDB) GetAll(id, t string, ctx *gin.Context) ([]model.Transaction, error) {
	query := "SELECT acT.*, tt.* FROM account_transaction AS acT JOIN transaction_type AS tt ON acT.fk_t_type = tt.id_transaction_type"

	if t == "sender" {
		query += " WHERE sender_id = ?;"
	} else if t == "recipient" {
		query += " WHERE recipient_id = ?;"
	} else {
		query += " WHERE sender_id = ? OR recipient_id = ?;"
	}

	stmt, err := receiver.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			_ = ctx.Error(errors.New(fmt.Sprintf("stmt.Close() error: %v", err)))
		}
	}(stmt)

	var rows *sql.Rows
	if t == "all" {
		rows, err = stmt.Query(id, id)
	} else {
		rows, err = stmt.Query(id)
	}

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			_ = ctx.Error(errors.New(fmt.Sprintf("rows.Close() error: %v", err)))
		}
	}(rows)

	var types []model.Transaction

	for rows.Next() {
		var result model.Transaction
		var tDate string
		var tt model.TransactionType

		if err := rows.Scan(&result.ID, &result.SenderID, &result.RecipientID, &result.Amount, &tDate,
			&tt.ID, &tt.ID, &tt.Type); err != nil {
			_ = ctx.Error(errors.New(fmt.Sprintf("rows.Scan() error: %v", err)))
			continue
		}

		result.Date, err = time.Parse("2006-01-02 15:04:05", tDate)
		if err != nil {
			_ = ctx.Error(errors.New(fmt.Sprintf("time.Parse() error: %v", err)))
			continue
		}

		result.Type = tt

		types = append(types, result)
	}
	if err := rows.Err(); err != nil {
		_ = ctx.Error(errors.New(fmt.Sprintf("rows.Err() error: %v", err)))
	}
	return types, nil
}

func (receiver TransactionDB) delete(id, query string, ctx *gin.Context) error {
	stmt, err := receiver.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			_ = ctx.Error(errors.New(fmt.Sprintf("stmt.Close() error: %v", err)))
		}
	}(stmt)

	_, err = stmt.Exec(id)
	return err
}

func (receiver TransactionDB) Delete(id string, ctx *gin.Context) error {
	return receiver.delete(id, "DELETE FROM account_transaction WHERE id_transaction = ?;", ctx)
}

func (receiver TransactionDB) DeleteForAccount(id string, ctx *gin.Context) error {
	return receiver.delete(id, "DELETE FROM account_transaction WHERE sender_id = ?;", ctx)
}
