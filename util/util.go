package util

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"log"
	"main/model"
	"net/http"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func GetAccount(userID, accountID string) (model.Account, error) {
	response, err := http.Get("http://account-api:8080/api/v1/account/" + userID + "/" + accountID)

	if err != nil {
		return model.Account{}, errors.New("Get error: " + err.Error())
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return model.Account{}, errors.New("ReadAll error: " + err.Error())
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			log.Printf("Close error: %s\n", err)
		}
	}(response.Body)

	if response.StatusCode != 200 {
		return model.Account{}, errors.New("error:" + string(data))
	}

	var acc model.Account
	if err := json.Unmarshal(data, &acc); err != nil {
		return model.Account{}, errors.New("encode data error: " + err.Error())
	}
	return acc, nil
}

func ValidateAccount(account model.Account) (bool, error) {
	if account.PK == "" {
		return false, errors.New("invalid account")
	}

	if account.CloseDate != nil {
		return false, errors.New("account is closed")
	}
	return true, nil
}
