//go:build integration
// +build integration

package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/snowflakedb/gosnowflake"
)

func TestSF(t *testing.T) {
	runQuery()

	t.Fail()
}
func runQuery() {
	// connect to DB
	dsn := "olegch:***@QKCBVSB-WDB07503/SNOWFLAKE_SAMPLE_DATA?schema=TPCH_SF1"
	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatalf("failed to connect. %v, err: %v", dsn, err)
	}
	defer db.Close()

	// query
	query := "SELECT C_NAME FROM CUSTOMER ORDER BY C_NAME LIMIT 10"
	rows, err := db.Query(query) // no cancel is allowed
	if err != nil {
		log.Fatalf("failed to run a query. %v, err: %v", query, err)
	}
	defer rows.Close()

	// process query results
	var v string
	for rows.Next() {
		err := rows.Scan(&v)
		if err != nil {
			log.Fatalf("failed to get result. err: %v", err)
			break
		}
		fmt.Println(v)
	}
	if rows.Err() != nil {
		fmt.Printf("ERROR: %v\n", rows.Err())
		return
	}
	fmt.Printf("Congrats! You have successfully run %v with Snowflake DB!\n", query)
}
