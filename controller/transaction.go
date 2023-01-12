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
func (receiver TransactionController) Create(ctx *gin.Context) {
	var req request.TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if !util.IsValidUUID(req.SenderAccountID) {
		err := ctx.Error(errors.New("invalid sender id"))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if !util.IsValidUUID(req.RecipientAccountID) {
		err := ctx.Error(errors.New("invalid recipient id"))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if req.Amount < 1 {
		err := ctx.Error(errors.New("invalid amount, minimum is 1"))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	acc, err := util.GetAccount(req.SenderAccountID, ctx.MustGet("token").(string),
		ctx.GetString("Correlation"))
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ok, err := util.ValidateAccount(acc)
	if !ok {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
	}

	if acc.Amount-req.Amount < float64(-1*acc.Limit) {
		err := ctx.Error(errors.New("insufficient funds"))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
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

	err = receiver.DB.Create(tr, ctx)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, tr)
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
func (receiver TransactionController) GetAll(ctx *gin.Context) {
	accountID := ctx.Param("accountID")

	if !util.IsValidUUID(accountID) {
		err := ctx.Error(errors.New("invalid account id"))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	t := ctx.Param("type")
	if !(t == "sender" || t == "recipient" || t == "all") {
		err := ctx.Error(errors.New("invalid type, supported: 'sender', 'recipient', 'all'"))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	res, err := receiver.DB.GetAll(accountID, t, ctx)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	if len(res) == 0 {
		ctx.Status(http.StatusNoContent)
		return
	}
	ctx.JSON(http.StatusOK, res)
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
func (receiver TransactionController) Delete(ctx *gin.Context) {
	transactionID := ctx.Param("transactionID")

	if !util.IsValidUUID(transactionID) {
		err := ctx.Error(errors.New("invalid transaction id"))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	err := receiver.DB.Delete(transactionID, ctx)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
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
func (receiver TransactionController) DeleteForAccount(ctx *gin.Context) {
	accountID := ctx.Param("accountID")

	if !util.IsValidUUID(accountID) {
		err := ctx.Error(errors.New("invalid account id"))
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	err := receiver.DB.DeleteForAccount(accountID, ctx)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
