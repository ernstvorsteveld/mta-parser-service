package mta

import (
	"github.com/ernstvorsteveld/mta-common/common"
	"mime/multipart"
)

var ch chan common.FilenameMessage

type BindFile struct {
	Name string                `form:"name" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func PublishFile(filename string, dest string) {
	ch <- common.FilenameMessage{
		Filename: filename,
		Dst:      dest,
	}
}
