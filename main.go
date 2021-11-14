package main

import (
	"flag"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

func init() {
	flag.Bool("debug", false, "Drop Log Level to DEBUG(-1)")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	
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
}

func main() {
	defer func() {
		_ = log.Sync()
	}()
}