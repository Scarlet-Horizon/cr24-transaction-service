package db

import (
	"database/sql"
	"log"
	"main/model"
	"time"
)

type TransactionDB struct {
	DB *sql.DB
}

func (receiver TransactionDB) Create(transaction model.Transaction) error {
	stmt, err := receiver.DB.Prepare("INSERT INTO account_transaction (id_transaction, sender_id, recipient_id, " +
		"amount, t_date, fk_t_type) VALUES (?,?,?,?,?,?);")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Println("stmt.Close() error", err)
		}
	}(stmt)

	_, err = stmt.Exec(transaction.ID, transaction.SenderID, transaction.RecipientID, transaction.Amount,
		transaction.GetDate(), transaction.Type)
	return err
}

func (receiver TransactionDB) GetAll(id, t string) ([]model.Transaction, error) {
	var query string
	if t == "sender" {
		query = "SELECT * FROM account_transaction WHERE sender_id = ?;"
	} else {
		query = "SELECT * FROM account_transaction WHERE recipient_id = ?;"
	}

	stmt, err := receiver.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Println("stmt.Close() error", err)
		}
	}(stmt)

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Println("rows.Close() error", err)
		}
	}(rows)

	var types []model.Transaction

	for rows.Next() {
		var result model.Transaction
		var tDate string

		if err := rows.Scan(&result.ID, &result.SenderID, &result.RecipientID, &result.Amount, &tDate,
			&result.Type); err != nil {
			log.Println("rows.Scan() error", err)
			continue
		}

		result.Date, err = time.Parse("2006-01-02 15:04:05", tDate)
		if err != nil {
			log.Println("time.Parse() error", err)
			continue
		}

		types = append(types, result)
	}
	return types, nil
}

func (receiver TransactionDB) delete(id, query string) error {
	stmt, err := receiver.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Println("stmt.Close() error", err)
		}
	}(stmt)

	_, err = stmt.Exec(id)
	return err
}

func (receiver TransactionDB) Delete(id string) error {
	return receiver.delete(id, "DELETE FROM account_transaction WHERE id_transaction = ?;")
}

func (receiver TransactionDB) DeleteForAccount(id string) error {
	return receiver.delete(id, "DELETE FROM account_transaction WHERE sender_id = ?;")
}
