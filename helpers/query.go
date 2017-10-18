package helpers

import (
	. "attendanceuploader/models"
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func GetData() {
	flag.Parse()

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
