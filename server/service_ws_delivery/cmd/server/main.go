package main

import (
	"fmt"
	"os"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/start"
)

// @title           Swagger Chat API
// @version         1.0
// @description     This is a sample server for a chat application.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/Youknow2509/
// @contact.email  lytranvinh.work@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	err := start.StartService()
	if err != nil {
		fmt.Printf("Service start failed: %+v\n", err)
		os.Exit(1)
	}
}
