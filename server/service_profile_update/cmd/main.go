package main

import (
	"log"

	_ "github.com/youknow2509/cio_verify_face/server/service_profile_update/docs"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/start"
)

// @title Face Profile Update Service API
// @version 1.0
// @description API for managing face profile updates and password resets
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	if err := start.StartService(); err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}
}
