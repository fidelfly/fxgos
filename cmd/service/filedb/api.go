package filedb

import (
	"context"
	"errors"

	"github.com/fidelfly/gox/pkg/filex"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/filedb/res"
	"github.com/fidelfly/fxgos/cmd/utilities/mctx"
)

func Save(ctx context.Context, name string, data []byte) (int64, error) {
	if len(data) == 0 {
		return 0, errors.New("data is empty")
	}
	md5 := filex.CalculateBytesMd5(data)
	resFile := &res.File{
		Md5: md5,
	}
	if find, err := db.Read(resFile); err != nil {
		return 0, err
	} else if !find {
		resFile.Name = name
		resFile.Data = data
		resFile.Size = int64(len(data))
		resFile.CreateUser = mctx.GetUserId(ctx)
		if id, err := db.Create(resFile); err != nil {
			return 0, err
		} else {
			return id, nil
		}
	} else {
		return resFile.Id, nil
	}
}

func Read(ctx context.Context, id int64) (*res.File, error) {
	if id <= 0 {
		return nil, errors.New("invalid value of id")
	}
	resFile := &res.File{Id: id}
	if find, err := db.Read(resFile); err != nil {
		return nil, err
	} else if find {
		return resFile, nil
	}
	return nil, nil
}
