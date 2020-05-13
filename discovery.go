package toy_dhcp_client

import (
	"fmt"
	"net"
	"toy_dhcp_client/message"
)

const (
	broadcastAddr   = "255.255.255.255:67"
	dhcpPortReceive = 68
)

func (cl *Client) discover() error {
	conn, err := net.Dial("udp", broadcastAddr)
	if err != nil {
		fmt.Printf("Closing Connecting: %v\n", err)
		return err
	}
	defer conn.Close()

	msg := message.NewDiscoverMsg(cl.xid, cl.iface.HardwareAddr, cl.ops)
	err = msg.WriteToConn(conn)
	if err != nil {
		fmt.Printf("Closing Connecting: %v\n", err)
		return err
	}

	cl.state = DHCP_CLIENT_DISCOVERING
	return nil
}
