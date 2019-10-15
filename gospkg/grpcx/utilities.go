package rpcx

import "google.golang.org/grpc/metadata"

func GetMDValue(md metadata.MD, key string) string {
	vs := md.Get(key)
	if vs != nil && len(vs) > 0 {
		return vs[0]
	}
	return ""
}
