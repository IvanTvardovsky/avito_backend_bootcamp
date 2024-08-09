package http

import (
	"avito_bootcamp/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type errorResponseStruct struct {
	Message   string `json:"message" example:"что-то пошло не так"`
	RequestID string `json:"request_id" example:"g12ugs67gqw67yu12fgeuqwd"`
	Code      int    `json:"code" example:"12345"`
}

func generateRequestID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return "unknown"
	}
	return id.String()
}

func errorResponse(c *gin.Context, l logger.Interface, code int, msg string, err error) {
	requestID := generateRequestID()

	l.Error(fmt.Sprintf("Error occurred: %s %d: %s because of: %s", requestID, code, msg, err.Error()))

	// тут можно определить и проверить кастомные ошибки при желании, в api такое есть только у 500
	// я решил, что code пока будет просто дублировать http код ответа
	/*
		errorCode := 0
		switch code {
		case http.StatusInternalServerError:
			errorCode = 12345
		default:
			errorCode = 0
		} */

	c.JSON(code, errorResponseStruct{
		Message:   msg,
		RequestID: requestID,
		Code:      code,
	})
}
