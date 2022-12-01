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

	if !util.IsValidUUID(req.SenderID) {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid sender id"})
		return
	}

	if !util.IsValidUUID(req.RecipientID) {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid recipient id"})
		return
	}

	if req.SenderID == req.RecipientID {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "can't transfer money to same accounts"})
		return
	}

	if req.Amount >= 0 && req.Amount < 1 || req.Amount <= 0 && req.Amount > -1 {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid amount, minimum is 1 or -1"})
		return
	}

	tr := model.Transaction{
		ID:          uuid.NewString(),
		SenderID:    req.SenderID,
		RecipientID: req.RecipientID,
		Amount:      req.Amount,
		Date:        time.Now(),
		Type:        req.Type,
	}

	err := receiver.DB.Create(tr)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.JSON(http.StatusCreated, tr)
}
