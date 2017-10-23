package helpers

import (
	. "attendanceuploader/models"
	"io/ioutil"
	"strings"
)

func GetConfig() (config AttendanceConfig) {

	file, err := ioutil.ReadFile("config.txt")
	check(err)

	configString := strings.Split(string(file), ";")
	config = AttendanceConfig{strings.TrimSpace(strings.Split(configString[0], "=")[1]), strings.TrimSpace(strings.Split(configString[1], "=")[1]), strings.TrimSpace(strings.Split(configString[2], "=")[1])}
	return config
}
