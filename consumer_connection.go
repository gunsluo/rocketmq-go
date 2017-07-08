package rocketmq

var ConsumeTypes = struct {
	Actively  string //主动方式消费
	Passively string //被动方式消费
}{
	"CONSUME_ACTIVELY",
	"CONSUME_PASSIVELY",
}

type ConsumerConnection struct {
	ConnectionSet     []*Connection                `json:"connectionSet"`
	SubscriptionTable map[string]*SubscriptionData `json:"subscriptionTable"`
	ConsumeType       string                       `json:"consumeType"`
	MessageModel      string                       `json:"messageModel"`
	ConsumeFromWhere  string                       `json:"consumeFromWhere"`
	//subscriptionTableLock sync.RWMutex                 `json:"-"`
}

type Connection struct {
	ClientId   string `json:"clientId"`
	ClientAddr string `json:"clientAddr"`
	Language   string `json:"language"`
	Version    int    `json:"version"`
}

type ConsumeStats struct {
	OffsetTable map[MessageQueue]OffsetWrapper `json:"offsetTable"`
	ConsumeTps  int64                          `json:"consumeTps"`
}

type OffsetWrapper struct {
	BrokerOffset   int64 `json:"brokerOffset"`
	ConsumerOffset int64 `json:"consumerOffset"`
	LastTimestamp  int64 `json:"lastTimestamp"` // 消费的最后一条消息对应的时间戳
}
