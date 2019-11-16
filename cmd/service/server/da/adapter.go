package da

import (
	"github.com/casbin/casbin/model"

	"github.com/fidelfly/gostool/db"
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
	//users, _, err := user.List(context.Background(), nil, fmt.Sprintf("super_admin = %d", 1))
	dbs := db.NewSession()
	defer dbs.Close()
	uids := make([]int64, 0)
	err := dbs.Find(&uids, db.Table("user"), db.Cols("id"), db.Where("super_admin = ? ", true))
	if err != nil {
		return err
	}
	if len(uids) > 0 {
		for _, userId := range uids {
			m.AddPolicy("g", "g", []string{encodeUserKey(userId), "sa"})

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
