package otk

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fidelfly/gox/pkg/randx"

	"github.com/fidelfly/fxgos/cmd/service/res"
	"github.com/fidelfly/fxgos/cmd/utilities/syserr"
	"github.com/fidelfly/gostool/db"
	"github.com/fidelfly/gostool/dbo"
)

type server struct {
}

func (s server) NewOtk(keyType string, typeId string, expired time.Duration, usage string, data string) (key string, err error) {
	otkData := &res.OneTimeKey{
		Type:   keyType,
		TypeId: typeId,
		Usage:  usage,
		Data:   data,
	}

	otkData.Key, err = randx.GetString(12)
	if err != nil {
		return
	}
	otkData.CreateTime = time.Now()

	dbSession := db.Engine.NewSession()
	defer dbSession.Close()
	err = dbSession.Begin()
	if err != nil {
		return
	}

	if _, err = dbSession.Insert(otkData); err != nil {
		return
	}
	if _, err = dbSession.Exec("update one_time_key set invalid = 1 where id != ? and type = ? and type_id = ?", otkData.Id, otkData.Type, otkData.TypeId); err != nil {
		_ = dbSession.Rollback()
		return
	}

	claims := &jwt.StandardClaims{
		Id:        strconv.FormatInt(otkData.Id, 10),
		IssuedAt:  otkData.CreateTime.Unix(),
		ExpiresAt: otkData.CreateTime.Add(expired).Unix(),
		Issuer:    otkData.Type,
		Subject:   otkData.Usage,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if key, err = token.SignedString([]byte(otkData.Key)); err != nil {
		_ = dbSession.Rollback()
	} else {
		_ = dbSession.Commit()
	}
	return
}

func (s server) update(input dbo.UpdateInfo) error {
	if input.Data == nil {
		return errors.New("data is empty")
	}
	var otk *res.OneTimeKey
	if t, ok := input.Data.(*res.OneTimeKey); ok {
		otk = t
	} else {
		otk = new(res.OneTimeKey)
	}
	opts := dbo.ApplyUpdateOption(otk, input)

	if rows, err := dbo.Update(context.Background(), otk, opts...); err != nil {
		return syserr.DatabaseErr(err)
	} else if rows == 0 {
		return syserr.ErrNotFound
	}
	return nil
}

func (s server) Consume(id int64) error {
	return s.update(dbo.UpdateInfo{
		Id:   id,
		Cols: []string{"consumed"},
		Data: &res.OneTimeKey{Id: id, Consumed: true},
	})
}

func (s server) Validate(key string) (*res.OneTimeKey, error) {
	otkv := OneTimeKeyVerifier{
		ResOtk: &res.OneTimeKey{},
	}
	if token, err := jwt.ParseWithClaims(key, &jwt.StandardClaims{}, otkv.KeyFunction); err != nil {
		return nil, err
	} else if token.Valid {
		return otkv.ResOtk, nil
	} else {
		return nil, errors.New("token is invalid")
	}
}

type OneTimeKeyVerifier struct {
	ResOtk *res.OneTimeKey
}

func (otkv *OneTimeKeyVerifier) KeyFunction(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		if id, err := strconv.ParseInt(claims.Id, 10, 64); err != nil {
			return nil, err
		} else {
			otkv.ResOtk.Id = id
		}
		if find, err := db.Read(otkv.ResOtk); err != nil {
			return nil, err
		} else if !find || otkv.ResOtk.Invalid || otkv.ResOtk.Consumed {
			return nil, errors.New("one time key is invalid")
		}
		return []byte(otkv.ResOtk.Key), nil
	}
	return nil, fmt.Errorf("one time key is invalid")
}
