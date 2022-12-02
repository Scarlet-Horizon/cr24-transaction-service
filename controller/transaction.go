package controller

import (
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
//	@router			/transaction [POST]
func (receiver TransactionController) Create(context *gin.Context) {
	var req request.TransactionRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if !util.IsValidUUID(req.SenderAccountID) {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid sender id"})
		return
	}

	if !util.IsValidUUID(req.RecipientAccountID) {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid recipient id"})
		return
	}

	if req.SenderAccountID == req.RecipientAccountID {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "can't transfer money to same accounts"})
		return
	}

	if req.Amount >= 0 && req.Amount < 1 || req.Amount <= 0 && req.Amount > -1 {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid amount, minimum is 1 or -1"})
		return
	}

	acc, err := util.GetAccount(req.SenderID, req.SenderAccountID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ok, err := util.ValidateAccount(acc)
	if !ok {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
	}

	if acc.Amount-req.Amount < float64(-1*acc.Limit) {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "insufficient funds"})
		return
	}

	acc, err = util.GetAccount(req.RecipientID, req.RecipientAccountID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ok, err = util.ValidateAccount(acc)
	if !ok {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
	}

	tr := model.Transaction{
		ID:          uuid.NewString(),
		SenderID:    req.SenderAccountID,
		RecipientID: req.RecipientAccountID,
		Amount:      req.Amount,
		Date:        time.Now(),
		Type:        req.Type,
	}

	err = receiver.DB.Create(tr)
	if err != nil {
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
//	@failure		400			{object}	response.ErrorResponse
//	@failure		500			{object}	response.ErrorResponse
//	@router			/transaction/{accountID}/{type} [GET]
func (receiver TransactionController) GetAll(context *gin.Context) {
	accountID := context.Param("accountID")

	if !util.IsValidUUID(accountID) {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid account id"})
		return
	}

	t := context.Param("type")
	if !(t == "sender" || t == "recipient" || t == "all") {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid type, supported: 'sender', 'recipient', 'all'"})
		return
	}

	res, err := receiver.DB.GetAll(accountID, t)
	if err != nil {
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
//	@router			/transaction/{transactionID}/ [DELETE]
func (receiver TransactionController) Delete(context *gin.Context) {
	transactionID := context.Param("transactionID")

	if !util.IsValidUUID(transactionID) {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid transaction id"})
		return
	}

	err := receiver.DB.Delete(transactionID)
	if err != nil {
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
//	@router			/transactions/{accountID}/ [DELETE]
func (receiver TransactionController) DeleteForAccount(context *gin.Context) {
	accountID := context.Param("accountID")

	if !util.IsValidUUID(accountID) {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid accountID id"})
		return
	}

	err := receiver.DB.DeleteForAccount(accountID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.Status(http.StatusNoContent)
}
