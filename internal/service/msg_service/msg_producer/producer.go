package msg_producer

import (
	"github.com/donscoco/gochat/internal/base/kafka/mkafka"
	"github.com/donscoco/gochat/pkg/iron_config"
)

var MsgProducer *mkafka.KafkaProducer

func MsgProducerStart() (err error) {

	conf := iron_config.Conf

	addrs := make([]string, 0, 16)
	conf.GetByScan("/kafka/addrs", &addrs)

	p, err := mkafka.CreateProducer(
		//[]string{"192.168.2.132:9092", "192.168.2.132:9093", "192.168.2.132:9094"},
		addrs,
		map[string]string{
			"ReturnSuc":    conf.GetString("/kafka/producer/return_suc"),
			"Version":      conf.GetString("/kafka/version"),
			"RequiredAcks": conf.GetString("/kafka/producer/required_acks"),
			"SetLog":       "true",
		})
	if err != nil {
		return err
	}

	p.Start()
	MsgProducer = p

	return nil
}

func MsgProducerStop() error {
	MsgProducer.Stop()
	return nil
}
