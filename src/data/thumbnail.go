package data

import (
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type Thumbnail interface {
	// TODO(john)
	GetMimeType() string
	GetSize() int64
	GetContent() []byte
	GetWidth() int64
	GetHeight() int64
	GetExtension() string
	GetError() error
}

type ThumbnailImpl struct {
	MimeType  string `json:"mimetype"`
	Size      int64  `json:"size"`
	Content   []byte `json:"content"`
	Width     int64  `json:"width"`
	Height    int64  `json:"height"`
	TextBody  string `json:"text_body"`
	Extension string `json:"extension"`
	Error     error  `json:"error,omitempty"`
}

func (t *ThumbnailImpl) GetMimeType() string {
	return t.MimeType
}

func (t *ThumbnailImpl) GetSize() int64 {
	return t.Size
}

func (t *ThumbnailImpl) GetContent() []byte {
	return t.Content
}

func (t *ThumbnailImpl) GetWidth() int64 {
	return t.Width
}

func (t *ThumbnailImpl) GetHeight() int64 {
	return t.Height
}

func (t *ThumbnailImpl) GetExtension() string {
	if len(t.Content) == 0 {
		return ""
	}

	mt := mimetype.Detect(t.Content)
	if mt == nil {
		return ""
	}
	return strings.ToLower(mt.Extension())
}

func (t *ThumbnailImpl) GetError() error {
	return t.Error
}
