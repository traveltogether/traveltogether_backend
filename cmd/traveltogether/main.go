package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/database"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver"
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

	if *debug {
		general.Log.SetLevel(logrus.DebugLevel)
	}

	database.OpenConnection(*databaseHostname, *databasePort, *databaseUser, *databasePassword, *databaseName,
		*databaseSSLMode)
	database.MustExec("CREATE TABLE IF NOT EXISTS users(" +
		"id BIGSERIAL PRIMARY KEY," +
		"username TEXT NOT NULL," +
		"mail TEXT NOT NULL," +
		"first_name TEXT," +
		"password TEXT NOT NULL," +
		"session_key VARCHAR(36)," +
		"profile_image TEXT," +
		"disabilities TEXT)")
	database.MustExec("CREATE TABLE IF NOT EXISTS journeys(" +
		"id BIGSERIAL PRIMARY KEY," +
		"user_id INTEGER NOT NULL," +
		"request BOOLEAN NOT NULL," +
		"offer BOOLEAN NOT NULL," +
		"start_lat_long VARCHAR(64) NOT NULL," +
		"end_lat_long VARCHAR(64) NOT NULL," +
		"approximate_start_lat_long VARCHAR(64) NOT NULL," +
		"approximate_end_lat_long VARCHAR(64) NOT NULL," +
		"start_address TEXT NOT NULL," +
		"end_address TEXT NOT NULL," +
		"approximate_start_address TEXT NOT NULL," +
		"approximate_end_address TEXT NOT NULL," +
		"time_value BIGINT NOT NULL," +
		"time_is_departure BOOLEAN NOT NULL," +
		"time_is_arrival BOOLEAN NOT NULL," +
		"open_for_requests BOOLEAN NOT NULL," +
		"pending_user_ids INTEGER[]," +
		"accepted_user_ids INTEGER[]," +
		"declined_user_ids INTEGER[]," +
		"cancelled_by_host BOOLEAN NOT NULL," +
		"cancelled_by_host_reason TEXT," +
		"cancelled_by_attendee_ids INTEGER[]," +
		"note TEXT)")
	database.MustExec("CREATE TABLE IF NOT EXISTS chat_rooms(" +
		"id BIGSERIAL PRIMARY KEY," +
		"participants INTEGER[]," +
		"group_chat BOOLEAN NOT NULL)")
	database.MustExec("CREATE TABLE IF NOT EXISTS chat_messages(" +
		"id BIGSERIAL PRIMARY KEY," +
		"chat_id INTEGER NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE," +
		"message TEXT NOT NULL," +
		"sender_id INTEGER NOT NULL," +
		"read_by INTEGER[]," +
		"time BIGINT NOT NULL)")

	webserver.Run(*webserverHostname, *webserverPort)
}
