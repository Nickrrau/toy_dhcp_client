package main

import (
	"bufio"
	"fmt"
	"net"
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

	Options []DHCPOption

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
		msg.Options,
	)
}

func NewDiscoverMsg(hwaddr []byte, ops []DHCPOption) *DHCPMsg {
	msg := NewDHCPMsg()
	msg.MsgType = BOOT_REQEUST
	msg.HardwareType = ETHERNET
	msg.HardwareLength = 6
	msg.XID = DHCP_XID
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
