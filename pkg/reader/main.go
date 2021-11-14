package reader

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Make(host, topic string) (*kafka.Reader) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{host},
		Topic:     topic,
		Partition: 0,
	})
	r.SetOffset(0)

	return r
}

func Start(r *kafka.Reader, log *zap.Logger, ch chan kafka.Message) {
	log = log.Named("KafkaReader")
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Debug("Error reading Message", zap.Error(err))
			break
		}
		log.Debug("Message read", zap.Any("msg", m))
		ch <- m
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", zap.Error(err))
	}
}