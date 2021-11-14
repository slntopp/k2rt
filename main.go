package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/slntopp/k2rt/pkg/reader"
	"github.com/slntopp/k2rt/pkg/timeseries"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger

	kafkaHost, kafkaTopic string
)

func init() {
	flag.Bool("debug", false, "Drop Log Level to DEBUG(-1)")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetDefault("KAFKA_HOST", "kafka:9092")
	viper.SetDefault("TOPIC", "shadow.reported-state.delta")

	viper.AutomaticEnv()
	
	level := 0
	if viper.GetBool("debug") {
		level = -1
	}

	atom := zap.NewAtomicLevel()
	atom.SetLevel(zapcore.Level(level))
	encoderCfg := zap.NewProductionEncoderConfig()
	log = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	kafkaHost  = viper.GetString("KAFKA_HOST")
	kafkaTopic = viper.GetString("TOPIC")
}

func main() {
	defer func() {
		_ = log.Sync()
	}()
	log.Info("Starting Service")

	r := reader.Make(kafkaHost, kafkaTopic)
	ch := make(chan kafka.Message)
	go reader.Start(r, log, ch)

	ts := timeseries.NewTSClient(log)

	for msg := range ch {
		var err error

		device := string(msg.Key)

		var data map[string]interface{}
		err = json.Unmarshal(msg.Value, &data)
		if err != nil {
			log.Error("Failed to Unmarshal message value", zap.ByteString("value", msg.Value), zap.Error(err))
			continue
		}

		ts.AddRecord(device, data)
	}
}