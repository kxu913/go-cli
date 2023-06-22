package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func setField(field reflect.Value, defaultVal string) error {

	if !field.CanSet() {
		return fmt.Errorf("Can't set value\n")
	}

	switch field.Kind() {

	case reflect.Int:
		if val, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(defaultVal).Convert(field.Type()))
	}

	return nil
}

func Set(ptr interface{}) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("Not a pointer")
	}

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if defaultVal := t.Field(i).Tag.Get("default"); defaultVal != "-" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}

		}
	}
	return nil
}

// 下划线写法转为驼峰写法
func SnakeToCaml(field string) string {
	children := strings.Split(field, "_")
	upperCaseChildren := []string{}
	for _, c := range children {
		upperCaseChildren = append(upperCaseChildren, Ucfirst(c))

	}
	return strings.Join(upperCaseChildren, "")
}

// 首字母大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}
func GetEnvWithDefault(key string, defValue string) string {
	val, err := os.LookupEnv(key)
	if !err {
		return defValue
	}
	return val

}

func WriteFile(fileName string) *os.File {
	os.OpenFile(fileName, os.O_CREATE, 0o666)
	file, err := os.OpenFile(fileName, os.O_RDWR, 0o666)
	if err != nil {
		panic(err)
	}
	return file
}
