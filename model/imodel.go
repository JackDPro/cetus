package model

type IModel interface {
	ToMap() (map[string]interface{}, error)
}
