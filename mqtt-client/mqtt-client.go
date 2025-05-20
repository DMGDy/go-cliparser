package mqtt-client

import (
	"net"
)

type Client interface {
	Subscribe()
	Publish()
	Set()
	Get()
}

type MqttPacket struct {
	ControlHeader byte
	ContentLength int
	Content []byte
}

type MqttClient struct {
	conn net.Conn
}

const (
	PARTITION1 = "@/Panel/Partition_/1"
	PORT = "1883"
	IP = "localhost:"

	CONNECT    = 0x00010000
	CONNACK    = 0x00100000
	DISCONNECT = 0x10000000
	QOS        = 0x00000010 // 2
)

func InitConnClient() (MqttClient, error) {
	conn, err := net.Dial("tcp", IP+PORT)
}


