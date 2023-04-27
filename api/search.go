package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-nlp-search-scrubber/db"
	"github.com/ONSdigital/dp-nlp-search-scrubber/models"
	"github.com/ONSdigital/log.go/v2/log"
)

func FindAllMatchingAreasAndIndustriesHandler(scrubberDB *db.ScrubberDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx := r.Context()

		log.Info(ctx, "api contains /scrubber/search endpoint which return a list of possible locations and industries based on OAC and SIC")

		start := time.Now()

		if len(scrubberDB.AreasPFM.Children) == 0 && len(scrubberDB.IndustriesPFM.Children) == 0 {
			log.Error(ctx, "There is no data to display due to a database issue", fmt.Errorf("missing raw data"))
			w.Header().Set("X-Error-Message", "There was an issue with the database")
			w.WriteHeader(http.StatusNoContent)

			return
		}

		scrubberParams := models.GetScrubberParams(r.URL.Query())

		matchingAreas := getAllMatchingAreas(scrubberParams.OAC, scrubberDB)
		matchingIndustries := getAllMatchingIndustries(scrubberParams.SIC, scrubberDB)

		scrubberResp := models.ScrubberResp{
			Time:  fmt.Sprint(time.Since(start).Microseconds(), "µs"),
			Query: scrubberParams.Query,
			Results: models.Results{
				Areas:      matchingAreas,
				Industries: matchingIndustries,
			},
		}

		if err := json.NewEncoder(w).Encode(scrubberResp); err != nil {
			log.Error(ctx, "Unable to encode the response data", err)

			w.WriteHeader(http.StatusInternalServerError)

			errObj := ErrorResp{
				errors: []Errors{
					{
						error_code: "", // to be added once Nathan finished the error-codes lib
						message:    "An unexpected error occurred while processing your request",
					},
				},
				trace_id: getRequestId(ctx),
			}

			if err := json.NewEncoder(w).Encode(errObj); err != nil {
				log.Fatal(ctx, "cannot encode errObj", err)
			}
		}
	}
}

func getAllMatchingAreas(querySl []string, ScrubberDB *db.ScrubberDB) []*models.AreaResp {
	var matchingAreas []*models.AreaResp

	areaRespMap := make(map[string]*models.AreaResp)

	for _, q := range querySl {
		matchingRecords := ScrubberDB.AreasPFM.GetByPrefix(strings.ToUpper(q))

		for _, rData := range matchingRecords {
			area := rData.(*db.Area)
			key := area.LAName + area.RegionName + area.RegionCode

			if _, found := areaRespMap[key]; found {
				areaRespMap[key].Codes[area.OutputAreaCode] = area.OutputAreaCode
			} else {
				areaResp := &models.AreaResp{
					Name:       area.LAName,
					Region:     area.RegionName,
					RegionCode: area.RegionCode,
					Codes: map[string]string{
						area.OutputAreaCode: area.OutputAreaCode,
					},
				}

				areaRespMap[key] = areaResp
				matchingAreas = append(matchingAreas, areaRespMap[key])
			}
		}
	}

	return matchingAreas
}

func getAllMatchingIndustries(querySl []string, ScrubberDB *db.ScrubberDB) []*models.IndustryResp {
	var matchingIndustries []*models.IndustryResp

	validation := make(map[string]string)

	for _, q := range querySl {
		matchingRecords := ScrubberDB.IndustriesPFM.GetByPrefix(strings.ToUpper(q))

		for _, rData := range matchingRecords {
			industry := rData.(*db.Industry)

			if _, valid := validation[industry.Code]; !valid {
				industryResp := &models.IndustryResp{
					Code: industry.Code,
					Name: industry.Name,
				}

				matchingIndustries = append(matchingIndustries, industryResp)
			}

			validation[industry.Code] = industry.Name
		}
	}

	return matchingIndustries
}

func getRequestId(ctx context.Context) string {
	requestID := ctx.Value(request.RequestIdKey)
	if requestID == nil {
		requestID = ctx.Value("request-id")
	}

	correlationID, ok := requestID.(string)
	if !ok {
		return ""
	}

	return correlationID
}
