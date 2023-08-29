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

	result.splitAllAcceptableCodesFromQuery()

	return &result, nil
}

func (sp *ScrubberParams) rmSpecialCharsFromQuery() {
	re := regexp.MustCompile("[^A-Za-z0-9]+")

	sp.Query = re.ReplaceAllString(sp.Query, " ")
}

func (sp *ScrubberParams) splitAllAcceptableCodesFromQuery() {
	querySl := strings.Split(sp.Query, " ")
	sp.Query = ""

	// regex for how a sic code looks like e.g. 12345
	sicCodeRe := regexp.MustCompile(`^\d{5}$`)

	// regex for how a output area code looks like e.g. E12345678
	oacCodeRe := regexp.MustCompile(`^[a-zA-Z]\d{8}$`)

	// cache is here to make sure we don't duplicate entries
	cache := make(map[string]string)
	for _, v := range querySl {
		// if it matches a SIC code
		if _, ok := cache[v]; !ok && sicCodeRe.MatchString(v) {
			cache[v] = v
			sp.SIC = append(sp.SIC, v)
			continue
		}

		// if it matches a OAC code
		if _, ok := cache[v]; !ok && oacCodeRe.MatchString(v) {
			cache[v] = v
			sp.OAC = append(sp.OAC, v)
			continue
		}

		// if it doesn't match a OAC or SIC code and isn't composed of 2 letters
		if _, ok := cache[v]; !ok && len(v) > 2 {
			cache[v] = v

			// first sp.Query is always empty
			if sp.Query == "" {
				sp.Query = v
				continue
			}

			sp.Query = sp.Query + " " + v
		}
	}
}
