package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorDetail struct {
	Field string `json:"field,omitempty"`
	Issue string `json:"issue,omitempty"`
}

type ErrorBody struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type Envelope struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  *ErrorBody  `json:"error,omitempty"`
}

func OK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Envelope{Status: "ok", Data: data})
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Envelope{Status: "ok", Data: data})
}

func Fail(c echo.Context, httpStatus int, code, message string, details ...ErrorDetail) error {
	b := &ErrorBody{Code: code, Message: message}
	if len(details) > 0 {
		b.Details = details
	}
	return c.JSON(httpStatus, Envelope{Status: "error", Error: b})
}
