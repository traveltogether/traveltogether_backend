package main

import (
	"flag"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/database"
)

func main() {
	databaseHostname := flag.String("db-hostname", "127.0.0.1", "hostname of the PostgreSQL instance")
	databasePort := flag.Int("db-port", 5432, "port of the PostgreSQL instance")
	databaseUser := flag.String("db-user", "traveltogether", "name of the PostgreSQL instance user")
	databasePassword := flag.String("db-password", "tr4v3lt0g3th3r!",
		"password of the PostgreSQL instance password")
	databaseName := flag.String("db-name", "traveltogether",
		"name of the PostgreSQL instance database")
	databaseSSLMode := flag.String("db-ssl", "disable",
		"enable/disable SSL connection to the PostgreSQL instance "+
			"(see PostgreSQL documentation of specific values to enable)")
	webserverHostname := flag.String("web-hostname", "127.0.0.1", "ip to bind the webserver to")
	webserverPort := flag.Int("web-port", 4269, "port to bind the webserver to")
	debug := flag.Bool("debug", false, "enable/disable debug logging")

	flag.Parse()

	database.OpenConnection(*databaseHostname, *databasePort, *databaseUser, *databasePassword, *databaseName)
	database.MustExec("CREATE TABLE IF NOT EXISTS users(" +
		"id BIGSERIAL PRIMARY KEY," +
		"name TEXT NOT NULL," +
		"mail TEXT NOT NULL," +
		"password TEXT NOT NULL," +
		"session_key VARCHAR(36))")
}
