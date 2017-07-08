package rocketmq

type kvTable struct {
	Table map[string]string `json:"table"`
}

type BrokerRuntimeInfo struct {
	BrokerVersionDesc           string  `json:"brokerVersionDesc"`
	BrokerVersion               string  `json:"brokerVersion"`
	MsgPutTotalYesterdayMorning string  `json:"msgPutTotalYesterdayMorning"`
	MsgPutTotalTodayMorning     string  `json:"msgPutTotalTodayMorning"`
	MsgPutTotalTodayNow         string  `json:"msgPutTotalTodayNow"`
	MsgGetTotalYesterdayMorning string  `json:"msgGetTotalYesterdayMorning"`
	MsgGetTotalTodayNow         string  `json:"msgGetTotalTodayNow"`
	SendThreadPoolQueueSize     string  `json:"sendThreadPoolQueueSize"`
	SendThreadPoolQueueCapacity string  `json:"sendThreadPoolQueueCapacity"`
	MsgGetTotalTodayMorning     string  `json:"msgGetTotalTodayMorning"`
	InTps                       float64 `json:"inTps"`
	OutTps                      float64 `json:"outTps"`
}
