package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-search-scrubber-api/db"
	"github.com/ONSdigital/dp-search-scrubber-api/models"
	"github.com/ONSdigital/log.go/v2/log"
)

func FindAllMatchingAreasAndIndustriesHandler(scrubberDB db.ScrubberDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx := r.Context()

		start := time.Now()

		if len(scrubberDB.AreasPFM.Children) == 0 && len(scrubberDB.IndustriesPFM.Children) == 0 {
			log.Error(ctx, "There is no data to display due to a database issue", fmt.Errorf("missing raw data"))

			w.Header().Set("X-Error-Message", "There was an issue with the database")

			w.WriteHeader(http.StatusNoContent)

			errObj := ErrorResp{
				Errors: []Errors{
					{
						ErrorCode: "", // to be added once Nathan finished the error-codes lib
						Message:   "An unexpected error occurred while processing your request",
					},
				},
				TraceID: getRequestID(ctx),
			}

			if err := json.NewEncoder(w).Encode(errObj); err != nil {
				log.Error(ctx, "Unable to encode the error response data", err)
			}

			return
		}

		scrubberParams, err := models.GetScrubberParams(r.URL.Query())
		if err != nil {
			log.Error(ctx, "Error getting scrubber query", err)

			w.WriteHeader(http.StatusBadRequest)

			errObj := ErrorResp{
				Errors: []Errors{
					{
						ErrorCode: "", // to be added once Nathan finished the error-codes lib
						Message:   "An unexpected error occurred while processing your request",
					},
				},
				TraceID: getRequestID(ctx),
			}

			if err := json.NewEncoder(w).Encode(errObj); err != nil {
				log.Error(ctx, "Unable to encode the error response data", err)
			}

			return
		}

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
				Errors: []Errors{
					{
						ErrorCode: "", // to be added once Nathan finished the error-codes lib
						Message:   "An unexpected error occurred while processing your request",
					},
				},
				TraceID: getRequestID(ctx),
			}

			if err := json.NewEncoder(w).Encode(errObj); err != nil {
				log.Error(ctx, "Unable to encode the error response data", err)
			}
		}
	}
}

func getAllMatchingAreas(querySl []string, scrubberDB db.ScrubberDB) []models.AreaResp {
	var matchingAreas []models.AreaResp

	areaRespMap := make(map[string]models.AreaResp)

	for _, q := range querySl {
		matchingRecords := scrubberDB.AreasPFM.GetByPrefix(strings.ToUpper(q))
		for _, rData := range matchingRecords {
			area := rData.(db.Area)
			key := area.LAName + area.RegionName + area.RegionCode

			if _, found := areaRespMap[key]; found {
				areaRespMap[key].Codes[area.OutputAreaCode] = area.OutputAreaCode
			} else {
				areaResp := models.AreaResp{
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

func getAllMatchingIndustries(querySl []string, scrubberDB db.ScrubberDB) []models.IndustryResp {
	var matchingIndustries []models.IndustryResp

	validation := make(map[string]string)

	for _, q := range querySl {
		matchingRecords := scrubberDB.IndustriesPFM.GetByPrefix(strings.ToUpper(q))

		for _, rData := range matchingRecords {
			industry := rData.(db.Industry)

			if _, valid := validation[industry.Code]; !valid {
				industryResp := models.IndustryResp{
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

func getRequestID(ctx context.Context) string {
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
