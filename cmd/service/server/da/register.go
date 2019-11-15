package da

import (
	"context"
	"reflect"

	"github.com/fidelfly/gox/logx"
	"github.com/fidelfly/gox/pubsubx"

	"github.com/fidelfly/fxgos/cmd/service/api/da"
	"github.com/fidelfly/fxgos/cmd/service/api/user"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/fxgos/cmd/utilities/pub"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

func (s server) RegisterData(data da.SecurityData) error {
	resType := data.GetResourceType()
	if resType == user.ResourceType {
		pub.Subscribe(pub.TopicResource, userSgHandler(data))
	}
	pub.Subscribe(pub.TopicResource, resSghandler(data))
	return nil
}

func userSgHandler(sd da.SecurityData) pubsubx.SubscriberHandler {
	return func(msg interface{}) error {
		if event, ok := msg.(pub.ResourceEvent); ok {
			if event.Type != sd.GetResourceType() {
				return nil
			}
			ctx, dbs := dbo.WithDBSession(context.Background(), dbo.DefaultSession)
			defer dbs.Close()
			if event.Action == pub.ResourceDelete {
				if _, err := dbo.Delete(ctx, new(res.UserSg), db.Where("user_id = ?", event.Id)); err != nil {
					return err
				}
				removeUserFromDaModel(event.Id)
				_ = sgCache.Delete(encodeUserKey(event.Id))
				return nil
			}
			if resUser, err := user.Read(ctx, event.Id); err == nil {
				groups := resUser.SecurityGroups
				updateDaModelByUser(resUser.Id, resUser.SuperAdmin)
				if event.Action == pub.ResourceCreate || !reflect.DeepEqual(groups, getUserSg(ctx, event.Id)) {
					if err := dbs.Begin(); err == nil {
						if _, err = dbo.Delete(ctx,
							new(res.UserSg),
							db.Where("user_id = ?", event.Id)); err != nil {
							logx.CaptureError(dbs.Rollback())
							return err
						}

						groups = append(groups, 0)
						sgs := make([]res.UserSg, len(groups))
						for i := 0; i < len(sgs); i++ {
							sgs[i] = res.UserSg{
								UserId:        event.Id,
								SecurityGroup: groups[i],
							}
						}

						if err = dbo.Create(ctx, &sgs); err != nil {
							logx.CaptureError(dbs.Rollback())
							return err
						}

						if err = dbs.Commit(); err != nil {
							return err
						}
						_ = sgCache.Set(encodeUserKey(event.Id), userSg{
							UserId:         event.Id,
							SecurityGroups: groups,
						})
					} else {
						return err
					}
				}
			} else {
				return err
			}
		}
		return nil
	}
}

func resSghandler(sd da.SecurityData) pubsubx.SubscriberHandler {
	return func(msg interface{}) error {
		if event, ok := msg.(pub.ResourceEvent); ok {
			if event.Type != sd.GetResourceType() {
				return nil
			}
			ctx, dbs := dbo.WithDBSession(context.Background(), dbo.DefaultSession)
			defer dbs.Close()
			if event.Action == pub.ResourceDelete {
				if _, err := dbo.Delete(ctx, new(res.ResourceSg),
					db.Where("res_type = ? and res_id = ?", sd.GetResourceType(), event.Id)); err != nil {
					return err
				}
				_ = sgCache.Delete(encodeResKey(sd.GetResourceType(), event.Id))
				return nil
			}

			if groups, err := sd.GetSecurityGroups(ctx, event.Id); err == nil {
				if event.Action == pub.ResourceCreate || !reflect.DeepEqual(groups, getResourceSg(ctx, sd.GetResourceType(), event.Id)) {
					if err := dbs.Begin(); err == nil {
						if _, err = dbo.Delete(ctx,
							new(res.ResourceSg),
							db.Where("res_type = ? and res_id = ?", sd.GetResourceType(), event.Id)); err != nil {
							logx.CaptureError(dbs.Rollback())
							return err
						}
						if len(groups) == 0 {
							groups = append(groups, 0)
						}

						sgs := make([]res.ResourceSg, len(groups))
						for i := 0; i < len(groups); i++ {
							sgs[i] = res.ResourceSg{
								ResType:       sd.GetResourceType(),
								ResId:         event.Id,
								SecurityGroup: groups[i],
							}
						}

						if err = dbo.Create(ctx, &sgs); err != nil {
							logx.CaptureError(dbs.Rollback())
							return err
						}

						if err = dbs.Commit(); err != nil {
							return err
						}

						_ = sgCache.Set(encodeResKey(sd.GetResourceType(), event.Id), resourceSg{
							ResType:        sd.GetResourceType(),
							ResId:          event.Id,
							SecurityGroups: groups,
						})
					} else {
						return err
					}

				}
			} else {
				return err
			}
		}

		return nil
	}
}

func getResourceSg(ctx context.Context, resType string, id int64) []int64 {
	sgs := make([]int64, 0)

	if err := dbo.Find(ctx, &sgs,
		db.Table(new(res.ResourceSg)),
		db.Cols("security_group"),
		db.Where("res_type = ? and res_id = ? and security_group > 0", resType, id)); err != nil {
		logx.Error(err)
	}
	return sgs
}

func getUserSg(ctx context.Context, userId int64) []int64 {
	sgs := make([]int64, 0)

	if err := dbo.Find(ctx, &sgs,
		db.Table(new(res.UserSg)),
		db.Cols("security_group"),
		db.Where("user_id = ? and security_group > 0", userId)); err != nil {
		logx.Error(err)
	}
	return sgs
}
