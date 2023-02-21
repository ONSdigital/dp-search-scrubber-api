package payloads

type ScrubberResp struct {
	Time    string  `json:"time"`
	Query   string  `json:"query"`
	Results Results `json:"results"`
}

type Results struct {
	Areas      []*AreaResp     `json:"areas"`
	Industries []*IndustryResp `json:"industries"`
}

type AreaResp struct {
	Name       string            `json:"name"`
	Region     string            `json:"region"`
	RegionCode string            `json:"region_code"`
	Codes      map[string]string `json:"codes"`
}

type IndustryResp struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
