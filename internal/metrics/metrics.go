package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	TotalMsgCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "total_msg_count",
		Help: "Count of total messages",
	})
)

func init() {
	prometheus.MustRegister(TotalMsgCount)
}

func StartMetricServer() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9176", nil)
}
