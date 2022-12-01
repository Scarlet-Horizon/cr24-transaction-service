package db

import (
	"database/sql"
	"main/model"
)

type TransactionDB struct {
	DB *sql.DB
}

func (receiver TransactionDB) Create(transaction model.Transaction) error {
	_, err := receiver.DB.Exec("INSERT INTO account_transaction VALUES (?,?,?,?,?,?);", transaction.ID,
		transaction.SenderID, transaction.RecipientID, transaction.Amount, transaction.GetDate(), transaction.Type)
	return err
}
