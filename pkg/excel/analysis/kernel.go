package excel_analysis

import (
	"fmt"
	"reflect"

	model "service-monitor/pkg/model"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
)

/* Тип предиката для сравнения каждой ячейки таблицы с определённым значением */
type PredicatType func(value string, row, column int) bool

/* Интерфейс для ядра анализа Excel-таблиц */
type IExcelAnalysis interface {
	CopyTo(filePath string) (bool, error)
	GetSheet() (*excelize.File, error)
	GetHeaderInfo() (model.HeaderInfoModel, error)
	GetIndexByValue(value string, sheet *excelize.File) (model.IndexCellModel, error)
	GetIndexNextRow(index model.IndexCellModel, sheet *excelize.File, predicat PredicatType) (model.IndexCellModel, error)
	GetIndexNextRowOffset(index model.IndexCellModel, sheet *excelize.File, offset int, predicat PredicatType) (model.IndexCellModel, error)
	GetLengthCells(index model.IndexCellModel, sheet *excelize.File, predicat PredicatType) int
	GetValueCells(data *model.HeaderInfoModel, sheet *excelize.File, index model.IndexCellModel, place string) model.IndexCellModel
}

/* Структура основного ядра анализа Excel-таблиц */
type ExcelAnalysis struct {
	Filepath string
}

/* Функция создания нового экземпляра ExAnalysisKernel */
func NewExccelAnalysis(filepath string) *ExcelAnalysis {
	return &ExcelAnalysis{
		Filepath: filepath,
	}
}

/* Получение длины ячеек начиная с определённого индекса */
func (k *ExcelAnalysis) GetLengthCells(index model.IndexCellModel, sheet *excelize.File, predicat PredicatType) int {
	if sheet == nil {
		var err error
		if sheet, err = k.GetSheet(); err != nil {
			return 0
		}
	}

	var result int
	result = 0

	rows, err := sheet.GetRows(viper.GetString("table.sheet"))
	if err != nil {
		return 0
	}

	for i := index.Row; i < len(rows); i++ {
		cell := rows[i][index.Column]
		if predicat(cell, i, index.Column) {
			result++
		} else {
			break
		}
	}

	return result
}

/* Получение подробной информации о ячейки в таблице по её значению */
func (k *ExcelAnalysis) GetIndexByValue(value string, sheet *excelize.File) (model.IndexCellModel, error) {
	if sheet == nil {
		var err error
		if sheet, err = k.GetSheet(); err != nil {
			return model.IndexCellModel{}, err
		}
	}

	var result model.IndexCellModel
	rows, err := sheet.GetRows(viper.GetString("table.sheet"))
	if err != nil {
		return model.IndexCellModel{}, err
	}

	for rInd, row := range rows {
		for cInd, cell := range row {
			if cell == value {
				pos, err := excelize.CoordinatesToCellName(rInd, cInd)
				if err != nil {
					return model.IndexCellModel{}, err
				}
				result = model.IndexCellModel{
					Pos:    pos,
					Row:    rInd,
					Column: cInd,
					Value:  cell,
				}
				break
			}
		}
	}

	return result, nil
}

/* Получение следующей строки, которая удовлетворяет некоторому условию, определённому в предикате */
func (k *ExcelAnalysis) GetIndexNextRowOffset(index model.IndexCellModel, sheet *excelize.File, offset int, predicat PredicatType) (model.IndexCellModel, error) {
	if sheet == nil {
		var err error
		if sheet, err = k.GetSheet(); err != nil {
			return model.IndexCellModel{}, err
		}
	}

	var result model.IndexCellModel
	rows, err := sheet.GetRows(viper.GetString("table.sheet"))
	if err != nil {
		return model.IndexCellModel{}, err
	}

	for i := (index.Row + offset); i < len(rows); i++ {
		cell := rows[i][index.Column]
		if predicat(cell, i, index.Column) {
			pos, err := excelize.CoordinatesToCellName(i, index.Column)
			if err != nil {
				return model.IndexCellModel{}, err
			}

			result = model.IndexCellModel{
				Pos:    pos,
				Row:    i,
				Column: index.Column,
				Value:  cell,
			}
			break
		}
	}

	return result, nil
}

/* Получение следующей строки, которая удовлетворяет некоторому условию, определённому в предикате */
func (k *ExcelAnalysis) GetIndexNextRow(index model.IndexCellModel, sheet *excelize.File, predicat PredicatType) (model.IndexCellModel, error) {
	if sheet == nil {
		var err error
		if sheet, err = k.GetSheet(); err != nil {
			return model.IndexCellModel{}, err
		}
	}

	var result model.IndexCellModel
	rows, err := sheet.GetRows(viper.GetString("table.sheet"))
	if err != nil {
		return model.IndexCellModel{}, err
	}

	for i := index.Row; i < len(rows); i++ {
		cell := rows[i][index.Column]
		if predicat(cell, i, index.Column) {
			pos, err := excelize.CoordinatesToCellName(i, index.Column)
			if err != nil {
				return model.IndexCellModel{}, err
			}

			result = model.IndexCellModel{
				Pos:    pos,
				Row:    i,
				Column: index.Column,
				Value:  cell,
			}
			break
		}
	}

	return result, nil
}

/* Получение указателя на объект */
func (k *ExcelAnalysis) GetSheet() (*excelize.File, error) {
	file, err := excelize.OpenFile(viper.GetString("paths.table"))
	if err != nil {
		logrus.Fatal(err.Error())
		return nil, err
	}

	return file, nil
}

/* Копирование данных из удалённой таблицы в локальную */
func (k *ExcelAnalysis) CopyTo(filePath string) (bool, error) {
	sheet, err := k.GetSheet()
	if err != nil {
		return false, err
	}

	// Создание нового файла Excel
	f := excelize.NewFile()
	rows, err := sheet.GetRows(viper.GetString("table.sheet"))
	if err != nil {
		return false, err
	}

	// Копирование данных из одной таблицы в другую
	for rInd, row := range rows {
		for cInd, cell := range row {
			pos, err := excelize.CoordinatesToCellName(rInd, cInd)
			if err != nil {
				return false, err
			}
			f.SetCellValue(viper.GetString("table.sheet"), pos, cell)
		}
	}

	// Сохранение нового файла
	if err := f.SaveAs(filePath); err != nil {
		fmt.Println(err)
	}

	return true, nil
}

/* Получение информации о ячейках и её загрузка в структуру */
func (k *ExcelAnalysis) GetValueCells(data *model.HeaderInfoModel, sheet *excelize.File, index model.IndexCellModel, place string) model.IndexCellModel {
	nextIndex, _ := k.GetIndexNextRowOffset(index, sheet, 1, func(value string, row, column int) bool { return (len(value) > 0) })
	nextIndex_RC := model.IndexCellModel{
		Row:    nextIndex.Row,
		Column: (nextIndex.Column + 1),
	}

	indexLen := k.GetLengthCells(nextIndex_RC, sheet, func(value string, row, column int) bool {
		if len(value) <= 0 {
			return false
		}

		if row == nextIndex.Row {
			return true
		}

		rows, err := sheet.GetRows(viper.GetString("table.sheet"))
		if err != nil {
			return false
		}

		return (len(rows[row][column-1]) <= 0)
	})

	ps := reflect.ValueOf(data)
	s := ps.Elem()
	var reflectValue reflect.Value

	if s.Kind() == reflect.Struct {
		reflectValue = s.FieldByName(place)

		if !(reflectValue.IsValid() && reflectValue.CanSet()) {
			return model.IndexCellModel{}
		}
	} else {
		return model.IndexCellModel{}
	}

	for i := 0; i < indexLen; i++ {
		next, _ := k.GetIndexNextRow(nextIndex_RC, sheet, func(value string, row, column int) bool {
			if len(value) <= 0 {
				return false
			}

			if row == nextIndex.Row {
				return true
			}

			rows, err := sheet.GetRows(viper.GetString("table.sheet"))
			if err != nil {
				return false
			}

			return (len(rows[row][column-1]) <= 0)
		})

		if reflectValue.Kind() == reflect.String {
			reflectValue.SetString(next.Value)
		} else if reflectValue.Kind() == reflect.Slice {
			reflectValue.Set(reflect.Append(reflectValue, reflect.ValueOf(next.Value)))
		}

		nextIndex_RC.Row += 1
	}

	return nextIndex
}

func (k *ExcelAnalysis) GetHeaderInfo() (model.HeaderInfoModel, error) {
	sheet, err := k.GetSheet()
	if err != nil {
		return model.HeaderInfoModel{}, err
	}

	var headerInfo model.HeaderInfoModel

	// Определение индекса ячейки с указанием на основную информацию
	index, _ := k.GetIndexByValue("Основная информация", sheet)
	var nextIndex model.IndexCellModel
	nextIndex, _ = k.GetIndexNextRowOffset(index, sheet, 1, func(value string, row, column int) bool { return (len(value) > 0) })

	headerInfo.IPv4 = nextIndex.Value

	// Заполнение полей структуры (начиная с адреса)
	t := reflect.TypeOf(headerInfo)

	for i := 1; i < t.NumField(); i++ {
		nextIndex = k.GetValueCells(&headerInfo, sheet, nextIndex, t.Field(i).Name)
	}

	return headerInfo, nil
}
