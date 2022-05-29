package metrics

import (
	"context"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Speakerkfm/iso/internal/pkg/router"
)

var (
	UnaryServerInterceptor = grpc_prometheus.UnaryServerInterceptor

	RequestProcessingTimeHistogramVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_processing_time_seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method"},
	)
)

func init() {
	prometheus.MustRegister(
		RequestProcessingTimeHistogramVec,
	)
	grpc_prometheus.EnableHandlingTimeHistogram()
}

func RegisterMetricsHandler(ctx context.Context, mux router.ServeMux) error {
	mux.Handle("/metrics", promhttp.Handler())
	return nil
}
