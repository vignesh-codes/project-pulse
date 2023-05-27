package constants

import "fmt"

var (
	ENV                 = GetEnvString("ENVIRONMENT", "dev")
	SERVICE_NAME        = fmt.Sprintf("%s:%s", "PULSE_SERVICE", GetEnvString("ENVIRONMENT", "dev"))
	PORT         string = GetEnvString("PORT", ":8080")
	REGION       string = GetEnvString("REGION", "xx")
	JWT_SECRET   []byte = []byte(GetEnvString("JWT_SECRET", "your-very-secret-key"))
)

// psql -h localhost -d logdb -U logman
var (
	REDIS_SERVER    string = GetEnvString("REDIS_SERVER", "localhost:6379")
	POSTGRESDB_HOST string = GetEnvString("POSTGRESDB_HOST", "localhost")
	POSTGRESDB_DB   string = GetEnvString("POSTGRESDB_DB", "postgres")
	POSTGRESDB_PORT string = GetEnvString("POSTGRESDB_PORT", "4432")
	POSTGRESDB_USER string = GetEnvString("POSTGRESDB_USER", "user")
	POSTGRESDB_PWD  string = GetEnvString("POSTGRESDB_PWD", "pass")
)

var (
	SERVICE_ALPHA_INTERNAL_CLIENT_ID_SECRET map[string]string = map[string]string{"id": "secret"}
)

// has to be loaded from a db if used in production
var (
	AVAILABLE_SERVICES map[string]string = map[string]string{
		"id:secret":   "service_1",
		"id2:secret2": "service_2",
	}
)
