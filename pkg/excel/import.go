package excel

import (
	"errors"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"

	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/utils"
)

type Importer struct {
	file           *excelize.File
	Sheets         []string
	Headers        []Header
	headerMap      map[int]string
	CheckHeader    bool
	StartRow       int
	EndRow         int
	rows           []map[string]string
	DateFormat     string
	DatetimeFormat string
	SkipFunc       ImportSkipFunc
	Filter         map[string]ImportFilter
	count          int
}

type ImportSkipFunc func(sheet string, rowNum int, rowMap map[string]string) bool
type ImportFilter func(content, sheet string, rowNum int, rowMap map[string]string) string

func NewImporter(headers []Header) *Importer {
	for headerIdx, header := range headers {
		if utils.IsEmpty(header.Column) {
			headers[headerIdx].Column = headerIdx + 1
		}
	}
	importer := &Importer{
		Headers:        headers,
		CheckHeader:    true,
		StartRow:       2,
		DateFormat:     config.GetCommonConfig().DateFormat,
		DatetimeFormat: config.GetCommonConfig().DatetimeFormat,
		Filter:         make(map[string]ImportFilter),
	}
	return importer
}

func (e *Importer) WithSheets(sheets []string) *Importer {
	e.Sheets = sheets
	return e
}

func (e *Importer) WithStartRow(num int) *Importer {
	e.StartRow = num
	return e
}

func (e *Importer) WithEndRow(num int) *Importer {
	e.EndRow = num
	return e
}

func (e *Importer) WithDateFormat(format string) *Importer {
	e.DateFormat = format
	return e
}

func (e *Importer) WithDatetimeFormat(format string) *Importer {
	e.DatetimeFormat = format
	return e
}

func (e *Importer) WithoutCheckHeader() *Importer {
	e.CheckHeader = false
	return e
}

func (e *Importer) RegisterSkipFunc(fn ImportSkipFunc) *Importer {
	e.SkipFunc = fn
	return e
}

func (e *Importer) RegisterFilter(key string, fn ImportFilter) *Importer {
	e.Filter[key] = fn
	return e
}

func (e *Importer) checkHeader(sheet string) bool {
	for _, header := range e.Headers {
		cell, _ := excelize.CoordinatesToCellName(header.Column, 1)
		value, _ := e.file.GetCellValue(sheet, cell)
		if value != header.Title {
			return false
		}
	}
	return true
}

func (e *Importer) Import(file io.Reader) error {
	var err error
	e.file, err = excelize.OpenReader(file)
	if utils.HasErr(err) {
		return err
	}
	sheets := e.Sheets
	if utils.IsEmpty(sheets) {
		sheets = []string{e.file.GetSheetName(e.file.GetActiveSheetIndex())}
	}
	for _, sheet := range sheets {
		if e.CheckHeader && !e.checkHeader(sheet) {
			return errors.New("xlsx template is invalid")
		}
		err := e.getContents(sheet)
		if utils.HasErr(err) {
			return err
		}
	}
	return nil
}

func (e *Importer) GetRows() []map[string]string {
	return e.rows
}

func (e *Importer) GetCount() int {
	return e.count
}

func (e *Importer) getContents(sheet string) error {
	e.rows = []map[string]string{}
	e.count = 0
	rowNum := 1
	rows, err := e.file.Rows(sheet)
	if utils.HasErr(err) {
		return err
	}
	for rows.Next() {
		if !utils.IsEmpty(e.EndRow) && rowNum > e.EndRow {
			break
		}
		row, err := rows.Columns()
		if utils.HasErr(err) {
			return err
		}
		e.getRowMap(sheet, rowNum, row)
		rowNum++
	}
	return nil
}

func (e *Importer) getRowMap(sheet string, rowNum int, row []string) {
	if rowNum < e.StartRow {
		return
	}
	rowMap := make(map[string]string)
	for _, header := range e.Headers {
		if len(row) >= header.Column {
			rowMap[header.Key] = strings.TrimSpace(row[header.Column-1])
		} else {
			rowMap[header.Key] = ""
		}
		if header.Required && utils.IsEmpty(rowMap[header.Key]) {
			return
		}
	}
	e.formatRowMap(sheet, rowNum, rowMap)
	if utils.IsEmpty(e.SkipFunc) || !e.SkipFunc(sheet, rowNum, rowMap) {
		e.rows = append(e.rows, rowMap)
		e.count++
	}
}

func (e *Importer) formatRowMap(sheet string, rowNum int, rowMap map[string]string) {
	for key, value := range rowMap {
		if filter, ok := e.Filter[key]; ok {
			rowMap[key] = filter(value, sheet, rowNum, rowMap)
		}
	}
}

func (e *Importer) BindData(data interface{}) {
	utils.ArrayMapStringToStruct(e.rows, data, e.DatetimeFormat, e.DateFormat)
}
