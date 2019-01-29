package monitor

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MoonighT/elastic"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var gatewayURL = "http://localhost:9091/"
var jobName = "golang_web_server_push"

var webRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "web_request_duration_seconds",
		Help:    "web request duration distribution",
		Buckets: []float64{1, 2, 5, 10, 20, 60},
	},
	[]string{"method", "endpoint"},
)

var webRequestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "web_reqeust_total",
		Help: "Number of web requests in total",
	},
	[]string{"method", "endpoint"},
)

var esRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "es_request_duration_seconds",
		Help:    "elasticsearch query duration distribution",
		Buckets: []float64{1, 2, 5, 10, 20, 60},
	},
	[]string{"query"},
)

var esRequestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "es_reqeust_total",
		Help: "Number of elasticsearch requests in total",
	},
	[]string{"query"},
)

func init() {
	// register exposed metrics
	prometheus.MustRegister(webRequestDuration)
	prometheus.MustRegister(webRequestCount)

	prometheus.MustRegister(esRequestDuration)
	prometheus.MustRegister(esRequestCount)
}

// Monitor ...
func Monitor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h(w, r)
		duration := time.Since(start)
		webRequestCount.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
		webRequestDuration.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Observe(duration.Seconds())
		// push webRequestCount and webRequestDuration to pushgateway.
		if err := push.Collectors(jobName, push.HostnameGroupingKey(), gatewayURL, webRequestCount, webRequestDuration); err != nil {
			log.Printf("can not push webRequestCount and webRequestDuration to pushgateway: %s", err)
		}
	}
}

// ESQuery ...
func ESQuery(es *elastic.Client, index string, typ string, query elastic.Query) ([]byte, error) {
	start := time.Now()

	searchResult, err := es.Search().
		Index(index).
		Type(typ).
		Query(query).
		Size(1).
		Do()

	duration := time.Since(start)

	if err != nil {
		return nil, err
	}

	source, err := query.Source()
	if err == nil {
		bytes, ee := json.Marshal(source)
		if ee == nil {
			esRequestCount.With(prometheus.Labels{"query": string(bytes)}).Inc()
			esRequestDuration.With(prometheus.Labels{"query": string(bytes)}).Observe(duration.Seconds())
			// push esRequestCount and esRequestDuration to pushgateway.
			if err := push.Collectors(jobName, push.HostnameGroupingKey(), gatewayURL, esRequestCount, esRequestDuration); err != nil {
				log.Printf("can not push esRequestCount and esRequestDuration to pushgateway: %s", err)
			}
		}
	}

	var line []byte
	line, err = searchResult.Hits.Hits[0].Source.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return line, nil
}
