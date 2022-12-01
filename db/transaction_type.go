package db

import (
	"main/model"
)

func (receiver TransactionDB) GetTypes() ([]model.TransactionType, error) {
	var types []model.TransactionType

	res, err := receiver.DB.Query("SELECT * FROM transaction_type;")
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
