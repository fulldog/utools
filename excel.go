package utools

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/fulldog/utools/timex"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const fileExt = ".xlsx"

type excelTool struct{}

var ExcelTool = &excelTool{}

type ExportModel struct {
	Path        string
	ExcelType   int
	FileName    string
	ShowExtFlag int
	SheetDatas  []*SheetData //[]sheet1,sheet2
}
type SheetData struct {
	SheetName string
	Data      []interface{}
	Header    []interface{}
}

func (e excelTool) SpecialFileName(s string) string {
	for _, v := range []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "%", "+"} {
		s = strings.ReplaceAll(s, v, "")
	}
	return s
}

// CreateExcel 生成Excel
// m.ObjData 是一个二维数组，每个数组表示一个sheet
func (e excelTool) CreateExcel(m ExportModel) (fileRoute string, err error) {
	if m.SheetDatas == nil {
		err = errors.New("数据不能为空")
		return
	}

	fileNameOrg := e.SpecialFileName(m.FileName)
	node, _ := snowflake.NewNode(16)
	fileName := fmt.Sprintf("%s_%s%s", fileNameOrg, node.Generate().String(), fileExt)
	filePath := e.GetFilePath(m.Path)
	//filePathUrl := strings.ReplaceAll(filePath, "\\", pathRoute)
	//relativeUrl := fmt.Sprint(filePathUrl, pathRoute, fileName)
	//fileSize := 0

	//生成文件路径
	err = e.CreateDir(filePath)
	if err != nil {
		return
	}
	excel := excelize.NewFile()
	styleFloat64, _ := excel.NewStyle(&excelize.Style{
		NumFmt: 2,
	})
	defer func() {
		if err := excel.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	var deleteSheet1 bool
	var streamWriter *excelize.StreamWriter
	for obi, _ := range m.SheetDatas {
		if m.SheetDatas[obi] == nil {
			continue
		}
		sheet := TernaryOperation(m.SheetDatas[obi].SheetName == "", "Sheet"+strconv.Itoa(obi+1), m.SheetDatas[obi].SheetName)
		_, err = excel.NewSheet(sheet)
		if err != nil {
			return
		}
		if deleteSheet1 == false && sheet == "Sheet1" {
			deleteSheet1 = true
		}
		streamWriter, err = excel.NewStreamWriter(sheet)
		if err != nil {
			return
		}
		desc := reflect.TypeOf(m.SheetDatas[obi].Data[0]).Elem()
		var numField = desc.NumField()
		if m.SheetDatas[obi].Header == nil {
			for j := 0; j < numField; j++ {
				de := desc.Field(j).Tag.Get("desc")
				if de == "" || de == "-" {
					continue
				}
				m.SheetDatas[obi].Header = append(m.SheetDatas[obi].Header, de)
			}
		}

		err = streamWriter.SetRow("A1", m.SheetDatas[obi].Header)
		if err != nil || len(m.SheetDatas[obi].Header) == 0 {
			err = errors.WithMessagef(err, "创建表头错误")
			return
		}

		i := 2
		for _, d := range m.SheetDatas[obi].Data {
			ref := reflect.ValueOf(d).Elem()
			var dt []interface{}
			for k := 0; k < numField; k++ {
				de := desc.Field(k).Tag.Get("desc")
				if de == "" || de == "-" {
					continue
				}
				switch ref.Field(k).Type().String() {
				case "int", "int8", "int64", "float32", "float64":
					dt = append(dt, excelize.Cell{
						StyleID: styleFloat64,
						Formula: "",
						Value:   ref.Field(k).Float(),
					})
					//fmt.Println(btype.ToString(ref.Field(k).Float()))
				default:
					dt = append(dt, excelize.Cell{
						StyleID: styleFloat64,
						Formula: "",
						Value:   ref.Field(k).String(),
					})
				}
			}
			err = streamWriter.SetRow("A"+strconv.Itoa(i), dt)
			if err != nil {
				err = errors.WithMessagef(err, fmt.Sprintf("创建数据错误错误:%s,A%d", sheet, i))
				return
			}
			i++
		}
		err = streamWriter.Flush()
		if err != nil {
			err = errors.WithMessagef(err, "写入文件错误")
			return
		}
	}
	if !deleteSheet1 {
		_ = excel.DeleteSheet("Sheet1")
	}
	savePath := fmt.Sprint(filePath, fileName)
	err = excel.SaveAs(savePath)
	if err != nil {
		err = errors.WithMessagef(err, "生成文件错误")
		return
	}
	fileRoute = fmt.Sprint(strings.TrimLeft(savePath, "."))
	return
}

func (e excelTool) GetFilePath(path string) string {
	return fmt.Sprint(time.Now().Format(timex.DateYearMonth), path, "xlsx")
}

func (e excelTool) CreateDir(s string) (err error) {
	if b, _ := e.PathExists(s); !b {
		err = os.MkdirAll(s, os.ModePerm)
	}
	return
}

func (e excelTool) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
