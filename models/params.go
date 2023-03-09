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

func GetScrubberParams(query url.Values) *ScrubberParams {
	result := ScrubberParams{
		Query: "",
		SIC:   []string{},
		OAC:   []string{},
	}

	if len(query["q"]) != 0 {
		result.Query = query["q"][0]
	}
	result.removeSpecialCharacters()
	result.populateOACandSICcodes()

	return &result
}

func (sp *ScrubberParams) removeSpecialCharacters() {
	re := regexp.MustCompile("[^a-zA-Z0-9_]+")

	sp.Query = re.ReplaceAllString(sp.Query, " ")
}

func (sp *ScrubberParams) populateOACandSICcodes() {
	querySl := strings.Split(sp.Query, " ")

	sicCodeRe := regexp.MustCompile(`^\d{5}$`)
	oacCodeRe := regexp.MustCompile(`^[a-zA-Z]\d{8}$`)

	cache := make(map[string]string)
	for _, v := range querySl {
		fmt.Printf("v: %s\n", v)

		if _, ok := cache[v]; !ok && sicCodeRe.MatchString(v) {
			cache[v] = v
			sp.SIC = append(sp.SIC, v)
		}

		if _, ok := cache[v]; !ok && oacCodeRe.MatchString(v) {
			cache[v] = v
			sp.OAC = append(sp.OAC, v)
		}
		fmt.Println(sp)
	}
}
