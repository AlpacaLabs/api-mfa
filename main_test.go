package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/AlpacaLabs/go-kontext"

	"github.com/AlpacaLabs/api-mfa/internal/app"
	"github.com/AlpacaLabs/api-mfa/internal/configuration"
	mfaV1 "github.com/AlpacaLabs/protorepo-mfa-go/alpacalabs/mfa/v1"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

var conn *grpc.ClientConn

func TestMain(m *testing.M) {
	c := configuration.LoadConfig()
	logrus.Infof("Loaded config: %s", c)

	a := app.NewApp(c)

	go a.Run()

	conn = createGRPCConn(c)
	waitUntilHealthy(conn)

	code := m.Run()

	os.Exit(code)
}

func Test_MFA_Flow(t *testing.T) {
	Convey("a user wants to get their MFA options", t, func(c C) {
		ctx := context.TODO()
		client := mfaV1.NewMFAServiceClient(conn)
		_, err := client.GetDeliveryOptions(ctx, &mfaV1.GetDeliveryOptionsRequest{AccountId: xid.New().String()})
		So(err, ShouldBeNil)
	})
}

func createGRPCConn(c configuration.Config) *grpc.ClientConn {
	grpcAddress := fmt.Sprintf("localhost:%d", c.GrpcPort)
	grpcConn, err := kontext.Dial(grpcAddress)
	if err != nil {
		logrus.Fatalf("failed to connect to our own gRPC server: %v", err)
	}
	return grpcConn
}

func waitUntilHealthy(grpcConn *grpc.ClientConn) {
	healthClient := health.NewHealthClient(grpcConn)

	ticker := time.Tick(time.Second * 1)
	timeout := time.After(time.Second * 5)

	var timedOut bool
	for {
		if timedOut {
			break
		}
		select {
		case <-ticker:
			// check health
			if res, err := healthClient.Check(context.TODO(), &health.HealthCheckRequest{}); err != nil {
				logrus.Warnf("got error while checking server health because it may still be starting up: %v", err)
			} else {
				if res.Status == health.HealthCheckResponse_SERVING {
					break
				}
			}
		case <-timeout:
			timedOut = true
		}
	}

	if timedOut {
		logrus.Fatal("timed out waiting for server to come alive")
	}
}
