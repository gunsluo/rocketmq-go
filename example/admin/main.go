package main

import (
	"fmt"

	"github.com/gunsluo/rocketmq-go"
)

var conf *rocketmq.Config

func main() {
	conf = &rocketmq.Config{
		Nameserver: "10.122.1.201:9876",
	}

	admin, err := rocketmq.NewDefaultAdmin(conf)
	if err != nil {
		panic(err)
	}

	searchClusterInfo(admin)
	createTopic(admin)
	searchTopic(admin)
	queryConsumerGroupByTopic(admin)
	queryConsumerConnection(admin)
	queryConsumerProgress(admin)
	queryMessageById(admin)
	fetchBrokerRuntimeStats(admin)

	admin.Close()
}

func searchTopic(admin rocketmq.Admin) {

	topicData, err := admin.SearchTopic()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nsearch cluser's topic from nameserver[%s]: %s\n", conf.Nameserver, topicData.List)
}

func searchClusterInfo(admin rocketmq.Admin) {

	clusterInfo, err := admin.SearchClusterInfo()
	if err != nil {
		panic(err)
	}

	fmt.Printf("search cluster info from nameserver[%s]: \n", conf.Nameserver)
	for clusterName, brokerNames := range clusterInfo.ClusterAddrTable {
		fmt.Printf("\tclusterName: [%s] brokerNames: %s\n", clusterName, brokerNames)
	}
	for _, brokerData := range clusterInfo.BrokerAddrTable {
		fmt.Printf("\t[%s]: ", brokerData.BrokerName)
		for bid, addr := range brokerData.BrokerAddrs {
			fmt.Printf("bid=%d, addr=%s\n", bid, addr)
		}
	}
}

func createTopic(admin rocketmq.Admin) {

	topic := "luoji"
	clusterName := "DefaultCluster"
	ret, err := admin.CreateTopic(topic, clusterName, 8, 8, false)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\ncreate or update topic[%s], result: %v\n", topic, ret)
}

func queryMessageById(admin rocketmq.Admin) {
	msgId := "0A7A01C800002A9F0000000000000779"

	msg, err := admin.QueryMessageById(msgId)
	if err != nil {
		panic(err)
	}
	if msg != nil {
		fmt.Printf("\nquery message by id: msgId=%s|topic=%s|storeTimestamp=%d|bornTimestamp=%d|storeHost=%s|bornHost=%s"+
			"|msgTag=%s|msgKey=%s|sysFlag=%d|storeSize=%d|queueId=%d|queueOffset=%d|commitOffset=%d|Properties=%s|body=%s\n",
			msg.MsgId, msg.Topic, msg.StoreTimestamp, msg.BornTimestamp, msg.StoreHost, msg.BornHost,
			msg.Tag(), msg.Key(), msg.SysFlag, msg.StoreSize, msg.QueueId, msg.QueueOffset, msg.CommitLogOffset, msg.Properties, msg.Body)
	}

	msgTracks, err := admin.MessageTrackDetail(msg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("msg track:\n")
	for i, track := range msgTracks {
		fmt.Printf("%d:\t ConsumerGroup=%s, TrackType=%d, ExceptionDesc=%s\n", i, track.ConsumerGroup, track.TrackType, track.ExceptionDesc)
	}
}

func queryConsumerGroupByTopic(admin rocketmq.Admin) {
	topic := "testTopic"
	groupNames, err := admin.QueryConsumerGroupByTopic(topic)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nquery consumer group by topic[%s]: %s\n", topic, groupNames)
}

func queryConsumerConnection(admin rocketmq.Admin) {
	consumerGroupId := "ConsumerGroupName"
	consumerConnection, err := admin.QueryConsumerConnection(consumerGroupId)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nquery consumer conection by consumeGroup[%s]: consumeType=%s, MessageModel=%s, consumeFromWhere=%s\n",
		consumerGroupId, consumerConnection.ConsumeType, consumerConnection.MessageModel, consumerConnection.ConsumeFromWhere)
	fmt.Printf("client list(%d)\t|index\t|ClientId\t\t|ClientAddr\t\t|Language\t|Version\n", len(consumerConnection.ConnectionSet))
	for i, conn := range consumerConnection.ConnectionSet {
		fmt.Printf("\t\t|%d\t|%s\t|%s\t|%s\t\t|%d\n", i+1, conn.ClientId, conn.ClientAddr, conn.Language, conn.Version)
	}
	fmt.Printf("Subscribed list\t|Topic\t\t|SubString\t\t|ClassFilterMode\t|TagsSet\t|CodeSet\t|SubVersion\n")
	for key, subData := range consumerConnection.SubscriptionTable {
		fmt.Printf("\t\t|%s\t|%s\t|%v\t|%s\t|%d\t|%d\n", key, subData.SubString, subData.ClassFilterMode, subData.TagsSet, subData.CodeSet, subData.SubVersion)
	}
}

func queryConsumerProgress(admin rocketmq.Admin) {
	consumerGroupId := "ConsumerGroupName"
	consumerProgress, err := admin.QueryConsumerProgress(consumerGroupId)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nquery consumer progress by consumeGroup[%s]: tps[%d] diff[%d] \n", consumerGroupId, consumerProgress.Tps, consumerProgress.Diff)
	fmt.Printf("Progress list(%d)\t|Topic\t\t|BrokerName\t\t|QueueId\t|BrokerOffset\t|ConsumerOffset\t|Diff\n", len(consumerProgress.List))
	for _, cg := range consumerProgress.List {
		fmt.Printf("\t|%s\t\t|%s\t\t|%d\t\t|%d\t\t|%d\t\t|%d\n", cg.Topic, cg.BrokerName, cg.QueueId, cg.BrokerOffset, cg.ConsumerOffset, cg.Diff)
	}
}

func fetchBrokerRuntimeStats(admin rocketmq.Admin) {
	brokerAddr := "10.122.1.200:10911"
	brokerRuntimeInfo, err := admin.FetchBrokerRuntimeStats(brokerAddr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nquery broker runtime info[%s]: [%#v] \n", brokerAddr, brokerRuntimeInfo)
}
