package main

import (
    "fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"strings"
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

func main() {
	b, err := ioutil.ReadFile("config.txt") // just pass the file name
    if err != nil {
        fmt.Print(err)
    }
	str := string(b) // convert content to a 'string'
	strs := strings.Split(str, "\n")
	
	//Initialize DB
	db, err := sql.Open("mysql", "root@127.0.0.1/dbname")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT squareNumber FROM squarenum WHERE number = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()
	
	excelFileName := strs[0]
    xlFile, err := xlsx.OpenFile(excelFileName)
    if err != nil {
		fmt.Printf("%s\n", err)
	}
    for _, sheet := range xlFile.Sheets {
        for _, row := range sheet.Rows {
            for _, cell := range row.Cells {
                text := cell.String()
                fmt.Printf("%s\n", text)
            }
        }
    }
}