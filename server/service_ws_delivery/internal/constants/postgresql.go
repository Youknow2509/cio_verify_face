package constants

// ================================================
//
//	Constants for PostgreSQL
//
// ================================================
const (
	// sslmode	Eavesdropping protection	MITM protection	Statement
	// disable	No	No	I don't care about security, and I don't want to pay the overhead of encryption.
	// allow	Maybe	No	I don't care about security, but I will pay the overhead of encryption if the server insists on it.
	// prefer	Maybe	No	I don't care about encryption, but I wish to pay the overhead of encryption if the server supports it.
	// require	Yes	No	I want my data to be encrypted, and I accept the overhead. I trust that the network will make sure I always connect to the server I want.
	// verify-ca	Yes	Depends on CA policy	I want my data encrypted, and I accept the overhead. I want to be sure that I connect to a server that I trust.
	// verify-full	Yes	Yes	I want my data encrypted, and I accept the overhead. I want to be sure that I connect to a server I trust, and that it's the one I specify.
	POSTGRESQL_SSL_MODE_DISABLE     = "disable"
	POSTGRESQL_SSL_MODE_ALLOW       = "allow"
	POSTGRESQL_SSL_MODE_PREFER      = "prefer"
	POSTGRESQL_SSL_MODE_REQUIRE     = "require"
	POSTGRESQL_SSL_MODE_VERIFY_CA   = "verify-ca"
	POSTGRESQL_SSL_MODE_VERIFY_FULL = "verify-full"
)
