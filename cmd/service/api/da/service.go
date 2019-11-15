package da

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/service"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/gostool/dbo"
)

const (
	ServiceName  = "service.data.access"
	ResourceType = "resource.da"
)

func RegisterServer(server Service) {
	service.Register(ServiceName, server)
	pub.Subscribe(pub.TopicResource, cacheSubscriber)
}

type SecurityData interface {
	GetResourceType() string
	GetSecurityGroups(ctx context.Context, id int64) ([]int64, error)
}

type Service interface {
	RegisterData(data SecurityData) error
	ValidateAccess(userId int64, resourceType string, resId int64) bool

	Create(ctx context.Context, input interface{}) (*res.SecurityGroup, error)
	Update(ctx context.Context, info dbo.UpdateInfo) error
	Read(ctx context.Context, id int64) (*res.SecurityGroup, error)
	ReadByCode(ctx context.Context, code string) (*res.SecurityGroup, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, input *dbo.ListInfo, conds ...string) ([]*res.SecurityGroup, int64, error)
}
