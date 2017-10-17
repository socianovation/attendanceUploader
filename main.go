package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

type AttendanceClock struct {
	UserId    string `json:"user_id"`
	Checktime string `json:"checktime"`
	CompanyId string `json:"company_id"`
}

var (
	debug         = flag.Bool("debug", false, "enable debugging")
	password      = flag.String("password", "12345678", "the database password")
	port     *int = flag.Int("port", 1433, "the database port")
	server        = flag.String("server", "192.168.1.32", "the database server")
	user          = flag.String("user", "fingerprint", "the database user")
	database      = flag.String("d", "fingerprint", "fingerprint")
)

func main() {
	getConfig()
	getData()
}

func getConfig() {

}

func getData() {
	flag.Parse()

	if *debug {
		fmt.Printf(" password:%s\n", *password)
		fmt.Printf(" port:%d\n", *port)
		fmt.Printf(" server:%s\n", *server)
		fmt.Printf(" user:%s\n", *user)
	}

	conn, err := sql.Open("mssql", "server=localhost;user id=fingerprint;password=12345678;port=1433;database=fingerprint")
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT COUNT(USERID) FROM CHECKINOUT")
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
	}

	rows, err = conn.Query("SELECT USERID, CHECKTIME, sn FROM CHECKINOUT")
	if err != nil {
		log.Fatal("Failed query:", err.Error())
	}

	var userid string
	var checktime string
	var sn string
	var attendanceClocks []AttendanceClock
	f, err := os.Create("result.json")
	w := bufio.NewWriter(f)
	defer f.Close()

	for rows.Next() {
		rows.Scan(&userid, &checktime, &sn)
		if userid == "" {
			continue
		}
		attendanceClocks = append(attendanceClocks, AttendanceClock{userid, checktime, sn})
	}

	output, err := json.Marshal(attendanceClocks)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(output))

	w.Flush()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
}
