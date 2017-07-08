package rocketmq

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strconv"
)

const (
	CompressedFlag = (0x1 << 0)
)

// Action 共享申请
var Action = struct {
	CommitMessage  int //消费成功
	ReconsumeLater int //消费失败(消息进入重试队列，过一会儿消费)
}{
	1,
	2,
}

type Message struct {
	Topic      string            `json:"topic"`
	Flag       int32             `json:"flag"`
	Body       []byte            `json:"body"`
	Properties map[string]string `json:"properties"`
}

type MessageExt struct {
	Message
	QueueId                   int32  `json:"queueId"`
	StoreSize                 int32  `json:"storeSize"`
	QueueOffset               int64  `json:"queueOffset"`
	SysFlag                   int32  `json:"sysFlag"`
	BornTimestamp             int64  `json:"bornTimestamp"`
	BornHost                  string `json:"bornHost"`
	StoreTimestamp            int64  `json:"storeTimestamp"`
	StoreHost                 string `json:"storeHost"`
	MsgId                     string `json:"msgId"`
	CommitLogOffset           int64  `json:"commitLogOffset"`
	BodyCRC                   int32  `json:"bodyCRC"`
	ReconsumeTimes            int32  `json:"reconsumeTimes"`
	PreparedTransactionOffset int64  `json:"preparedTransactionOffset"`
}

func decodeMessage(data []byte) []*MessageExt {
	buf := bytes.NewBuffer(data)
	var storeSize, magicCode, bodyCRC, queueId, flag, sysFlag, reconsumeTimes, bodyLength, bornPort, storePort int32
	var queueOffset, physicOffset, preparedTransactionOffset, bornTimeStamp, storeTimestamp int64
	var topicLen byte
	var topic, body, properties, bornHost, storeHost []byte
	var propertiesLength int16

	msgs := make([]*MessageExt, 0, 32)
	for buf.Len() > 0 {
		msg := new(MessageExt)
		binary.Read(buf, binary.BigEndian, &storeSize)
		binary.Read(buf, binary.BigEndian, &magicCode)
		binary.Read(buf, binary.BigEndian, &bodyCRC)
		binary.Read(buf, binary.BigEndian, &queueId)
		binary.Read(buf, binary.BigEndian, &flag)
		binary.Read(buf, binary.BigEndian, &queueOffset)
		binary.Read(buf, binary.BigEndian, &physicOffset)
		binary.Read(buf, binary.BigEndian, &sysFlag)
		binary.Read(buf, binary.BigEndian, &bornTimeStamp)
		bornHost = make([]byte, 4)
		binary.Read(buf, binary.BigEndian, &bornHost)
		binary.Read(buf, binary.BigEndian, &bornPort)
		binary.Read(buf, binary.BigEndian, &storeTimestamp)
		storeHost = make([]byte, 4)
		binary.Read(buf, binary.BigEndian, &storeHost)
		binary.Read(buf, binary.BigEndian, &storePort)
		binary.Read(buf, binary.BigEndian, &reconsumeTimes)
		binary.Read(buf, binary.BigEndian, &preparedTransactionOffset)
		binary.Read(buf, binary.BigEndian, &bodyLength)
		if bodyLength > 0 {
			body = make([]byte, bodyLength)
			binary.Read(buf, binary.BigEndian, body)

			if (sysFlag & CompressedFlag) == CompressedFlag {
				b := bytes.NewReader(body)
				z, err := zlib.NewReader(b)
				if err != nil {
					logger.Error("%s", err)
					return nil
				}
				defer z.Close()
				body, err = ioutil.ReadAll(z)
				if err != nil {
					logger.Error("%s", err)
					return nil
				}
			}

		}
		binary.Read(buf, binary.BigEndian, &topicLen)
		topic = make([]byte, topicLen)
		binary.Read(buf, binary.BigEndian, &topic)
		binary.Read(buf, binary.BigEndian, &propertiesLength)
		if propertiesLength > 0 {
			properties = make([]byte, propertiesLength)
			binary.Read(buf, binary.BigEndian, &properties)

			// 解析消息属性 Add: tianyuliang Since: 2017/4/19
			msg.Properties = convertProperties(properties)
		}

		if magicCode != -626843481 {
			logger.Info("magic code is error %d", magicCode)
			return nil
		}

		msg.Topic = string(topic)
		msg.QueueId = queueId
		msg.SysFlag = sysFlag
		msg.QueueOffset = queueOffset
		msg.BodyCRC = bodyCRC
		msg.StoreSize = storeSize
		// 解析消息BornHost字段 Add: tianyuliang Since: 2017/5/3
		msg.BornHost = convertHostString(bornHost, bornPort)
		msg.BornTimestamp = bornTimeStamp
		msg.ReconsumeTimes = reconsumeTimes
		msg.Flag = flag
		msg.CommitLogOffset = physicOffset
		// 解析消息ID字段 Add: tianyuliang Since: 2017/4/19
		msg.MsgId = convertMessageId(storeHost, storePort, physicOffset)
		msg.StoreHost = convertHostString(storeHost, storePort)
		msg.StoreTimestamp = storeTimestamp
		msg.PreparedTransactionOffset = preparedTransactionOffset
		msg.Body = body

		msgs = append(msgs, msg)
	}

	return msgs
}

// 解析消息msgId字段(ip + port + commitOffset，其中ip、port长度分别是4位，offset占用8位长度)
// Author: tianyuliang
// Since: 2017/5/4
func convertMessageId(storeHost []byte, storePort int32, offset int64) string {
	buffMsgId := make([]byte, MSG_ID_LENGTH)
	input := bytes.NewBuffer(buffMsgId)
	input.Reset()
	input.Grow(MSG_ID_LENGTH)
	input.Write(storeHost)

	storePortBytes := int32ToBytes(storePort)
	input.Write(storePortBytes)

	offsetBytes := int64ToBytes(offset)
	input.Write(offsetBytes)

	return bytesToHexString(input.Bytes())
}

// 根据ip和port，解析host
// Author: tianyuliang
// Since: 2017/5/4
func convertHostString(ipBytes []byte, port int32) string {
	ip := bytesToIPv4String(ipBytes)

	return fmt.Sprintf("%s:%s", ip, strconv.FormatInt(int64(port), 10))
}

func convertProperties(buf []byte) map[string]string {

	tbuf := buf
	properties := make(map[string]string)
	for len(tbuf) > 0 {
		pi := bytes.IndexByte(tbuf, PROPERTY_SEPARATOR)
		if pi == -1 {
			break
		}

		propertie := tbuf[0:pi]

		ni := bytes.IndexByte(propertie, NAME_VALUE_SEPARATOR)
		if ni == -1 || ni > pi {
			break
		}

		key := string(propertie[0:ni])
		properties[key] = string(propertie[ni+1:])

		tbuf = tbuf[pi+1:]
	}

	return properties
}

func (msg *Message) Propertie(key string) string {
	if msg.Properties == nil || key == "" {
		return ""
	}

	if val, ok := msg.Properties[key]; ok {
		return val
	}

	return ""
}

func (msg *Message) Key() string {
	return msg.Propertie(PROPERTY_KEYS)
}

func (msg *Message) Tag() string {
	return msg.Propertie(PROPERTY_TAGS)
}

func decodeMessageId(msgId string) (string, uint64, error) {

	if len(msgId) != 32 {
		return "", 0, fmt.Errorf("msgid length[%d] invalid.", len(msgId))
	}

	buf, err := hex.DecodeString(msgId)
	if err != nil {
		return "", 0, err
	}

	ip := bytesToIPv4String(buf[0:4])
	port := binary.BigEndian.Uint32(buf[4:8])
	offset := binary.BigEndian.Uint64(buf[8:16])

	return fmt.Sprintf("%s:%d", ip, port), offset, nil
}
