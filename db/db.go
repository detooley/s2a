package db

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"s2a.app/idm"
	"s2a.app/utils/structs"
	"s2a.app/utils/text"
)

func connectToDatabase() *sql.DB {
	// check environment variable DB_PASSWORD set with your database password
	if os.Getenv("DB_PASSWORD") == "" {
		panic("DB_PASSWORD environment variable is not set")
	}
	// connection string for PostgreSQL
	connStr := "user=s2a password=" + os.Getenv("DB_PASSWORD") + " dbname=s2a sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return conn
}

// create a new pulls table
func CreatePulls() {
	conn := connectToDatabase()
	defer conn.Close()
	// create SQL commands
	backupRedTable := "CREATE TABLE pulls" +
		string(time.Now().Format("20060102150405")) +
		" AS SELECT * FROM pulls;"
	dropRedTable := "DROP TABLE pulls;"
	newRedTable := "CREATE TABLE pulls(" +
		"ser_no INT PRIMARY KEY," +
		"mark VARCHAR(255)," +
		"country VARCHAR(255)," +
		"examiner VARCHAR(255)," +
		"pull_date TIMESTAMP DEFAULT NOW());"
	// execute sequential SQL commands
	data := [][]string{
		{backupRedTable, "Failed to backup table", "Successfully backed up table"},
		{dropRedTable, "Failed to drop table", "Successfully dropped table"},
		{newRedTable, "Failed to create new Table", "Successfully created new table"},
	}
	seqSql(data, conn)
}

// create new id manual table
func CreateIdManual() {
	conn := connectToDatabase()
	defer conn.Close()
	// create SQL commands
	backupIdTable := "CREATE TABLE id_manual" +
		string(time.Now().Format("20060102150405")) +
		" AS SELECT * FROM id_manual;"
	dropIdTable := "DROP TABLE id_manual;"
	newIdTable := "CREATE TABLE id_manual(" +
		"id INT PRIMARY KEY," +
		"id_tx VARCHAR(10)," +
		"class_id CHAR(3)," +
		"description_tx VARCHAR(1000)," +
		"notes TEXT," +
		"version CHAR(7)," +
		"tm5 CHAR(1)," +
		"begin_effective_date TIMESTAMP WITH TIME ZONE," +
		"status CHAR(1));"
	// execute sequential SQL commands
	data := [][]string{
		{backupIdTable, "Failed to backup table", "Successfully backed up table"},
		{dropIdTable, "Failed to drop table", "Successfully dropped table"},
		{newIdTable, "Failed to create new Table", "Successfully created new table"},
	}
	seqSql(data, conn)
	populateIdManual(conn)
}

// Search the id manual table
func SearchIdManual(idm structs.IdmPageData) []structs.IdManual {
	var results []structs.IdManual
	if idm.Search == "" {
		return results
	}
	conn := connectToDatabase()
	defer conn.Close()
	// create SQL commands
	sqlStart := "SELECT " + idManualFields + " FROM id_manual WHERE "
	sqlSearch := parseIdmSearch(idm.Search)
	sqlSort := parseOrderBy(idm.OrderBy)
	sql := sqlStart + sqlSearch + sqlSort
	fmt.Print(sql + "\n")
	// execute sequential SQL commands
	rows, err := conn.Query(sql)
	if err != nil {
		fmt.Printf("Error executing query: %s\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var id_tx string
		var class_id string
		var description_tx string
		var notes string
		var version string
		var tm5 string
		var begin_effective_date time.Time
		var status string

		err := rows.Scan(&id,
			&id_tx,
			&class_id,
			&description_tx,
			&notes,
			&version,
			&tm5,
			&begin_effective_date,
			&status,
		)

		if err != nil {
			fmt.Printf("Error scanning row: %s\n", err)
			continue
		}
		results = append(results, structs.IdManual{
			Id:                 id,
			IdTx:               id_tx,
			ClassId:            class_id,
			DescriptionTx:      description_tx,
			Notes:              notes,
			Version:            version,
			Tm5:                tm5,
			BeginEffectiveDate: begin_effective_date,
			Status:             status,
		})
	}
	return results
}

// Generic search function for the database
func SearchDb(sql string) []map[string]any {
	conn := connectToDatabase()
	defer conn.Close()
	sql = text.SqlSanitize(sql)
	rows, err := conn.Query(sql)
	if err != nil {
		fmt.Printf("Error executing query: %s\n", err)
	}
	defer rows.Close()
	results := convertRows(rows)
	return results
}

// execute sequential sql commands and handle errors
func seqSql(data [][]string, conn *sql.DB) string {
	for _, sql := range data {
		_, err := conn.Exec(sql[0])
		if err != nil {
			fmt.Printf("%s: %s\n", sql[1], err)
			continue
		} else {
			fmt.Println(sql[2])
		}
	}
	results := ""
	return results
}

func convertRows(rows *sql.Rows) []map[string]any {
	// Written by Visual Stuidio Code Copilot (don't totally understand it)
	// Get the column names
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("Error getting columns:", err)
		return nil
	}
	// Create a slice of maps to hold the results
	results := make([]map[string]any, 0)
	// Create a slice to hold the values for each row
	values := make([]any, len(columns))
	for i := range values {
		values[i] = new(any) // Use new(any) to get a pointer to an empty interface
	}
	// Iterate over the rows
	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		rowMap := make(map[string]any)
		for i, colName := range columns {
			rowMap[colName] = *(values[i].(*any)) // Dereference the pointer to get the value
		}
		results = append(results, rowMap)
	}
	return results
}

// populate id manual table
func populateIdManual(conn *sql.DB) {
	// loop through the 45 classes and A, B, and 200
	for i := 0; i < 49; i++ {
		cl := strconv.Itoa(i)
		switch {
		case cl == "46":
			cl = "200"
		case cl == "47":
			cl = "A"
		case cl == "48":
			cl = "B"
		}
		results := idm.GetEntries("class-num=" + cl)
		x, err := strconv.Atoi(fmt.Sprint(results["numFound"]))
		if err != nil {
			fmt.Println("Error converting numFound to int:", err)
		}
		// poorly formed JSON made me do this
		for i := 1; i < x+1; i++ {
			// conver the int to string
			ti := strconv.Itoa(i)
			// build the variables for the SQL statement
			id := results["docs"].(map[string]any)[ti].(map[string]any)["id"]
			id_tx := results["docs"].(map[string]any)[ti].(map[string]any)["id_tx"]
			class_id := results["docs"].(map[string]any)[ti].(map[string]any)["class_id"]
			// sanatize the description
			description_tx := text.SqlSanitize(fmt.Sprint(results["docs"].(map[string]any)[ti].(map[string]any)["description_tx"]))
			notes, okay := results["docs"].(map[string]any)[ti].(map[string]any)["notes"]
			// add a empty string if notes is missing or sanitize it if its not
			if !okay {
				notes = ""
			} else {
				notes = text.StripHTMLTags(text.SqlSanitize(fmt.Sprint(notes)))
			}
			version := results["docs"].(map[string]any)[ti].(map[string]any)["version"]
			tm5 := results["docs"].(map[string]any)[ti].(map[string]any)["TM5"]
			begin_effective_date := (results["docs"].(map[string]any)[strconv.Itoa(i)].(map[string]any)["begin_effective_date"])
			status := (results["docs"].(map[string]any)[strconv.Itoa(i)].(map[string]any)["status"])
			// build the SQL statement
			populateIds := fmt.Sprintf("INSERT INTO id_manual (id, id_tx, class_id, description_tx, notes, version, tm5, begin_effective_date, status) "+
				"VALUES (%s, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');",
				id, id_tx, class_id, description_tx, notes, version, tm5, begin_effective_date, status,
			)
			// fmt.Print(description_tx)
			// insert the entries and handle errors
			_, err := conn.Exec(populateIds)
			if err != nil {
				fmt.Printf("Error inserting '%s': %s\n", id_tx, err)
			}
		}
		fmt.Print("Class " + cl + " Complete\n")
	}
	createSearchVector := "ALTER TABLE id_manual ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english', coalesce(description_tx, '') || ' ' || coalesce(notes, ''))) STORED;"
	createFullTextIndex := "CREATE INDEX idx_search_vector ON id_manual USING GIN (search_vector);"
	data := [][]string{
		{createSearchVector, "Failed to create search vector", "Successfully created search vector"},
		{createFullTextIndex, "Failed to create full text index", "Successfully created full text index"},
	}
	seqSql(data, conn)
}

func parseIdmSearch(search string) string {
	paras := strings.Fields(search)
	var sql_search string
	var sql_bit string
	class := regexp.MustCompile(`\b(00[1-9]|0[1-3]\d|04[0-5])\b`)
	punc := regexp.MustCompile(`[{}()\[\],-.:]`)
	// Parse Full Text Search
	if strings.Contains(search, "/f") {
		search = strings.ReplaceAll(search, "/f", "")
		search = strings.Join(strings.Fields(punc.ReplaceAllString(search, "")), " & ")
		sql_search = "search_vector @@ to_tsquery('english', '" + search + "')"
	} else {
		// Parse Manual Search
		for _, para := range paras {
			// Remove relevent punctuation from para
			para = punc.ReplaceAllString(para, "")
			// Check if search is a class id
			if class.MatchString(para) {
				sql_bit = "class_id ILIKE '" + para + "'"
				// Assume search of description
			} else {
				runes := []rune(para)
				parax := para
				paray := para
				if runes[0] != 42 {
					parax = "% " + parax
				} else {
					parax = "%" + parax
					paray = "%" + paray
				}
				if runes[len(runes)-1] != 42 {
					parax = parax + " %"
				} else {
					parax = parax + "%"
					paray = paray + "%"
				}
				// Get rid of asterisks for SQL ILIKE
				parax = strings.ReplaceAll(parax, "*", "")
				paray = strings.ReplaceAll(paray, "*", "")
				paraz := strings.ReplaceAll(para, "*", "")
				//sql_bit = "(description_tx ILIKE '" + parax + "' OR description_tx ILIKE '" + paray + "' OR description_tx ILIKE '" + paraz + " %' OR description_tx ILIKE '% " + paraz + "')"
				sql_bit = "(regexp_replace(description_tx, '[[:punct:]]', '', 'g') ILIKE '" +
					parax + "' OR regexp_replace(description_tx, '[[:punct:]]', '', 'g') ILIKE '" +
					paray + "' OR regexp_replace(description_tx, '[[:punct:]]', '', 'g') ILIKE '" +
					paraz + " %' OR regexp_replace(description_tx, '[[:punct:]]', '', 'g') ILIKE '% " +
					paraz + "')"
			}
			sql_search += sql_bit + " AND "
		}
	}
	sql_search = strings.TrimSuffix(sql_search, " AND ")
	//fmt.Print(sql_search)
	return sql_search
}

func parseOrderBy(orderBy structs.OrderBy) string {
	var sql string
	var dir string
	if orderBy.Direction == "desc" {
		dir = "DESC"
	} else {
		dir = "ASC"
	}
	switch orderBy.Element {
	case "class_id", "version", "tm5", "begin_effective_date", "status":
		sql = " ORDER BY " + orderBy.Element + " " + dir + ", description_tx " + dir + ";"
	case "description_tx":
		sql = " ORDER BY description_tx " + dir + ";"
	default:
		sql = " ORDER BY class_id " + dir + ", description_tx " + dir + ";"
	}
	return sql
}

// Define string with id_manual fields to pull
var idManualFields = "id, id_tx, class_id, description_tx, notes, version, tm5, begin_effective_date, status"
