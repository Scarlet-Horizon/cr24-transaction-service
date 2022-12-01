package controller

import (
	"github.com/gin-gonic/gin"
	"main/response"
	"net/http"
)

func (receiver TransactionController) GetTypes(context *gin.Context) {
	types, err := receiver.DB.GetTypes()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	context.JSON(http.StatusOK, types)
}
