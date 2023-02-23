package msg_consumer

import (
	"github.com/donscoco/gochat/internal/base/kafka/mkafka"
	"github.com/donscoco/gochat/internal/module/msg_process"
	"github.com/donscoco/gochat/pkg/iron_config"
)

var MsgConsumer *mkafka.KafkaConsumer

func MsgConsumerStart() (err error) { // todo 文件配置

	conf := iron_config.Conf

	addrs := make([]string, 0, 16)
	conf.GetByScan("/kafka/addrs", &addrs)
	topics := make([]string, 0, 16)
	conf.GetByScan("/kafka/consumer/topics", &topics)

	groupId := conf.GetString("/kafka/consumer/group_id")
	clientId := conf.GetString("/kafka/consumer/client_id")

	c, err := mkafka.CreateConsumer(
		//[]string{"192.168.2.132:9092", "192.168.2.132:9093", "192.168.2.132:9094"},
		addrs,
		groupId,
		topics,
		map[string]string{
			"ClientId":           clientId,
			"AutoCommitInterval": "100",
			"ReturnSuc":          conf.GetString("/kafka/producer/return_suc"),
			"Version":            conf.GetString("/kafka/version"),
			"RequiredAcks":       conf.GetString("/kafka/producer/required_acks"),
			"SetLog":             "true",
		},
		msg_process.MsgHandleFunc,
		msg_process.SuccessFunc,
		msg_process.FailureFunc,
	)

	c.Start()

	MsgConsumer = c

	return
}
func MsgConsumerStop() (err error) {
	MsgConsumer.Stop()
	return
}
