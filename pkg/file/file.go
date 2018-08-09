package file

import (
	"io"
	"bufio"
	"crypto/md5"
	"fmt"
	"encoding/hex"
)

func GetReaderMd5(reader io.Reader) string {
	r := bufio.NewReader(reader)
	h := md5.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetBytesMd5(data []byte) string {
	mh := md5.New()
	mh.Write(data)
	return hex.EncodeToString(mh.Sum(nil))
}
