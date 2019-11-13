package otk

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fidelfly/fxgos/cmd/service"
	"github.com/fidelfly/fxgos/cmd/service/res"
)

type ResourceKey struct {
	Id int64 `json:"id"`
}

func NewResourceKey(id int64) string {
	if jsonData, err := json.Marshal(&ResourceKey{Id: id}); err == nil {
		return string(jsonData)
	}
	return ""
}

func getServer() Service {
	if v, ok := service.GetService(ServiceName); ok {
		if server, ok := v.(Service); ok {
			return server
		}
	}
	panic(fmt.Sprintf("Service(%s) is not registered", ServiceName))
}

func NewOtk(keyType string, typeId string, expired time.Duration, usage string, data string) (string, error) {
	return getServer().NewOtk(keyType, typeId, expired, usage, data)
}
func Consume(id int64) error {
	return getServer().Consume(id)
}
func Validate(key string) (*res.OneTimeKey, error) {
	return getServer().Validate(key)
}
