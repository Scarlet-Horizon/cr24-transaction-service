package db

import (
	"database/sql"
	"log"
	"main/model"
)

func (receiver TransactionDB) GetTypes() ([]model.TransactionType, error) {
	stmt, err := receiver.DB.Prepare("SELECT * FROM transaction_type;")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Println("stmt.Close() error", err)
		}
	}(stmt)

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Println("rows.Close() error", err)
		}
	}(rows)

	var types []model.TransactionType

	for rows.Next() {
		var result model.TransactionType
		if err := rows.Scan(&result.ID, &result.Type); err != nil {
			log.Println("rows.Scan() error", err)
			continue
		}
		types = append(types, result)
	}
	if err := rows.Err(); err != nil {
		log.Println("rows.Err() error", err)
	}

	return types, nil
}
