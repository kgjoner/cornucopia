package dbhandler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/kgjoner/cornucopia/v2/utils/structop"
)

// Implements Scanner for anything coming from json. It normalizes keys and timestamps.
// For example: created_at, createdAt, CreatedAt are considered all the same, normalized to createdat.
// Timestamps without timezone are normalized to have Z timezone.
func FromJSON[K any](v *K) interface{ sql.Scanner } {
	return &fromJSONScan[K]{
		value: v,
	}
}

type fromJSONScan[K any] struct {
	value *K
}

func (s *fromJSONScan[K]) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if s.value == nil {
		s.value = new(K)
	}

	if v, ok := value.([]byte); ok {
		var dataAsInterface interface{}
		err := json.Unmarshal(v, &dataAsInterface)
		if err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}

		// Recursively normalize all keys
		normalizedData := normalizeMapKeys(dataAsInterface)
		normalizedData = normalizeTimestamps(normalizedData)

		v, err = json.Marshal(normalizedData)
		if err != nil {
			return fmt.Errorf("failed to re-marshal data: %w", err)
		}

		if err := json.Unmarshal(v, s.value); err != nil {
			return fmt.Errorf("failed to unmarshal json: %w", err)
		}

		return nil
	}

	return fmt.Errorf("failed to scan from json")
}

// Implements Scanner for structs
func Struct[K any](v *K) interface{ sql.Scanner } {
	return &structScan[K]{
		value: v,
	}
}

type structScan[K any] struct {
	value *K
}

func (s *structScan[K]) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if s.value == nil {
		s.value = new(K)
	}

	if v, ok := value.([]byte); ok {
		var data map[string]any
		err := json.Unmarshal(v, &data)
		if err != nil {
			return err
		}

		if reflect.Indirect(reflect.ValueOf(s.value)).Kind() == reflect.Pointer {
			pointer := reflect.Indirect(reflect.ValueOf(s.value))
			if pointer.IsNil() {
				pointer.Set(reflect.New(reflect.TypeOf(*s.value).Elem()))
			}

			err = structop.New(*s.value).UpdateViaMap(data)
			return err
		}

		err = structop.New(s.value).UpdateViaMap(data)
		return err
	}

	return fmt.Errorf("failed to scan struct")
}

// Implements Scanner for struct arrays
func StructArray[K any](v *[]K) interface{ sql.Scanner } {
	return (*structArrayScan[K])(v)
}

type structArrayScan[K any] []K

func (s *structArrayScan[K]) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return s.scanBytes(src)
	case nil:
		*s = nil
		return nil
	}

	return fmt.Errorf("cannot convert %T to structArrayScan", src)
}

func (a *structArrayScan[K]) scanBytes(src []byte) error {
	// fmt.Printf("Src: %s", src)
	elems, err := parseJSONArray(src)
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(structArrayScan[K], len(elems))
		for i, v := range elems {
			var data map[string]any
			err := json.Unmarshal(v, &data)
			if err != nil {
				return err
			}

			err = structop.New(&b[i]).UpdateViaMap(data)
			if err != nil {
				return err
			}
		}
		*a = b
	}
	return nil
}

// Implements Scanner for map type
func Map[K comparable, V any](m *map[K]V) interface{ sql.Scanner } {
	return (*mapScan[K, V])(m)
}

type mapScan[K comparable, V any] map[K]V

func (p *mapScan[K, V]) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, p)
	default:
		return fmt.Errorf("failed to scan map")
	}
}

// Implements Scanner for int arrays
func IntArray(a *[]int) interface{ sql.Scanner } {
	return (*intArray)(a)
}

type intArray []int

func (a *intArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to intArray", src)
}

func (a *intArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "intArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(intArray, len(elems))
		for i, v := range elems {
			if b[i], err = strconv.Atoi(string(v)); err != nil {
				return fmt.Errorf("pq: parsing array element index %d: %v", i, err)
			}
		}
		*a = b
	}
	return nil
}

// Implements Scanner for array of enum type
func EnumArray[K ~string](m *[]K) interface{ sql.Scanner } {
	return (*enumArrayScan[K])(m)
}

type enumArrayScan[K ~string] []K

func (a *enumArrayScan[K]) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("cannot convert %T to enumArrayScan", src)
}

func (a *enumArrayScan[K]) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "enumArrayScan")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(enumArrayScan[K], len(elems))
		for i, v := range elems {
			if b[i] = K(v); v == nil {
				return fmt.Errorf("parsing array element index %d: cannot convert nil to string", i)
			}
		}
		*a = b
	}
	return nil
}

// parseArray extracts the dimensions and elements of an array represented in
// text format. Only representations emitted by the backend are supported.
// Notably, whitespace around brackets and delimiters is significant, and NULL
// is case-sensitive.
//
// See http://www.postgresql.org/docs/current/static/arrays.html#ARRAYS-IO
func parseArray(src, del []byte) (dims []int, elems [][]byte, err error) {
	var depth, i int

	if len(src) < 1 || src[0] != '{' {
		return nil, nil, fmt.Errorf("pq: unable to parse array; expected %q at offset %d", '{', 0)
	}

Open:
	for i < len(src) {
		switch src[i] {
		case '{':
			depth++
			i++
		case '}':
			elems = make([][]byte, 0)
			goto Close
		default:
			break Open
		}
	}
	dims = make([]int, i)

Element:
	for i < len(src) {
		switch src[i] {
		case '{':
			if depth == len(dims) {
				break Element
			}
			depth++
			dims[depth-1] = 0
			i++
		case '"':
			var elem = []byte{}
			var escape bool
			for i++; i < len(src); i++ {
				if escape {
					elem = append(elem, src[i])
					escape = false
				} else {
					switch src[i] {
					default:
						elem = append(elem, src[i])
					case '\\':
						escape = true
					case '"':
						elems = append(elems, elem)
						i++
						break Element
					}
				}
			}
		default:
			for start := i; i < len(src); i++ {
				if bytes.HasPrefix(src[i:], del) || src[i] == '}' {
					elem := src[start:i]
					if len(elem) == 0 {
						return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
					}
					if bytes.Equal(elem, []byte("NULL")) {
						elem = nil
					}
					elems = append(elems, elem)
					break Element
				}
			}
		}
	}

	for i < len(src) {
		if bytes.HasPrefix(src[i:], del) && depth > 0 {
			dims[depth-1]++
			i += len(del)
			goto Element
		} else if src[i] == '}' && depth > 0 {
			dims[depth-1]++
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}

Close:
	for i < len(src) {
		if src[i] == '}' && depth > 0 {
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}
	if depth > 0 {
		err = fmt.Errorf("pq: unable to parse array; expected %q at offset %d", '}', i)
	}
	if err == nil {
		for _, d := range dims {
			if (len(elems) % d) != 0 {
				err = fmt.Errorf("pq: multidimensional arrays must have elements with matching dimensions")
			}
		}
	}
	return
}

func scanLinearArray(src, del []byte, typ string) (elems [][]byte, err error) {
	dims, elems, err := parseArray(src, del)
	if err != nil {
		return nil, err
	}
	if len(dims) > 1 {
		return nil, fmt.Errorf("pq: cannot convert ARRAY%s to %s", strings.Replace(fmt.Sprint(dims), " ", "][", -1), typ)
	}
	return elems, err
}

func parseJSONArray(src []byte) (elems [][]byte, err error) {
	if len(src) < 1 || src[0] != '[' {
		return nil, fmt.Errorf("pq: unable to parse json array; expected %q at offset %d", '[', 0)
	}

	var elem []byte
	depth := 0
	scape := false
	inString := false
	for i, b := range src {
		switch b {
		case '{', '[':
			depth++
			if i != 0 {
				elem = append(elem, b)
			}
		case '}', ']':
			depth--
			if i != len(src)-1 {
				elem = append(elem, b)
			} else if len(elem) > 0 {
				elems = append(elems, elem)
			}
		case '\\':
			scape = !scape
			elem = append(elem, b)
		case '"':
			if !scape {
				inString = !inString
			}
			elem = append(elem, b)
		case ',':
			if depth == 1 && !inString {
				elems = append(elems, elem)
				elem = []byte{}
			} else {
				elem = append(elem, b)
			}
		default:
			elem = append(elem, b)
		}

		if b != '\\' {
			scape = false
		}
	}

	return elems, nil
}
