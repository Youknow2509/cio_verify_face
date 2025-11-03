package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    int         `json:"code"`    // Ma status code
	Message string      `json:"message"` // Thong bao loi
	Data    interface{} `json:"data"`    // Du lieu duoc return
}

type ErrResponseData struct {
	Code   int         `json:"code"`   // Ma status code
	Error  string      `json:"error"`  //
	Detail interface{} `json:"detail"` // Thong bao loi
}

// func success response
func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: msg[code],
		Data:    data,
	})
}

// func bad response
func BadRequestResponse(c *gin.Context, code int, message string) {
	if message == "" {
		message = msg[code]
	}
	c.JSON(http.StatusBadRequest, ErrResponseData{
		Code:   code,
		Error:  message,
		Detail: nil,
	})
}

// func error response
func ErrorResponse(c *gin.Context, code int, message string) {
	if message == "" {
		message = msg[code]
	}
	c.JSON(http.StatusInternalServerError, ErrResponseData{
		Code:   code,
		Error:  message,
		Detail: nil,
	})
}

// func not found response
func NotFoundResponse(c *gin.Context, code int, message string) {
	if message == "" {
		message = msg[code]
	}
	c.JSON(http.StatusNotFound, ErrResponseData{
		Code:   code,
		Error:  message,
		Detail: nil,
	})
}

// func unauthorized response
func UnauthorizedResponse(c *gin.Context, code int, message string) {
	if message == "" {
		message = msg[code]
	}
	c.JSON(http.StatusUnauthorized, ErrResponseData{
		Code:   code,
		Error:  message,
		Detail: nil,
	})
}

// func forbidden response
func ForbiddenResponse(c *gin.Context, code int, message string) {
	if message == "" {
		message = msg[code]
	}
	c.JSON(http.StatusForbidden, ErrResponseData{
		Code:   code,
		Error:  message,
		Detail: nil,
	})
}

// func too many requests response
func TooManyRequestsResponse(c *gin.Context, code int, message string) {
	if message == "" {
		message = msg[code]
	}
	c.JSON(http.StatusTooManyRequests, ErrResponseData{
		Code:   code,
		Error:  message,
		Detail: nil,
	})
}
