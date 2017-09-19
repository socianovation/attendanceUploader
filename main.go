package main

import (
    "fmt"
	"github.com/tealeg/xlsx"
)

func main() {
	excelFileName := "C:/Users/yosis/Desktop/excel/absensi.xlsx"
    xlFile, err := xlsx.OpenFile(excelFileName)
    if err != nil {
		fmt.Printf("%s\n", err)
	}
    for _, sheet := range xlFile.Sheets {
		fmt.Printf("%s\n", "asdasd")
        for _, row := range sheet.Rows {
            for _, cell := range row.Cells {
                text := cell.String()
                fmt.Printf("%s\n", text)
            }
        }
    }
}
