// The go.mod file defines the module's properties, dependencies, and the Go version being used.
// It specifies the module's path, required dependencies, and their versions, ensuring consistent builds across different environments.

module tiktok-webapp // Declares the module name, which is typically the path to the repository.

go 1.23.0 // Specifies the version of Go to be used.

require (
	filippo.io/edwards25519 v1.1.0 // indirect dependency: elliptic curve cryptography package.
	github.com/go-sql-driver/mysql v1.8.1 // Direct dependency: MySQL driver for Go, used for database interactions.
)