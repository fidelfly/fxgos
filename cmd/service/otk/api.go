package otk

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fidelfly/gox/pkg/randx"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/otk/res"
)

type ResourceKey struct {
	Id int64 `json:"id"`
}

func NewResourceKey(id int64) string {
	if jsonData, err := json.Marshal(&ResourceKey{Id: id}); err == nil {
		return string(jsonData)
	}
	return ""
}

func NewOtk(keyType string, typeId string, expired time.Duration, usage string, data string) (key string, err error) {
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

type updateInput struct {
	db.UpdateInfo
	Data *res.OneTimeKey
}

func update(input updateInput) (int64, error) {
	if input.Data == nil {
		return 0, errors.New("data is empty")
	}
	opts := make([]db.QueryOption, 0)

	if input.Id > 0 {
		opts = append(opts, db.ID(input.Id))
	}
	if len(input.Cols) > 0 {
		opts = append(opts, db.Cols(input.Cols...))
	}

	if rows, err := db.Update(input.Data, opts...); err != nil {
		return 0, err
	} else if rows > 0 {
		return input.Data.Id, nil
	}
	return 0, nil
}

func Consume(id int64) error {
	_, err := update(updateInput{
		UpdateInfo: db.UpdateInfo{Id: id, Cols: []string{"consumed"}},
		Data:       &res.OneTimeKey{Id: id, Consumed: true},
	})
	return err
}

func Validate(key string) (*res.OneTimeKey, error) {
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
