package util

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tealeg/xlsx"
)

var xlsxPath = flag.String("f", "", "Path to an XLSX file")
var csvpath = flag.String("r", "", "Path to store csv file")
var sheetIndex = flag.Int("i", 0, "Index of sheet to convert, zero based")
var delimiter = flag.String("d", ";", "Delimiter to use between fields")

type outputer func(s string)

func generateCSVFromXLSXFile(excelFileName string, csvFileName string, sheetIndex int, outputf outputer) (string, error) {

	xlFile, error := xlsx.OpenFile(excelFileName)
	if error != nil {
		return "", error
	}

	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		return "", errors.New("This XLSX file contains no sheets.")
	case sheetIndex >= sheetLen:
		return "", fmt.Errorf("No sheet %d available, please select a sheet between 0 and %d\n", sheetIndex, sheetLen-1)
	}

	file, err := os.Create(csvFileName)
	if err != nil {
		log.Fatal("can not create file")
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	sheet := xlFile.Sheets[sheetIndex]
	for _, row := range sheet.Rows {
		var vals []string
		// line := ""
		if row != nil {
			for _, cell := range row.Cells {
				str, err := cell.FormattedValue()
				if err != nil {
					vals = append(vals, err.Error())
				}
				vals = append(vals, fmt.Sprintf("%q", str))
			}
			// outputf(strings.Join(vals, *delimiter) + "\n")
			// line = strings.Join(vals, *delimiter)
		}
		err := writer.Write(vals)
		if err != nil {
			log.Fatal("can not write")
		}
	}
	totalString := ""
	return totalString, nil
}

// func main() {
// 	flag.Parse()
// 	if len(os.Args) < 3 {
// 		flag.PrintDefaults()
// 		return
// 	}
// 	flag.Parse()
// 	printer := func(s string) { fmt.Printf("%s", s) }
// 	if totalString, err := generateCSVFromXLSXFile(*xlsxPath, *csvpath, *sheetIndex, printer); err != nil {
// 		fmt.Println(totalString)
// 	}
// }
