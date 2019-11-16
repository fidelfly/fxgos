package otk

import (
	"time"

	"github.com/fidelfly/fxgos/cmd/service"
	"github.com/fidelfly/fxgos/cmd/service/res"
)

const (
	ServiceName = "service.otk"
)

func RegisterServer(server Service, dependencies ...string) {
	service.Register(ServiceName, server, dependencies...)
}

type Service interface {
	NewOtk(keyType string, typeId string, expired time.Duration, usage string, data string) (string, error)
	Consume(id int64) error
	Validate(key string) (*res.OneTimeKey, error)
}
