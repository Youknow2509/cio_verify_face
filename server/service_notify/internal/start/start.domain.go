package start

import (
	infraConn "github.com/youknow2509/cio_verify_face/server/service_notify/internal/infrastructure/conn"
)

func initDomain() error {
	// ============================================
	// 			Get client connection
	// ============================================
	postgres, err := infraConn.GetPostgresqlClient()
	if err != nil {
		return err
	}
	// ============================================
	// 			Initialize domain components
	// ============================================
	_ = postgres // TODO: use postgres connection here

	// v.v
	return nil
}
