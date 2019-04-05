package main

import (
	"fmt"
	"strconv"
	"time"
)

var nilDateTimeIF bool = false
var nilDateTime time.Time

func NilDateTime() time.Time {
	if !nilDateTimeIF {
		nilDateTime, _ = time.Parse("2006-01-02T15:04:05", "1900-01-01T01:01:01")
		nilDateTimeIF = true
	}
	return nilDateTime

}

func InterfaceToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

/*
func StrToDate(s string, format string) (time.Time, error) {
	var t time.Time
	var err error
	if format == "YYYY-MM-DD" {
		t, err = time.Parse(time.RFC3339, s+"T11:45:26.371Z")
	} else if format == "YYYY-MM-DD hh:mm:ss" {
		t, err = time.Parse("2006-01-02T15:04:05", s)
	} else if format == "DD-MM-YYYY" {
		//t, err = time.Parse("02-01-2006 15:04:05", s+" 00:00:00") - тоже работает
		t, err = time.Parse("02-01-2006", s)
	} else if format == "DD-MM-YYYY hh:mm:ss" {
		t, err = time.Parse("02-01-2006 15:04:05", s)
	}

	return t, err
}
*/

func InterfaceToDateTime(v interface{}) time.Time {
	switch v.(type) {
	case time.Time:
		return v.(time.Time)
	case string:
		s := InterfaceToString(v)
		if s == "" {
			s = "1900-01-01T01:01:01"
		}
		if len(s) > 20 { ///"1970-10-28T15:04:05"
			s = s[0:19]
		}
		t, _ := time.Parse("2006-01-02T15:04:05", s)
		return t
	default:
		return NilDateTime()

	}
	return NilDateTime()
}

func InterfaceToFloat64(v interface{}) float64 {
	switch v.(type) {
	case float64:
		return v.(float64)
	case int:
		return float64(v.(int))
	case int64:
		return float64(v.(int64))
	case string:
		s, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return 0
		}
		return s
	default:
		return 0
	}
	return 0
}
