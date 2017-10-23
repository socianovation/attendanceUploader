package main

import (
	"attendanceuploader/helpers"
	. "attendanceuploader/models"
)

var config AttendanceConfig

func main() {
	config = helpers.GetConfig()
	helpers.GetAttendanceData(config.DatabaseName, config.AuthToken)
	helpers.Compress(config.AuthToken+".json", config.AuthToken+".zip")
	helpers.SendFile(config.MiddlewareUrl, config.AuthToken+".zip")
}
