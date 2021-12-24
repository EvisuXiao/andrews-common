package excel

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/utils"
)

type Exporter struct {
	file           *excelize.File
	Title          string
	Sheet          string
	Headers        []Header
	SkipHeader     bool
	Contents       []map[string]interface{}
	DateFormat     string
	DatetimeFormat string
	Filter         map[string]ExportFilter
	UniqueName     bool
	AutoSave       bool
	count          int
	actionTime     time.Time
}
type Header struct {
	Key      string
	Title    string
	Column   int
	Required bool
}
type ExportFilter func(content interface{}, rowNum int, rowMap map[string]interface{}) string

const defaultSheet = "Sheet1"

func NewExporter(title string, headers []Header) *Exporter {
	exporter := &Exporter{
		file:           excelize.NewFile(),
		Title:          title,
		Sheet:          defaultSheet,
		DateFormat:     config.GetCommonConfig().DateFormat,
		DatetimeFormat: config.GetCommonConfig().DatetimeFormat,
		Filter:         make(map[string]ExportFilter),
		UniqueName:     true,
		AutoSave:       true,
	}
	return exporter.WithHeaders(headers)
}

func (e *Exporter) WithHeaders(headers []Header) *Exporter {
	for headerIdx, header := range headers {
		if utils.IsEmpty(header.Column) {
			headers[headerIdx].Column = headerIdx + 1
		}
	}
	e.Headers = headers
	return e
}

func (e *Exporter) WithSheet(sheet string) *Exporter {
	e.Sheet = sheet
	return e
}

func (e *Exporter) WithDateFormat(format string) *Exporter {
	e.DateFormat = format
	return e
}

func (e *Exporter) WithDatetimeFormat(format string) *Exporter {
	e.DatetimeFormat = format
	return e
}

func (e *Exporter) WithoutUniqueName() *Exporter {
	e.UniqueName = false
	return e
}

func (e *Exporter) WithoutAutoSave() *Exporter {
	e.AutoSave = false
	return e
}

func (e *Exporter) WithoutHeaderRow() *Exporter {
	e.SkipHeader = true
	return e
}

func (e *Exporter) WithContents(contents interface{}) *Exporter {
	var newContents []map[string]interface{}
	rv := reflect.ValueOf(contents)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			row := rv.Index(i).Interface()
			if v, ok := row.(map[string]interface{}); ok {
				newContents = append(newContents, v)
			} else {
				newContents = append(newContents, utils.StructToMap(row))
			}
		}
		e.Contents = newContents
	}
	return e
}

func (e *Exporter) RegisterFilter(key string, fn ExportFilter) *Exporter {
	e.Filter[key] = fn
	return e
}

func (e *Exporter) Export() error {
	e.file.NewSheet(e.Sheet)
	e.file.DeleteSheet(defaultSheet)
	e.count = 0
	if !e.SkipHeader {
		e.putHeaders()
	}
	e.putContents()
	if e.AutoSave {
		return e.SaveFile()
	}
	return nil
}

func (e *Exporter) GetFilename() string {
	if utils.InvalidTime(e.actionTime) {
		return ""
	}
	suffix := ""
	if e.UniqueName {
		suffix = e.actionTime.Format(config.GetCommonConfig().SerialDatetimeFormat)
	}
	filename := fmt.Sprintf("%s%s.xlsx", e.Title, suffix)
	return url.QueryEscape(filename)
}

func (e *Exporter) GetCount() int {
	return e.count
}

func (e *Exporter) putHeaders() {
	for _, header := range e.Headers {
		cell, _ := excelize.CoordinatesToCellName(header.Column, 1)
		_ = e.file.SetCellStr(e.Sheet, cell, header.Title)
	}
}

func (e *Exporter) putContents() {
	startRow := 1
	if !e.SkipHeader {
		startRow++
	}
	var text interface{}
	for contentIdx, content := range e.Contents {
		for _, header := range e.Headers {
			if v, ok := content[header.Key]; ok {
				text = v
			} else {
				text = nil
			}
			rowNum := contentIdx + startRow
			cell, _ := excelize.CoordinatesToCellName(header.Column, rowNum)
			_ = e.file.SetCellStr(e.Sheet, cell, e.formatContent(header.Key, text, rowNum, content))
			e.count++
		}
	}
}

func (e *Exporter) formatContent(key string, content interface{}, rowNum int, rowMap map[string]interface{}) string {
	if filter, ok := e.Filter[key]; ok {
		return filter(content, rowNum, rowMap)
	}
	if v, ok := content.(time.Time); ok {
		if utils.InvalidTime(v) {
			return ""
		}
		dateFmt := utils.If(strings.HasSuffix(strings.ToLower(key), "date"), e.DateFormat, e.DatetimeFormat).(string)
		return v.Local().Format(dateFmt)
	}
	if v, ok := content.(bool); ok {
		return utils.If(v, "是", "否").(string)
	}
	return fmt.Sprint(content)
}

func (e *Exporter) SaveFile() error {
	e.file.SetActiveSheet(0)
	e.actionTime = time.Now()
	err := e.file.SaveAs(config.TempFilePath(e.GetFilename()))
	if utils.HasErr(err) {
		return err
	}
	e.file = nil
	return nil
}
