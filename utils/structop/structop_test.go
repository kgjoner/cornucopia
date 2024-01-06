package structop_test

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/kgjoner/cornucopia/utils/structop"
	"github.com/stretchr/testify/assert"
)

type OriginalStruct struct {
	Name   string
	Number int
	Kind   Kind
	Time   time.Time
	Nested NestedStruct
	Arr    []string
}

type NestedStruct struct {
	Name string
	Time time.Time
}

type Kind string

type EditableFields struct {
	Name   string
	Number int
	Nested NestedStruct
}

type SimilarStruct struct {
	Name   string
	Number int
	Kind   string
	Time   time.Time
	Nested struct {
		Name sql.NullString
		Time time.Time
	}
}

var time1, _ = time.Parse("2006-Jan-02", "2020-Jul-14")
var time2, _ = time.Parse("2006-Jan-02", "2020-Sep-22")

var mockedOriginal = OriginalStruct{
	"OrigName",
	10,
	Kind("orig-kind"),
	time1,
	NestedStruct{
		"OrigNestedName",
		time1,
	},
	[]string{},
}

var mockedEdited = EditableFields{
	"EditName",
	20,
	NestedStruct{
		Name: "EditNestedName",
	},
}

var mockedSimilar = SimilarStruct{
	"SimilarName",
	30,
	"similar-kind",
	time2,
	struct {
		Name sql.NullString
		Time time.Time
	}{
		sql.NullString{String: "SimilarNestedName", Valid: true},
		time.Now(),
	},
}

func TestUpdate(t *testing.T) {
	target := mockedOriginal
	err := structop.New(&target).Update(mockedEdited)
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, mockedEdited.Name, target.Name)
	assert.Equal(t, mockedEdited.Number, target.Number)
	assert.Equal(t, mockedOriginal.Kind, target.Kind)
	assert.Equal(t, mockedOriginal.Time, target.Time)
	assert.Equal(t, mockedEdited.Nested.Name, target.Nested.Name)
	assert.Equal(t, mockedOriginal.Nested.Time, target.Nested.Time)
}

func TestUpdateViaMap(t *testing.T) {
	editedMap := structop.New(mockedEdited).Map()

	target := mockedOriginal
	err := structop.New(&target).UpdateViaMap(editedMap)
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, mockedEdited.Name, target.Name)
	assert.Equal(t, mockedEdited.Number, target.Number)
	assert.Equal(t, mockedOriginal.Kind, target.Kind)
	assert.Equal(t, mockedOriginal.Time, target.Time)
	assert.Equal(t, mockedEdited.Nested.Name, target.Nested.Name)
	assert.Equal(t, mockedOriginal.Nested.Time, target.Nested.Time)
}

func TestCopy(t *testing.T) {
	target := mockedOriginal
	err := structop.New(mockedSimilar).Copy(&target)
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, mockedSimilar.Name, target.Name)
	assert.Equal(t, mockedSimilar.Number, target.Number)
	assert.Equal(t, Kind(mockedSimilar.Kind), target.Kind)
	assert.Equal(t, mockedSimilar.Time, target.Time)
	assert.Equal(t, mockedSimilar.Nested.Name.String, target.Nested.Name)
	assert.Equal(t, mockedSimilar.Nested.Time, target.Nested.Time)
}

func TestKeys(t *testing.T) {
	keys := structop.New(mockedOriginal).Keys()

	expected := []string{"Name", "Number", "Kind", "Time", "Nested", "Arr"}
	assert.Equal(t, expected, keys)
}

func TestCopySlice(t *testing.T) {
	var target []OriginalStruct
	from := []SimilarStruct{mockedSimilar}
	err := structop.CopySlice(from, &target)
	if err != nil {
		t.Errorf(err.Error())
	}

	var expectedEl OriginalStruct
	err = structop.New(mockedSimilar).Copy(&expectedEl)
	if err != nil {
		t.Errorf(err.Error())
	}

	expected := []OriginalStruct{expectedEl}
	assert.Equal(t, expected, target)
}

type SimilarWithJson struct {
	Name   string
	Number int
	Kind   string
	Time   json.RawMessage
	Nested json.RawMessage
	Arr    json.RawMessage
}

var mockedSimilarWithJson = SimilarWithJson{
	Name:   "JsonName",
	Number: 31,
	Kind:   "json-kind",
	Time:   json.RawMessage("2023-12-12T09:10:11.2341"),
	Nested: json.RawMessage(`{"name":"NestedJsonName","time":"2023-12-12T09:10:11.2341"}`),
	Arr:    json.RawMessage(`["opt1","opt2"]`),
}

var timeFormat = "2006-01-02T15:04:05.9"

func TestCopyWithJson(t *testing.T) {
	target := mockedOriginal
	err := structop.New(mockedSimilarWithJson).Copy(&target)
	if err != nil {
		t.Errorf(err.Error())
	}

	expectedTime, _ := time.Parse(timeFormat, "2023-12-12T09:10:11.2341")

	assert.Equal(t, mockedSimilarWithJson.Name, target.Name)
	assert.Equal(t, mockedSimilarWithJson.Number, target.Number)
	assert.Equal(t, Kind(mockedSimilarWithJson.Kind), target.Kind)
	assert.Equal(t, expectedTime, target.Time)
	assert.Equal(t, "NestedJsonName", target.Nested.Name)
	assert.Equal(t, expectedTime, target.Nested.Time)
	assert.Equal(t, []string{"opt1", "opt2"}, target.Arr)
}
