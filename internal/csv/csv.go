package csv

import (
	"cmp"
	"encoding/csv"
	"errors"
	"io"
	"reflect"
	"slices"
	"strconv"

	"github.com/labstack/gommon/log"
)

func Parse[T any](r io.Reader, row *T) ([]T, error) {
	rowType := reflect.TypeOf(*row)
	if rowType.Kind() != reflect.Struct {
		return nil, errors.New("row must be a struct")
	}
	reader := csv.NewReader(r)
	header, err := reader.Read()
	if err != nil {
		return nil, errors.Join(errors.New("error reading header line"), err)
	}
	headerMap := GetHeaderMap(header, rowType)

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, errors.Join(errors.New("error reading csv body"), err)
	}

	output := []T{}
	for _, line := range lines {
		if len(line) != len(headerMap) {
			log.Warn("line did not match expected length. Expected ", len(headerMap), " but found ", len(line), " ", line)
			continue
		}

		outputStruct := new(T)
		for i, f := range headerMap {
			if f == nil {
				continue
			}
			structValue := reflect.ValueOf(outputStruct).Elem().FieldByName(f.Name)
			if !structValue.CanSet() {
				log.Warn("cannot set unexported field " + structValue.Type().Name())
				continue
			}
			switch structValue.Kind() {
			case reflect.String:
				structValue.SetString(line[i])
			case reflect.Int8:
				fallthrough
			case reflect.Int16:
				fallthrough
			case reflect.Int32:
				fallthrough
			case reflect.Int64:
				fallthrough
			case reflect.Int:
				if intVal, err := strconv.Atoi(line[i]); err != nil {
					log.Warn("Invalid int value ", line[i], " on line ", i)
				} else if structValue.OverflowInt(int64(intVal)) {
					log.Warn("Int value overflow ", line[i], " on line ", i)
				} else {
					structValue.SetInt(int64(intVal))
				}
			case reflect.Float32:
				fallthrough
			case reflect.Float64:
				if floatVal, err := strconv.ParseFloat(line[i], 64); err != nil {
					log.Warn("Invalid float value ", line[i], " on line ", i)
				} else if structValue.OverflowFloat(floatVal) {
					log.Warn("Float value overflow ", line[i], " on line ", i)
				} else {
					structValue.SetFloat(floatVal)
				}
			default:
				return nil, errors.New("Unsupported csv parse type " + structValue.Kind().String())

			}
		}
		output = append(output, *outputStruct)
	}

	return output, nil
}

func GetHeaderMap(header []string, rowType reflect.Type) []*reflect.StructField {
	headerMap := make([]*reflect.StructField, len(header))

	for i := 0; i < rowType.NumField(); i++ {
		field := rowType.Field(i)
		csvTag := field.Tag.Get("csv")
		idx := slices.Index(header, cmp.Or(csvTag, field.Name))
		if idx >= 0 && headerMap[idx] == nil {
			headerMap[idx] = &field
		}
	}
	return headerMap
}
