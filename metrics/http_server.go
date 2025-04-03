package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/kodenkai-labs/go-lib/httplib"
)

func NewMetricsServer(port, path string, shutdownTimeout time.Duration) httplib.Server {
	router := gin.Default()
	router.GET(path, gin.WrapH(promhttp.Handler()))

	prometheus.DefaultRegisterer.Unregister(collectors.NewGoCollector())
	prometheus.DefaultRegisterer.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return httplib.NewHTTPServer(router, port, shutdownTimeout)
}
