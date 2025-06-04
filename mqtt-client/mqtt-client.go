package mqtt_client

import (
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
	PORT = ":1883"
	//IP = "localhost:"

	CONNECT    = 1
	CONNACK    = 2
	PUBLISH    = 3
	PUBACK     = 4
	SUBSCRIBE  = 8
	SUBACK     = 9
	DISCONNECT = 14
	QOS        = 2
	// dont need DUP or RET flags

	KEEPALIVE  = 5

	CLIENT_ID  = "cli_client"
	MAX_PAYLOAD_LEN = 5096
)

var IP = "localhost"
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

	n, err := conn.Read(response)
	if err != nil {
		return err
	}
	decode(response, n)


	return nil
}

func CloseClient() {
	conn.Close()
}


func sendConnect(conn net.Conn) error {
	// fixed header
	fixed_header := []byte{CONNECT << 4, 0} // Length will be updated later

	protocol_name := []byte{0, 4, 'M', 'Q', 'T', 'T'}
	// Variable header
	

	protocol_level := []byte{4}

	conn_flags := []byte{QOS}

	keep_alive := []byte{0, KEEPALIVE}

	client_id_len := len(CLIENT_ID)
	client_id_len_bytes := []byte{byte(client_id_len >> 8), byte(client_id_len)}
	client_id_bytes := []byte(CLIENT_ID)

	remaining_len := len(protocol_name) + len(protocol_level) + len(conn_flags) + len(keep_alive) + len(client_id_len_bytes) + len(client_id_bytes)


	fixed_header[1] = byte(remaining_len)

	packet := append(fixed_header, protocol_name...)
	
	packet = append(packet, protocol_level...)
	packet = append(packet, conn_flags...)
	packet = append(packet, keep_alive...)

	
	packet = append(packet, client_id_len_bytes...)
	packet = append(packet, client_id_bytes...)
	
	

	_, err := conn.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func subscribe(sub_topic string) error {
	fixed_header := []byte{0x82}

	packet_id := []byte{0,1}

	topic_len := len(sub_topic)
	topic_len_bytes := []byte{byte(topic_len>>8), byte(topic_len)}
	topic_bytes := []byte(sub_topic)

	qos_bytes := []byte{2}


	remaining_len := len(packet_id) + len(topic_len_bytes) + len(topic_bytes) + len(qos_bytes)
	fixed_header = append(fixed_header, byte(remaining_len))

	packet := append(fixed_header, packet_id...)
	packet = append(packet, topic_len_bytes...)
	packet = append(packet, topic_bytes...)
	packet = append(packet, qos_bytes...)

	_, err := conn.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func publish(topic string, payload string) error {
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
		return err
	}
	return nil
}

/* deconstruct packet bytes
	determine what kind of message (CONNACK, SUBACK, PUBLISH)
	get message len (end of control header OR end of variable header)
	n being the amount of bytes read so we can determine whent to stop reading
*/
func decode(response []byte, n int) responsePacket {
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
	*/

	// only case when we care about topic and payload
	case(PUBLISH):
		// if its not zero, its a length
		var topic_len = 0
		var offset = 0
		for i := 2; i < 5; i++ {
			if int(response[i]) == 0 {
				offset = i+1
				topic_len = int(response[offset])
				break
			}
		}
		topic := string(response[offset+1:offset+1+topic_len])
		payload := string(response[offset+1+topic_len:n])
		return responsePacket {
			message_type: PUBLISH,
			topic: topic,
			payload: payload,
		}

	default:
		return responsePacket{}
	}

}

func SendCommand(sub_topic string, pub_topic string, payload string)(string, error) {

	err := subscribe(sub_topic)
	if err != nil {
		return "", errors.New("Error sending SUBSCRIBE:" + err.Error())
	}

	sub_response := make([]byte, 8)

	n, err := conn.Read(sub_response)
	if err != nil || n == 0{
		return "", errors.New("Error reading SUBACK"+err.Error())
	}
	decode(sub_response, n)


	err = publish(pub_topic, payload)
	if err != nil {
		return "", errors.New("Error sending PUBLISH with: "+err.Error())
	}
	pub_response := make([]byte, MAX_PAYLOAD_LEN)
	n, err = conn.Read(pub_response)
	if err != nil || n == 0{
		return "", err
	}

	decoded_response := decode(pub_response, n)

	return decoded_response.payload, nil
}
