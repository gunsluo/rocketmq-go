package rocketmq

type GetRouteInfoRequestHeader struct {
	Topic string `json:"topic"`
}

type QueryTopicConsumeByWhoRequestHeader struct {
	Topic string `json:"topic"`
}

type GetConsumerConnectionListRequestHeader struct {
	ConsumerGroup string `json:"consumerGroup"`
}

type GetConsumeStatsRequestHeader struct {
	ConsumerGroup string `json:"consumerGroup"`
	Topic         string `json:"topic"`
}

type GetKVConfigRequestHeader struct {
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
}

type GetConsumerListByGroupRequestHeader struct {
	ConsumerGroup string `json:"consumerGroup"`
}

type QueryConsumerOffsetRequestHeader struct {
	ConsumerGroup string `json:"consumerGroup"`
	Topic         string `json:"topic"`
	QueueId       int32  `json:"queueId"`
}

type ViewMessageRequestHeader struct {
	Offset uint64 `json:"offset"`
}
