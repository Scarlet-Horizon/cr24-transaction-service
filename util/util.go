package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"io"
	"log"
	"main/model"
	"main/response"
	"net/http"
	"os"
	"strings"
	"time"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func GetAccount(accountID, token, correlation string) (model.Account, error) {
	req, err := http.NewRequest(http.MethodGet, "http://account-api:8080/api/v1/account/"+accountID, nil)
	if err != nil {
		return model.Account{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Correlation", correlation)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return model.Account{}, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return model.Account{}, err
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			log.Printf("Close error: %s\n", err)
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return model.Account{}, errors.New("error: " + string(data))
	}

	var acc model.Account
	if err := json.Unmarshal(data, &acc); err != nil {
		return model.Account{}, err
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

func ValidateToken(context *gin.Context) {
	token := context.GetHeader("Authorization")
	if token == "" {
		context.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "unauthorized"})
		context.Abort()
		return
	}

	values := strings.Split(token, "Bearer ")
	if len(values) != 2 {
		context.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "token is not set properly"})
		context.Abort()
		return
	}

	token = values[1]

	to, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		context.Abort()
		return
	}

	if !to.Valid {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid token"})
		context.Abort()
		return
	}

	if claims, ok := to.Claims.(jwt.MapClaims); ok {
		if claims["sub"] == "" {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid id"})
			context.Abort()
			return
		}

		if claims["iat"] == "" || claims["exp"] == "" {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "iat or exp not set"})
			context.Abort()
			return
		}

		tokenIat := time.Unix(int64(claims["iat"].(float64)), 0)
		if tokenIat.After(time.Now()) {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "iat can't be in the future"})
			context.Abort()
			return
		}

		tokenExp := time.Unix(int64(claims["exp"].(float64)), 0)
		if tokenExp.Before(time.Now()) {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "expired token"})
			context.Abort()
			return
		}

		context.Set("ID", claims["sub"])
		context.Set("token", token)
		context.Next()
		return
	}
	context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid token"})
}

func CORS(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Credentials", "true")
	context.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, Origin, Accept, Cache-Control")
	context.Header("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PATCH, DELETE")
	context.Header("Access-Control-Max-Age", "86400")

	if context.Request.Method == http.MethodOptions {
		context.AbortWithStatus(http.StatusOK)
		return
	}
	context.Next()
}
