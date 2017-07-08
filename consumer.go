package rocketmq

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	BrokerSuspendMaxTimeMillis       = 1000 * 15
	FLAG_COMMIT_OFFSET         int32 = 0x1 << 0
	FLAG_SUSPEND               int32 = 0x1 << 1
	FLAG_SUBSCRIPTION          int32 = 0x1 << 2
	FLAG_CLASS_FILTER          int32 = 0x1 << 3
)

type MessageListener func(msgs []*MessageExt) (int, error)

var DEFAULT_IP = GetLocalIp4()

type Config struct {
	Nameserver   string
	ClientIp     string
	InstanceName string
}

type Consumer interface {
	//Admin
	Start() error
	Shutdown()
	RegisterMessageListener(listener MessageListener)
	Subscribe(topic string, subExpression string)
	UnSubcribe(topic string)
	SendMessageBack(msg MessageExt, delayLevel int) error
	SendMessageBack1(msg MessageExt, delayLevel int, brokerName string) error
	fetchSubscribeMessageQueues(topic string) error
}

type DefaultConsumer struct {
	conf             *Config
	consumerGroup    string
	consumeFromWhere string
	consumerType     string
	messageModel     string
	unitMode         bool

	subscription    map[string]string
	messageListener MessageListener
	offsetStore     OffsetStore
	brokers         map[string]net.Conn

	rebalance      *Rebalance
	remotingClient RemotingClient
	mqClient       *MqClient
}

func NewDefaultConsumer(name string, conf *Config) (Consumer, error) {
	if conf == nil {
		conf = &Config{
			Nameserver:   os.Getenv("ROCKETMQ_NAMESVR"),
			InstanceName: "DEFAULT",
		}
	}

	if conf.ClientIp == "" {
		conf.ClientIp = DEFAULT_IP
	}

	remotingClient := NewDefaultRemotingClient()
	mqClient := NewMqClient()

	rebalance := NewRebalance()
	rebalance.groupName = name
	rebalance.mqClient = mqClient

	offsetStore := new(RemoteOffsetStore)
	offsetStore.mqClient = mqClient
	offsetStore.groupName = name
	offsetStore.offsetTable = make(map[MessageQueue]int64)

	pullMessageService := NewPullMessageService()

	consumer := &DefaultConsumer{
		conf:             conf,
		consumerGroup:    name,
		consumeFromWhere: "CONSUME_FROM_LAST_OFFSET",
		subscription:     make(map[string]string),
		offsetStore:      offsetStore,
		brokers:          make(map[string]net.Conn),
		rebalance:        rebalance,
		remotingClient:   remotingClient,
		mqClient:         mqClient,
	}

	mqClient.consumerTable[name] = consumer
	mqClient.remotingClient = remotingClient
	mqClient.conf = conf
	mqClient.clientId = fmt.Sprintf("%s@%d#%d#%d", conf.ClientIp, os.Getpid(), HashCode(conf.Nameserver), UnixNano())
	mqClient.pullMessageService = pullMessageService

	rebalance.consumer = consumer
	pullMessageService.consumer = consumer

	return consumer, nil
}

func (self *DefaultConsumer) Start() error {
	self.mqClient.start()
	return nil
}

func (self *DefaultConsumer) Shutdown() {
}

func (self *DefaultConsumer) RegisterMessageListener(messageListener MessageListener) {
	self.messageListener = messageListener
}

func (self *DefaultConsumer) Subscribe(topic string, subExpression string) {
	self.subscription[topic] = subExpression

	subData := &SubscriptionData{
		Topic:     topic,
		SubString: subExpression,
	}
	self.rebalance.subscribe(topic, subData)
}

func (self *DefaultConsumer) UnSubcribe(topic string) {
	delete(self.subscription, topic)
	self.rebalance.unSubscribe(topic)
}

func (self *DefaultConsumer) SendMessageBack(msg MessageExt, delayLevel int) error {
	return nil
}

func (self *DefaultConsumer) SendMessageBack1(msg MessageExt, delayLevel int, brokerName string) error {
	return nil
}

func (self *DefaultConsumer) fetchSubscribeMessageQueues(topic string) error {
	return nil
}

func (self *DefaultConsumer) pullMessage(pullRequest *PullRequest) {

	commitOffsetEnable := false
	commitOffsetValue := int64(0)

	commitOffsetValue = self.offsetStore.readOffset(pullRequest.messageQueue, READ_FROM_MEMORY)
	if commitOffsetValue > 0 {
		commitOffsetEnable = true
	}

	var sysFlag int32 = 0
	if commitOffsetEnable {
		sysFlag |= FLAG_COMMIT_OFFSET
	}

	sysFlag |= FLAG_SUSPEND

	subscriptionData, ok := self.rebalance.subscriptionInner[pullRequest.messageQueue.Topic]
	var subVersion int64
	var subString string
	if ok {
		subVersion = subscriptionData.SubVersion
		subString = subscriptionData.SubString

		sysFlag |= FLAG_SUBSCRIPTION
	}

	requestHeader := new(PullMessageRequestHeader)
	requestHeader.ConsumerGroup = pullRequest.consumerGroup
	requestHeader.Topic = pullRequest.messageQueue.Topic
	requestHeader.QueueId = pullRequest.messageQueue.QueueId
	requestHeader.QueueOffset = pullRequest.nextOffset

	requestHeader.SysFlag = sysFlag
	requestHeader.CommitOffset = commitOffsetValue
	requestHeader.SuspendTimeoutMillis = BrokerSuspendMaxTimeMillis

	if ok {
		requestHeader.SubVersion = subVersion
		requestHeader.Subscription = subString
	}

	pullCallback := func(responseFuture *ResponseFuture) {
		var nextBeginOffset int64 = pullRequest.nextOffset

		if responseFuture != nil {
			responseCommand := responseFuture.responseCommand
			if responseCommand.Code == SUCCESS && len(responseCommand.Body) > 0 {
				var err error
				pullResult, ok := responseCommand.ExtFields.(map[string]interface{})
				if ok {
					if nextBeginOffsetInter, ok := pullResult["nextBeginOffset"]; ok {
						if nextBeginOffsetStr, ok := nextBeginOffsetInter.(string); ok {
							nextBeginOffset, err = strconv.ParseInt(nextBeginOffsetStr, 10, 64)
							if err != nil {
								logger.Error("%s", err)
								return
							}

						}

					}

				}

				// 解析消息所有属性
				msgs := decodeMessage(responseFuture.responseCommand.Body)
				// 执行消息消费业务，返回消息被消费的结果
				result, err := self.messageListener(msgs)
				if err != nil {
					logger.Error("consume msg error: %s", err)
					//TODO retry
				} else {
					if result == Action.CommitMessage {
						// 业务处理完毕并返回CommitMessage，消费成功，更新offset
						self.offsetStore.updateOffset(pullRequest.messageQueue, nextBeginOffset, false)
					} else {
						//  业务处理完毕并返回ReconsumeLater，消费失败，消息进入重试队列
						logger.Info("consume failed, it's return ReconsumeLater.")
					}
				}

			} else if responseCommand.Code == PULL_NOT_FOUND {
			} else if responseCommand.Code == PULL_RETRY_IMMEDIATELY || responseCommand.Code == PULL_OFFSET_MOVED {
				logger.Error("pull message error,code=%d,request=%v", responseCommand.Code, requestHeader)
				var err error
				pullResult, ok := responseCommand.ExtFields.(map[string]interface{})
				if ok {
					if nextBeginOffsetInter, ok := pullResult["nextBeginOffset"]; ok {
						if nextBeginOffsetStr, ok := nextBeginOffsetInter.(string); ok {
							nextBeginOffset, err = strconv.ParseInt(nextBeginOffsetStr, 10, 64)
							if err != nil {
								logger.Error("%s", err)
							}

						}

					}

				}

				//time.Sleep(1 * time.Second)
			} else {
				logger.Error(fmt.Sprintf("pull message error,code=%d,body=%s", responseCommand.Code, string(responseCommand.Body)))
				logger.Error("%#v", pullRequest.messageQueue)
				time.Sleep(1 * time.Second)
			}
		} else {
			logger.Error("responseFuture is nil")
		}

		nextPullRequest := &PullRequest{
			consumerGroup: pullRequest.consumerGroup,
			nextOffset:    nextBeginOffset,
			messageQueue:  pullRequest.messageQueue,
		}

		self.mqClient.pullMessageService.pullRequestQueue <- nextPullRequest
	}

	brokerAddr, _, found := self.mqClient.findBrokerAddressInSubscribe(pullRequest.messageQueue.BrokerName, 0, false)

	if found {
		currOpaque := atomic.AddInt32(&opaque, 1)
		remotingCommand := new(RemotingCommand)
		remotingCommand.Code = PULL_MESSAGE
		remotingCommand.Opaque = currOpaque
		remotingCommand.Flag = 0
		remotingCommand.Language = LANGUAGE
		remotingCommand.Version = 79

		remotingCommand.ExtFields = requestHeader

		self.remotingClient.invokeAsync(brokerAddr, remotingCommand, 1000, pullCallback)
	}
}

func (self *DefaultConsumer) updateTopicSubscribeInfo(topic string, info []*MessageQueue) {
	if self.rebalance.subscriptionInner != nil {
		self.rebalance.subscriptionInnerLock.RLock()
		_, ok := self.rebalance.subscriptionInner[topic]
		self.rebalance.subscriptionInnerLock.RUnlock()
		if ok {
			self.rebalance.subscriptionInnerLock.Lock()
			self.rebalance.topicSubscribeInfoTable[topic] = info
			self.rebalance.subscriptionInnerLock.Unlock()
		}
	}
}

// all suscribe(topic-->tags)
func (self *DefaultConsumer) subscriptions() []*SubscriptionData {
	subscriptions := make([]*SubscriptionData, 0)
	for _, subscription := range self.rebalance.subscriptionInner {
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions
}

func (self *DefaultConsumer) doRebalance() {
	self.rebalance.doRebalance()
}
