package toy_dhcp_client

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"time"
	"toy_dhcp_client/message"
)

//Client.listen() listens on port 68 for broadcasted messages, messages are checked for correct message type and xid before commiting to parsing the full bytes.
func (cl *Client) listen() (*message.DHCPMsg, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: dhcpPortReceive})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	reader := bufio.NewReader(conn)
	buff := make([]byte, reader.Size())
	_, err = reader.Read(buff)
	if err != nil {
		return nil, err
	}
	switch cl.state {
	case DHCP_CLIENT_DISCOVERING:
		if message.MessageType(buff[0]) == message.BOOT_REPLY && bytes.Compare(buff[4:8], cl.xid) == 0 { //TODO: Move this check into parseMsg so that we can get rid of the duped code
			return message.BytesToDHCPMsg(buff)
		} else {
			return nil, errors.New("Incorrect Message Type, Ignoring") //TODO: Make this more informative
		}
	case DHCP_CLIENT_REQUESTING:
		if message.MessageType(buff[0]) == message.BOOT_REPLY && bytes.Compare(buff[4:8], cl.xid) == 0 { //TODO: Move this check into parseMsg so that we can get rid of the duped code
			return message.BytesToDHCPMsg(buff)
		} else {
			return nil, errors.New("Incorrect Message Type, Ignoring") //TODO: Make this more informative
		}
	default:
		return nil, nil
	}
}
