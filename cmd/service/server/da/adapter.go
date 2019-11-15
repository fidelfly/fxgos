package da

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/model"

	"github.com/fidelfly/fxgos/cmd/service/api/user"
)

type sgAdapter struct {
}

func (a *sgAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

func (a *sgAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

func (a *sgAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}

func (a *sgAdapter) LoadPolicy(m model.Model) (err error) {

	err = a.LoadDbPolicy(m)

	return
}

func (a *sgAdapter) LoadDbPolicy(m model.Model) error {
	users, _, err := user.List(context.Background(), nil, fmt.Sprintf("super_admin = %d", 1))
	//dbs := db.NewSession()
	//defer dbs.Close()
	//err := dbs.Find(&userRoles, db.Table("user"), db.Where("super_admin = ? ", true))
	if err != nil {
		return err
	}
	if len(users) > 0 {
		for _, resUser := range users {
			if resUser.SuperAdmin {
				m.AddPolicy("g", "g", []string{encodeUserKey(resUser.Id), "sa"})
			}
		}
	}
	return nil
}

func (a *sgAdapter) SavePolicy(model model.Model) error {
	return nil
}

func newAdapter() *sgAdapter {
	return &sgAdapter{}
}
