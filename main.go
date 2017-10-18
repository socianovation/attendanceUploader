package main

import (
	"attendanceuploader/helpers"
	"flag"
)

var (
	debug         = flag.Bool("debug", false, "enable debugging")
	password      = flag.String("password", "12345678", "the database password")
	port     *int = flag.Int("port", 1433, "the database port")
	server        = flag.String("server", "192.168.1.32", "the database server")
	user          = flag.String("user", "fingerprint", "the database user")
	database      = flag.String("d", "fingerprint", "fingerprint")
)

func main() {
	helpers.GetConfig()
	helpers.GetData()
	helpers.Compress("result.json", "result.zip")
	helpers.SendFile("http://localhost:8082/test/upload.php", "result.zip")
}
