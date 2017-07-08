package rocketmq

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/pquerna/ffjson/ffjson"
)

const (
	defaultTimeoutMillis     = 3000 * 1000
	RETRY_GROUP_TOPIC_PREFIX = "%RETRY%"
	MASTER_ID                = 0
)

// TopicData topic data
type TopicData struct {
	List []string `json:"topicList"`
}

// ClusterInfo cluster info
type ClusterInfo struct {
	BrokerAddrTable  map[string]BrokerData `json:"brokerAddrTable"`
	ClusterAddrTable map[string][]string   `json:"clusterAddrTable"`
}

type Admin interface {
	SearchTopic() (*TopicData, error)
	SearchClusterInfo() (*ClusterInfo, error)
	CreateTopic(topicName, clusterName string, readQueueNum int, writeQueueNum int, order bool) (bool, error)
	QueryMessageById(msgId string) (*MessageExt, error)
	MessageTrackDetail(msg *MessageExt) ([]*MessageTrack, error)
	QueryConsumerGroupByTopic(topic string) ([]string, error)
	QueryConsumerProgress(consumeGroupId string) (*ConsumerProgress, error)
	QueryConsumerConnection(consumeGroupId string) (*ConsumerConnection, error)
	FetchBrokerRuntimeStats(brokerAddr string) (*BrokerRuntimeInfo, error)
	Close() error
}

type DefaultAdmin struct {
	conf           *Config
	remotingClient RemotingClient
	mqClient       *MqClient
}

func NewDefaultAdmin(conf *Config) (Admin, error) {

	if conf.ClientIp == "" {
		conf.ClientIp = DEFAULT_IP
	}

	remotingClient := NewDefaultRemotingClient()
	mqClient := NewMqClient()
	mqClient.remotingClient = remotingClient
	mqClient.conf = conf
	mqClient.clientId = fmt.Sprintf("%s@%d#%d#%d", conf.ClientIp, os.Getpid(), HashCode(conf.Nameserver), UnixNano())

	err := mqClient.initInfo()
	if err != nil {
		return nil, err
	}

	admin := &DefaultAdmin{
		conf:           conf,
		remotingClient: remotingClient,
		mqClient:       mqClient,
	}

	return admin, nil
}

// GetTopic 查询topic列表接口
func (admin *DefaultAdmin) SearchTopic() (*TopicData, error) {
	remotingCommand := createRemotingCommand(GET_ALL_TOPIC_LIST_FROM_NAMESERVER)
	response, err := admin.remotingClient.invokeSync(admin.conf.Nameserver, remotingCommand, defaultTimeoutMillis)
	if err != nil {
		return nil, err
	}

	if response != nil && response.Code == SUCCESS {
		topicData := new(TopicData)
		err = ffjson.Unmarshal(response.Body, topicData)
		if err != nil {
			logger.Error("%s", err)
			return nil, err
		}
		return topicData, nil
	}

	return &TopicData{}, nil
}

// SearchClusterInfo 查询集群的信息
func (admin *DefaultAdmin) SearchClusterInfo() (*ClusterInfo, error) {
	remotingCommand := createRemotingCommand(GET_BROKER_CLUSTER_INFO)
	response, err := admin.remotingClient.invokeSync(admin.conf.Nameserver, remotingCommand, defaultTimeoutMillis)
	if err != nil {
		return nil, err
	}

	if response != nil && response.Code == SUCCESS {
		clusterInfo := new(ClusterInfo)
		err = ffjson.Unmarshal(response.Body, clusterInfo)
		if err != nil {
			logger.Error("%s", err)
			return nil, err
		}
		return clusterInfo, nil
	}

	return &ClusterInfo{}, nil
}

// CreateTopic 创建topic
func (admin *DefaultAdmin) CreateTopic(topicName, clusterName string, readQueueNum int, writeQueueNum int, order bool) (bool, error) {
	clusterInfo, _ := admin.SearchClusterInfo()
	if clusterInfo == nil {
		return false, errors.New("CreateTopic get clusterInfo failed")
	}

	remotingCommand := createRemotingCommand(UPDATE_AND_CREATE_TOPIC)
	extFields := make(map[string]string)
	extFields["topic"] = topicName
	extFields["defaultTopic"] = "MY_DEFAULT_TOPIC"
	extFields["readQueueNums"] = strconv.Itoa(readQueueNum)
	extFields["writeQueueNums"] = strconv.Itoa(writeQueueNum)
	extFields["perm"] = "6"
	extFields["topicFilterType"] = "SINGLE_TAG"
	extFields["topicSysFlag"] = "0"
	extFields["order"] = fmt.Sprintf("%t", order)
	remotingCommand.ExtFields = extFields

	var flag bool
	for clsName, brokerNames := range clusterInfo.ClusterAddrTable {
		if clsName != clusterName {
			continue
		}
		for _, brName := range brokerNames {
			if brokerData, ok := clusterInfo.BrokerAddrTable[brName]; ok {
				for k, v := range brokerData.BrokerAddrs {
					if k == 0 {
						flag = true
						cmd, err := admin.remotingClient.invokeSync(v, remotingCommand, defaultTimeoutMillis)
						if err != nil {
							return false, err
						}
						if cmd.Code != SUCCESS {
							return false, fmt.Errorf("创建Topic失败，错误码:%d", cmd.Code)
						}
					}
				}
			}
		}
	}

	return flag, nil
}

// QueryMessageById 根据msgid查询消息详情
func (admin *DefaultAdmin) QueryMessageById(msgId string) (*MessageExt, error) {

	addr, offset, err := decodeMessageId(msgId)
	if err != nil {
		return nil, err
	}

	requestHeader := new(ViewMessageRequestHeader)
	requestHeader.Offset = offset
	remotingCommand := createRemotingCommand(VIEW_MESSAGE_BY_ID)
	remotingCommand.ExtFields = requestHeader

	response, err := admin.remotingClient.invokeSync(addr, remotingCommand, defaultTimeoutMillis)
	if err != nil {
		return nil, err
	}

	if response != nil && response.Code == SUCCESS {
		msgs := decodeMessage(response.Body)
		if len(msgs) > 0 {
			msg := msgs[0]
			msg.MsgId = msgId
			return msg, nil
		}
	}

	return nil, nil
}

// MessageTrackDetail 根据msgid查询消息轨迹
func (admin *DefaultAdmin) MessageTrackDetail(msg *MessageExt) ([]*MessageTrack, error) {

	groupList, err := admin.queryTopicConsumeByWho(msg.Topic)
	if err != nil {
		return nil, err
	}

	if groupList == nil {
		return []*MessageTrack{}, nil
	}

	var tracks []*MessageTrack
	for _, group := range groupList.GroupList {
		track := new(MessageTrack)
		track.ConsumerGroup = group
		track.TrackType = TrackTypes.UnknowExeption
		consumerConnection, err := admin.getConsumerConnectionList(group)
		if err != nil {
			track.ExceptionDesc = err.Error()
			tracks = append(tracks, track)
			continue
		}
		if len(consumerConnection.ConnectionSet) == 0 {
			track.ExceptionDesc = fmt.Sprintf("CODE: %d DESC: the consumer group[%s] not online.", CONSUMER_NOT_ONLINE, group)
			tracks = append(tracks, track)
			continue
		}

		switch consumerConnection.ConsumeType {
		case ConsumeTypes.Actively:
			track.TrackType = TrackTypes.SubscribedButPull
		case ConsumeTypes.Passively:
			flag, err := admin.consumed(msg, group)
			if err != nil {
				track.ExceptionDesc = err.Error()
				break
			}

			if flag {
				track.TrackType = TrackTypes.SubscribedAndConsumed
				// 查看订阅关系是否匹配
				for topic, subscriptionData := range consumerConnection.SubscriptionTable {
					if topic != msg.Topic {
						continue
					}
					for _, tag := range subscriptionData.TagsSet {
						if tag != msg.Tag() && tag != "*" {
							track.TrackType = TrackTypes.SubscribedButFilterd
						}
					}
				}

			} else {
				track.TrackType = TrackTypes.SubscribedAndNotConsumeYet
			}
		default:
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (admin *DefaultAdmin) consumed(msg *MessageExt, group string) (bool, error) {
	consumeStats, err := admin.getConsumeStats(group, "")
	if err != nil {
		return false, err
	}

	clusterInfo, err := admin.SearchClusterInfo()
	if err != nil {
		return false, err
	}

	for mq, ow := range consumeStats.OffsetTable {
		if mq.Topic == msg.Topic && mq.QueueId == msg.QueueId {
			if bd, ok := clusterInfo.BrokerAddrTable[mq.BrokerName]; ok {
				if addr, yes := bd.BrokerAddrs[MASTER_ID]; yes {
					if addr == msg.StoreHost && ow.ConsumerOffset > msg.QueueOffset {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}

// 在线consumerGroup
func (admin *DefaultAdmin) queryTopicConsumeByWho(topic string) (*GroupList, error) {

	topicRouteData, err := admin.mqClient.getTopicRouteInfoFromNameServer(topic, defaultTimeoutMillis)
	if err != nil {
		return nil, err
	}

	for _, bd := range topicRouteData.BrokerDatas {
		for _, ba := range bd.BrokerAddrs {
			groupList, err := admin.mqClient.queryTopicConsumeByWho(ba, topic, defaultTimeoutMillis)
			if err != nil {
				return nil, err
			}

			return groupList, nil
		}
	}

	return &GroupList{}, nil
}

// 在线consume
func (admin *DefaultAdmin) getConsumerConnectionList(group string) (*ConsumerConnection, error) {

	topic := getRetryTopic(group)
	topicRouteData, err := admin.mqClient.getTopicRouteInfoFromNameServer(topic, defaultTimeoutMillis)
	if err != nil {
		return nil, err
	}

	for _, bd := range topicRouteData.BrokerDatas {
		for _, ba := range bd.BrokerAddrs {
			consumerConnection, err := admin.mqClient.getConsumerConnectionList(ba, group, defaultTimeoutMillis)
			if err != nil {
				return nil, err
			}

			return consumerConnection, nil
		}
	}

	return &ConsumerConnection{}, nil
}

// QueryConsumerGroupByTopic 根据topic查询消费组
func (admin *DefaultAdmin) QueryConsumerGroupByTopic(topic string) ([]string, error) {
	groupList, err := admin.queryTopicConsumeByWho(topic)
	if err != nil {
		return nil, err
	}

	if groupList == nil {
		return []string{}, nil
	}

	return groupList.GroupList, nil
}

// 消费进度
func (admin *DefaultAdmin) getConsumeStats(consumeGroupId, topic string) (*ConsumeStats, error) {
	retryTopic := getRetryTopic(consumeGroupId)
	topicRouteData, err := admin.mqClient.getTopicRouteInfoFromNameServer(retryTopic, defaultTimeoutMillis)
	if err != nil {
		return nil, err
	}

	result := new(ConsumeStats)
	result.OffsetTable = make(map[MessageQueue]OffsetWrapper)
	for _, bd := range topicRouteData.BrokerDatas {
		for _, ba := range bd.BrokerAddrs {
			// 由于查询时间戳会产生IO操作，可能会耗时较长，所以超时时间设置为15s
			consumeStats, err := admin.mqClient.getConsumeStats(ba, consumeGroupId, topic, defaultTimeoutMillis*5)
			if err != nil {
				return nil, err
			}

			for mq, ow := range consumeStats.OffsetTable {
				// BrokerOffset < ConsumerOffset的broke是未消费groupId的节点
				if ow.BrokerOffset < ow.ConsumerOffset {
					continue
				}
				result.OffsetTable[mq] = ow
			}

			result.ConsumeTps += consumeStats.ConsumeTps
		}
	}

	return result, nil
}

// QueryConsumerProgress 根据consumeGroupId查询消费进度
func (admin *DefaultAdmin) QueryConsumerProgress(consumeGroupId string) (*ConsumerProgress, error) {
	consumeStats, err := admin.getConsumeStats(consumeGroupId, "")
	if err != nil {
		return nil, err
	}

	length := len(consumeStats.OffsetTable)
	if length == 0 {
		return &ConsumerProgress{}, nil
	}

	mqs := make([]*MessageQueue, 0, len(consumeStats.OffsetTable))
	for mq, _ := range consumeStats.OffsetTable {
		cmq := mq.clone()
		mqs = append(mqs, cmq)
	}
	var messageQueues MessageQueues = mqs
	sort.Sort(messageQueues)

	consumerProgress := new(ConsumerProgress)
	for _, mq := range messageQueues {
		if ow, ok := consumeStats.OffsetTable[*mq]; ok {
			cg := new(ConsumerGroup)
			cg.Topic = mq.Topic
			cg.BrokerName = mq.BrokerName
			cg.BrokerOffset = ow.BrokerOffset
			cg.ConsumerOffset = ow.ConsumerOffset
			cg.QueueId = mq.QueueId
			cg.Diff = ow.BrokerOffset - ow.ConsumerOffset
			consumerProgress.Diff += cg.Diff
			consumerProgress.List = append(consumerProgress.List, cg)
		}
	}
	consumerProgress.Tps = consumeStats.ConsumeTps
	consumerProgress.GroupId = consumeGroupId

	return consumerProgress, nil
}

// QueryConsumerConnection 根据consumeGroupId查询消费进程
func (admin *DefaultAdmin) QueryConsumerConnection(consumeGroupId string) (*ConsumerConnection, error) {
	return admin.getConsumerConnectionList(consumeGroupId)
}

// FetchBrokerRuntimeStats 根据broker addr查询在线broker运行统计信息
func (admin *DefaultAdmin) FetchBrokerRuntimeStats(brokerAddr string) (*BrokerRuntimeInfo, error) {
	return admin.mqClient.getBrokerRuntimeInfo(brokerAddr, defaultTimeoutMillis)
}

func (admin *DefaultAdmin) Close() error {
	return admin.mqClient.Close()
}

func getRetryTopic(consumerGroup string) string {
	return RETRY_GROUP_TOPIC_PREFIX + consumerGroup
}
