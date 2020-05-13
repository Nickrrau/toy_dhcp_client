package message

import (
	"bufio"
	"fmt"
	"net"
)

var (
	DHCP_MAGIC = []byte{99, 130, 83, 99}
)

//TODO: Unexport these fields
type DHCPMsg struct {
	MsgType        MessageType
	HardwareType   HardwareType
	HardwareLength byte
	Hops           byte

	XID []byte

	ElapsedTime []byte

	Flags []byte

	//TODO: Use Go IP struct and preform the conversions to byte arrays later
	ClientAddr  []byte
	YourAddr    []byte
	ServerAddr  []byte
	GatewayAddr []byte

	ClientHardwareAddr []byte

	ServerName []byte
	File       []byte

	Magic []byte

	Options []Option

	RawBody    []byte
	RawOptions []byte
}

func NewDHCPMsg() *DHCPMsg {
	return &DHCPMsg{
		XID:                make([]byte, 4),
		ElapsedTime:        make([]byte, 2),
		Flags:              make([]byte, 2),
		ClientAddr:         make([]byte, 4),
		YourAddr:           make([]byte, 4),
		ServerAddr:         make([]byte, 4),
		GatewayAddr:        make([]byte, 4),
		ClientHardwareAddr: make([]byte, 16),
		ServerName:         make([]byte, 64),
		File:               make([]byte, 128),
		Magic:              make([]byte, 4),
	}
}

func (msg *DHCPMsg) String() string {
	return fmt.Sprintf("Message type: %v\nHardware Type: %v\nHardware Length: %v\nXID: %v\nClient IP: %v\nYour IP: %v\nServer IP: %v\nGateway IP: %v\nClient Hardware Address: %v\nOptions: %v\n",
		msg.MsgType,
		msg.HardwareType,
		msg.HardwareLength,
		msg.XID,
		msg.ClientAddr,
		msg.YourAddr,
		msg.ServerAddr,
		msg.GatewayAddr,
		msg.ClientHardwareAddr,
		PrettyOptions(msg.Options),
	)
}

// NewDiscoverMsg() is a helper function that sets up a DHCPMsg struct for a Discover Message being broadcasted
func NewBroadcastMsg(xid, hwaddr []byte, ops []Option) *DHCPMsg {
	msg := NewDHCPMsg()
	msg.MsgType = BOOT_REQEUST
	msg.HardwareType = ETHERNET
	msg.HardwareLength = 6
	msg.XID = xid
	msg.Magic = DHCP_MAGIC
	msg.Options = ops

	copy(msg.ClientHardwareAddr, hwaddr)

	return msg
}

func (msg *DHCPMsg) WriteToConn(conn net.Conn) error {
	writer := bufio.NewWriter(conn)

	err := writer.WriteByte(byte(msg.MsgType))
	err = writer.WriteByte(byte(msg.HardwareType))
	err = writer.WriteByte(msg.HardwareLength)
	err = writer.WriteByte(msg.Hops)

	_, err = writer.Write(msg.XID)
	_, err = writer.Write(msg.ElapsedTime)
	_, err = writer.Write(msg.Flags)
	_, err = writer.Write(msg.ClientAddr)
	_, err = writer.Write(msg.YourAddr)
	_, err = writer.Write(msg.ServerAddr)
	_, err = writer.Write(msg.GatewayAddr)
	_, err = writer.Write(msg.ClientHardwareAddr)
	_, err = writer.Write(msg.ServerName)
	_, err = writer.Write(msg.File)
	_, err = writer.Write(msg.Magic)

	if err != nil {
		writer.Reset(conn)
		return err
	}
	for i := range msg.Options {
		if _, err := writer.Write(msg.Options[i].Bytes()); err != nil {
			return err
		}
	}

	writer.Flush()
	return err
}

func BytesToDHCPMsg(data []byte) (*DHCPMsg, error) {
	msg := &DHCPMsg{}
	var err error

	msg.MsgType = MessageType(data[0])
	msg.HardwareType = HardwareType(data[1])
	msg.HardwareLength = data[2]
	msg.Hops = data[3]
	msg.XID = data[4:8]
	msg.ElapsedTime = data[8:10]
	msg.Flags = data[10:12]
	msg.ClientAddr = data[12:16]
	msg.YourAddr = data[16:20]
	msg.ServerAddr = data[20:24]
	msg.GatewayAddr = data[24:28]
	msg.ClientHardwareAddr = data[28:44]
	msg.ServerName = data[44:108]
	msg.File = data[108:236]
	msg.Magic = data[236:240]
	msg.RawBody = data[:240]

	msg.Options, err = OptionsFromBytes(data[240:])
	if err != nil {
		return msg, err
	}
	msg.RawOptions = data[240:]

	return msg, nil
}
