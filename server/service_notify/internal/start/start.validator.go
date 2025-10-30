package start

import (
	"github.com/go-playground/validator/v10"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/global"
)

// init validator
func initValidator() error {
	global.Validator = validator.New()
	return nil
}