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

	XID [4]byte

	ElapsedTime [2]byte

	Flags [2]byte

	//TODO: Use Go IP struct and preform the conversions to byte arrays later
	ClientAddr  [4]byte
	YourAddr    [4]byte
	ServerAddr  [4]byte
	GatewayAddr [4]byte

	ClientHardwareAddr [16]byte

	ServerName [64]byte
	File       [128]byte

	Magic [4]byte

	Options []DHCPOption
}

func NewDiscoverMsg(hwaddr []byte, ops []DHCPOption) DHCPMsg {
	var mac = [16]byte{}
	copy(mac[:], hwaddr)
	fmt.Println(hwaddr)
	return DHCPMsg{
		MsgType:            BOOT_REQEUST,
		HardwareType:       ETHERNET,
		HardwareLength:     6,
		Hops:               0,
		XID:                DHCP_XID,
		Flags:              [2]byte{0, 0},
		ClientHardwareAddr: mac,
		Magic:              DHCP_MAGIC,
		Options:            ops,
	}
}

func (msg *DHCPMsg) WriteToConn(conn net.Conn) error {
	writer := bufio.NewWriter(conn)

	// err := binary.Write(writer, binary.BigEndian, msg)

	//TODO: Can't use the binary encoder package with slice for variable length options, manually writing fields solves the issue but at the cost of not being very friendly code...
	err := writer.WriteByte(byte(msg.MsgType))
	err = writer.WriteByte(byte(msg.HardwareType))
	err = writer.WriteByte(msg.HardwareLength)
	err = writer.WriteByte(msg.Hops)

	_, err = writer.Write(msg.XID[:]) // Converting fixed array into slice? Probably not good
	_, err = writer.Write(msg.ElapsedTime[:])
	_, err = writer.Write(msg.Flags[:])
	_, err = writer.Write(msg.ClientAddr[:])
	_, err = writer.Write(msg.YourAddr[:])
	_, err = writer.Write(msg.ServerAddr[:])
	_, err = writer.Write(msg.GatewayAddr[:])
	_, err = writer.Write(msg.ClientHardwareAddr[:])
	_, err = writer.Write(msg.ServerName[:])
	_, err = writer.Write(msg.File[:])
	_, err = writer.Write(msg.Magic[:])

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
