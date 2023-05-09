package mock

import (
	"github.com/ONSdigital/dp-nlp-search-scrubber/db"
	"github.com/alediaferia/prefixmap"
)

func Inds() []*db.Industry {
	industries := []*db.Industry{
		{Code: "IND1", Name: "Industry 1"},
		{Code: "IND2", Name: "Industry 2"},
		{Code: "IND3", Name: "Industry 3"},
	}

	return industries
}

func Areas() []*db.Area {
	areas := []*db.Area{
		{
			RegionCode:         "RC1",
			OutputAreaCode:     "OAC1",
			LocalAuthorityCode: "LAC1",
			LAName:             "LAN1",
			RegionName:         "RN1",
		},
		{
			RegionCode:         "RC2",
			OutputAreaCode:     "OAC2",
			LocalAuthorityCode: "LAC2",
			LAName:             "LAN2",
			RegionName:         "RN2",
		},
		{
			RegionCode:         "RC3",
			OutputAreaCode:     "OAC3",
			LocalAuthorityCode: "LAC3",
			LAName:             "LAN3",
			RegionName:         "RN3",
		},
	}

	return areas
}

func DB() *db.ScrubberDB {
	areaData := Areas()
	industryData := Inds()

	areasMap := prefixmap.New()
	for _, area := range areaData {
		areasMap.Insert(area.OutputAreaCode, area)
	}

	industryMap := prefixmap.New()
	for _, industry := range industryData {
		industryMap.Insert(industry.Code, industry)
	}

	return &db.ScrubberDB{
		AreasPFM:      areasMap,
		IndustriesPFM: industryMap,
	}
}
