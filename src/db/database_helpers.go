package db

import(
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"strings"
)

var database *sql.DB

func init() {
	//get db file path dynamically from the server.exe running location
	ExePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error getting executing path")
	}
	dbPath := strings.Replace(ExePath, "\\src", "", 1)
	dbPath = dbPath + "\\data\\parkinglot.db"
	fmt.Println(dbPath)

	dbConn, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			fmt.Println(err)
		}
	database = dbConn
}

func Close(){
	fmt.Println("Closing db")
	database.Close()
	
}


func QueryForString(query string) string{// to return strings

	var returnStr string
	fmt.Println("Query: ", query)
	rows, err := database.Query(query)
	defer rows.Close()
	if err != nil {
		fmt.Println("error executing query", err)
		return ""
	} else {
		
		for rows.Next() {
			err = rows.Scan(&returnStr)
			fmt.Println(returnStr)
			if err != nil {
				fmt.Println(err)
			}
			
		}
	}
	return returnStr
}

func Query(query string) int {// to return integers

	//selectQuery:= "select company, floorNum from lot"
	var returnNum int
	fmt.Println("Query: ", query)
	rows, err := database.Query(query)
	defer rows.Close()
	if err != nil {
		fmt.Println("error executing query", err)
		return 0
	} else {
		
		for rows.Next() {
			err = rows.Scan(&returnNum)
			fmt.Println(returnNum)
			if err != nil {
				fmt.Println(err)
			}
			
		}
	}
	return returnNum
}


func ExecuteStatement(queryStmt string) {
	fmt.Println(queryStmt)
	_, err := database.Exec(queryStmt)
	if err != nil {
	fmt.Println(err)
	}

}

