package sta

import (
	"github.com/ernstvorsteveld/mta-common/common"
	"mime/multipart"
)

var ch chan common.FilenameMessage

type BindFile struct {
	Name  string                `form:"name" binding:"required"`
	File  *multipart.FileHeader `form:"file" binding:"required"`
}

