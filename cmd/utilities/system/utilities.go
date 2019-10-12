package system

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fidelfly/fxgo/logx"
)

//todo remove
func GetProtectedPath(path string) string {
	if len(path) > 0 {
		if strings.HasPrefix(path, "/") {
			return ProtectedPrefix + path
		}
		return ProtectedPrefix + "/" + path
	}
	return path
}

//todo remove
func GetPublicPath(path string) string {
	if len(path) > 0 {
		if strings.HasPrefix(path, "/") {
			return PublicPrefix + path
		}
		return PublicPrefix + "/" + path
	}
	return path
}

func CreateTemporaryFile(data []byte, file string) (path string, err error) {
	tempFile, err := os.Create(filepath.Join(Runtime.TemporaryPath, file))
	if err != nil {
		return
	}

	err = os.Chmod(tempFile.Name(), os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		logx.CaptureError(tempFile.Close())
	}()

	_, err = tempFile.Write(data)
	if err != nil {
		return
	}

	path = tempFile.Name()
	return
}

func GetImagePath(name string) string {
	return filepath.Join(Runtime.AssetPath, "image", name)
}

func GetAssetPath(relativePath string) string {
	return filepath.Join(Runtime.AssetPath, relativePath)
}
