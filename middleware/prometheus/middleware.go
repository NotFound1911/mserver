package prometheus

import (
	"github.com/NotFound1911/mserver"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
}

func (m MiddlewareBuilder) Build() mserver.Middleware {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      m.Name,
		Subsystem: m.Subsystem,
		Help:      m.Help,
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"pattern", "method", "status"})
	prometheus.MustRegister(vector)

	return func(next mserver.HandleFunc) mserver.HandleFunc {
		return func(ctx *mserver.Context) error {
			startTime := time.Now()
			defer func() {
				duration := time.Now().Sub(startTime).Milliseconds()
				patern := ctx.MatchedRoute
				if patern == "" {
					patern = "unknown"
				}
				vector.WithLabelValues(patern, ctx.GetRequest().
					Method, strconv.Itoa(ctx.GetRespStatusCode())).Observe(float64(duration))
			}()
			next(ctx)
			return nil
		}
	}
}
