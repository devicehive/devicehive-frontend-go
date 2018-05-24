package dhkafka

type KafkaRequestMessage struct {
	PartionKey          string      `json:"pK"`
	SingleReplyExpected bool        `json:"sre"`
	ReplyTo             string      `json:"rTo"`
	Body                interface{} `json:"body"`
	CorrelationId       string      `json:"cId"`
}
