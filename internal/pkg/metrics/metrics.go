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

	RequestProcessingTimeSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "request_processing_time_summary_ms",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"service", "method"},
	)
)

func init() {
	prometheus.MustRegister(
		RequestProcessingTimeSummary,
	)
}

func RegisterMetricsHandler(ctx context.Context, mux router.ServeMux) error {
	mux.Handle("/metrics", promhttp.Handler())
	return nil
}
