package main

import (
	"fmt"

	"github.com/gunsluo/rocketmq-go"
)

func main() {
	conf := &rocketmq.Config{
		Nameserver:   "127.0.0.1:9876",
		InstanceName: "DEFAULT",
	}
	consumer, err := rocketmq.NewDefaultConsumer("ConsumerGroupName", conf)
	if err != nil {
		panic(err)
	}
	consumer.Subscribe("testTopic", "*")

	var count int
	consumer.RegisterMessageListener(func(msgs []*rocketmq.MessageExt) (int, error) {
		for _, msg := range msgs {
			count++
			fmt.Printf("count=%d|msgId=%s|topic=%s|storeTimestamp=%d|bornTimestamp=%d|storeHost=%s|bornHost=%s"+
				"|msgTag=%s|msgKey=%s|sysFlag=%d|storeSize=%d|queueId=%d|queueOffset=%d|body=%s\n",
				count, msg.MsgId, msg.Topic, msg.StoreTimestamp, msg.BornTimestamp, msg.StoreHost, msg.BornHost,
				msg.Tag(), msg.Key(), msg.SysFlag, msg.StoreSize, msg.QueueId, msg.QueueOffset, msg.Body)
		}
		return rocketmq.Action.CommitMessage, nil
	})
	consumer.Start()

	select {}
}
