package model

type Pagination struct {
	Count       int   `json:"count"`
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}

func (p *Pagination) ToJsonMap() map[string]interface{} {
	modelMap := make(map[string]interface{})
	modelMap["count"] = p.Count
	modelMap["current_page"] = p.CurrentPage
	modelMap["page_size"] = p.PageSize
	modelMap["total"] = p.Total
	modelMap["total_pages"] = p.TotalPages
	return modelMap
}
