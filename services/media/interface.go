package media

import (
	"bytes"
)

type MediaService interface {
	Store(file *bytes.Buffer, kind, id string) (string, error)
}
