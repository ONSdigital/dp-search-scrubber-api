package service_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/ONSdigital/dp-search-scrubber-api/config"
	"github.com/ONSdigital/dp-search-scrubber-api/service"
	serviceMock "github.com/ONSdigital/dp-search-scrubber-api/service/mock"

	"github.com/pkg/errors"
	c "github.com/smartystreets/goconvey/convey"
)

var (
	ctx           = context.Background()
	testBuildTime = "BuildTime"
	testGitCommit = "GitCommit"
	testVersion   = "Version"
	errServer     = errors.New("HTTP Server error")
)

var (
	errHealthcheck = errors.New("healthCheck error")
)

var funcDoGetHealthcheckErr = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return nil, errHealthcheck
}

var funcDoGetHTTPServerNil = func(bindAddr string, router http.Handler) service.HTTPServer {
	return nil
}

func TestRun(t *testing.T) {
	c.Convey("Having a set of mocked dependencies", t, func() {
		cfg, err := config.Get()
		c.So(err, c.ShouldBeNil)

		hcMock := &serviceMock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
		}

		serverWg := &sync.WaitGroup{}
		serverMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return nil
			},
		}

		failingServerMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return errServer
			},
		}

		funcDoGetHealthcheckOk := func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		funcDoGetHTTPServer := func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		funcDoGetFailingHTTPSerer := func(bindAddr string, router http.Handler) service.HTTPServer {
			return failingServerMock
		}

		c.Convey("Given that initialising healthcheck returns an error", func() {
			// setup (run before each `c.Convey` at this scope / indentation):
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:  funcDoGetHTTPServerNil,
				DoGetHealthCheckFunc: funcDoGetHealthcheckErr,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			c.Convey("Then service Run fails with the same error and the flag is not set", func() {
				c.So(err, c.ShouldResemble, errHealthcheck)
				c.So(svcList.HealthCheck, c.ShouldBeFalse)
			})

			c.Reset(func() {
				// This c.reset is run after each `c.Convey` at the same scope (indentation)
			})
		})

		c.Convey("Given that all dependencies are successfully initialised", func() {
			// setup (run before each `c.Convey` at this scope / indentation):
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:  funcDoGetHTTPServer,
				DoGetHealthCheckFunc: funcDoGetHealthcheckOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			c.Convey("Then service Run succeeds and all the flags are set", func() {
				c.So(err, c.ShouldBeNil)
				c.So(svcList.HealthCheck, c.ShouldBeTrue)
			})

			c.Convey("The checkers are registered and the healthcheck and http server started", func() {
				c.So(len(hcMock.AddCheckCalls()), c.ShouldEqual, 0)
				c.So(len(initMock.DoGetHTTPServerCalls()), c.ShouldEqual, 1)
				c.So(initMock.DoGetHTTPServerCalls()[0].BindAddr, c.ShouldEqual, ":28700")
				c.So(len(hcMock.StartCalls()), c.ShouldEqual, 1)
				//!!! a call needed to stop the server, maybe ?
				serverWg.Wait() // Wait for HTTP server go-routine to finish
				c.So(len(serverMock.ListenAndServeCalls()), c.ShouldEqual, 1)
			})

			c.Reset(func() {
				// This c.reset is run after each `c.Convey` at the same scope (indentation)
			})
		})

		c.Convey("Given that all dependencies are successfully initialised but the http server fails", func() {
			// setup (run before each `c.Convey` at this scope / indentation):
			initMock := &serviceMock.InitialiserMock{
				DoGetHealthCheckFunc: funcDoGetHealthcheckOk,
				DoGetHTTPServerFunc:  funcDoGetFailingHTTPSerer,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			c.So(err, c.ShouldBeNil)

			c.Convey("Then the error is returned in the error channel", func() {
				sErr := <-svcErrors
				c.So(sErr.Error(), c.ShouldResemble, fmt.Sprintf("failure in http listen and serve: %s", errServer.Error()))
				c.So(len(failingServerMock.ListenAndServeCalls()), c.ShouldEqual, 1)
			})

			c.Reset(func() {
				// This c.reset is run after each `c.Convey` at the same scope (indentation)
			})
		})
	})
}

func TestClose(t *testing.T) {
	c.Convey("Having a correctly initialised service", t, func() {
		cfg, err := config.Get()
		c.So(err, c.ShouldBeNil)

		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &serviceMock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("Server stopped before healthcheck")
				}
				return nil
			},
		}

		c.Convey("Closing the service results in all the dependencies being closed in the expected order", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer { return serverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, serviceErr := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			c.So(serviceErr, c.ShouldBeNil)

			closingErr := svc.Close(context.Background())
			c.So(closingErr, c.ShouldBeNil)
			c.So(len(hcMock.StopCalls()), c.ShouldEqual, 1)
			c.So(len(serverMock.ShutdownCalls()), c.ShouldEqual, 1)
		})

		c.Convey("If services fail to stop, the Close operation tries to close all dependencies and returns an error", func() {
			failingserverMock := &serviceMock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					return errors.New("Failed to stop http server")
				},
			}

			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer { return failingserverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, serviceErr := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			c.So(serviceErr, c.ShouldBeNil)

			closingErr := svc.Close(context.Background())
			c.So(closingErr, c.ShouldNotBeNil)
			c.So(len(hcMock.StopCalls()), c.ShouldEqual, 1)
			c.So(len(failingserverMock.ShutdownCalls()), c.ShouldEqual, 1)
		})

		c.Convey("If service times out while shutting down, the Close operation fails with the expected error", func() {
			cfg.GracefulShutdownTimeout = 1 * time.Millisecond
			timeoutServerMock := &serviceMock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					time.Sleep(2 * time.Millisecond)
					return nil
				},
			}

			svcList := service.NewServiceList(nil)
			svcList.HealthCheck = true
			svc := service.Service{
				Config:      cfg,
				ServiceList: svcList,
				Server:      timeoutServerMock,
				HealthCheck: hcMock,
			}

			closingErr := svc.Close(context.Background())
			c.So(closingErr, c.ShouldNotBeNil)
			c.So(closingErr.Error(), c.ShouldResemble, "context deadline exceeded")
			c.So(len(hcMock.StopCalls()), c.ShouldEqual, 1)
			c.So(len(timeoutServerMock.ShutdownCalls()), c.ShouldEqual, 1)
		})
	})
}
