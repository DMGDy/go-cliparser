package mqtt_client

import (
	"fmt"
	"net"
	"errors"
)

type responsePacket struct {
	message_type int
	topic string
	payload string
}

const (
	SUBSCRIBE_TOPIC = "@/Panel/Partition_/#"
	PORT = "1883"
	//IP = "localhost:"
	IP = "localhost:"

	CONNECT    = 1
	CONNACK    = 2
	PUBLISH    = 3
	PUBACK     = 4
	SUBSCRIBE  = 8
	SUBACK     = 9
	DISCONNECT = 14
	QOS        = 2
	// dont need DUP or RET flags

	KEEPALIVE  = 30

	CLIENT_ID  = "cli_client"
	MAX_PAYLOAD_LEN = 5096
)

var conn net.Conn

func InitClient() error {

	var err error
	//var packet MqttPacket
	conn, err = net.Dial("tcp", IP+PORT)

	if err != nil {
		return err
	}

	sendConnect(conn)

	response := make([]byte, 24)

	_, err = conn.Read(response)
	if err != nil {
		return err
	}
	fmt.Println(response)
	decode(response)


	return nil
}

func CloseClient() {
	conn.Close()
}


func sendConnect(conn net.Conn) {
	// fixed header
	fixedHeader := []byte{CONNECT << 4, 0} // Length will be updated later

	protocolName := []byte{0, 4, 'M', 'Q', 'T', 'T'}
	// Variable header
	

	protocolLevel := []byte{4}

	connectFlags := []byte{QOS}

	keepAlive := []byte{0, KEEPALIVE}

	clientIDLen := len(CLIENT_ID)
	clientIDLenBytes := []byte{byte(clientIDLen >> 8), byte(clientIDLen)}
	clientIDBytes := []byte(CLIENT_ID)

	remainingLength := len(protocolName) + len(protocolLevel) + len(connectFlags) + len(keepAlive) + len(clientIDLenBytes) + len(clientIDBytes)


	fixedHeader[1] = byte(remainingLength)

	packet := append(fixedHeader, protocolName...)
	
	packet = append(packet, protocolLevel...)
	packet = append(packet, connectFlags...)
	packet = append(packet, keepAlive...)

	
	packet = append(packet, clientIDLenBytes...)
	packet = append(packet, clientIDBytes...)
	
	

	_, err := conn.Write(packet)
	if err != nil {
		fmt.Println("Error sending CONNECT packet:", err)
	}
}

func subscribe(conn net.Conn) {
	fixed_header := []byte{0x82}

	packet_id := []byte{0,1}

	topic_len := len(SUBSCRIBE_TOPIC)
	topic_len_bytes := []byte{byte(topic_len>>8), byte(topic_len)}
	topic_bytes := []byte(SUBSCRIBE_TOPIC)

	qos_bytes := []byte{2}


	remaining_len := len(packet_id) + len(topic_len_bytes) + len(topic_bytes) + len(qos_bytes)
	fixed_header = append(fixed_header, byte(remaining_len))

	packet := append(fixed_header, packet_id...)
	packet = append(packet, topic_len_bytes...)
	packet = append(packet, topic_bytes...)
	packet = append(packet, qos_bytes...)

	_, err := conn.Write(packet)
	if err != nil {
		fmt.Println("Error sending SUBSCRIBE packet,", err)
		return
	}
}

func publish(conn net.Conn, topic string, payload string) {
	fixed_header := []byte{PUBLISH<< 4, 0}

	topic_len := len(topic)
	topic_len_bytes := []byte{byte(topic_len>>8), byte(topic_len)}
	topic_bytes := []byte(topic)

	payload_bytes := []byte(payload)

	remaining_len := len(topic_len_bytes) + len(topic_bytes) + len(payload_bytes)
	fixed_header[1] = byte(remaining_len)

	packet := append(fixed_header, topic_len_bytes...)
	packet = append(packet, topic_bytes...)
	packet = append(packet, payload_bytes...)

	_, err := conn.Write(packet)
	if err != nil {
		fmt.Println("Error sending PUBLISH packet,", err)
		return
	}
}

/* deconstruct packet bytes
	determine what kind of message (CONNACK, SUBACK, PUBLISH)
	get message len (end of control header OR end of variable header)
*/
func decode(response []byte) responsePacket {
	// first byte is control header
	control_header := response[0]
	// upper 4 bytes is message type
	msg_type := control_header>>4
	// need to only know about CONNACK and SUBACK (maybe PUBACK?)

	switch(msg_type) {
	case(CONNACK):
		return responsePacket {
			message_type: CONNACK,
		}
	case(SUBACK): 
		return responsePacket {
			message_type: SUBACK,
		}
	/*
	case(PUBACK): 
		fmt.Println("PUBACK received")
	*/
	case(PUBLISH): 
		topic_info_len := response[2]
		// if its 2, perfect
		var topic_len = 0
		if topic_info_len == 2 {
			topic_len = int(response[4])
		}
		topic := string(response[4:4+topic_len])
		payload := string(response[4+topic_len:])
		return responsePacket {
			message_type: PUBLISH,
			topic: topic,
			payload: payload,
		}

	default:
		return responsePacket{}
	}

}

func SendCommand(topic string, payload string)(string, error) {

	subscribe(conn)

	sub_response := make([]byte, 8)

	_, err := conn.Read(sub_response)
	if err != nil {
		return "", errors.New("could not Subscribe")
	}
	fmt.Println(sub_response)
	decode(sub_response)


	publish(conn, topic, payload)
	pub_response := make([]byte, MAX_PAYLOAD_LEN)
	_, err = conn.Read(pub_response)
	if err != nil {
		return "", err
	}

	decoded_response := decode(pub_response)

	return decoded_response.payload, nil
}
