package test

import (
	"github.com/JackDPro/cetus/model"
	"testing"
)

type User struct {
	model.BaseModel
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string
}

func TestBaseModelToJson(t *testing.T) {
	user := User{
		Id:       123,
		Name:     "Jack",
		Password: "password",
	}
	json, err := user.ToJson(&user)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(json)
}
