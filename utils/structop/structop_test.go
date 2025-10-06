package structop_test

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/kgjoner/cornucopia/v2/utils/structop"
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

type SimilarWithJSON struct {
	Name   string
	Number int
	Kind   string
	Time   json.RawMessage
	Nested json.RawMessage
	Arr    json.RawMessage
}

var mockedSimilarWithJSON = SimilarWithJSON{
	Name:   "JSONName",
	Number: 31,
	Kind:   "json-kind",
	Time:   json.RawMessage(`"2023-12-12T09:10:11.234Z"`),
	Nested: json.RawMessage(`{"Name":"NestedJSONName","Time":"2023-12-12T09:10:11.234Z"}`),
	Arr:    json.RawMessage(`["opt1","opt2"]`),
}

func TestCopyWithJSON(t *testing.T) {
	target := mockedOriginal
	err := structop.New(mockedSimilarWithJSON).Copy(&target)
	if err != nil {
		t.Errorf(err.Error())
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2023-12-12T09:10:11.234Z")

	assert.Equal(t, mockedSimilarWithJSON.Name, target.Name)
	assert.Equal(t, mockedSimilarWithJSON.Number, target.Number)
	assert.Equal(t, Kind(mockedSimilarWithJSON.Kind), target.Kind)
	assert.Equal(t, expectedTime, target.Time)
	assert.Equal(t, "NestedJSONName", target.Nested.Name)
	assert.Equal(t, expectedTime, target.Nested.Time)
	assert.Equal(t, []string{"opt1", "opt2"}, target.Arr)
}

type SimilarWithPointer struct {
	Name   string
	Number int
	Kind   Kind
	Time   time.Time
	Nested *NestedStruct
}

func TestUpdateViaMapWithPointer(t *testing.T) {
	editedMap := map[string]any{
		"name":   "SimilarName",
		"number": 30,
		"kind":   Kind("similar-kind"),
		"time":   time2,
		"nested": map[string]any{
			"name": "SimilarNestedName",
			"time": time.Now(),
		},
	}

	target := SimilarWithPointer{}
	err := structop.New(&target).UpdateViaMap(editedMap)
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equal(t, editedMap["name"], target.Name)
	assert.Equal(t, editedMap["number"], target.Number)
	assert.Equal(t, editedMap["kind"], target.Kind)
	assert.Equal(t, editedMap["time"], target.Time)
	assert.Equal(t, editedMap["nested"].(map[string]any)["name"], target.Nested.Name)
	assert.Equal(t, editedMap["nested"].(map[string]any)["time"], target.Nested.Time)
}

type StructEmbedding struct {
	ID string
	NestedStruct
}

func TestStructEmbeddingUpdate(t *testing.T) {
	target := StructEmbedding{
		ID: "123",
		NestedStruct: NestedStruct{
			Name: "Name",
			Time: time.Now(),
		},
	}
	err := structop.New(&target.NestedStruct).Update(NestedStruct{
		Name: "NewName",
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.NotZero(t, target.Name, "NewName")
}

// Custom type that implements json.Unmarshaler
type CustomUnmarshaler struct {
	Value          string
	WasUnmarshaled bool
}

func (c *CustomUnmarshaler) UnmarshalJSON(data []byte) error {
	c.WasUnmarshaled = true
	// Remove quotes if present
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	c.Value = "unmarshaled:" + str
	return nil
}

// Test struct with a field that implements json.Unmarshaler
type StructWithUnmarshaler struct {
	ID     string
	Custom CustomUnmarshaler
}

func TestJSONUnmarshalerPrecedence(t *testing.T) {
	target := StructWithUnmarshaler{
		ID: "test-id",
		Custom: CustomUnmarshaler{
			Value:          "original",
			WasUnmarshaled: false,
		},
	}

	// Test with []byte input - should use UnmarshalJSON
	editedWithBytes := struct {
		Custom []byte
	}{
		Custom: []byte(`"test-data"`),
	}

	err := structop.New(&target).Update(editedWithBytes)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	assert.True(t, target.Custom.WasUnmarshaled, "UnmarshalJSON should have been called")
	assert.Equal(t, "unmarshaled:test-data", target.Custom.Value, "Value should be set by UnmarshalJSON")
}

func TestJSONUnmarshalerWithJSONMarshaler(t *testing.T) {
	target := StructWithUnmarshaler{
		ID: "test-id",
		Custom: CustomUnmarshaler{
			Value:          "original",
			WasUnmarshaled: false,
		},
	}

	// Test with json.RawMessage (which implements json.Marshaler) - should also use UnmarshalJSON
	editedWithRawMessage := struct {
		Custom json.RawMessage
	}{
		Custom: json.RawMessage(`"from-raw-message"`),
	}

	err := structop.New(&target).Update(editedWithRawMessage)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	assert.True(t, target.Custom.WasUnmarshaled, "UnmarshalJSON should have been called")
	assert.Equal(t, "unmarshaled:from-raw-message", target.Custom.Value, "Value should be set by UnmarshalJSON")
}

func TestJSONUnmarshalerSkippedWhenNotBytes(t *testing.T) {
	target := StructWithUnmarshaler{
		ID: "test-id",
		Custom: CustomUnmarshaler{
			Value:          "original",
			WasUnmarshaled: false,
		},
	}

	// Test with string input - should NOT use UnmarshalJSON, should try other transformations
	editedWithString := struct {
		Custom string
	}{
		Custom: "direct-string",
	}

	err := structop.New(&target).Update(editedWithString)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Since string cannot be converted to CustomUnmarshaler and it doesn't implement json.Unmarshaler
	// for string inputs, the original value should remain unchanged
	assert.False(t, target.Custom.WasUnmarshaled, "UnmarshalJSON should NOT have been called for string input")
	assert.Equal(t, "original", target.Custom.Value, "Value should remain original since no conversion is possible")
}

func TestIsolatedTimeUnmarshal(t *testing.T) {
	// Test if our JSON unmarshaling works for time.Time specifically
	source := struct {
		Time json.RawMessage
	}{
		Time: json.RawMessage(`"2023-12-12T09:10:11.234Z"`),
	}

	target := struct {
		Time time.Time
	}{
		Time: time.Date(2020, time.July, 14, 0, 0, 0, 0, time.UTC),
	}

	t.Logf("Before copy - Target time: %v", target.Time)

	err := structop.New(source).Copy(&target)
	if err != nil {
		t.Errorf("Copy failed: %v", err)
	}

	t.Logf("After copy - Target time: %v", target.Time)

	expectedTime, _ := time.Parse(time.RFC3339, "2023-12-12T09:10:11.234Z")
	if !target.Time.Equal(expectedTime) {
		t.Errorf("Expected %v, got %v", expectedTime, target.Time)
	}
}

func TestIsolatedNestedUnmarshal(t *testing.T) {
	// Test nested structure with time field
	source := struct {
		Nested json.RawMessage
	}{
		Nested: json.RawMessage(`{"Name":"NestedJSONName","Time":"2023-12-12T09:10:11.234Z"}`),
	}

	target := struct {
		Nested NestedStruct
	}{
		Nested: NestedStruct{
			Name: "OrigNestedName",
			Time: time.Date(2020, time.July, 14, 0, 0, 0, 0, time.UTC),
		},
	}

	t.Logf("Before copy - Target nested: %+v", target.Nested)

	err := structop.New(source).Copy(&target)
	if err != nil {
		t.Errorf("Copy failed: %v", err)
	}

	t.Logf("After copy - Target nested: %+v", target.Nested)

	// Let's also check what happens when we unmarshal the JSON manually
	var jsonData map[string]interface{}
	json.Unmarshal(source.Nested, &jsonData)
	t.Logf("Unmarshaled JSON data: %+v", jsonData)
	for k, v := range jsonData {
		t.Logf("Key: %s, Value: %v, Type: %T", k, v, v)
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2023-12-12T09:10:11.234Z")
	if target.Nested.Name != "NestedJSONName" {
		t.Errorf("Expected name 'NestedJSONName', got '%s'", target.Nested.Name)
	}
	if !target.Nested.Time.Equal(expectedTime) {
		t.Errorf("Expected time %v, got %v", expectedTime, target.Nested.Time)
	}
}
