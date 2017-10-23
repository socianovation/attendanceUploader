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

	_ "github.com/mattn/go-adodb"
)

func GetData() {
	flag.Parse()

	flag.Parse()

	if _, err := os.Stat("att2000.mdb"); err != nil {
		fmt.Println("put here empty database named 'example.mdb'.")
		return
	}
	conn, err := sql.Open("adodb", "Provider=Microsoft.Jet.OLEDB.4.0;Data Source=att2000.mdb;")
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := conn.Query("SELECT COUNT(USERID) FROM CHECKINOUT")
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
	}

	rows, err = conn.Query("SELECT USERID, CHECKTIME, sn  FROM CHECKINOUT")
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
