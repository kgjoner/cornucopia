package media

import (
	"bytes"
	"net/http"
	"reflect"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

type Media struct {
	url          string
	file         *bytes.Buffer
	mediaService MediaService
	mime         string
	// May be used as reference in URL by media service
	id string
	// May be used for adding custom settings by media service
	kind string
}

func New(file *bytes.Buffer, mediaService MediaService) *Media {
	return &Media{
		file:         file,
		mediaService: mediaService,
		mime:         http.DetectContentType(file.Bytes()),
	}
}

func (m *Media) IsValid() error {
	if m.url == "" && (reflect.ValueOf(m.file).IsZero() || reflect.ValueOf(m.mediaService).IsZero()) {
		return normalizederr.NewValidationError("Missing fields in media type")
	}

	return nil
}

func (m Media) IsZero() bool {
	return reflect.ValueOf(m.file).IsZero()
}

// Set optional props that may be used for media service.
func (m *Media) Config(id string, kind string) {
	m.id = id
	m.kind = kind
}

// Save media in the cloud, if not already, and return their URL.
func (m *Media) URL() (string, error) {
	if m.url != "" {
		return m.url, nil
	}

	if reflect.ValueOf(m.file).IsZero() || reflect.ValueOf(m.mediaService).IsZero() {
		return "", normalizederr.NewValidationError("Missing fields in media type")
	}

	url, err := m.mediaService.Store(m.file, m.kind, m.id)
	if err != nil {
		return "", err
	}

	m.url = url
	return url, nil
}
