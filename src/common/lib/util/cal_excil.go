package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/tealeg/xlsx"
)

func MkExcelFile() (file *xlsx.File) {
	file = xlsx.NewFile()
	return file
}

func MkExcelSheet(name string, file *xlsx.File) (sheet *xlsx.Sheet, err error) {
	sheet, err = file.AddSheet(name)
	if err != nil {
		fmt.Printf(err.Error())
	}
	return
}

func FillExcelSheet(cells []string, sheet *xlsx.Sheet, is_header bool, is_even bool) (err error) {
	row := sheet.AddRow()
	for _, cellv := range cells {
		var style *xlsx.Style
		if is_header {
			style = mkStyle(16, "黑体", true, "00FFF2DF", true)
		} else if is_even {
			style = mkStyle(14, "Verdana", false, "FFFFFFFF", true)
		} else {
			style = mkStyle(14, "Verdana", false, "00FFF2DF", true)
		}
		cell := row.AddCell()
		cell.SetStyle(style)
		cell.Value = cellv

	}
	return
}

func mkStyle(fontSize int, fontName string, isbold bool, bgcolor string, is_bord bool) (style *xlsx.Style) {
	style = xlsx.NewStyle()
	style.Fill.PatternType = "solid"
	//	style.Fill.BgColor = bgcolor
	style.Fill.FgColor = bgcolor
	style.Font.Size = fontSize
	style.Font.Name = fontName
	style.Font.Bold = isbold
	style.Border = xlsx.Border{
		Left:   "thin",
		Right:  "thin",
		Top:    "thin",
		Bottom: "thin",
	}
	return
}

func GetExcelData(file *xlsx.File) (data []byte, err error) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	parts, err := file.MarshallParts()
	if err != nil {
		return
	}
	zipWriter := zip.NewWriter(buffer)
	for partName, part := range parts {
		w, err := zipWriter.Create(partName)
		if err != nil {
			return data, err
		}
		_, err = w.Write([]byte(part))
		if err != nil {
			return data, err
		}
	}
	err = zipWriter.Close()
	if err != nil {
		return data, err
	}
	data = buffer.Bytes()
	return

}
