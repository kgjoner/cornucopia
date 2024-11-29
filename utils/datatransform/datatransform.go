package datatransform

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"

	"github.com/kgjoner/cornucopia/helpers/htypes"
)

func ToRawMessage(obj interface{}) json.RawMessage {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return json.RawMessage(data)
}

func ToRawMessageSlice[T any](objSlc []T) []json.RawMessage {
	res := []json.RawMessage{}
	for _, obj := range objSlc {
		data, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}

		res = append(res, json.RawMessage(data))
	}

	return res
}

func ToNullRawMessage(obj interface{}) htypes.NullRawMessage {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return htypes.NullRawMessage{
		RawMessage: json.RawMessage(data),
		Valid:      len(data) != 0,
	}
}

func ToNullString(str string) sql.NullString {
	isNull := str == ""
	return sql.NullString{
		String: str,
		Valid:  !isNull,
	}
}

func ToNullInt(num int) sql.NullInt32 {
	isNull := num == 0
	return sql.NullInt32{
		Int32: int32(num),
		Valid: !isNull,
	}
}

func ToNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

func ToStringArray[T any](arr []T) []string {
	res := []string{}

	for _, v := range arr {
		refV := reflect.ValueOf(v)
		if refV.Kind() == reflect.String {
			res = append(res, refV.String())
		}
	}

	return res
}

func ToInt64Array(arr []int) []int64 {
	res := []int64{}

	for _, v := range arr {
		res = append(res, int64(v))
	}

	return res
}
