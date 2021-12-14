package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"
)

var (
	db      *sql.DB
	logpath string
	debug   string
)

type Event struct {
	Id         int    `json:"id"`
	Content    string `json:"content"`
	Created_at string `json:"created_at"`
}

func main() {
	app := &cli.App{
		Name:  "api",
		Usage: "simple go api",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "pg-user",
				Aliases:  []string{"u"},
				Usage:    "Postgres username",
				EnvVars:  []string{"PG_USER"},
				FilePath: "/var/config/pg_user",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "pg-password",
				Aliases:  []string{"p"},
				Usage:    "Postgres password",
				EnvVars:  []string{"PG_PASSWORD"},
				FilePath: "/var/config/pg_password",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "pg-database",
				Aliases:  []string{"d"},
				Usage:    "Postgres database",
				EnvVars:  []string{"PG_DATABASE"},
				FilePath: "/var/config/pg_database",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "pg-host",
				Aliases:  []string{"c"},
				Usage:    "Postgres host",
				EnvVars:  []string{"PG_HOST"},
				FilePath: "/var/config/pg_host",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "pg-port",
				Aliases:  []string{"o"},
				Usage:    "Postgres port",
				Value:    "5432",
				EnvVars:  []string{"PG_PORT"},
				FilePath: "/var/config/pg_port",
			},
			&cli.StringFlag{
				Name:     "log-path",
				Aliases:  []string{"l"},
				Usage:    "Logfile path",
				Value:    "/var/log/logs.log",
				FilePath: "/var/config/log_path",
				EnvVars:  []string{"LOG_PATH"},
			},
			&cli.StringFlag{
				Name:     "debug",
				Aliases:  []string{"i"},
				Usage:    "Debug setting",
				Value:    "false",
				EnvVars:  []string{"DEBUG"},
				FilePath: "/var/config/debug",
			},
		},
		Action: func(c *cli.Context) error {
			logpath = c.String("log-path")
			debug = c.String("debug")
			connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", c.String("pg-user"), c.String("pg-database"), c.String("pg-password"), c.String("pg-host"), c.String("pg-port"))
			var err error
			db, err = sql.Open("postgres", connStr)
			if err != nil {
				CustomLogging("Database connection could not be established!")
				CustomLogging("Shutting down...")
				os.Exit(1)
			}
			defer db.Close()
			err = db.Ping()
			if err != nil {
				CustomLogging("Database is not reachable!")
				CustomLogging("Shutting down...")
				os.Exit(1)
			}
			CustomLogging("Successfully connected to database!")

			r := newRouter()
			CustomLogging("Start listening on port :8000")
			http.ListenAndServe(":8000", r)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/time", handler).Methods("GET")
	r.HandleFunc("/healthz", healthz).Methods("GET")
	return r
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := createTable(db)
	if err != nil {
		CustomLogging("Table creation has failed!")
		return
	}

	err = insertTimestamp(db, r.Host)
	if err != nil {
		CustomLogging("Insert timestamp failed!")
		return
	}

	events := queryTime(db)

	err = JSONResponse(http.StatusOK, events)(w)
	if err != nil {
		CustomLogging("failed to send OK response for reset")
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	err := db.Ping()
	if err != nil {
		CustomLogging("Database is not reachable!")
		CustomLogging("Shutting down...")
		os.Exit(1)
	}
	fmt.Fprintln(w, "All good")
}

func queryTime(db *sql.DB) []Event {
	CustomLogging("Querying Time on Database will be executed")

	rows, err := db.Query("SELECT * FROM logs")
	if err != nil {
		CustomLogging("Database query has failed!")
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		event := Event{}
		err = rows.Scan(&event.Id, &event.Content, &event.Created_at)
		if err != nil {
			CustomLogging(err.Error())
		}
		events = append(events, event)
	}

	return events
}

func CustomLogging(message string) {
	valBool, _ := strconv.ParseBool(debug)
	if valBool {
		fmt.Println(message)
	} else {
		WriteToLoggingFile(message)
	}
}

func WriteToLoggingFile(message string) {
	f, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(message + "\n"); err != nil {
		log.Println(err)
	}
}

func createTable(db *sql.DB) (err error) {
	const q = `
CREATE TABLE IF NOT EXISTS logs (
	id serial PRIMARY KEY,
	content varchar(255) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)`

	// Exec executes a query without returning any rows.
	_, err = db.Exec(q)

	return err
}

func insertTimestamp(db *sql.DB, host string) (err error) {
	_, err = db.Exec("INSERT INTO logs (content) VALUES ($1)", host)

	return err
}

type APIResponse func(http.ResponseWriter) error

func JSONResponse(code int, payload interface{}) APIResponse {
	return func(w http.ResponseWriter) error {
		response, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(code)
		_, err = w.Write(response)
		return err
	}
}
