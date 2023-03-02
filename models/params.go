package params

import "net/url"

type ScrubberParams struct {
	Query string `url:"q,omitempty"`
}

func GetScrubberParams(query url.Values) *ScrubberParams {
	result := ScrubberParams{
		Query: "",
	}

	if len(query["q"]) != 0 {
		result.Query = query["q"][0]
	}

	return &result
}
