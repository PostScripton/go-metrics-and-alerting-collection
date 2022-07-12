package server

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/service"
	"net/http"
	"regexp"
)

func UpdateMetric(service *service.MetricService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}
		if r.Header.Get("Content-Type") != "plain/text" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Content-Type"))
			return
		}

		regex, regErr := regexp.Compile(`^/update/(\w+)/(\w+)/([\-\d.,]+)$`)
		if regErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Bad regex"))
			return
		}
		matches := regex.FindSubmatch([]byte(r.URL.Path))

		metricType := string(matches[1])
		if !(metricType == metrics.StringCounterType || metricType == metrics.StringGaugeType) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Wrong metric type"))
			return
		}
		metricName := string(matches[2])
		metricValue := string(matches[3])

		service.UpdateMetric(metricType, metricName, metricValue)
	}
}
