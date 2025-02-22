package structop

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kgjoner/cornucopia/helpers/htypes"
	"github.com/kgjoner/cornucopia/services/media"
)

type structop struct {
	value any
}

func New(strct any) *structop {
	s := strct
	if r, ok := strct.(reflect.Value); ok {
		s = r.Interface()
	}

	return &structop{s}
}

// Overwrite original struct fields with edited fields of similar name and type, if not zero value.
// Field names are case insensitive and if types don't match, a conversion attempt will occur.
func (s *structop) Update(editedStruct interface{}) error {
	targetFields := s.ReflectMap()
	editedFields := New(editedStruct).ReflectMap()

	for fieldName, editedFieldV := range editedFields {
		targetFieldV := targetFields[fieldName]
		err := setValue(targetFieldV, editedFieldV, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// Overwrite original struct fields with edited fields of similar name and type, if not zero value.
// Field names are case insensitive and if types don't match, a conversion attempt will occur.
//
// Struct nested fields don't need to be nested in the map.
func (s *structop) UpdateViaMap(editedMap map[string]any) error {
	targetFields := s.ReflectMap()

	editedFields := map[string]reflect.Value{}
	for fieldName, fieldValue := range editedMap {
		normalizedName := normalizeFieldName(fieldName)
		editedFields[normalizedName] = reflect.ValueOf(fieldValue)
	}

	for fieldName, targetFieldV := range targetFields {
		editedFieldV, exists := editedFields[fieldName]
		if !exists && targetFieldV.Kind() == reflect.Struct {
			err := New(targetFieldV.Addr()).UpdateViaMap(editedMap)
			if err != nil {
				return err
			}
			continue
		}

		err := setValue(targetFieldV, editedFieldV, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// Copy data from original struct into target struct. Target struct must be a pointer.
func (s structop) Copy(targetStruct any) error {
	targetFields := New(targetStruct).ReflectMap()
	originalFields := s.ReflectMap()

	for fieldName, originalFieldV := range originalFields {
		targetFieldV := targetFields[fieldName]
		err := setValue(targetFieldV, originalFieldV, &setValueOption{ShouldSetZeroValue: true})
		if err != nil {
			return err
		}
	}

	return nil
}

// List struct fields names.
func (s structop) Keys() []string {
	keys := []string{}
	v := reflect.ValueOf(s.value)
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		keys = append(keys, fi.Name)
	}

	return keys
}

// List struct fields names accordingly to their json tag. If marked with "-", field is skipped. 
// If no json tag is provided, field name appears.
func (s structop) JsonKeys() []string {
	keys := []string{}
	v := reflect.ValueOf(s.value)
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		jsonTag := typ.Field(i).Tag.Get("json")
		jsonKey := strings.Split(jsonTag, ",")[0]
		if jsonKey == "-" {
			continue
		} else if jsonKey == "" {
			jsonKey = typ.Field(i).Name
		}

		keys = append(keys, jsonKey)
	}

	return keys
}

// Create a map based on struct using fields as keys. Keep original field name as it is.
func (s structop) Map() map[string]interface{} {
	maps := make(map[string]interface{})
	v := reflect.ValueOf(s.value)
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if !fi.IsExported() {
			continue
		}

		maps[fi.Name] = v.Field(i).Interface()
	}

	return maps
}

// Create a map based on struct using fields as keys. Normalize field names (let all
// lower case and remove underscores). They lead to the reflect value of their original
// value.
func (s structop) ReflectMap() map[string]reflect.Value {
	maps := make(map[string]reflect.Value)
	v := reflect.Indirect(reflect.ValueOf(s.value))
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if !fi.IsExported() {
			continue
		}
		normalizedName := normalizeFieldName(fi.Name)
		maps[normalizedName] = v.Field(i)
	}

	return maps
}

type setValueOption struct {
	ShouldSetZeroValue bool
}

func setValue(target reflect.Value, edited reflect.Value, opt *setValueOption) error {
	if !target.IsValid() || !edited.IsValid() {
		return nil
	}

	shouldSkipZeroValue := opt == nil || !opt.ShouldSetZeroValue
	if shouldSkipZeroValue && edited.IsZero() {
		return nil
	}

	if !target.CanSet() {
		return fmt.Errorf("target cannot be set, it may be unadressable or a private field")
	}

	shouldSetDeeply := target.Kind() == reflect.Struct && shouldSkipZeroValue &&
		target.Type() != reflect.TypeOf(time.Now()) && target.Type() != reflect.TypeOf(media.Media{})

	if edited.Type() == target.Type() {
		if shouldSetDeeply {
			return New(target.Addr()).Update(edited)
		}

		target.Set(edited)
		return nil
	}

	if edited.Kind() == reflect.Pointer {
		edited = reflect.Indirect(edited)
	}

	if edited.Kind() == reflect.Interface {
		switch v := edited.Interface().(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			edited = reflect.ValueOf(v)
		case string:
			edited = reflect.ValueOf(v)
		}
	}

	if edited.Type() == reflect.TypeOf(htypes.NullRawMessage{}) {
		edited = edited.FieldByName("RawMessage")
		return setValue(target, edited, opt)
	}

	if edited.CanConvert(target.Type()) {
		edited = edited.Convert(target.Type())

	} else if target.Kind() == edited.Kind() &&
		!(edited.Kind() == reflect.Slice && edited.Type().Elem().Kind() == reflect.Uint8) {

		switch target.Kind() {
		case reflect.Struct:
			if target.Type() == reflect.TypeOf(htypes.NullTime{}) &&
				edited.Type() == reflect.TypeOf(time.Time{}) {
				timeV := target.FieldByName("Time")
				timeV.Set(edited)
				return nil
			}

			if opt.ShouldSetZeroValue {
				return New(edited).Copy(target.Addr())
			} else {
				return New(target.Addr()).Update(edited)
			}

		case reflect.Slice, reflect.Array:
			return copyReflectedSlice(edited, target.Addr())

		case reflect.Map:
			return copyMap(edited, target.Addr())
		}

	} else if target.Kind() == reflect.Array && edited.Kind() == reflect.Slice {
		return copyReflectedSlice(edited, target.Addr())

	} else if edited.Kind() == reflect.String {
		str := edited.String()

		switch target.Kind() {
		case reflect.String:
			return setValue(target, reflect.ValueOf(str), opt)
		case reflect.Int:
			num, err := strconv.ParseInt(str, 10, 32)
			if err != nil {
				return err
			}
			edited = reflect.ValueOf(int(num))
		case reflect.Uint:
			num, err := strconv.ParseUint(str, 10, 32)
			if err != nil {
				return err
			}
			edited = reflect.ValueOf(uint(num))
		case reflect.Bool:
			bl, err := strconv.ParseBool(str)
			if err != nil {
			 return err
			}
			edited = reflect.ValueOf(bl)
		}

		switch target.Type() {
		case reflect.TypeOf(uuid.New()):
			uuid, err := uuid.Parse(str)
			if err != nil {
				return err
			}

			edited = reflect.ValueOf(uuid)

		case reflect.TypeOf(time.Now()), reflect.TypeOf(htypes.NullTime{}):
			v, err := time.Parse("2006-01-02T15:04:05.9", str)
			if err != nil {
				v, err = time.Parse("2006-01-02", str)
				if err != nil {
					return err
				}
			}
			if target.Type() == reflect.TypeOf(htypes.NullTime{}) {
				edited = reflect.ValueOf(htypes.NullTime{Time: v})
			} else {
				edited = reflect.ValueOf(v)
			}
		default:
			raw := []byte(str)

			switch target.Kind() {
			case reflect.Struct:
				var v map[string]any
				json.Unmarshal(raw, &v)
				return New(target.Addr()).UpdateViaMap(v)

			case reflect.Slice, reflect.Array:
				if target.Type().Elem().Kind() == reflect.Struct {
					var v []map[string]any
					json.Unmarshal(raw, &v)
					return copyReflectedSlice(reflect.ValueOf(v), target.Addr())
				} else {
					var v []string
					json.Unmarshal(raw, &v)
					return copyReflectedSlice(reflect.ValueOf(v), target.Addr())
				}

			case reflect.Map:
				var v map[string]any
				json.Unmarshal(raw, &v)
				return copyMap(reflect.ValueOf(v), target.Addr())
			}
		}

	} else if v, ok := edited.Interface().(map[string]any); ok &&
		(target.Kind() == reflect.Struct || target.Kind() == reflect.Map) {

		switch target.Kind() {
		case reflect.Struct:
			return New(target.Addr()).UpdateViaMap(v)
		case reflect.Map:
			return copyMap(reflect.ValueOf(v), target.Addr())
		}

	} else if byteArr, ok := edited.Interface().([]byte); ok || edited.CanConvert(reflect.TypeOf([]byte{})) {

		if !ok {
			marshaller, ok := edited.Interface().(json.Marshaler)
			if !ok {
				return nil
			}

			data, err := marshaller.MarshalJSON()
			if err != nil {
				return err
			}
			byteArr = data
		}

		rawString := bytes.NewBuffer(byteArr).String()
		return setValue(target, reflect.ValueOf(rawString), opt)

	} else if nullStr, ok := edited.Interface().(sql.NullString); ok {
		if nullStr.Valid {
			setValue(target, reflect.ValueOf(nullStr.String), opt)
		} else {
			setValue(target, reflect.ValueOf(""), opt)
		}

	} else if nullRawMsg, ok := edited.Interface().(htypes.NullRawMessage); ok {
		rawString := bytes.NewBuffer(nullRawMsg.RawMessage).String()
		return setValue(target, reflect.ValueOf(rawString), opt)
	}

	if edited.Type() == target.Type() {
		if shouldSetDeeply {
			return New(target.Addr()).Update(edited)
		}

		target.Set(edited)
		return nil
	}

	return nil
}

// Copy slice of type T into slice of type K. Target data must be a pointer.
// T and K must be structs with similar fields.
func CopySlice[T, K any](original []T, target *[]K) error {
	originalV := reflect.Indirect(reflect.ValueOf(original))
	targetV := reflect.Indirect(reflect.ValueOf(target))

	if targetV.Kind() != originalV.Kind() {
		return fmt.Errorf("kinds must be the same to copy")
	}

	if originalV.Kind() == reflect.Array || originalV.Kind() == reflect.Slice {
		for i := 0; i < originalV.Len(); i++ {
			var newData K
			originalEl := originalV.Index(i)
			err := New(originalEl.Addr()).Copy(&newData)
			if err != nil {
				return err
			}
			*target = append(*target, newData)
		}
	}

	return nil
}

func copyReflectedSlice(original reflect.Value, target reflect.Value) error {
	originalV := reflect.Indirect(original)
	targetV := reflect.Indirect(target)

	if originalV.Kind() == reflect.Array || originalV.Kind() == reflect.Slice {
		if targetV.Kind() == reflect.Array {
			t := reflect.SliceOf(targetV.Type().Elem())
			targetV = reflect.New(t).Elem()
		}

		for i := 0; i < originalV.Len(); i++ {
			newData := reflect.New(targetV.Type().Elem())
			originalEl := originalV.Index(i)

			err := setValue(reflect.Indirect(newData), originalEl, &setValueOption{ShouldSetZeroValue: true})
			if err != nil {
				return err
			}

			targetV.Set(reflect.Append(targetV, reflect.Indirect(newData)))
		}

		if reflect.Indirect(target).Kind() == reflect.Array {
			reflect.Copy(reflect.Indirect(target), targetV)
		}
	}

	return nil
}

func copyMap(original reflect.Value, target reflect.Value) error {
	originalV := reflect.Indirect(original)
	targetV := reflect.Indirect(target)

	if targetV.IsNil() {
		if originalV.IsNil() {
			return nil
		}

		targetV.Set(reflect.MakeMapWithSize(targetV.Type(), 0))
	}

	for _, key := range originalV.MapKeys() {
		newData := reflect.New(targetV.Type().Elem())
		originalEl := originalV.MapIndex(key)

		err := setValue(reflect.Indirect(newData), originalEl, &setValueOption{ShouldSetZeroValue: true})
		if err != nil {
			return err
		}

		originalKeyType := originalV.Type().Key()
		targetKeyType := targetV.Type().Key()
		if originalKeyType != targetKeyType {

			if key.CanConvert(targetKeyType) {
				key = key.Convert(targetKeyType)

			} else if key.Kind() == reflect.String {
				str := key.String()
				switch targetKeyType.Kind() {
				case reflect.Int:
					num, err := strconv.ParseInt(str, 10, 32)
					if err != nil {
						return err
					}
					key = reflect.ValueOf(int(num))
				case reflect.Uint:
					num, err := strconv.ParseUint(str, 10, 32)
					if err != nil {
						return err
					}
					key = reflect.ValueOf(uint(num))
				}

			} else {
				return fmt.Errorf("cannot convert keys")
			}
		}

		targetV.SetMapIndex(key, reflect.Indirect(newData))
	}

	return nil
}

func normalizeFieldName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), "_", "")
}
