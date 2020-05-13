package message

type MessageType byte

const (
	BOOT_REQEUST MessageType = iota + 1
	BOOT_REPLY
)
