package user

import (
	"context"

	"github.com/fidelfly/fxgos/cmd/service"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/gostool/dbo"
)

const ServiceName = "service.user"

const ResourceType = "resource.user"

func RegisterServer(server Service) {
	service.Register(ServiceName, server)

	pub.Subscribe(pub.TopicResource, cacheSubscriber)
}

type ValidateInput struct {
	Id       int64
	Code     string
	Email    string
	Password string
}

type Service interface {
	Create(ctx context.Context, input interface{}) (user *res.User, err error)
	Update(ctx context.Context, info dbo.UpdateInfo) error
	Read(ctx context.Context, id int64) (*res.User, error)
	ReadByCode(ctx context.Context, code string) (*res.User, error)
	ReadByEmail(ctx context.Context, email string) (*res.User, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, input *dbo.ListInfo, conds ...string) (results []*res.User, count int64, err error)
	Validate(ctx context.Context, input ValidateInput) (*res.User, error)
}
