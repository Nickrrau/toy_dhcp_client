package toy_dhcp_client

import (
	"errors"
	"fmt"
	"net"
	"time"
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
	state     ClientState
	iface     net.Interface
	retry     int
	ip        net.IP
	leaseTime time.Duration
	sip       net.IP
	xid       []byte
	ops       []message.Option
}

// NewClient() returns a pointer to a new DHCP Client initialized with the Interface, Transaction ID, and Options provided
func NewClient(iface net.Interface, xid []byte, ops []message.Option) *Client {
	return &Client{
		iface: iface,
		state: DHCP_CLIENT_UNINITIALIZED,
		ops:   ops,
		xid:   xid,
		retry: 5,
	}
}

// Client.State() returns the current/last state the client reported.
func (cl *Client) State() ClientState {
	return cl.state
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

func (cl *Client) acceptAck(msg *message.DHCPMsg) error {
	cl.state = DHCP_CLIENT_ACCEPTACK
	msgType := message.FindOption(message.OPTION_MSG_TYPE, msg.Options)
	if msgType == nil {
		return errors.New("Malformed Response from Server: No DHCP Message Type Option")
	} else if message.DHCPMessageType(msgType.Data[0]) != message.MSGTYPE_ACK {
		if message.DHCPMessageType(msgType.Data[0]) == message.MSGTYPE_NACK {
			return errors.New("Server NACK Response")
		}
		return errors.New("Wrong Response from Server: Incorrect Message Type")
	}
	return nil
}

//Client.Run() starts the client and begins attempting to connect with a local DHCP server.
func (cl *Client) Run() error {

	cl.StatePrint()
	err := cl.discover(cl.ops)
	if err != nil {
		return errors.New(fmt.Sprintf("Discovery Failed, closing Client\nErr:%v", err))
	}
	cl.StatePrint()

	offerMsg := &message.DHCPMsg{}
	for counter := 0; ; counter++ {
		offerMsg, err = cl.listen()
		if err == nil {
			break
		}
		if err != nil && counter >= cl.retry {
			return err
		}
		fmt.Println("Failed, Retrying")
		time.Sleep(time.Second)
	}
	cl.StatePrint()

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
	for counter := 0; ; counter++ {
		ackMsg, err = cl.listen()
		if err == nil {
			break
		}
		if err != nil && counter >= cl.retry {
			return err
		}
		fmt.Println("Failed, Retrying")
		time.Sleep(time.Second)
	}
	cl.StatePrint()
	cl.acceptAck(ackMsg)
	cl.StatePrint()
	cl.state = DHCP_CLIENT_ACKED
	cl.StatePrint()

	fmt.Println(ackMsg)

	return nil

}

func (cl *Client) StatePrint() {
	switch cl.state {
	case DHCP_CLIENT_UNINITIALIZED:
		fmt.Println("Client Is Uninitialized")
	case DHCP_CLIENT_INITIALIZING:
		fmt.Println("Client Is Initializing")
	case DHCP_CLIENT_INITIALIZED:
		fmt.Println("Client Is Initialized")
	case DHCP_CLIENT_DISCOVERING:
		fmt.Println("Client Is Broadcasting a Discover Message")
	case DHCP_CLIENT_OFFERED:
		fmt.Println("Client has recieved an Offer from server")
	case DHCP_CLIENT_REQUESTING:
		fmt.Println("Client Is Requesting an IP")
	case DHCP_CLIENT_ACCEPTACK:
		fmt.Println("Client has recieved an ACK from server")
	case DHCP_CLIENT_ACKED:
		fmt.Println("Client has approved an ACK from server")
	}
}
