package model

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type BaseModel struct {
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	model := tx.Statement.ReflectValue
	if model.Kind() == reflect.Ptr {
		model = model.Elem()
	}
	if model.Kind() != reflect.Struct {
		return nil
	}
	idField := model.FieldByName("Id")
	if idField.IsValid() && idField.CanSet() && idField.Kind() == reflect.String && idField.String() == "" {
		idField.SetString(ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String())
	}
	return nil
}

func (b *BaseModel) ToJson(model interface{}) (string, error) {
	modelMap, err := b.ToMap(model)
	if err != nil {
		return "", err
	}
	dataStr, err := json.Marshal(modelMap)
	if err != nil {
		return "", err
	}
	return string(dataStr), nil
}

func (b *BaseModel) ToMap(model interface{}) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct or pointer to struct")
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name, opts, _ := strings.Cut(tag, ",")
		if name == "" || name == "-" {
			continue
		}
		fv := val.Field(i)
		if strings.Contains(opts, "omitempty") && fv.IsZero() {
			continue
		}
		modelMap[name] = fv.Interface()
	}
	return modelMap, nil
}

func (b *BaseModel) Include(model interface{}, includeStr string, db *gorm.DB) *gorm.DB {
	includeStr = strings.TrimSpace(includeStr)
	if includeStr == "" {
		return db
	}
	relations := strings.Split(includeStr, ",")
	for i := 0; i < len(relations); i++ {
		relation := relations[i]
		relation = strings.ToUpper(relation[:1]) + relation[1:]
		rType := reflect.TypeOf(model)
		if rType.Kind() == reflect.Ptr {
			rType = rType.Elem()
		}
		_, ok := rType.FieldByName(relation)
		if ok {
			db = db.Preload(relation)
		}
	}
	return db
}
