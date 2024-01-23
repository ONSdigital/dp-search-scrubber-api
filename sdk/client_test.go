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
	"github.com/ONSdigital/dp-search-scrubber-api/models"
	c "github.com/smartystreets/goconvey/convey"
)

const testHost = "http://localhost:23900"

var (
	initialTestState = healthcheck.CreateCheckState(service)

	scrubberResults = models.ScrubberResp{
		Time:  "10",
		Query: "sth",
		Results: models.Results{
			Areas: []models.AreaResp{
				{
					Name:       "name1",
					Region:     "region1",
					RegionCode: "regioncode1",
					Codes: map[string]string{
						"code1": "code1",
					},
				},
			},
			Industries: []models.IndustryResp{
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

	c.Convey("Given clienter.Do returns an error", t, func() {
		clientError := errors.New("unexpected error")
		httpClient := newMockHTTPClient(&http.Response{}, clientError)
		scrubberAPIClient := newScrubberAPIClient(httpClient)
		check := initialTestState

		c.Convey("When scrubber API client Checker is called", func() {
			err := scrubberAPIClient.Checker(ctx, &check)
			c.So(err, c.ShouldBeNil)

			c.Convey("Then the expected check is returned", func() {
				c.So(check.Name(), c.ShouldEqual, service)
				c.So(check.Status(), c.ShouldEqual, health.StatusCritical)
				c.So(check.StatusCode(), c.ShouldEqual, 0)
				c.So(check.Message(), c.ShouldEqual, clientError.Error())
				c.So(*check.LastChecked(), c.ShouldHappenAfter, timePriorHealthCheck)
				c.So(check.LastSuccess(), c.ShouldBeNil)
				c.So(*check.LastFailure(), c.ShouldHappenAfter, timePriorHealthCheck)
			})

			c.Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				c.So(doCalls, c.ShouldHaveLength, 1)
				c.So(doCalls[0].Req.URL.Path, c.ShouldEqual, path)
			})
		})
	})

	c.Convey("Given a 500 response for health check", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusInternalServerError}, nil)
		scrubberAPIClient := newScrubberAPIClient(httpClient)
		check := initialTestState

		c.Convey("When scrubber API client Checker is called", func() {
			err := scrubberAPIClient.Checker(ctx, &check)
			c.So(err, c.ShouldBeNil)

			c.Convey("Then the expected check is returned", func() {
				c.So(check.Name(), c.ShouldEqual, service)
				c.So(check.Status(), c.ShouldEqual, health.StatusCritical)
				c.So(check.StatusCode(), c.ShouldEqual, 500)
				c.So(check.Message(), c.ShouldEqual, service+healthcheck.StatusMessage[health.StatusCritical])
				c.So(*check.LastChecked(), c.ShouldHappenAfter, timePriorHealthCheck)
				c.So(check.LastSuccess(), c.ShouldBeNil)
				c.So(*check.LastFailure(), c.ShouldHappenAfter, timePriorHealthCheck)
			})

			c.Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				c.So(doCalls, c.ShouldHaveLength, 1)
				c.So(doCalls[0].Req.URL.Path, c.ShouldEqual, path)
			})
		})
	})
}

func TestGetScrubber(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c.Convey("Given request to find scrubber API results", t, func() {
		body, err := json.Marshal(scrubberResults)
		if err != nil {
			t.Errorf("failed to setup test data, error: %v", err)
		}

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(bytes.NewReader(body)),
			},
			nil)

		scrubberAPIClient := newScrubberAPIClient(httpClient)

		c.Convey("When GetScrubber is called", func() {
			query := url.Values{}
			query.Add("q", "sic code")
			resp, err := scrubberAPIClient.GetScrubber(ctx, Options{Query: query})

			c.Convey("Then the expected response body is returned", func() {
				c.So(*resp, c.ShouldResemble, scrubberResults)

				c.Convey("And no error is returned", func() {
					c.So(err, c.ShouldBeNil)

					c.Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						c.So(doCalls, c.ShouldHaveLength, 1)
						c.So(doCalls[0].Req.Method, c.ShouldEqual, "GET")
						c.So(doCalls[0].Req.URL.Path, c.ShouldEqual, "/v1/scrubber")
						c.So(doCalls[0].Req.URL.Query().Get("q"), c.ShouldEqual, "sic code")
						c.So(doCalls[0].Req.Header["Authorization"], c.ShouldBeEmpty)
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

func newScrubberAPIClient(httpClient *dphttp.ClienterMock) *Client {
	healthClient := healthcheck.NewClientWithClienter(service, testHost, httpClient)
	return NewWithHealthClient(healthClient)
}
