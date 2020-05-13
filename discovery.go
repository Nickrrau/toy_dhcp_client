package toy_dhcp_client

import (
	"fmt"
	"net"
	"time"
	"toy_dhcp_client/message"
)

const (
	broadcastAddr   = "255.255.255.255:67"
	dhcpPortReceive = 68
)

// Client.discovery() broadcasts a DISCOVERY message over the client's interface.
// The dial has a 10 second timeout.
// Once completed the Client state is changed to DHCP_CLIENT_DISCOVERING
func (cl *Client) discover(ops []message.Option) error {
	conn, err := net.DialTimeout("udp", broadcastAddr, 10*time.Second)
	if err != nil {
		//fmt.Printf("Closing Connection: %v\n", err)
		return err
	}
	defer conn.Close()

	msg := message.NewBroadcastMsg(cl.xid, cl.iface.HardwareAddr, ops)
	err = msg.WriteToConn(conn)
	if err != nil {
		//fmt.Printf("Closing Connection: %v\n", err)
		return err
	}

	cl.state = DHCP_CLIENT_DISCOVERING
	fmt.Println("=== Discover -> Server ===")
	fmt.Print(msg.String())
	return nil
}
