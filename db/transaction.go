package db

import (
	"database/sql"
	"log"
	"main/model"
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

func (receiver TransactionDB) GetAll(id string) ([]model.Transaction, error) {
	var types []model.TransactionType

	res, err := receiver.DB.Query("SELECT * FROM account_transaction WHERE sender_id = ?;")
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var result model.TransactionType
		if err := res.Scan(&result.ID, &result.Type); err != nil {
			return nil, err
		}
		types = append(types, result)
	}
	return types, nil
}
