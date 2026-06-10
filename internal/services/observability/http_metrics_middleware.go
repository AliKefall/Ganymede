package observability

import (
	"bufio"
	"net"
	"net/http"
	"time"
)

type statusRecorder struct{
	http.ResponseWriter
	status int
}

func (r *statusRecorder) Hijack()(net.Conn, *bufio.ReadWriter, error){
	hj, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hj.Hijack()
}

func (r *statusRecorder) Flush(){
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (r *statusRecorder) WriteHeader(code int){
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func MiddlewareHTTPMetrics(m *Metrics) func(http.Handler) http.Handler{
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if m == nil {
				h.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			m.IncInFlight()
			defer m.DecInFlight()

			rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			h.ServeHTTP(rw, r)
			m.ObserveHTTP(r.Method, r.URL.Path, rw.status, time.Since(start))
		})
	}
}
