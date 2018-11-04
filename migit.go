package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xwb1989/sqlparser"
	"strings"
)

func main() {
	fmt.Printf("migit\n")

	db, err := sql.Open("mysql", "root@/referenceum")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//rows, err := db.Query("SHOW CREATE TABLE accounts")
	//rows, err := db.Query("SHOW TABLES")
	//rows, err := db.Query("SELECT * FROM accounts")

	var rows = QuerySingleColumn(db, "SHOW TABLES", 0)
	fmt.Printf("%s\n", strings.Join(rows, ","))

	ddl := QuerySingleColumn(db, "SHOW CREATE TABLE accounts", 1)[0]
	stmt, err := sqlparser.Parse(ddl)
	_ = stmt
	if err != nil {
		panic(err.Error())
	}

	// Otherwise do something with stmt
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		_ = stmt
	case *sqlparser.DDL:
		_ = stmt
		//switch action := stmt.Action.(type) {
		//}
	}
}

func QuerySingleColumn(db *sql.DB, statement string, idx int) []string {
	rows, err := db.Query(statement)
	if err != nil {
		panic(err.Error())
	}
	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	data := make([]string, 0)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	fmt.Printf("%s\n", strings.Join(columns, ","))
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}
		for i, value := range values {
			if i == idx {
				data = append(data, string(value))
			}
		}
	}
	return data
}
