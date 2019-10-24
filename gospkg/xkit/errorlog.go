package xkit

import (
	"io"

	"github.com/fidelfly/gox/logx"
)

func Close(t io.Closer) {
	logx.CaptureError(t.Close())
}
