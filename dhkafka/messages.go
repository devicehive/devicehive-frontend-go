package dhkafka

import "github.com/devicehive/devicehive-frontend-go/request"

type KafkaRequestMessage struct {
	PartionKey          string       `json:"pK"`
	SingleReplyExpected bool         `json:"sre"`
	ReplyTo             string       `json:"rTo"`
	Body                interface{}  `json:"b"`
	CorrelationId       string       `json:"cId"`
	Type                request.Type `json:"t"`
}
