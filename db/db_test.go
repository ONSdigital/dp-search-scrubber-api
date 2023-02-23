package db

import (
	"testing"

	"github.com/ONSdigital/dp-nlp-search-scrubber/config"
	"github.com/ONSdigital/dp-nlp-search-scrubber/db/mock"
	"github.com/stretchr/testify/assert"
)

func TestLoadCsvData(t *testing.T) {
	// Split API tests from Unit tests
	skipUnitTests(t)

	m := mock.CreateFiles(t)
	defer m.CloseFiles()

	cfg := config.Config{
		AreaDataFile:     "area.csv",
		IndustryDataFile: "industry.csv",
	}

	sr := LoadCsvData(&cfg)

	expectedAreas := []struct {
		OutputAreaCode     string
		LocalAuthorityCode string
		LAName             string
		RegionCode         string
		RegionName         string
	}{
		{
			OutputAreaCode:     "Test Output Area Code1",
			LocalAuthorityCode: "Test LAC1",
			LAName:             "Test LAN1",
			RegionCode:         "Test RC1",
			RegionName:         "Test RN 1",
		},
		{
			OutputAreaCode:     "Test Output Area Code2",
			LocalAuthorityCode: "Test LAC2",
			LAName:             "Test LAN2",
			RegionCode:         "Test RC2",
			RegionName:         "Test RN 2",
		},
	}

	// check if the function returns the expected result
	assert.NotNil(t, sr.AreasPFM)
	assert.NotNil(t, sr.IndustriesPFM)

	for _, e := range expectedAreas {
		matchingRecords := sr.AreasPFM.GetByPrefix(e.OutputAreaCode)
		assert.NotEqual(t, len(matchingRecords), 0)
		for _, mr := range matchingRecords {
			area := mr.(*Area)
			assert.Equal(t, area.RegionCode, e.RegionCode)
			assert.Equal(t, area.LocalAuthorityCode, e.LocalAuthorityCode)
			assert.Equal(t, area.LAName, e.LAName)
			assert.Equal(t, area.RegionName, e.RegionName)
			assert.Equal(t, area.OutputAreaCode, e.OutputAreaCode)
		}
	}

	expectedIndustries := []struct {
		code string
		name string
	}{
		{
			code: "TestCode1",
			name: "TestName1",
		},
		{
			code: "TestCode2",
			name: "TestName2",
		},
	}

	for _, e := range expectedIndustries {
		matchingRecords := sr.IndustriesPFM.GetByPrefix(e.code)
		assert.NotEqual(t, len(matchingRecords), 0)
		for _, mr := range matchingRecords {
			industry := mr.(*Industry)
			assert.Equal(t, industry.Code, e.code)
			assert.Equal(t, industry.Name, e.name)
		}
	}
}
