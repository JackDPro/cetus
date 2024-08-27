package jwt

import (
	"gorm.io/gorm"
	"time"
)

type ValidToken struct {
	gorm.Model
	Id        string    `json:"id"`
	UserId    uint64    `json:"user_id"`
	Token     string    `json:"token"`
	Type      string    `json:"type"`
	Audience  string    `json:"audience"`
	Now       time.Time `json:"now"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (token *ValidToken) ToJsonMap() map[string]interface{} {
	modelMap := make(map[string]interface{})
	modelMap["id"] = token.Id
	modelMap["user_id"] = token.UserId
	modelMap["token"] = token.Token
	modelMap["type"] = token.Type
	modelMap["expired_at"] = token.ExpiredAt
	modelMap["now"] = token.Now
	modelMap["audience"] = token.Audience
	return modelMap
}
