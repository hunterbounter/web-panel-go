package database

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"strings"
	"time"
)

var (
	// SQL wrapper
	SQL *sqlx.DB
	// Database info
	databases Info
)

// Type is the type of database from a Type* constant
type Type string

const (
	// TypeBolt is BoltDB
	TypeBolt Type = "Bolt"
	// TypeMongoDB is MongoDB
	TypeMongoDB Type = "MongoDB"
	// TypeMySQL is MySQL
	TypeMySQL Type = "MySQL"

	TypePostgreSQL Type = "PostgreSQL"
)

// Info contains the database configurations
type Info struct {
	// Database type
	Type Type
	// MySQL info if used
	MySQL MySQLInfo
	// Bolt info if used
	Bolt BoltInfo
	// MongoDB info if used
	MongoDB MongoDBInfo

	PostgreSQL PostgreSQLInfo
}

// Postgre

type PostgreSQLInfo struct {
	Username  string
	Password  string
	Hostname  string
	Port      int
	Name      string
	Parameter string
}

// MySQLInfo is the details for the database connection
type MySQLInfo struct {
	Username  string
	Password  string
	Name      string
	Hostname  string
	Port      int
	Parameter string
}

// BoltInfo is the details for the database connection
type BoltInfo struct {
	Path string
}

// MongoDBInfo is the details for the database connection
type MongoDBInfo struct {
	URL      string
	Database string
}

// DSN returns the Data Source Name
func DSN(ci MySQLInfo) string {
	// Example: root:@tcp(localhost:3306)/test
	return ci.Username +
		":" +
		ci.Password +
		"@tcp(" +
		ci.Hostname +
		":" +
		fmt.Sprintf("%d", ci.Port) +
		")/" +
		ci.Name + ci.Parameter
}

func PostgreDSN(ci PostgreSQLInfo) string {
	// Example: root:@tcp(localhost:3306)/test
	var connectionString = "user=" + ci.Username +
		" password=" + ci.Password +
		" host=" + ci.Hostname +
		" port=" + fmt.Sprintf("%d", ci.Port) +
		" dbname=" + ci.Name + " " + ci.Parameter
	log.Println("Database Connection DSN -> " + connectionString)
	return connectionString
}

// Connect to the database
func Connect(d Info) {
	var err error

	// Store the config
	databases = d

	log.Println("Database Connecting -> " + d.Type)

	switch d.Type {
	case TypeMySQL:
		// Connect to MySQL
		if SQL, err = sqlx.Connect("mysql", DSN(d.MySQL)); err != nil {
			log.Println("SQL Driver Error", err)
		}

		// Check if is alive
		if err = SQL.Ping(); err != nil {
			log.Println("Database Error", err)
		}
	case TypePostgreSQL:
		// Connect to PostgreSQL
		log.Println("Connecting PostgreSQL..")
		if SQL, err = sqlx.Connect("postgres", PostgreDSN(d.PostgreSQL)); err != nil {
			log.Println("SQL Driver Error", err)
		}
		if err = SQL.Ping(); err != nil {
			log.Println("Database Error", err)
		}
		log.Println("Connected PostgreSQL..")
	default:
		log.Println("No registered database in config")
	}
}

func logQuery(query string, args ...interface{}) {
	queryWithValues := query
	for _, arg := range args {
		// Convert argument to string and replace in the query
		queryWithValues = strings.Replace(queryWithValues, "?", fmt.Sprintf("%v", arg), 1)
	}
	log.Println("Full SQL Query:", queryWithValues)
}

func Insert(tableName string, data map[string]interface{}, generateUUID bool) (interface{}, error) {
	var err error

	if generateUUID {
		query := "SELECT uuid_generate_v4()"

		var uuid string
		err = SQL.QueryRow(query).Scan(&uuid)
		if err != nil {
			log.Fatal(err)
		}

		data["id"] = uuid
	}

	// Map interface içinden gelen key'ler ile Query oluşturuluyor.
	var insertQuery = "INSERT INTO " + tableName + " ("
	var valuesData = " VALUES ("
	i := 0
	for key := range data {
		if i == 0 {
			insertQuery += key
			valuesData += fmt.Sprintf("'%v'", data[key])
		} else {
			insertQuery += "," + key
			valuesData += fmt.Sprintf(",'%v'", data[key])
		}
		i++
	}
	valuesData += ")"
	insertQuery += ")" + valuesData

	if i == 0 {
		return nil, errors.New("insert interface is null")
	}

	log.Printf("SQL : %s\n", insertQuery)

	insertResult, err := SQL.Exec(insertQuery)
	if err != nil {
		return nil, err
	}

	if generateUUID {
		// UUID oluşturulmuşsa, yeni oluşturulan UUID'yi döndür
		return data["id"], nil
	}

	lastInsertedID, _ := insertResult.LastInsertId()
	return lastInsertedID, nil
}

func RowArray(query string, queryData map[string]interface{}) (map[string]interface{}, error) {

	row := map[string]interface{}{}
	var rows *sqlx.Rows
	var err error

	if queryData != nil {
		rows, err = SQL.NamedQuery(query, queryData)
		defer rows.Close()
	} else {
		rows, err = SQL.Queryx(query)
		defer rows.Close()
	}

	if err != nil {
		fmt.Errorf(err.Error())
		return row, err
	}

	for rows.Next() {
		err = rows.MapScan(row)
		if err != nil {
			fmt.Errorf(err.Error())
			return row, err
		}

		for k, v := range row {
			if _, ok := v.([]byte); ok {
				// Damn. Byte. Arrays. Sqlx.
				row[k] = string(v.([]byte))
			}
		}
		break
	}

	return row, nil
}

func ResultArray(query string, query_data map[string]interface{}) ([]map[string]interface{}, error) {

	var rows *sqlx.Rows
	var err error
	result := []map[string]interface{}{}

	if len(query_data) > 0 {

		rows, err = SQL.NamedQuery(query, query_data)
		defer rows.Close()
		if err != nil {

			fmt.Errorf(err.Error())
			return result, err
		}
	} else {
		//https://github.com/go-sql-driver/mysql/issues/407
		//log.Println(query)
		smtp, err := SQL.Preparex(query)

		defer smtp.Close()
		if err == nil {
		}
		rows, err = smtp.Queryx()
		defer rows.Close()

		if err != nil {
			return nil, err
		}
	}

	for rows.Next() {

		row := make(map[string]interface{})
		err = rows.MapScan(row)
		if err != nil {
			fmt.Errorf(err.Error())
			return result, err
		}

		for k, v := range row {
			if _, ok := v.([]byte); ok {
				// Damn. Byte. Arrays. Sqlx.
				row[k] = string(v.([]byte))
			}
		}

		result = append(result, row)

	}

	return result, nil
}

// Big Query
func ExecuteSql(query string) ([]map[string]interface{}, error) {
	sqlString := strings.ReplaceAll(query, "\n", "\\n")
	//remove \n to blank
	sqlString = strings.ReplaceAll(sqlString, "\\n", " ")
	log.Println(sqlString)

	if strings.Contains(query, "truncate") {

		a, b := SQL.Exec(query)
		fmt.Println(a)
		fmt.Println(b)
		return nil, nil
	}
	return ResultArray(query, map[string]interface{}{})
}
func Select(tableName string, whereData map[string]interface{}) ([]map[string]interface{}, error) {

	var selectQuery = "SELECT * FROM " + tableName + " "

	if len(whereData) > 0 {
		selectQuery += "WHERE "
	}
	i := 0

	orderby := false
	limit := false
	for key := range whereData {
		// check is order by
		if key == "order_by" {
			orderby = true
			continue
		}
		// check is limit
		if key == "limit" {
			limit = true
			continue
		}
		if i == 0 {
			selectQuery += key + "= :" + key
		} else {
			selectQuery += " AND " + key + "= :" + key
		}
		i++
	}

	if orderby {
		selectQuery += " ORDER BY " + whereData["order_by"].(string)
	}
	if limit {
		selectQuery += " LIMIT " + whereData["limit"].(string)
	}

	// Log the SQL query with values, wrapping string values in single quotes
	//log.Printf("Executing SQL Query: %s\n", formatSQLQuery(selectQuery, whereData))

	return ResultArray(selectQuery, whereData)
}

// formatSQLQuery wraps string values in single quotes for logging
func formatSQLQuery(query string, data map[string]interface{}) string {
	for key, value := range data {
		switch value := value.(type) {
		case string:
			// Wrap string values in single quotes
			quotedValue := fmt.Sprintf("'%s'", value)
			query = strings.Replace(query, ":"+key, quotedValue, -1)
		}
	}
	return query
}

func Update(tableName string, updateData map[string]interface{}, whereData map[string]interface{}) (bool, error) {
	var err error
	newUpdateData := map[string]interface{}{}
	var updateQuery = "UPDATE " + tableName + " set "
	i := 0
	for key, value := range updateData {
		if i == 0 {
			updateQuery += key + "= :u_" + key
		} else {
			updateQuery += "," + key + "= :u_" + key
		}
		i++

		newUpdateData["u_"+key] = value

	}

	if i == 0 {
		return false, errors.New("update set interface is null")
	}

	i = 0
	whereQuery := ""
	if len(whereData) > 0 {
		whereQuery = " WHERE "

		for key, _ := range whereData {
			if i == 0 {
				whereQuery += key + "= :" + key
			} else {
				whereQuery += " AND " + key + "= :" + key
			}
			i++

		}
	}

	updateQuery += whereQuery

	updateParameters := mergeMaps(newUpdateData, whereData)

	_, err = SQL.NamedExec(updateQuery, updateParameters)

	if err != nil {
		log.Println("Update Query Error -- Query : ", updateQuery, updateParameters)
		log.Println("Error:" + err.Error())
		return false, err
	}

	return true, err
}

func Delete(tableName string, whereData map[string]interface{}) (bool, error) {
	var err error

	var deleteQuery = "DELETE FROM " + tableName + " "
	i := 0
	whereQuery := ""
	if len(whereData) > 0 {
		deleteQuery += " WHERE "

		for key, _ := range whereData {
			if i == 0 {
				deleteQuery += key + "= :" + key
			} else {
				deleteQuery += " AND " + key + "= :" + key
			}
			i++

		}
	}

	deleteQuery += whereQuery

	_, err = SQL.NamedExec(deleteQuery, whereData)

	if err != nil {
		return false, err
	}

	return true, err
}

// For Update Query
func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// ReadConfig returns the database information
func ReadConfig() Info {
	return databases
}

func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
