package model

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"

	"github.com/fidelfly/fxgos/cmd/service/api/iam"
	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/gostool/db"
)

/*LoadPolicy(model model.Model) error
// SavePolicy saves all policy rules to the storage.
SavePolicy(model model.Model) error

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
AddPolicy(sec string, ptype string, rule []string) error
// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
RemovePolicy(sec string, ptype string, rule []string) error
// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error*/

type Adapter struct {
	ResourceType string
	InitData     []byte
}

func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}

func (a *Adapter) LoadPolicy(m model.Model) (err error) {
	err = a.LoadInitPolicy(m)

	if err != nil {
		return
	}

	err = a.LoadDbPolicy(m)

	return
}

func (a *Adapter) LoadInitPolicy(m model.Model) error {
	if len(a.InitData) > 0 {
		scanner := bufio.NewScanner(bytes.NewReader(a.InitData))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			persist.LoadPolicyLine(line, m)
		}
		return scanner.Err()
	}
	return nil
}

type userRole struct {
	Id         int64
	SuperAdmin bool
	Roles      []int64 `xorm:"json"`
}

func (a *Adapter) LoadDbPolicy(m model.Model) error {
	userRoles := make([]*userRole, 0)
	err := db.Find(&userRoles, db.Table("user"))
	if err != nil {
		return err
	}
	if len(userRoles) > 0 {
		for _, user := range userRoles {
			if len(user.Roles) > 0 {
				for _, role := range user.Roles {
					m.AddPolicy("g", "g", []string{iam.EncodeUserSubject(user.Id), iam.EncodeRoleSubject(role)})
				}
			}
			if user.SuperAdmin {
				m.AddPolicy("g", "g", []string{iam.EncodeUserSubject(user.Id), "sa"})
			}
		}
	}

	policies := make([]*res.Policy, 0)
	err = db.Find(&policies, db.Where("resource_type = ?", a.ResourceType))
	if err != nil {
		return err
	}

	if len(policies) > 0 {
		for _, p := range policies {
			if len(p.Act) > 0 {
				for _, act := range p.Act {
					m.AddPolicy("p", "p", []string{p.Sub, p.Obj, act})
				}
			}
		}
	}
	return nil
}

func (a *Adapter) SavePolicy(model model.Model) error {
	return nil
}

func NewAdapter(resType string, data []byte) *Adapter {
	a := &Adapter{
		resType,
		data,
	}
	return a
}
