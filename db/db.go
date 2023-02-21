package db

import (
	"context"

	"github.com/ONSdigital/dp-nlp-search-scrubber/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/alediaferia/prefixmap"
)

type ScrubberDB struct {
	AreasPFM      *prefixmap.PrefixMap
	IndustriesPFM *prefixmap.PrefixMap
}

func LoadCsvData(cfg *config.Config) *ScrubberDB {
	// setting up a logger
	log.Namespace = "Scrubber DB"
	ctx := context.Background()

	// gets area data
	areaData, err := getArea(cfg)
	if err != nil {
		log.Fatal(ctx, "Error loading Area data: ", err)
	}

	log.Info(ctx, "Successfully loaded Area data")

	// gets industry data
	industryData, err := getIndustry(cfg)
	if err != nil {
		log.Info(ctx, err.Error())
		log.Fatal(ctx, "Error loading Industry data: ", err)
	}

	log.Info(ctx, "Successfully loaded Industry data")

	// creates a new area prefixmap and populates it
	areasMap := prefixmap.New()
	for _, area := range areaData {
		areasMap.Insert(area.OutputAreaCode, area)
	}

	// creates a new industry prefixmap and populates it
	industryMap := prefixmap.New()
	for _, industry := range industryData {
		industryMap.Insert(industry.Code, industry)
	}

	return &ScrubberDB{
		AreasPFM:      areasMap,
		IndustriesPFM: industryMap,
	}
}
