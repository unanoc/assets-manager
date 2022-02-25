package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/trustwallet/assets-manager/internal/config"
)

const (
	namespace = "assets_manager"
	subsystem = "worker"
)

// Prometheus is a struct for prometheus metrics.
type Prometheus struct {
	PullRequestsOpen           prometheus.Gauge
	PullRequestsToPay          prometheus.Gauge
	CounterPullRequestsCreated prometheus.Counter
	CounterPaymentsDetected    prometheus.Counter
}

// NewPrometheus return an instance of Prometheus with registered metrics.
func NewPrometheus() *Prometheus {
	constLabels := prometheus.Labels{"service": config.Default.ServiceName}

	p := Prometheus{
		PullRequestsOpen: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name:        prometheus.BuildFQName(namespace, subsystem, "state_pull_requests_open"),
				Help:        "Current number of open pull requests",
				ConstLabels: constLabels,
			},
		),
		PullRequestsToPay: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name:        prometheus.BuildFQName(namespace, subsystem, "state_pull_requests_to_pay"),
				Help:        "Current number of pull requests expecting a payment",
				ConstLabels: constLabels,
			},
		),
		CounterPullRequestsCreated: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:        prometheus.BuildFQName(namespace, subsystem, "event_pull_requests_created"),
				Help:        "Number of PR created notications recevied",
				ConstLabels: constLabels,
			},
		),
		CounterPaymentsDetected: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name:        prometheus.BuildFQName(namespace, subsystem, "action_payments_detected"),
				Help:        "Number of incoming payments detected so far",
				ConstLabels: constLabels,
			},
		),
	}

	// Register metrics.
	prometheus.MustRegister(
		p.PullRequestsOpen,
		p.PullRequestsToPay,
		p.CounterPullRequestsCreated,
		p.CounterPaymentsDetected,
	)

	prometheus.DefaultRegisterer.Unregister(collectors.NewGoCollector())
	prometheus.DefaultRegisterer.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return &p
}

func (p *Prometheus) SetPullRequestsOpen(n int) {
	p.PullRequestsOpen.Set(float64(n))
}

func (p *Prometheus) SetPullRequestsToPay(n int) {
	p.PullRequestsToPay.Set(float64(n))
}

func (p *Prometheus) IncCounterPullRequestsCreated() {
	p.CounterPullRequestsCreated.Inc()
}

func (p *Prometheus) IncCounterPaymentsDetected() {
	p.CounterPaymentsDetected.Inc()
}
