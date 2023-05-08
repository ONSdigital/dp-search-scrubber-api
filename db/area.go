package db

import (
	"os"

	"github.com/ONSdigital/dp-nlp-search-scrubber/config"
	"github.com/gocarina/gocsv"
)

type Area struct {
	OutputAreaCode     string `csv:"Output Area Code"`
	LocalAuthorityCode string `csv:"Local Authority Code"`
	LAName             string `csv:"Local Authority Name"`
	RegionCode         string `csv:"Region/Country Code"`
	RegionName         string `csv:"Region/Country Name"`
}

func getArea(cfg *config.Config) ([]*Area, error) {
	file, err := os.Open(cfg.AreaDataFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	ar := []*Area{}

	if err := gocsv.UnmarshalFile(file, &ar); err != nil {
		return nil, err
	}

	return ar, nil
}
