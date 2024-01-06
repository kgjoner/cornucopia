package validator

type Enum interface {
	Enumerate() any // it must be a struct with field-value as string-validvalue
}
