package timeseries

import (
	"go.uber.org/zap"

	rts "github.com/RedisTimeSeries/redistimeseries-go"
)

type TSClient struct {
	log *zap.Logger
	ts *rts.Client
}

func NewTSClient(log *zap.Logger, redisHost string) TSClient {
	return TSClient{
		log: log.Named("TSClient"),
		ts: rts.NewClient(redisHost, "", nil),
	}
}

func (c *TSClient) AddRecord(prefix string, stamp int64, data map[string]interface{}) {
	log := c.log.Named("AddRecord")

	values := make(map[string]float64)
	opts := rts.CreateOptions{
		DuplicatePolicy: rts.FirstDuplicatePolicy,
		Labels: make(map[string]string),
	}

	for key, val := range data {
		switch value := val.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			log.Debug("result", zap.String("key", key), zap.Float64("value", value.(float64)))
			values[key] = value.(float64)
		case bool:
			var r float64 = 0
			if value {
				r = 1
			}
			log.Debug("result", zap.String("key", key), zap.Float64("value", r))
			values[key] = r
		case string:
			log.Debug("result", zap.String("key", key), zap.String("value", value))
			opts.Labels[key] = value
		default:
			log.Debug("Skipped value with unhandled type", zap.String("key", key), zap.Any("value", value))
			continue
		}
	}

	for key, value := range values {
		_, err := c.ts.AddWithOptions(prefix + ":" + key, stamp, value, opts)
		if err != nil {
			log.Error("Failed to write key to timeseries",
				zap.String("key", prefix + ":" + key),
				zap.Int64("ts", stamp), zap.Float64("value", value),
				zap.Any("labels", opts.Labels),
			)
		}
	}
}