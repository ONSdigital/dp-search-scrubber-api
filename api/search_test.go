package api

import (
	"testing"

	"github.com/ONSdigital/dp-nlp-search-scrubber/api/mock"
	"github.com/ONSdigital/dp-nlp-search-scrubber/models"
	"github.com/stretchr/testify/assert"
)

func TestEmptyDB(t *testing.T) {
	mockDB := mock.EmptyDB()

	tests := []struct {
		name          string
		query         []string
		expectedCodes []string
	}{
		{
			name:          "query with empty db",
			query:         []string{"ind1"},
			expectedCodes: []string{},
		},
		{
			name:          "empty query with empty db",
			query:         []string{},
			expectedCodes: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchingIndustries := getAllMatchingIndustries(tt.query, &mockDB)
			assert.Equal(t, len(tt.expectedCodes), len(matchingIndustries), "expected %d matching industries, got %d", len(tt.expectedCodes), len(matchingIndustries))
			for i, industryResp := range matchingIndustries {
				assert.Equal(t, tt.expectedCodes[i], industryResp.Code, "expected industry with code %s, got %s", tt.expectedCodes[i], industryResp.Code)
			}
		})
	}
}

func TestGetAllMatchingIndustries(t *testing.T) {
	// get a mock ScrubberDB with some industries
	mockDB := mock.DB()

	tests := []struct {
		name          string
		query         []string
		expectedCodes []string
	}{
		{
			name:          "matching single query",
			query:         []string{"ind1"},
			expectedCodes: []string{"IND1"},
		},
		{
			name:          "matching multiple queries",
			query:         []string{"ind1", "ind2"},
			expectedCodes: []string{"IND1", "IND2"},
		},
		{
			name:          "no matching queries",
			query:         []string{"foo", "bar"},
			expectedCodes: []string{},
		},
		// algorithm of PrefixMap is depth first search
		// so It will get the data in reverse
		// keep that in mind when updating tests
		{
			name:          "matching partial query",
			query:         []string{"ind"},
			expectedCodes: []string{"IND3", "IND2", "IND1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchingIndustries := getAllMatchingIndustries(tt.query, mockDB)
			assert.Equal(t, len(tt.expectedCodes), len(matchingIndustries), "expected %d matching industries, got %d", len(tt.expectedCodes), len(matchingIndustries))
			for i, industryResp := range matchingIndustries {
				assert.Equal(t, tt.expectedCodes[i], industryResp.Code, "expected industry with code %s, got %s", tt.expectedCodes[i], industryResp.Code)
			}
		})
	}
}

func TestGetAllMatchingAreas(t *testing.T) {
	// get a mock ScrubberDB with some areas
	mockDB := mock.DB()

	tests := []struct {
		name          string
		query         []string
		expectedNames []*models.AreaResp
	}{
		{
			name:  "matching single query",
			query: []string{"OAC1"},
			expectedNames: []*models.AreaResp{
				{
					Name:       "LAN1",
					Region:     "RN1",
					RegionCode: "RC1",
					Codes: map[string]string{
						"OAC1": "OAC1",
					},
				},
			},
		},
		{
			name:  "matching multiple queries",
			query: []string{"OAC1", "OAC2"},
			expectedNames: []*models.AreaResp{
				{
					Name:       "LAN1",
					Region:     "RN1",
					RegionCode: "RC1",
					Codes: map[string]string{
						"OAC1": "OAC1",
					},
				},
				{
					Name:       "LAN2",
					Region:     "RN2",
					RegionCode: "RC2",
					Codes: map[string]string{
						"OAC2": "OAC2",
					},
				},
			},
		},
		{
			// PrefixMap algorithm is depth first search
			// so when running partial queries it will get
			// the last area first, keep that in mind when updating tests
			name:  "matching partial queries",
			query: []string{"OAC"},
			expectedNames: []*models.AreaResp{
				{
					Name:       "LAN3",
					Region:     "RN3",
					RegionCode: "RC3",
					Codes: map[string]string{
						"OAC3": "OAC3",
					},
				},
				{
					Name:       "LAN2",
					Region:     "RN2",
					RegionCode: "RC2",
					Codes: map[string]string{
						"OAC2": "OAC2",
					},
				},
				{
					Name:       "LAN1",
					Region:     "RN1",
					RegionCode: "RC1",
					Codes: map[string]string{
						"OAC1": "OAC1",
					},
				},
			},
		},
		{
			name:          "no matching queries",
			query:         []string{"foo", "bar"},
			expectedNames: []*models.AreaResp{},
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchingAreas := getAllMatchingAreas(tt.query, mockDB)

			assert.Equal(t, len(tt.expectedNames), len(matchingAreas),
				"expected %d matching areas, got %d", len(tt.expectedNames), len(matchingAreas))

			for i, areaResp := range matchingAreas {
				assert.Equal(t, tt.expectedNames[i].Name, areaResp.Name)
				assert.Equal(t, tt.expectedNames[i].Region, areaResp.Region)
				assert.Equal(t, tt.expectedNames[i].RegionCode, areaResp.RegionCode)
			}
		})
	}
}
