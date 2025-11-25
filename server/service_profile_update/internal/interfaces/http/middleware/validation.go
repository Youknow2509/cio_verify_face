package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/interfaces/response"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateDTO validates the DTO struct using validator tags
func ValidateDTO(dto interface{}) error {
	return validate.Struct(dto)
}

// BindAndValidate binds request data and validates the DTO
func BindAndValidate(c *gin.Context, dto interface{}, bindType string) bool {
	var err error

	switch bindType {
	case "json":
		err = c.ShouldBindJSON(dto)
	case "query":
		err = c.ShouldBindQuery(dto)
	case "uri":
		err = c.ShouldBindUri(dto)
	case "form":
		err = c.ShouldBind(dto)
	default:
		err = c.ShouldBind(dto)
	}

	if err != nil {
		response.ValidationError(c, err.Error())
		return false
	}

	return true
}

// ValidateStruct validates a struct and returns error via response if validation fails
func ValidateStruct(c *gin.Context, dto interface{}) bool {
	if err := ValidateDTO(dto); err != nil {
		response.ValidationError(c, err.Error())
		return false
	}
	return true
}
