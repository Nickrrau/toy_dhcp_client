package toy_dhcp_client

import (
	"toy_dhcp_client/message"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
)

func (cl *DHCPClient) listen() (*message.DHCPMsg, error) {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{Port: 68}) //TODO: Set Timeout
	if err != nil {
		return nil,err
	}
	defer listener.Close()
	for {
		reader := bufio.NewReader(listener)
		status := make([]byte, reader.Size())
		_, err := reader.Read(status)
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch cl.state {
		case DHCP_CLIENT_DISCOVERING:
			if message.MessageType(status[0]) == message.BOOT_REPLY && bytes.Compare(status[4:8], cl.xid) == 0 { //TODO: Move this check into parseMsg so that we can get rid of the duped code
				return message.BytesToDHCPMsg(status)
			} else {
				return nil, errors.New("Incorrect Message Type, Ignoring") //TODO: Make this more informative
			}
		case DHCP_CLIENT_REQUESTING:
			fmt.Println(status[0])
			if message.MessageType(status[0]) == message.BOOT_REPLY && bytes.Compare(status[4:8], cl.xid) == 0 { //TODO: Move this check into parseMsg so that we can get rid of the duped code
				return message.BytesToDHCPMsg(status)
			} else {
				return nil, errors.New("Incorrect Message Type, Ignoring") //TODO: Make this more informative
			}
		default:
			continue
		}

	}
}
