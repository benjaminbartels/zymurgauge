package metrics

type Metrics interface {
	Count(bucket string, n interface{})
	Increment(bucket string)
	Gauge(bucket string, value interface{})
	Timing(bucket string, value interface{})
	Histogram(bucket string, value interface{})
}
