package role

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/service"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/gostool/dbo"
)

const (
	ResourceType = "resource.role"
	ServiceName  = "service.role"
)

func RegisterServer(server Service, dependencies ...string) {
	service.Register(ServiceName, server, dependencies...)

	pub.Subscribe(pub.TopicResource, cacheSubscriber)
}

type Service interface {
	Create(ctx context.Context, input interface{}) (*res.Role, error)
	Update(ctx context.Context, info dbo.UpdateInfo) error
	Read(ctx context.Context, id int64) (*res.Role, error)
	ReadByCode(ctx context.Context, code string) (*res.Role, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, input *dbo.ListInfo, conds ...string) ([]*res.Role, int64, error)
}

type Form struct {
	Code        string  `json:"code"`
	Roles       []int64 `json:"roles"`
	Description string  `json:"description"`
}
