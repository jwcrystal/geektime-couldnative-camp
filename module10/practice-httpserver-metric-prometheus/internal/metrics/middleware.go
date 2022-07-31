package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"metrics/internal/handlers"
	"net/http"
)

func Use(r *handlers.Engine) {
	// should integrate with engine (maybe be a context for use)
	r.GET(metricPath, func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})
}
