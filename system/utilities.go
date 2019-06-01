package system

import "strings"

func GetProtectedPath(path string) string {
	if len(path) > 0 {
		if strings.HasPrefix(path, "/") {
			return ProtectedPrefix + path
		}
		return ProtectedPrefix + "/" + path
	}
	return path
}

func GetPublicPath(path string) string {
	if len(path) > 0 {
		if strings.HasPrefix(path, "/") {
			return PublicPrefix + path
		}
		return PublicPrefix + "/" + path
	}
	return path
}
