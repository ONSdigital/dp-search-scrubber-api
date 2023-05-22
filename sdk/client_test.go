package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-nlp-search-scrubber/models"
	. "github.com/smartystreets/goconvey/convey"
)

const testHost = "http://localhost:23900"

var (
	initialTestState = healthcheck.CreateCheckState(service)

	searchResults = models.ScrubberResp{
		Time:  "10",
		Query: "sth",
		Results: models.Results{
			Areas: []*models.AreaResp{
				{
					Name:       "name1",
					Region:     "region1",
					RegionCode: "regioncode1",
					Codes: map[string]string{
						"code1": "code1",
					},
				},
			},
			Industries: []*models.IndustryResp{
				{
					Code: "indcode1",
					Name: "indname1",
				},
			},
		},
	}
)

func TestHealthCheckerClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	timePriorHealthCheck := time.Now().UTC()
	path := "/health"

	Convey("Given clienter.Do returns an error", t, func() {
		clientError := errors.New("unexpected error")
		httpClient := newMockHTTPClient(&http.Response{}, clientError)
		searchAPIClient := newSearchAPIClient(t, httpClient)
		check := initialTestState

		Convey("When search API client Checker is called", func() {
			err := searchAPIClient.Checker(ctx, &check)
			So(err, ShouldBeNil)

			Convey("Then the expected check is returned", func() {
				So(check.Name(), ShouldEqual, service)
				So(check.Status(), ShouldEqual, health.StatusCritical)
				So(check.StatusCode(), ShouldEqual, 0)
				So(check.Message(), ShouldEqual, clientError.Error())
				So(*check.LastChecked(), ShouldHappenAfter, timePriorHealthCheck)
				So(check.LastSuccess(), ShouldBeNil)
				So(*check.LastFailure(), ShouldHappenAfter, timePriorHealthCheck)
			})

			Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				So(doCalls, ShouldHaveLength, 1)
				So(doCalls[0].Req.URL.Path, ShouldEqual, path)
			})
		})
	})

	Convey("Given a 500 response for health check", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusInternalServerError}, nil)
		searchAPIClient := newSearchAPIClient(t, httpClient)
		check := initialTestState

		Convey("When search API client Checker is called", func() {
			err := searchAPIClient.Checker(ctx, &check)
			So(err, ShouldBeNil)

			Convey("Then the expected check is returned", func() {
				So(check.Name(), ShouldEqual, service)
				So(check.Status(), ShouldEqual, health.StatusCritical)
				So(check.StatusCode(), ShouldEqual, 500)
				So(check.Message(), ShouldEqual, service+healthcheck.StatusMessage[health.StatusCritical])
				So(*check.LastChecked(), ShouldHappenAfter, timePriorHealthCheck)
				So(check.LastSuccess(), ShouldBeNil)
				So(*check.LastFailure(), ShouldHappenAfter, timePriorHealthCheck)
			})

			Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				So(doCalls, ShouldHaveLength, 1)
				So(doCalls[0].Req.URL.Path, ShouldEqual, path)
			})
		})
	})
}

func TestGetSearch(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	Convey("Given request to find search results", t, func() {
		body, err := json.Marshal(searchResults)
		if err != nil {
			t.Errorf("failed to setup test data, error: %v", err)
		}

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(bytes.NewReader(body)),
			},
			nil)

		searchAPIClient := newSearchAPIClient(t, httpClient)

		Convey("When GetSearch is called", func() {
			query := url.Values{}
			query.Add("q", "census")
			resp, err := searchAPIClient.GetSearch(ctx, Options{Query: query})

			Convey("Then the expected response body is returned", func() {
				So(*resp, ShouldResemble, searchResults)

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.Method, ShouldEqual, "GET")
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/v1/search")
						So(doCalls[0].Req.URL.Query().Get("q"), ShouldEqual, "census")
						So(doCalls[0].Req.Header["Authorization"], ShouldBeEmpty)
					})
				})
			})
		})
	})
}

func newMockHTTPClient(r *http.Response, err error) *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		SetPathsWithNoRetriesFunc: func(paths []string) {
		},
		DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return r, err
		},
		GetPathsWithNoRetriesFunc: func() []string {
			return []string{"/healthcheck"}
		},
	}
}

func newSearchAPIClient(t *testing.T, httpClient *dphttp.ClienterMock) *Client {
	healthClient := healthcheck.NewClientWithClienter(service, testHost, httpClient)
	return NewWithHealthClient(healthClient)
}
