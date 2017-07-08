package rocketmq

type ConsumerGroup struct {
	BrokerOffset   int64 `json:"brokerOffset"`
	ConsumerOffset int64 `json:"consumerOffset"`
	Diff           int64 `json:"diff"`
	MessageQueue
}

type ConsumerProgress struct {
	GroupId string           `json:"groupId"`
	Tps     int64            `json:"tps"`
	Diff    int64            `json:"diff"`
	List    []*ConsumerGroup `json:"list"`
}
