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
	TotalCommandsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "total_commands_count",
		Help: "Count of total commands sent",
	}, []string{"command"})
)

func init() {
	prometheus.MustRegister(TotalMsgCount)
	prometheus.MustRegister(TotalCommandsCount)
}

func StartMetricServer() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9176", nil)
}
