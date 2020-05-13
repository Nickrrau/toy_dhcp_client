package toy_dhcp_client

import (
	"errors"
	"fmt"
	"net"
	"toy_dhcp_client/message"
)

type ClientState int

const (
	DHCP_CLIENT_UNINITIALIZED ClientState = iota
	DHCP_CLIENT_INITIALIZING
	DHCP_CLIENT_INITIALIZED
	DHCP_CLIENT_DISCOVERING
	DHCP_CLIENT_OFFERED
	DHCP_CLIENT_REQUESTING
	DHCP_CLIENT_ACCEPTACK
	DHCP_CLIENT_ACKED
	DHCP_CLIENT_FATAL
)

type Client struct {
	state    ClientState
	iface    net.Interface
	ip       net.IP
	serverIP net.IP
	xid      []byte
	ops      []message.Option
}

func NewClient(iface net.Interface, xid []byte, ops []message.Option) *Client {
	return &Client{
		iface: iface,
		state: DHCP_CLIENT_UNINITIALIZED,
		ops:   ops,
		xid:   xid,
	}
}

func (cl *Client) IP() net.IP {
	//TODO: Copy
	return cl.ip //Possible race?
}
func (cl *Client) SIP() net.IP {
	//TODO: Copy
	return cl.serverIP //Possible race?
}
func (cl *Client) XID() []byte {
	var b []byte
	copy(b, cl.xid[:])
	return b
}

//Connection related
func (cl *Client) reply() {

}

func (cl *Client) acceptAck() {

}

func (cl *Client) Run() error {
	fmt.Println("Sending Discovery Message")

	err := cl.discover()
	if err != nil {
		return errors.New(fmt.Sprintf("Discovery Failed, closing Client\nErr:%v", err))
	}

	offerMsg := &message.DHCPMsg{}
	for offerMsg, err = cl.listen(); ; {
		if err == nil {
			break
		}
		if err != nil {
			return err //TODO: Fail on N Retries
		}
		//TODO: Wait before Retry
	}

	//fmt.Println(offerMsg)

	msgType := message.FindOption(message.OPTION_MSG_TYPE, offerMsg.Options)
	if msgType == nil {
		return errors.New("Malformed Response from Server: No DHCP Message Type Option")
	} else if message.DHCPMessageType(msgType.Data[0]) != message.MSGTYPE_OFFER {
		return errors.New("Wrong Response from Server: Incorrect Message Type")
	}

	return nil

}
