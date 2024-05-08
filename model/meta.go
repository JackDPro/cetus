package model

type Meta struct {
	Pagination *Pagination `json:"pagination"`
}

func (p *Meta) IsNull() bool {
	if p.Pagination != nil {
		return false
	}
	return true
}

func (p *Meta) ToJsonMap() map[string]interface{} {
	modelMap := make(map[string]interface{})
	modelMap["pagination"] = p.Pagination
	return modelMap
}
