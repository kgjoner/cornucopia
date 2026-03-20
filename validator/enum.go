package validator

type Enum interface {
	Enumerate() any // it must be a slice with all possible values with the same type as the field
}
