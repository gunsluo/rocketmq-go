package rocketmq

var TrackTypes = struct {
	SubscribedAndConsumed      int //订阅了，而且消费了（Offset越过了）
	SubscribedButFilterd       int //订阅了，但是被过滤掉了
	SubscribedButPull          int //订阅了，但是PULL，结果未知
	SubscribedAndNotConsumeYet int //订阅了，但是没有消费（Offset小）
	UnknowExeption             int //未知异常
}{
	0,
	1,
	2,
	3,
	4,
}

type MessageTrack struct {
	ConsumerGroup string `json:"consumerGroup"`
	TrackType     int    `json:"trackType"`
	ExceptionDesc string `json:"exceptionDesc"`
}
