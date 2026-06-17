package observability

import (
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	once sync.Once
	instance *Metrics
)

type Metrics struct{
	httpRequestTotal *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpInFlight prometheus.Gauge
	wsConnectionsActive prometheus.Gauge
	wsConnectionsTotal *prometheus.CounterVec
	wsMessageTotal *prometheus.CounterVec
}

func New() *Metrics{
	once.Do(func() {
		m := &Metrics{
			httpRequestTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: "Ganymede",
					Subsystem: "http",
					Name: "request_total",
					Help: "Total HTTP requests",
				},
				[]string{"method", "route", "status"},
			),
			httpRequestDuration: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Namespace: "Ganymede",
					Subsystem: "http",
					Name: "request_duration_seconds",
					Help: "HTTP request latency in seconds",
					Buckets: prometheus.DefBuckets,
				},
				[]string{"method", "route"},
			),
			httpInFlight: prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: "Ganymede",
					Subsystem: "http",
					Name: "request_in_flight",
					Help: "Current in-flight HTTP requests",
				},

			),
			wsConnectionsActive: prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: "Ganymede",
					Subsystem: "ws",
					Name: "connections_active",
					Help: "Current active websocket connections",
				},

			),
			wsConnectionsTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: "Ganymede",
					Subsystem: "ws",
					Name: "connections_total",
					Help: "Total WebSocket connections by result",

				},
				[]string{"result"},
			),
			wsMessageTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: "Ganymede",
					Subsystem: "ws",
					Name: "messages_total",
					Help: "Total WebSocket messages",
				},
				[]string{"direction", "type"},
			),

		}
		prometheus.MustRegister(
			m.httpRequestTotal,
			m.httpRequestDuration,
			m.httpInFlight,
			m.wsConnectionsActive,
			m.wsConnectionsTotal,
			m.wsMessageTotal,
		)
		instance = m
	})
	return instance
}

func CounterInc(vec *prometheus.CounterVec, labels ...string){
	if metric, err := vec.GetMetricWithLabelValues(labels...); err == nil {
		metric.Inc()
	}
}

func Observe(hist *prometheus.HistogramVec, value float64, labels ...string){
	if metric, err := hist.GetMetricWithLabelValues(labels...); err == nil {
		metric.Observe(value)
	}
}

func NormalizeRoute(path string) string{
	if path == ""{
		return "unknown"
	}
	return path
}

func (m *Metrics) ObserveHTTP(method, route string, status int, duration time.Duration){
	if m == nil {
		return
	}

	route = NormalizeRoute(route)
	statusStr := strconv.Itoa(status)

	CounterInc(m.httpRequestTotal, method, route, statusStr)
	Observe(m.httpRequestDuration, duration.Seconds(), method, route)
}

func (m *Metrics) IncInFlight(){
	if m != nil {
		m.httpInFlight.Inc()
	}
}

func (m *Metrics) DecInFlight(){
	if m != nil {
		m.httpInFlight.Dec()
	}
}


func (m *Metrics) ObserveWSConnection(result string){
	if m == nil {
		return
	}
	if metric, err := m.wsConnectionsTotal.GetMetricWithLabelValues(result); err == nil {
		metric.Inc()
	}
}

func (m *Metrics) IncWSConnections(){
	if m != nil {
		m.wsConnectionsActive.Inc()
	}
}

func (m *Metrics) DecWSConnections(){
	if m != nil {
		m.wsConnectionsActive.Dec()
	}
}

func sanitizeMessageType(t string) string{
	switch t{
	case "chat", "join", "leave", "offer", "answer", "candidate":
		return t
	default:
		return "unknown"
	}
}

func (m *Metrics) ObserveWSMessage(direction, msgType string){
	if m == nil {
		return
	}
	msgType = sanitizeMessageType(msgType)

	if metric, err := m.wsMessageTotal.GetMetricWithLabelValues(direction, msgType); err == nil {
		metric.Inc()
	}
}
