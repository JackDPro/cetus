package model

type AccessToken struct {
	BaseModel
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"type"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (m *AccessToken) Fields() []string {
	return []string{"AccessToken", "RefreshToken", "Type", "ExpiresIn"}
}

func (m *AccessToken) ToMap() (map[string]interface{}, error) {
	return m.BaseModel.ToMap(m)
}
