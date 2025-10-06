package media

import (
	"bytes"
	"fmt"
	"image"
	"net/http"
	"reflect"
	"strings"

	"github.com/kgjoner/cornucopia/helpers/apperr"
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
	if m.url == "" && (m.IsEmpty() || reflect.ValueOf(m.mediaService).IsZero()) {
		return apperr.NewValidationError("missing fields in media type")
	}

	return nil
}

func (m Media) IsZero() bool {
	return reflect.ValueOf(m).IsZero()
}

func (m Media) IsEmpty() bool {
	return m.file == nil || m.file.Len() == 0
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

	if m.IsEmpty() || reflect.ValueOf(m.mediaService).IsZero() {
		return "", apperr.NewValidationError("missing fields in media type")
	}

	url, err := m.mediaService.Store(m.file, m.kind, m.id)
	if err != nil {
		return "", err
	}

	m.url = url
	return url, nil
}

// Check if media is an image type.
func (m *Media) IsImage() bool {
	return strings.Contains(m.mime, "image")
}

// Get dimensions of a image. Only works for image types.
func (m *Media) Shape() (width int, height int, err error) {
	if !m.IsImage() {
		return 0, 0, fmt.Errorf("must be an image to get its shape")
	}

	fileCopy := bytes.NewReader(m.file.Bytes())
	img, _, err := image.Decode(fileCopy)
	if err != nil {
		return 0, 0, err
	}

	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}
