package toy_dhcp_client

import (
	"fmt"
	"net"
	"toy_dhcp_client/message"
)

func (cl *Client) request(ops []message.Option) error {
	conn, err := net.Dial("udp", broadcastAddr)
	if err != nil {
		fmt.Printf("Closing Connection: %v\n", err)
		return err
	}
	defer conn.Close()

	msg := message.NewBroadcastMsg(cl.xid, cl.iface.HardwareAddr, ops)
	msg.ServerAddr = cl.SIP()

	err = msg.WriteToConn(conn)
	if err != nil {
		fmt.Printf("Closing Connection: %v\n", err)
		return err
	}

	cl.state = DHCP_CLIENT_REQUESTING

	return nil
}
