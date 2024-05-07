package model

type Probe struct {
	BaseModel
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	ConfigPath string `json:"config_path"`
}

func (m *Probe) Fields() []string {
	return []string{"AppName", "AppVersion", "ConfigPath"}
}

func (m *Probe) ToMap() (map[string]interface{}, error) {
	return m.BaseModel.ToMap(m)
}
