package rocketmq

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pquerna/ffjson/ffjson"
)

type GroupList struct {
	GroupList []string `json:"groupList"`
}

type QueueData struct {
	BrokerName     string
	ReadQueueNums  int32
	WriteQueueNums int32
	Perm           int32
	TopicSynFlag   int32
}

type BrokerData struct {
	BrokerName      string           `json:"brokerName"`
	BrokerAddrs     map[int64]string `json:"brokerAddrs"`
	BrokerAddrsLock sync.RWMutex     `json:"-"`
}

type TopicRouteData struct {
	OrderTopicConf string
	QueueDatas     []*QueueData
	BrokerDatas    []*BrokerData
}

type MqClient struct {
	clientId            string
	conf                *Config
	brokerAddrTable     map[string]map[int64]string //map[brokerName]map[bokerId]addrs
	brokerAddrTableLock sync.RWMutex
	consumerTable       map[string]*DefaultConsumer
	consumerTableLock   sync.RWMutex
	topicRouteTable     map[string]*TopicRouteData
	topicRouteTableLock sync.RWMutex
	remotingClient      RemotingClient
	pullMessageService  *PullMessageService
	projectGroupPrefix  string
}

func NewMqClient() *MqClient {
	return &MqClient{
		brokerAddrTable: make(map[string]map[int64]string),
		consumerTable:   make(map[string]*DefaultConsumer),
		topicRouteTable: make(map[string]*TopicRouteData),
	}
}
func (self *MqClient) findBrokerAddressInSubscribe(brokerName string, brokerId int64, onlyThisBroker bool) (brokerAddr string, slave bool, found bool) {
	slave = false
	found = false
	self.brokerAddrTableLock.RLock()
	brokerMap, ok := self.brokerAddrTable[brokerName]
	self.brokerAddrTableLock.RUnlock()
	if ok {
		brokerAddr, ok = brokerMap[brokerId]
		slave = (brokerId != 0)
		found = ok

		if !found && !onlyThisBroker {
			var id int64
			for id, brokerAddr = range brokerMap {
				slave = (id != 0)
				found = true
				break
			}
		}
	}

	return
}

func (self *MqClient) findBrokerAddressInAdmin(brokerName string) (addr string, found, slave bool) {
	found = false
	slave = false
	self.brokerAddrTableLock.RLock()
	brokers, ok := self.brokerAddrTable[brokerName]
	self.brokerAddrTableLock.RUnlock()
	if ok {
		for brokerId, addr := range brokers {

			if addr != "" {
				found = true
				if brokerId == 0 {
					slave = false
				} else {
					slave = true
				}
				break
			}
		}
	}

	return
}

func (self *MqClient) findBrokerAddrByTopic(topic string) (addr string, ok bool) {
	self.topicRouteTableLock.RLock()
	topicRouteData, ok := self.topicRouteTable[topic]
	self.topicRouteTableLock.RUnlock()
	if !ok {
		return "", ok
	}

	brokers := topicRouteData.BrokerDatas
	if brokers != nil && len(brokers) > 0 {
		brokerData := brokers[0]
		if ok {
			brokerData.BrokerAddrsLock.RLock()
			addr, ok = brokerData.BrokerAddrs[0]
			brokerData.BrokerAddrsLock.RUnlock()

			if ok {
				return
			}
			for _, addr = range brokerData.BrokerAddrs {
				return addr, ok
			}
		}
	}
	return
}

func (self *MqClient) findConsumerIdList(topic string, groupName string) ([]string, error) {
	brokerAddr, ok := self.findBrokerAddrByTopic(topic)
	if !ok {
		err := self.updateTopicRouteInfoFromNameServerByTopic(topic)
		logger.Error("%s", err)
		brokerAddr, ok = self.findBrokerAddrByTopic(topic)
	}

	if ok {
		return self.getConsumerIdListByGroup(brokerAddr, groupName, 3000)
	}

	return nil, errors.New("can't find broker")

}

type GetConsumerListByGroupResponseBody struct {
	ConsumerIdList []string `json:"consumerIdList"`
}

func (self *MqClient) getConsumerIdListByGroup(addr string, consumerGroup string, timeoutMillis int64) ([]string, error) {
	requestHeader := new(GetConsumerListByGroupRequestHeader)
	requestHeader.ConsumerGroup = consumerGroup

	currOpaque := atomic.AddInt32(&opaque, 1)
	request := &RemotingCommand{
		Code:      GET_CONSUMER_LIST_BY_GROUP,
		Language:  LANGUAGE,
		Version:   79,
		Opaque:    currOpaque,
		Flag:      0,
		ExtFields: requestHeader,
	}

	response, err := self.remotingClient.invokeSync(addr, request, timeoutMillis)
	if err != nil {
		logger.Error("%s", err)
		return nil, err
	}

	if response.Code == SUCCESS {
		getConsumerListByGroupResponseBody := new(GetConsumerListByGroupResponseBody)
		//bodyjson := strings.Replace(string(response.Body), "0:", "\"0\":", -1)
		//bodyjson = strings.Replace(bodyjson, "1:", "\"1\":", -1)
		//err := ffjson.Unmarshal([]byte(bodyjson), getConsumerListByGroupResponseBody)
		err := ffjson.Unmarshal(response.Body, getConsumerListByGroupResponseBody)
		if err != nil {
			logger.Error("%s", err)
			return nil, err
		}
		return getConsumerListByGroupResponseBody.ConsumerIdList, nil
	}

	return nil, errors.New("getConsumerIdListByGroup error")
}

func (self *MqClient) getTopicRouteInfoFromNameServer(topic string, timeoutMillis int64) (*TopicRouteData, error) {
	requestHeader := &GetRouteInfoRequestHeader{
		Topic: topic,
	}

	remotingCommand := new(RemotingCommand)
	remotingCommand.Code = GET_ROUTEINTO_BY_TOPIC
	currOpaque := atomic.AddInt32(&opaque, 1)
	remotingCommand.Opaque = currOpaque
	remotingCommand.Flag = 0
	remotingCommand.Language = LANGUAGE
	remotingCommand.Version = 79

	remotingCommand.ExtFields = requestHeader
	response, err := self.remotingClient.invokeSync(self.conf.Nameserver, remotingCommand, timeoutMillis)
	if err != nil {
		return nil, err
	}
	if response.Code == SUCCESS {
		topicRouteData := new(TopicRouteData)
		err = ffjson.Unmarshal(response.Body, topicRouteData)
		if err != nil {
			logger.Error("%s", err)
			return nil, err
		}
		return topicRouteData, nil
	} else {
		return nil, errors.New(fmt.Sprintf("get topicRouteInfo from nameServer error[code:%d,topic:%s]", response.Code, topic))
	}
}

func (self *MqClient) updateTopicRouteInfoFromNameServer() {
	for _, consumer := range self.consumerTable {
		subscriptions := consumer.subscriptions()
		for _, subData := range subscriptions {
			self.updateTopicRouteInfoFromNameServerByTopic(subData.Topic)
		}
	}
}

func (self *MqClient) updateTopicRouteInfoFromNameServerByTopic(topic string) error {

	topicRouteData, err := self.getTopicRouteInfoFromNameServer(topic, 3000*1000)
	if err != nil {
		logger.Error("%s", err)
		return err
	}

	for _, bd := range topicRouteData.BrokerDatas {
		self.brokerAddrTableLock.Lock()
		self.brokerAddrTable[bd.BrokerName] = bd.BrokerAddrs
		self.brokerAddrTableLock.Unlock()
	}

	mqList := make([]*MessageQueue, 0)
	for _, queueData := range topicRouteData.QueueDatas {
		var i int32
		for i = 0; i < queueData.ReadQueueNums; i++ {
			mq := &MessageQueue{
				Topic:      topic,
				BrokerName: queueData.BrokerName,
				QueueId:    i,
			}

			mqList = append(mqList, mq)
		}
	}

	for _, consumer := range self.consumerTable {
		consumer.updateTopicSubscribeInfo(topic, mqList)
	}
	self.topicRouteTableLock.Lock()
	self.topicRouteTable[topic] = topicRouteData
	self.topicRouteTableLock.Unlock()

	return nil
}

type ConsumerData struct {
	GroupName           string
	ConsumerType        string
	MessageModel        string
	ConsumeFromWhere    string
	SubscriptionDataSet []*SubscriptionData
	UnitMode            bool
}

type HeartbeatData struct {
	ClientId        string
	ConsumerDataSet []*ConsumerData
}

func (self *MqClient) prepareHeartbeatData() *HeartbeatData {
	heartbeatData := new(HeartbeatData)
	heartbeatData.ClientId = self.clientId
	heartbeatData.ConsumerDataSet = make([]*ConsumerData, 0)
	for group, consumer := range self.consumerTable {
		consumerData := new(ConsumerData)
		consumerData.GroupName = group
		consumerData.ConsumerType = consumer.consumerType
		consumerData.ConsumeFromWhere = consumer.consumeFromWhere
		consumerData.MessageModel = consumer.messageModel
		consumerData.SubscriptionDataSet = consumer.subscriptions()
		consumerData.UnitMode = consumer.unitMode

		heartbeatData.ConsumerDataSet = append(heartbeatData.ConsumerDataSet, consumerData)
	}
	return heartbeatData
}

func (self *MqClient) sendHeartbeatToAllBrokerWithLock() error {
	heartbeatData := self.prepareHeartbeatData()
	if len(heartbeatData.ConsumerDataSet) == 0 {
		return errors.New("send heartbeat error")
	}

	self.brokerAddrTableLock.RLock()
	for _, brokerTable := range self.brokerAddrTable {
		for brokerId, addr := range brokerTable {
			if addr == "" || brokerId != 0 {
				continue
			}
			currOpaque := atomic.AddInt32(&opaque, 1)
			remotingCommand := &RemotingCommand{
				Code:     HEART_BEAT,
				Language: LANGUAGE,
				Version:  79,
				Opaque:   currOpaque,
				Flag:     0,
			}

			data, err := ffjson.Marshal(*heartbeatData)
			if err != nil {
				logger.Error("%s", err)
				return err
			}
			remotingCommand.Body = data
			logger.Info("send heartbeat to broker[%s]", addr)
			response, err := self.remotingClient.invokeSync(addr, remotingCommand, 3000)
			if err != nil {
				logger.Error("%s", err)
			} else {
				if response == nil || response.Code != SUCCESS {
					logger.Error("send heartbeat response  error")
				}
			}
		}
	}
	self.brokerAddrTableLock.RUnlock()
	return nil
}

func (self *MqClient) startScheduledTask() {
	go func() {
		updateTopicRouteTimer := time.NewTimer(5 * time.Second)
		for {
			<-updateTopicRouteTimer.C
			self.updateTopicRouteInfoFromNameServer()
			updateTopicRouteTimer.Reset(5 * time.Second)
		}
	}()

	go func() {
		heartbeatTimer := time.NewTimer(10 * time.Second)
		for {
			<-heartbeatTimer.C
			self.sendHeartbeatToAllBrokerWithLock()
			heartbeatTimer.Reset(5 * time.Second)
		}
	}()

	go func() {
		rebalanceTimer := time.NewTimer(15 * time.Second)
		for {
			<-rebalanceTimer.C
			self.doRebalance()
			rebalanceTimer.Reset(30 * time.Second)
		}
	}()

	go func() {
		timeoutTimer := time.NewTimer(3 * time.Second)
		for {
			<-timeoutTimer.C
			self.remotingClient.ScanResponseTable()
			timeoutTimer.Reset(time.Second)
		}
	}()

}

func (self *MqClient) doRebalance() {
	for _, consumer := range self.consumerTable {
		consumer.doRebalance()
	}
}

func (self *MqClient) start() {
	self.startScheduledTask()
	go self.pullMessageService.start()
}

func (self *MqClient) queryConsumerOffset(addr string, requestHeader *QueryConsumerOffsetRequestHeader, timeoutMillis int64) (int64, error) {
	currOpaque := atomic.AddInt32(&opaque, 1)
	remotingCommand := &RemotingCommand{
		Code:     QUERY_CONSUMER_OFFSET,
		Language: LANGUAGE,
		Version:  79,
		Opaque:   currOpaque,
		Flag:     0,
	}

	remotingCommand.ExtFields = requestHeader
	reponse, err := self.remotingClient.invokeSync(addr, remotingCommand, timeoutMillis)

	if err != nil {
		logger.Error("%s", err)
		return 0, err
	}

	if reponse.Code == QUERY_NOT_FOUND {
		return 0, nil
	}

	if extFields, ok := (reponse.ExtFields).(map[string]interface{}); ok {
		if offsetInter, ok := extFields["offset"]; ok {
			if offsetStr, ok := offsetInter.(string); ok {
				offset, err := strconv.ParseInt(offsetStr, 10, 64)
				if err != nil {
					logger.Error("%s", err)
					return 0, err
				}
				return offset, nil

			}
		}
	}
	logger.Error("%#v %#v", requestHeader, reponse)
	return 0, errors.New("query offset error")
}

func (self *MqClient) updateConsumerOffsetOneway(addr string, header *UpdateConsumerOffsetRequestHeader, timeoutMillis int64) {

	currOpaque := atomic.AddInt32(&opaque, 1)
	remotingCommand := &RemotingCommand{
		Code:      QUERY_CONSUMER_OFFSET,
		Language:  LANGUAGE,
		Version:   79,
		Opaque:    currOpaque,
		Flag:      0,
		ExtFields: header,
	}

	self.remotingClient.invokeSync(addr, remotingCommand, timeoutMillis)
}

func (self *MqClient) queryTopicConsumeByWho(addr, topic string, timeoutMillis int64) (*GroupList, error) {
	requestHeader := &QueryTopicConsumeByWhoRequestHeader{
		Topic: topic,
	}

	remotingCommand := createRemotingCommand(QUERY_TOPIC_CONSUME_BY_WHO)
	remotingCommand.ExtFields = requestHeader

	response, err := self.remotingClient.invokeSync(addr, remotingCommand, timeoutMillis)
	if err != nil {
		return nil, err
	}
	if response.Code == SUCCESS {
		groupList := new(GroupList)

		err = ffjson.Unmarshal(response.Body, groupList)
		if err != nil {
			logger.Error("%s", err)
			return nil, err
		}
		return groupList, nil
	}

	return nil, errors.New(fmt.Sprintf("query topic consume from brokerServer error[code:%d,topic:%s]", response.Code, topic))
}

func (self *MqClient) getConsumerConnectionList(addr, group string, timeoutMillis int64) (*ConsumerConnection, error) {
	// TODO: group
	requestHeader := &GetConsumerConnectionListRequestHeader{
		ConsumerGroup: group,
	}

	remotingCommand := createRemotingCommand(GET_CONSUMER_CONNECTION_LIST)
	remotingCommand.ExtFields = requestHeader

	response, err := self.remotingClient.invokeSync(addr, remotingCommand, timeoutMillis)
	if err != nil {
		return nil, err
	}

	if response.Code == CONSUMER_NOT_ONLINE {
		return &ConsumerConnection{}, nil
	}
	if response.Code == SUCCESS {
		consumerConnection := new(ConsumerConnection)
		err = ffjson.Unmarshal(response.Body, consumerConnection)
		if err != nil {
			logger.Error("%s", err)
			return nil, err
		}
		return consumerConnection, nil
	}

	return nil, errors.New(fmt.Sprintf("get consume connection from brokerServer error[code:%d,group:%s]", response.Code, group))
}

func (self *MqClient) getConsumeStats(addr, group, topic string, timeoutMillis int64) (*ConsumeStats, error) {
	// TODO: group
	requestHeader := &GetConsumeStatsRequestHeader{
		ConsumerGroup: group,
		Topic:         topic,
	}

	remotingCommand := createRemotingCommand(GET_CONSUME_STATS)
	remotingCommand.ExtFields = requestHeader

	response, err := self.remotingClient.invokeSync(addr, remotingCommand, timeoutMillis)
	if err != nil {
		return nil, err
	}
	if response.Code == SUCCESS {
		consumeStats := new(ConsumeStats)
		err = ffjson.Unmarshal(response.Body, consumeStats)
		if err != nil {
			logger.Error("%s", err)
			return nil, err
		}
		return consumeStats, nil
	}

	return nil, errors.New(fmt.Sprintf("get consume connection from brokerServer error[code:%d,group:%s]", response.Code, group))

}
func (self *MqClient) getBrokerRuntimeInfo(addr string, timeoutMillis int64) (*BrokerRuntimeInfo, error) {

	remotingCommand := createRemotingCommand(GET_BROKER_RUNTIME_INFO)
	response, err := self.remotingClient.invokeSync(addr, remotingCommand, timeoutMillis)
	if err != nil {
		return nil, err
	}
	if response.Code == SUCCESS {
		table := &kvTable{}
		err = ffjson.Unmarshal(response.Body, table)
		if err != nil {
			logger.Error("%s", err)
			return nil, err
		}

		brokerRuntimeInfo := new(BrokerRuntimeInfo)
		if v, ok := table.Table["brokerVersionDesc"]; ok {
			brokerRuntimeInfo.BrokerVersionDesc = v
		}
		if v, ok := table.Table["brokerVersion"]; ok {
			brokerRuntimeInfo.BrokerVersion = v
		}
		if v, ok := table.Table["msgPutTotalYesterdayMorning"]; ok {
			brokerRuntimeInfo.MsgPutTotalYesterdayMorning = v
		}
		if v, ok := table.Table["msgPutTotalTodayMorning"]; ok {
			brokerRuntimeInfo.MsgPutTotalTodayMorning = v
		}
		if v, ok := table.Table["msgPutTotalTodayNow"]; ok {
			brokerRuntimeInfo.MsgPutTotalTodayNow = v
		}
		if v, ok := table.Table["msgGetTotalYesterdayMorning"]; ok {
			brokerRuntimeInfo.MsgGetTotalYesterdayMorning = v
		}
		if v, ok := table.Table["msgGetTotalTodayMorning"]; ok {
			brokerRuntimeInfo.MsgGetTotalTodayMorning = v
		}
		if v, ok := table.Table["msgGetTotalTodayNow"]; ok {
			brokerRuntimeInfo.MsgGetTotalTodayNow = v
		}
		if v, ok := table.Table["sendThreadPoolQueueSize"]; ok {
			brokerRuntimeInfo.SendThreadPoolQueueSize = v
		}
		if v, ok := table.Table["sendThreadPoolQueueCapacity"]; ok {
			brokerRuntimeInfo.SendThreadPoolQueueCapacity = v
		}
		if v, ok := table.Table["putTps"]; ok {
			vs := parseTpsString(v)
			brokerRuntimeInfo.InTps = vs[0]
		}
		if v, ok := table.Table["getTransferedTps"]; ok {
			vs := parseTpsString(v)
			brokerRuntimeInfo.OutTps = vs[0]
		}

		return brokerRuntimeInfo, nil
	}

	return nil, errors.New(fmt.Sprintf("get broker runtime from brokerServer error[code:%d,addr:%s]", response.Code, addr))
}

func (self *MqClient) getKVConfigValue(addr, namespace, key string, timeoutMillis int64) (string, error) {
	requestHeader := &GetKVConfigRequestHeader{
		Namespace: namespace,
		Key:       key,
	}

	remotingCommand := createRemotingCommand(GET_KV_CONFIG)
	remotingCommand.ExtFields = requestHeader

	response, err := self.remotingClient.invokeSync(addr, remotingCommand, timeoutMillis)
	if err != nil {
		return "", err
	}
	if response.Code == QUERY_NOT_FOUND {
		return "", nil
	}
	if response.Code == SUCCESS {
		getKVConfigResponseHeader := new(GetKVConfigResponseHeader)
		err := ffjson.Unmarshal(response.Body, getKVConfigResponseHeader)
		if err != nil {
			logger.Error("%s", err)
			return "", err
		}
		return getKVConfigResponseHeader.Value, nil
	}

	return "", errors.New(fmt.Sprintf("get kv config error[code:%d,namespace:%s, key:%s]", response.Code, namespace, key))
}

func (self *MqClient) Close() error {
	return self.remotingClient.close()
}

func (self *MqClient) initInfo() error {
	pjgp, err := self.getKVConfigValue(self.conf.Nameserver, NAMESPACE_PROJECT_CONFIG, self.conf.ClientIp, defaultTimeoutMillis)
	if err != nil {
		return err
	}

	self.projectGroupPrefix = pjgp
	return nil
}

func parseTpsString(str string) []float64 {
	var tpsf []float64
	tpss := strings.Split(str, " ")

	for _, v := range tpss {
		vi, err := strconv.ParseFloat(v, 64)
		if err != nil {
			continue
		}
		vi = math.Floor(vi*1e2) * 1e-2
		tpsf = append(tpsf, vi)
	}

	return tpsf
}
