package main

import (
	"./message"
	"fmt"
	"net"
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

type DHCPClient struct {
	state    ClientState
	iface    net.Interface
	ip       net.IP
	serverIP net.IP
	xid      []byte
	ops      []message.DHCPOption
}

func NewClient(iface net.Interface, ops []message.DHCPOption) *DHCPClient {
	return &DHCPClient{
		iface: iface,
		state: DHCP_CLIENT_UNINITIALIZED,
		ops:   ops,
		xid:   DHCP_XID,
	}
}

func (cl *DHCPClient) IP() net.IP {
	//TODO: Copy
	return cl.ip //Possible race?
}
func (cl *DHCPClient) SIP() net.IP {
	//TODO: Copy
	return cl.serverIP //Possible race?
}
func (cl *DHCPClient) XID() []byte {
	var b []byte
	copy(b, cl.xid[:])
	return b
}

//Connection related
func (cl *DHCPClient) reply() {

}

func (cl *DHCPClient) acceptAck() {

}

func (cl *DHCPClient) Run() {
	cl.discover()

	for parsed, err := cl.listen(); ; {
		if err == nil {
			fmt.Println(parsed.String())
			break
		} //TODO: Wait before Retry
		//TODO: Fail on N Retries
	}
}
