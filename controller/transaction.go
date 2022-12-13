package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main/db"
	"main/model"
	"main/request"
	"main/response"
	"main/util"
	"net/http"
	"time"
)

type TransactionController struct {
	DB *db.TransactionDB
}

//	@description	Create new transaction.
//	@summary		Create new transaction
//	@accept			json
//	@produce		json
//	@tags			transaction
//	@param			requestBody	body		request.TransactionRequest	true	"Transaction data"
//	@success		201			{object}	model.Transaction
//	@failure		400			{object}	response.ErrorResponse
//	@failure		500			{object}	response.ErrorResponse
//	@security		JWT
//	@param			Authorization	header	string	true	"Authorization"
//	@router			/transaction [POST]
func (receiver TransactionController) Create(context *gin.Context) {
	var req request.TransactionRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		_ = context.Error(err)
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if !util.IsValidUUID(req.SenderAccountID) {
		err := context.Error(errors.New("invalid sender id"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if !util.IsValidUUID(req.RecipientAccountID) {
		err := context.Error(errors.New("invalid recipient id"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if req.SenderAccountID == req.RecipientAccountID {
		err := context.Error(errors.New("can't transfer money to same accounts"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if req.Amount < 1 {
		err := context.Error(errors.New("invalid amount, minimum is 1"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	acc, err := util.GetAccount(req.SenderAccountID, context.MustGet("token").(string))
	if err != nil {
		_ = context.Error(err)
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ok, err := util.ValidateAccount(acc)
	if !ok {
		_ = context.Error(err)
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
	}

	if acc.Amount-req.Amount < float64(-1*acc.Limit) {
		err := context.Error(errors.New("insufficient funds"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	tr := model.Transaction{
		ID:          uuid.NewString(),
		SenderID:    req.SenderAccountID,
		RecipientID: req.RecipientAccountID,
		Amount:      req.Amount,
		Date:        time.Now(),
		Type: model.TransactionType{
			ID:   req.Type,
			Type: "",
		},
	}

	err = receiver.DB.Create(tr)
	if err != nil {
		_ = context.Error(err)
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.JSON(http.StatusCreated, tr)
}

//	@description	Get all transactions for a specific account, where that account was sender or recipient.
//	@summary		Get all transactions for a specific account, where that account was sender or recipient
//	@accept			json
//	@produce		json
//	@tags			transaction
//	@param			accountID	path		string				true	"Account ID"
//	@param			type		path		string				true	"Specifies type of returned transactions: ingoing, outgoing or both. Supported values: 'sender', 'recipient', 'all'"
//	@success		200			{object}	[]model.Transaction	"An array of model.Transaction"
//	@success		204			"No Content"
//	@failure		400			{object}	response.ErrorResponse
//	@failure		500			{object}	response.ErrorResponse
//	@security		JWT
//	@param			Authorization	header	string	true	"Authorization"
//	@router			/transaction/{accountID}/{type} [GET]
func (receiver TransactionController) GetAll(context *gin.Context) {
	accountID := context.Param("accountID")

	if !util.IsValidUUID(accountID) {
		err := context.Error(errors.New("invalid account id"))
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	t := context.Param("type")
	if !(t == "sender" || t == "recipient" || t == "all") {
		err := context.Error(errors.New("invalid type, supported: 'sender', 'recipient', 'all'"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	res, err := receiver.DB.GetAll(accountID, t)
	if err != nil {
		_ = context.Error(err)
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	if len(res) == 0 {
		context.Status(http.StatusNoContent)
		return
	}
	context.JSON(http.StatusOK, res)
}

//	@description	Delete transaction.
//	@summary		Delete transaction
//	@accept			json
//	@produce		json
//	@tags			transaction
//	@param			transactionID	path	string	true	"Transaction ID"
//	@success		204				"No Content"
//	@failure		400				{object}	response.ErrorResponse
//	@failure		500				{object}	response.ErrorResponse
//	@security		JWT
//	@param			Authorization	header	string	true	"Authorization"
//	@router			/transaction/{transactionID}/ [DELETE]
func (receiver TransactionController) Delete(context *gin.Context) {
	transactionID := context.Param("transactionID")

	if !util.IsValidUUID(transactionID) {
		err := context.Error(errors.New("invalid transaction id"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	err := receiver.DB.Delete(transactionID)
	if err != nil {
		_ = context.Error(err)
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.Status(http.StatusNoContent)
}

//	@description	Delete all transactions for the given sender.
//	@summary		Delete all transactions for the given sender
//	@accept			json
//	@produce		json
//	@tags			transaction
//	@param			accountID	path	string	true	"Account ID"
//	@success		204			"No Content"
//	@failure		400			{object}	response.ErrorResponse
//	@failure		500			{object}	response.ErrorResponse
//	@security		JWT
//	@param			Authorization	header	string	true	"Authorization"
//	@router			/transactions/{accountID}/ [DELETE]
func (receiver TransactionController) DeleteForAccount(context *gin.Context) {
	accountID := context.Param("accountID")

	if !util.IsValidUUID(accountID) {
		err := context.Error(errors.New("invalid account id"))
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	err := receiver.DB.DeleteForAccount(accountID)
	if err != nil {
		_ = context.Error(err)
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.Status(http.StatusNoContent)
}
