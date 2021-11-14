package timeseries

import "go.uber.org/zap"

type TSClient struct {
	log *zap.Logger
}

func NewTSClient(log *zap.Logger) TSClient {
	return TSClient{log: log.Named("TSClient")}
}

func (c *TSClient) AddRecord(prefix string, data map[string]interface{}) {
	log := c.log.Named("AddRecord")
	for key, val := range data {
		switch value := val.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			log.Debug("result", zap.String("key", key), zap.Float64("value", value.(float64)))
		case bool:
			var r float64 = 0
			if value {
				r = 1
			}
			log.Debug("result", zap.String("key", key), zap.Float64("value", r))
		case string:
			log.Debug("result", zap.String("key", key), zap.String("value", value))
		default:
			log.Debug("Skipped value with unhandled type", zap.String("key", key), zap.Any("value", value))
			continue
		}
	}
}