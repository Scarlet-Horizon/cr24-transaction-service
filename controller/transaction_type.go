package controller

import (
	"github.com/gin-gonic/gin"
	_ "main/model"
	"main/response"
	"net/http"
)

//	@description	Get all transaction types.
//	@summary		Get all transaction types
//	@accept			json
//	@produce		json
//	@tags			transaction
//	@success		200	{object}	[]model.TransactionType	"An array of TransactionType"
//	@failure		500	{object}	response.ErrorResponse
//	@router			/types [GET]
func (receiver TransactionController) GetTypes(ctx *gin.Context) {
	types, err := receiver.DB.GetTypes(ctx)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, types)
}
