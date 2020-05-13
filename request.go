package toy_dhcp_client

import (
	"fmt"
	"net"
	"time"
	"toy_dhcp_client/message"
)

// Client.request() broadcasts a OFFER message over the client's interface.
// The dial has a 10 second timeout.
// Once completed the Client state is changed to DHCP_CLIENT_REQUESTING
func (cl *Client) request(ops []message.Option) error {
	conn, err := net.DialTimeout("udp", broadcastAddr, 10*time.Second)
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
	fmt.Println("=== Request -> Server ===")
	fmt.Print(msg.String())
	return nil
}
