package model

import (
	"io"
)

type Upload struct {
	FileName string
	File     io.Reader
}

type ReadChunk struct {
	UploadID string
	Offset   int64 `form:"offset"`
	Limit    int64 `form:"limit" binding:"required"`
}
