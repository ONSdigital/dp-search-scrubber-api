package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/ONSdigital/dp-nlp-search-scrubber/api/mock"
	"github.com/ONSdigital/dp-nlp-search-scrubber/payloads"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestGetAllMatchingIndustries(t *testing.T) {
	// Split API tests from Unit tests
	skipUnitTests(t)
	// get a mock ScrubberDB with some industries
	mockDB := mock.Db()

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
	// Split API tests from Unit tests
	skipUnitTests(t)
	// get a mock ScrubberDB with some areas
	mockDB := mock.Db()
	tests := []struct {
		name          string
		query         []string
		expectedNames []*payloads.AreaResp
	}{
		{
			name:  "matching single query",
			query: []string{"OAC1"},
			expectedNames: []*payloads.AreaResp{
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
			expectedNames: []*payloads.AreaResp{
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
			expectedNames: []*payloads.AreaResp{
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
			expectedNames: []*payloads.AreaResp{},
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

func TestPrefixSearchHandler(t *testing.T) {
	// Create a new scrubberDB with mock data
	scrubberDB := mock.Db()

	// Create a new request with the "q" query parameter
	query := url.Values{}
	query.Set("q", "OAC IND")
	req := httptest.NewRequest("GET", "/scrubber/search?"+query.Encode(), nil)

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Call the handler function with the mock data and request
	handler := PrefixSearchHandler(context.Background(), scrubberDB)
	handler(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body into a ScrubberResp struct
	var resp payloads.ScrubberResp
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.Nil(t, err)

	// Check the response data
	assert.Equal(t, "OAC IND", resp.Query)
	assert.Len(t, resp.Results.Areas, 3)
	assert.Len(t, resp.Results.Industries, 3)
}
func skipUnitTests(t *testing.T) {
	if os.Getenv("UNIT") != "" {
		t.Skip("Skipping Unit tests in CI environment")
	}
}
