package model

import (
	"encoding/json"
)

type DataWrapper struct {
	Data interface{} `json:"data"`
	Meta *Meta       `json:"meta"`
}

func (data *DataWrapper) ToString() ([]byte, error) {
	dataStr, err := json.Marshal(data.Data)
	if err != nil {
		return nil, err
	}
	return dataStr, nil
}
