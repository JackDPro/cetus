package model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type BaseModel struct {
	gorm.Model
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
		if tag != "" {
			modelMap[tag] = val.Field(i).Interface()
		}
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
