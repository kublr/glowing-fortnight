package cmd

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	sf "github.com/snowflakedb/gosnowflake"
	"github.com/spf13/cobra"
)

const RECORDS_LIMIT = 15

var rootCmd = &cobra.Command{
	Use:   "snowflake-poc",
	Short: "snowflake-poc application",
	Run: func(cmd *cobra.Command, args []string) {
		// setup db connection pool
		err := setupDB()
		if err != nil {
			log.Fatal(err)
		}

		// register http(s) handler
		http.HandleFunc("/", handler)

		// start serving requests
		log.Fatal(http.ListenAndServe(":4080", nil))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//go:embed templates/page.html
var pageTemplateStr string
var pageTemplate *template.Template

func init() {
	var err error
	pageTemplate, err = template.New("page").Parse(pageTemplateStr)
	if err != nil {
		panic(err)
	}
}

var db *sql.DB

func setupDB() error {
	// get DSN from env var
	dsn := os.Getenv("SNOWFLAKE_DSN")
	if dsn == "" {
		return errors.New("SNOWFLAKE_DSN env var is not specified")
	}

	// parse DSN
	cfg, err := sf.ParseDSN(dsn)
	if err != nil {
		return fmt.Errorf("SNOWFLAKE_DSN cannot be parsed: %w", err)
	}

	// use env vars to override username/password if specified
	username := os.Getenv("SNOWFLAKE_USERNAME")
	password := os.Getenv("SNOWFLAKE_PASSWORD")
	if username != "" {
		cfg.User = username
		cfg.Password = password
	}

	databaseName := os.Getenv("SNOWFLAKE_DATABASE")
	if databaseName != "" {
		cfg.Database = databaseName
	}

	schema := os.Getenv("SNOWFLAKE_SCHEMA")
	if schema != "" {
		cfg.Schema = schema
	}

	// reassemble DSN
	dsn, err = sf.DSN(cfg)
	if err != nil {
		return fmt.Errorf("snowflake DSN cannot be generated: %w", err)
	}

	// open DB connection pool
	db, err = sql.Open("snowflake", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to %q, err: %w", dsn, err)
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	// process healthz requests
	if r.URL.Path == "/healthz" {
		fmt.Fprint(w, "ok")
		return
	}

	// select data
	data, total, err := getData(r.Context())

	// render data
	pageTemplate.Execute(w, map[string]any{
		"Data":     data,
		"Total":    total,
		"PageSize": RECORDS_LIMIT,
		"Error":    err,
	})
}

type Record struct {
	Name string
}

func getData(ctx context.Context) ([]Record, int, error) {
	// context
	ctx, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()

	// connection
	con, err := db.Conn(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get connection : %w", err)
	}
	defer con.Close()

	// query
	query := "SELECT C_NAME FROM CUSTOMER ORDER BY C_NAME LIMIT ?"
	rows, err := con.QueryContext(ctx, query, RECORDS_LIMIT)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to run a query %q: %w", query, err)
	}
	defer rows.Close()

	// process query results
	res := make([]Record, 0, RECORDS_LIMIT)
	var v Record
	for rows.Next() {
		err = rows.Scan(&v.Name)
		if err != nil {
			return nil, 0, fmt.Errorf("error reading DB query result: %w", err)
		}
		res = append(res, v)
	}
	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("DB query result set error: %w", rows.Err())
	}

	// total count query
	query2 := "SELECT COUNT(C_NAME) FROM CUSTOMER"
	row2 := con.QueryRowContext(ctx, query2)
	var total int
	err = row2.Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to run a query %q: %w", query2, err)
	}

	return res, total, nil
}
