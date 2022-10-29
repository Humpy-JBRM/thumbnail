package data

type Thumbnail interface {
	// TODO(john)
	GetMimeType() string
	GetSize() int64
	GetContent() []byte
	GetWidth() int64
	GetHeight() int64
}

type ThumbnailImpl struct {
	MimeType string `json:"mimetype"`
	Size     int64  `json:"size"`
	Content  []byte `json:"content"`
	Width    int64  `json:"width"`
	Height   int64  `json:"height"`
	TextBody string `json:"text_body"`
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
