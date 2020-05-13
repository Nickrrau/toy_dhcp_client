package toy_dhcp_client

import (
	"errors"
	"fmt"
	"net"
	"toy_dhcp_client/message"
)

//ClientState represents the current state of the DHCP Client
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

// Client maintains the needed state for communicating with a remote DHCP server
type Client struct {
	state ClientState
	iface net.Interface
	ip    net.IP
	sip   net.IP
	xid   []byte
	ops   []message.Option
}

// NewClient() returns a pointer to a new DHCP Client initialized with the Interface, Transaction ID, and Options provided
func NewClient(iface net.Interface, xid []byte, ops []message.Option) *Client {
	return &Client{
		iface: iface,
		state: DHCP_CLIENT_UNINITIALIZED,
		ops:   ops,
		xid:   xid,
	}
}

// Client.IP() returns the IP the client has configured.
func (cl *Client) IP() net.IP {
	//TODO: Copy
	return cl.ip
}

// Client.SIP() returns the IP of the server which has responded first to the client
func (cl *Client) SIP() net.IP {
	//TODO: Copy
	return cl.sip
}

// Client.XID() returns the transaction ID for the client. This is fixed when the client is created and is persisted for as long as the client is.
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

	err := cl.discover(cl.ops)
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

	serverIP := message.FindOption(message.OPTION_SERVER_ID, offerMsg.Options)
	if msgType == nil {
		return errors.New("Malformed Response from Server: No DHCP Server ID")
	}
	cl.ip = offerMsg.YourAddr
	cl.sip = serverIP.Data

	reqOptions := []message.Option{*serverIP, message.NewOption(message.OPTION_REQUEST_IP, cl.ip)}
	reqOptions = append(reqOptions, message.DefaultRequestOps...)
	cl.request(reqOptions)

	ackMsg := &message.DHCPMsg{}
	for ackMsg, err = cl.listen(); ; {
		if err == nil {
			break
		}
		if err != nil {
			return err //TODO: Fail on N Retries
		}
		//TODO: Wait before Retry
	}

	msgType = message.FindOption(message.OPTION_MSG_TYPE, offerMsg.Options)
	if msgType == nil {
		return errors.New("Malformed Response from Server: No DHCP Message Type Option")
	} else if message.DHCPMessageType(msgType.Data[0]) != message.MSGTYPE_ACK {
		if message.DHCPMessageType(msgType.Data[0]) == message.MSGTYPE_NACK {
			return errors.New("Server NACK Response")
		}
		return errors.New("Wrong Response from Server: Incorrect Message Type")
	}

	cl.state = DHCP_CLIENT_ACCEPTACK
	fmt.Println(ackMsg)

	cl.state = DHCP_CLIENT_ACKED

	return nil

}
