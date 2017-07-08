package rocketmq

import (
	"encoding/binary"
	"encoding/hex"
	"net"
	"strings"
	"time"
)

func int64ToBytes(value int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(value))
	return buf
}

func int32ToBytes(value int32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(value))
	return buf
}

func bytesToHexString(src []byte) string {
	//if src == nil || len(src) == 0 {
	//	return ""
	//}
	return strings.ToUpper(hex.EncodeToString(src))
}

// IPv4 address a.b.c.d. src is BigEndian buffer
func bytesToIPv4String(src []byte) string {
	return net.IPv4(src[0], src[1], src[2], src[3]).String()
}

// HashCode return stirng hashcode(equal java)
func HashCode(s string) int32 {
	var h int32
	for i := 0; i < len(s); i++ {
		h = 31*h + int32(s[i])
	}
	return h
}

// UnixNano return current time unix
func UnixNano() int64 {
	return time.Now().UnixNano()
}
