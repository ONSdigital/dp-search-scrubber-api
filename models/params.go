package models

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type ScrubberParams struct {
	Query string
	SIC   []string
	OAC   []string
}

func GetScrubberParams(query url.Values) (*ScrubberParams, error) {
	result := ScrubberParams{
		Query: "",
		SIC:   []string{},
		OAC:   []string{},
	}

	if len(query) != 1 {
		return nil, fmt.Errorf("one query expected, found multiple queries ")
	}

	if len(query["q"]) > 1 {
		return nil, fmt.Errorf("one query expected, found multiple queries with the same name ")
	}

	if len(query["q"]) == 0 {
		return nil, fmt.Errorf("no query provided or wrong query name")
	}

	result.Query = query["q"][0]
	result.rmSpecialCharsFromQuery()
	result.getAllAcceptableCodesFromQuery()

	return &result, nil
}

func (sp *ScrubberParams) rmSpecialCharsFromQuery() {
	re := regexp.MustCompile("[^a-zA-Z0-9]+")

	sp.Query = re.ReplaceAllString(sp.Query, " ")
}

func (sp *ScrubberParams) getAllAcceptableCodesFromQuery() {
	querySl := strings.Split(sp.Query, " ")

	// regex for how a sic code looks like e.g. 12345
	sicCodeRe := regexp.MustCompile(`^\d{5}$`)

	// regex for how a output area code looks like e.g. E12345678
	oacCodeRe := regexp.MustCompile(`^[a-zA-Z]\d{8}$`)

	// cache is here to make sure we don't duplicate entries
	cache := make(map[string]string)
	for _, v := range querySl {
		if _, ok := cache[v]; !ok && sicCodeRe.MatchString(v) {
			cache[v] = v
			sp.SIC = append(sp.SIC, v)
		}

		if _, ok := cache[v]; !ok && oacCodeRe.MatchString(v) {
			cache[v] = v
			sp.OAC = append(sp.OAC, v)
		}
	}
}
