package metrics

import (
	"github.com/kamontat/fthelper/metric/v4/src/collectors"
	"github.com/kamontat/fthelper/metric/v4/src/connection"
	"github.com/kamontat/fthelper/metric/v4/src/constants"
	"github.com/kamontat/fthelper/shared/caches"
	"github.com/kamontat/fthelper/shared/commandline/commands"
	"github.com/prometheus/client_golang/prometheus"
)

func callerBuilder(desc *prometheus.Desc, cacheKey string) []prometheus.Metric {
	var cache = caches.Global
	var data = cache.Get(cacheKey)
	var metric = 0
	if data.IsExist() {
		metric = data.Data.(int)
	}

	return []prometheus.Metric{prometheus.MustNewConstMetric(
		desc,
		prometheus.CounterValue,
		float64(metric),
	)}
}

var Internal = collectors.NewMetrics(
	collectors.NewMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName("fthelper", "internal", "ft_call"),
			"How many time do we call freqtrade apis",
			nil,
			prometheus.Labels{
				"type": "total",
			},
		),
		func(desc *prometheus.Desc, conn connection.Http, param *commands.ExecutorParameter) []prometheus.Metric {
			return callerBuilder(desc, constants.FTCONN_CALL)
		},
	),
	collectors.NewMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName("fthelper", "internal", "ft_call"),
			"How many time do we call freqtrade apis",
			nil,
			prometheus.Labels{
				"type": "success",
			},
		),
		func(desc *prometheus.Desc, conn connection.Http, param *commands.ExecutorParameter) []prometheus.Metric {
			return callerBuilder(desc, constants.FTCONN_CALL_SUCCESS)
		},
	),
	collectors.NewMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName("fthelper", "internal", "ft_call"),
			"How many time do we call freqtrade apis",
			nil,
			prometheus.Labels{
				"type": "failure",
			},
		),
		func(desc *prometheus.Desc, conn connection.Http, param *commands.ExecutorParameter) []prometheus.Metric {
			return callerBuilder(desc, constants.FTCONN_CALL_FAILURE)
		},
	),
	collectors.NewMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName("fthelper", "internal", "cache_size"),
			"Total keys we store on cache service, including expired ones",
			nil,
			nil,
		),
		func(desc *prometheus.Desc, conn connection.Http, param *commands.ExecutorParameter) []prometheus.Metric {
			return []prometheus.Metric{prometheus.MustNewConstMetric(
				desc,
				prometheus.CounterValue,
				float64(param.Cache.Size()),
			)}
		},
	),
	collectors.NewMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName("fthelper", "internal", "warmup_error"),
			"Counter of how many warmup is fail",
			nil,
			nil,
		),
		func(desc *prometheus.Desc, conn connection.Http, param *commands.ExecutorParameter) []prometheus.Metric {
			return callerBuilder(desc, constants.WARMUP_ERROR)
		},
	),
	collectors.NewRawMetric(
		func(desc *prometheus.Desc, conn connection.Http, param *commands.ExecutorParameter) []prometheus.Metric {
			var cache = caches.Global
			var data = cache.Get(constants.WARMUP_DURATIONS)
			if !data.IsExist() {
				return []prometheus.Metric{}
			}

			var durations = data.Data.([]int64)
			var histogram = prometheus.NewHistogram(prometheus.HistogramOpts{
				Namespace: "fthelper",
				Subsystem: "internal",
				Name:      "warmup_duration",
				Help:      "Warmup duration in milliseconds. (kept only last 1000)",
				Buckets:   prometheus.ExponentialBuckets(20, 2, 8),
			})

			for _, ms := range durations {
				histogram.Observe(float64(ms))
			}

			return []prometheus.Metric{histogram}
		},
	),
)
