package iam

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/fidelfly/gox/cachex/bcache"
	"github.com/fidelfly/gox/logx"
	"github.com/fidelfly/gox/pkg/filex"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/iam/res"
)

type Resource struct {
	Type    string   `json:"type"`
	Code    string   `json:"code"`
	Actions []string `json:"actions"`
}

type ResourceACL struct {
	Resource
	ACL []string `json:"acl"`
}

func initIam(resourceType string, model string, policy string, resources ...string) (err error) {

	modelData, err := ioutil.ReadFile(model)
	if err != nil {
		return
	}

	policyData, err := ioutil.ReadFile(policy)
	if err != nil {
		return
	}

	m := &res.Model{
		ResourceType: resourceType,
	}

	if ok, err := db.Read(m); err != nil {
		return err
	} else if ok {
		updateCols := make([]string, 0)
		if filex.CalculateBytesMd5(m.Data) != filex.CalculateBytesMd5(modelData) {
			m.Data = modelData
			updateCols = append(updateCols, "data")

		}
		if filex.CalculateBytesMd5(m.Policy) != filex.CalculateBytesMd5(policyData) {
			m.Policy = policyData
			updateCols = append(updateCols, "policy")
		}
		if len(updateCols) > 0 {
			_, err := db.Update(m, db.ID(m.Id), db.Cols(updateCols...))
			if err != nil {
				return err
			}
		}
	} else {
		m.Data = modelData
		m.Policy = policyData
		if _, err := db.Create(m); err != nil {
			return err
		}
	}

	err = initIamResource(resources...)
	return
}

func initIamResource(resources ...string) error {
	if len(resources) > 0 {
		type IamResource struct {
			Type      string
			Resources []Resource
		}
		for _, resource := range resources {
			resFile := &IamResource{}
			if _, err0 := toml.DecodeFile(resource, resFile); err0 == nil {
				if len(resFile.Resources) > 0 {
					for _, res0 := range resFile.Resources {
						if len(res0.Type) == 0 {
							res0.Type = resFile.Type
						}
						err := resDB.Set(bcache.NewKey(res0.Type, res0.Code), res0)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

//type ResourceData struct {
//	Type    string
//	Code    string
//	Actions []string
//}

var modelRegex, _ = regexp.Compile(`(.*)\.model\.conf`)
var policyRegex, _ = regexp.Compile(`(.*)\.policy\.csv`)
var iamResRegex, _ = regexp.Compile(`(.*)\.resource\.toml`)

//export
func ScanIam(folder string) {
	if _, err := os.Stat(folder); err != nil {
		return
	}

	type IamConfig struct {
		Type   string
		Model  string
		Policy string
	}

	iamFiles := make(map[string]*IamConfig)
	resFiles := make([]string, 0)
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		iamType := isModelFile(info.Name())
		if len(iamType) > 0 {
			if cfg, ok := iamFiles[iamType]; ok {
				cfg.Model = path
			} else {
				cfg = &IamConfig{
					Type:  iamType,
					Model: path,
				}
				iamFiles[iamType] = cfg
			}
		} else {
			iamType = isPolicyFile(info.Name())
			if len(iamType) > 0 {
				if cfg, ok := iamFiles[iamType]; ok {
					cfg.Policy = path
				} else {
					cfg = &IamConfig{
						Type:   iamType,
						Policy: path,
					}
					iamFiles[iamType] = cfg
				}
			} else if isIamResource(info.Name()) {
				resFiles = append(resFiles, path)
			}
		}
		return nil
	})
	if err != nil {
		logx.Error(err)
		return
	}

	if len(iamFiles) > 0 {
		for k, v := range iamFiles {
			err = initIam(k, v.Model, v.Policy)
			if err != nil {
				return
			}
		}
	}
	if len(resFiles) > 0 {
		err = initIamResource(resFiles...)
		if err != nil {
			return
		}
	}
}

func isModelFile(name string) (resType string) {
	m := modelRegex.FindSubmatch([]byte(name))
	if len(m) == 2 {
		resType = string(m[1])
		return
	}
	return ""
}

func isIamResource(name string) bool {
	return iamResRegex.Match([]byte(name))
}

func isPolicyFile(name string) (resType string) {
	m := policyRegex.FindSubmatch([]byte(name))
	if len(m) == 2 {
		resType = string(m[1])
		return
	}
	return ""
}
