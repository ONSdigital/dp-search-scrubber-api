package models

type ScrubberResp struct {
	Time    string  `json:"time"`
	Query   string  `json:"query"`
	Results Results `json:"results,omitempty"`
}

type Results struct {
	Areas      []AreaResp     `json:"areas,omitempty"`
	Industries []IndustryResp `json:"industries,omitempty"`
}

type AreaResp struct {
	Name       string            `json:"name,omitempty"`
	Region     string            `json:"region,omitempty"`
	RegionCode string            `json:"region_code,omitempty"`
	Codes      map[string]string `json:"codes,omitempty"`
}

type IndustryResp struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}
