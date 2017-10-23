package helpers

import (
	. "attendanceuploader/models"
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-adodb"
)

func GetAttendanceData(DatabaseName string, AuthToken string) {

	//Define Variables
	var userid string
	var checktime string
	var checktype string
	var sensorid string
	var attendanceClocks []AttendanceClock

	flag.Parse()

	//Check if the Database MS Access File available as configured
	if _, err := os.Stat(DatabaseName); err != nil {
		check(err)
	}

	//Connect the MS Access DB using ADODB Jet OLEDB
	conn, err := sql.Open("adodb", "Provider=Microsoft.Jet.OLEDB.4.0;Data Source="+DatabaseName+";")
	check(err)

	//Get Number of row(s) found
	rows, err := conn.Query("SELECT COUNT(USERID) FROM CHECKINOUT")
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		check(err)
	}

	//Get Result of queried DB
	rows, err = conn.Query("SELECT USERID, CHECKTIME, CHECKTYPE, SENSORID FROM CHECKINOUT")
	check(err)

	//Store into Our Predefined Variable
	f, err := os.Create(AuthToken + ".json")
	check(err)
	w := bufio.NewWriter(f)
	defer f.Close()

	for rows.Next() {
		rows.Scan(&userid, &checktime, &checktype, &sensorid)
		if userid == "" {
			continue
		}

		if checktype == "I" {
			checktype = "IN"
		} else {
			checktype = "OUT"
		}

		attendanceClocks = append(attendanceClocks, AttendanceClock{userid, checktime, checktype, sensorid})
	}
	output, err := json.Marshal(attendanceClocks)
	check(err)

	fmt.Fprintf(w, string(output))
	w.Flush()
	err = rows.Err()
	check(err)

	defer rows.Close()
}
