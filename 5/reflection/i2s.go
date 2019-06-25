package main

import (
	"errors"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	// todo
	dataValue := reflect.ValueOf(data)
	outValue := reflect.ValueOf(out)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}
	if outValue.Kind() == reflect.Ptr {
		outValue = outValue.Elem()
	} else {
		return errors.New("NOT POINTER, CAN`T WRITE DATA")
	}
	if outValue.Kind() == reflect.Slice {
		outValue.Set(reflect.MakeSlice(outValue.Type(), dataValue.Len(), dataValue.Len()))
		for j := 0; j < dataValue.Len(); j++ {
			ptr := outValue.Index(j).Addr()
			err := i2s(dataValue.Index(j).Interface(), ptr.Interface())
			if err != nil {
				return err
			}
		}
	} else {
		if dataValue.Kind() == reflect.Slice {
			return errors.New("EXPECTED STRUCTURE, NOT SLICE")
		}
		for i := 0; i < outValue.NumField(); i++ {
			outValueType := outValue.Type().Field(i)
			switch outValueType.Type.Kind() {
			case reflect.Int:
				k := reflect.ValueOf(outValueType.Name)
				toSet := dataValue.MapIndex(k).Elem()
				if toSet.Kind() != reflect.Float64 {
					return errors.New("BAD FIELD TYPE")
				}
				outValue.Field(i).SetInt(int64(toSet.Float()))
			case reflect.String:
				k := reflect.ValueOf(outValueType.Name)
				toSet := dataValue.MapIndex(k).Elem()
				if toSet.Kind() != reflect.String {
					return errors.New("BAD FIELD TYPE")
				}
				outValue.Field(i).SetString(toSet.String())
			case reflect.Bool:
				k := reflect.ValueOf(outValueType.Name)
				toSet := dataValue.MapIndex(k).Elem()
				if toSet.Kind() != reflect.Bool {
					return errors.New("BAD FIELD TYPE")
				}
				outValue.Field(i).SetBool(toSet.Bool())
			case reflect.Struct:
				k := reflect.ValueOf(outValueType.Name)
				toSet := dataValue.MapIndex(k).Elem()
				if toSet.Kind() != reflect.Map {
					return errors.New("EXPECTED STRUCTURE")
				}
				ptr := outValue.Field(i).Addr()
				err := i2s(toSet.Interface(), ptr.Interface())
				if err != nil {
					return err
				}
			case reflect.Slice:
				k := reflect.ValueOf(outValueType.Name)
				newDataValue := dataValue.MapIndex(k).Elem()
				if newDataValue.Kind() != reflect.Slice {
					return errors.New("NOT SLICE")
				}
				outValue.Field(i).Set(reflect.MakeSlice(outValue.Field(i).Type(), newDataValue.Len(), newDataValue.Len()))
				for j := 0; j < newDataValue.Len(); j++ {
					ptr := outValue.Field(i).Index(j).Addr()
					err := i2s(newDataValue.Index(j).Elem().Interface(), ptr.Interface())
					if err != nil {
						return err
					}
				}
			default:
				return errors.New("NOT SUPPORTED")

			}
		}
	}
	return nil
}
