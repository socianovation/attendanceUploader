package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	zipFile("result.json", "result.zip")
	sendFile("http://localhost/test/upload.php", "result.zip")
}

func getConfig() {

}

func sendFile(url, file string) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("image", file)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}
	// Add the other fields
	if fw, err = w.CreateFormField("key"); err != nil {
		return
	}
	if _, err = fw.Write([]byte("KEY")); err != nil {
		return
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
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

func zipFile(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
