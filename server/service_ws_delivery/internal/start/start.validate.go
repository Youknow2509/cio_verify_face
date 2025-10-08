package start

import (
	"github.com/go-playground/validator/v10"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
)


func initValidator() error{
	global.Validate = validator.New()
	return nil
}