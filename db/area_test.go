package db

import (
	"testing"

	"github.com/ONSdigital/dp-search-scrubber-api/config"
	"github.com/ONSdigital/dp-search-scrubber-api/db/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetArea(t *testing.T) {
	// Split API tests from Unit tests
	skipUnitTests(t)

	// mock data files
	m := mock.CreateFiles(t)
	defer m.CloseFiles()

	cfg := config.Config{
		AreaDataFile: "area.csv",
	}

	// create a test file with test data and write test data to the test file
	ar, err := getArea(&cfg)
	if err != nil {
		t.Fatalf("there was an error geting the area: %v ", err.Error())
	}

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

	for i, expected := range expectedAreas {
		assert.Equal(t, expected.OutputAreaCode, ar[i].OutputAreaCode, "OutputAreaCode does not match expected value")
		assert.Equal(t, expected.LocalAuthorityCode, ar[i].LocalAuthorityCode, "LocalAuthorityCode does not match expected value")
		assert.Equal(t, expected.LAName, ar[i].LAName, "LAName does not match expected value")
		assert.Equal(t, expected.RegionCode, ar[i].RegionCode, "RegionCode does not match expected value")
		assert.Equal(t, expected.RegionName, ar[i].RegionName, "RegionName does not match expected value")
	}
}
