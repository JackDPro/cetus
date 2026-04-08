package jwt

import (
	"time"
)

type ValidToken struct {
	Id        string    `json:"id" gorm:"primaryKey;type:char(26)"`
	UserId    string    `json:"user_id"`
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
