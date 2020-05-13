package main

import (
	"./message"
	"fmt"
	"net"
)

func (cl *DHCPClient) discover() error {
	conn, err := net.Dial("udp", "255.255.255.255:67")
	defer conn.Close()
	if err != nil {
		fmt.Printf("Closing Connecting: %v\n", err)
		return err
	}

	msg := message.NewDiscoverMsg(cl.xid, cl.iface.HardwareAddr, cl.ops)
	err = msg.WriteToConn(conn)
	if err != nil {
		fmt.Printf("Closing Connecting: %v\n", err)
		return err
	}

	cl.state = DHCP_CLIENT_DISCOVERING
	return nil
}
