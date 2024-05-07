package model

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type BaseModel struct{}

func (b *BaseModel) Fields() []string {
	return []string{}
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
	modelVal := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model).Elem()

	fieldsMethod := modelVal.MethodByName("Fields")
	if fieldsMethod.IsValid() {
		result := fieldsMethod.Call([]reflect.Value{})
		if len(result) > 0 {
			fields := result[0]
			for i := 0; i < fields.Len(); i++ {
				fieldName := fields.Index(i).String()
				key := fieldName
				field, ok := modelType.FieldByName(fieldName)
				if ok {
					jsonTag := field.Tag.Get("json")
					if jsonTag != "" {
						key = jsonTag
					}
				}
				fieldValue := modelVal.Elem().FieldByName(fieldName)
				if fieldValue.IsValid() && fieldValue.CanInterface() {
					modelMap[key] = fieldValue.Interface()
				}
			}
		}
	} else {
		return nil, errors.New("model need Fields method")
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
