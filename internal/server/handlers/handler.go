package handlers

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/service"
	"net/http"
	"regexp"
	"strconv"
)

func UpdateMetricHandler(service *service.MetricService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}
		//if r.Header.Get("Content-Type") != "text/plain" {
		//	w.WriteHeader(http.StatusBadRequest)
		//	w.Write([]byte("Invalid Content-Type"))
		//	return
		//}

		regex, regErr := regexp.Compile(`^/update/(\w+)/?(\w+)?/?([\-\d.,]+)?`)
		if regErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Bad regex"))
			return
		}
		matches := regex.FindSubmatch([]byte(r.URL.Path))

		metricType := string(matches[1])
		if !(metricType == metrics.StringCounterType || metricType == metrics.StringGaugeType) {
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("Wrong metric type"))
			return
		}
		metricName := string(matches[2])
		if metricName == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No metric ID specified"))
			return
		}
		metricValue := string(matches[3])
		if _, err := strconv.ParseFloat(metricValue, 64); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid metric value"))
			return
		}

		service.UpdateMetric(metricType, metricName, metricValue)

		w.WriteHeader(http.StatusOK)
	}
}
