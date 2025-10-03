package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/start"
)

// @title           Swagger Chat Service Auth REST API
// @version         1.0
// @description     This is a REST API for the Chat Service Auth.
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
	// init wait group
	global.WaitGroup = &sync.WaitGroup{}
	// start
	err := start.StartService()
	if err != nil {
		fmt.Printf("Service start failed: %+v\n", err)
		os.Exit(1)
	}
	global.WaitGroup.Wait()
}
