package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/inconshreveable/log15"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	healthClient healthpb.HealthClient
	healthConn   *grpc.ClientConn
	remoteUrl    string        = "0.0.0.0:50050"
	serviceName  string        = ""
	port         string        = ":8080"
	timeoutDur   time.Duration = time.Second
	logger       log15.Logger
)

func connectToRemote() {
	if url := os.Getenv("REMOTE_URL"); url != "" {
		remoteUrl = url
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	c, err := grpc.Dial(remoteUrl, opts...)
	if err != nil {
		glog.Errorf("failed to dial grpc server, %v", err)
		return
	}
	healthConn = c
	healthClient = healthpb.NewHealthClient(healthConn)
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if healthClient == nil {
		connectToRemote()
	}
	if healthClient == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "NOT_HEALTHY")
		logger.Error("server is not healthy, unable to connect remote grpc server")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), timeoutDur)
	resp, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{Service: serviceName})
	if err == nil && resp.Status == healthpb.HealthCheckResponse_SERVING {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "HEALTHY")
		logger.Info("health check is OK")
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "NOT_HEALTHY")
	logger.Error("server is not healthy err=%v response=%v", err, resp)
}

func main() {
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}

	if pod := os.Getenv("POD_NAME"); pod != "" {
		serviceName = pod
	}

	logger = log15.New("context", "HealthCheck", "podName", serviceName)
	connectToRemote()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHealthCheck)

	httpsrv := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	glog.Infof("Binding to %s...", httpsrv.Addr)
	glog.Fatal(httpsrv.ListenAndServe())
}
