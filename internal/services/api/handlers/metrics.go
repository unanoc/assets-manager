package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func SetupMetrics(router *gin.Engine) {
	m := ginmetrics.GetMonitor()

	m.SetMetricPath("/metrics")
	m.SetSlowTime(3)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	m.Use(router)

	prometheus.DefaultRegisterer.Unregister(collectors.NewGoCollector())
	prometheus.DefaultRegisterer.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}
