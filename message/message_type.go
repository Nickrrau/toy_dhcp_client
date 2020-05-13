package message

type MessageType byte

const (
	BOOT_REQEUST MessageType = iota + 1
	BOOT_REPLY
)

type DHCPMessageType byte

const (
	MSGTYPE_DISCOVER DHCPMessageType = iota + 1
	MSGTYPE_OFFER
	MSGTYPE_REQUEST
	//TODO: 4?
	MSGTYPE_ACK  DHCPMessageType = 5
	MSGTYPE_NACK DHCPMessageType = 6
)
