package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
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
)

type DHCPClient struct {
	status   chan ClientState
	iface    net.Interface
	IP       net.IP
	ServerIP net.IP
	XID      [4]byte
}

func NewClient(iface net.Interface) *DHCPClient {
	return &DHCPClient{
		iface:  iface,
		status: make(chan ClientState),
	}
}

func (cl *DHCPClient) listener() {

}

func (cl *DHCPClient) discover() {

}

func (cl *DHCPClient) parseMsg(data []byte) (DHCPMsg, error) {
	var msg DHCPMsg
	var err error

	reader := bufio.NewReader(bytes.NewReader(data))
	body := make([]byte)
}

func (cl *DHCPClient) reply() {

}

func (cl *DHCPClient) acceptAck() {

}

func (cl *DHCPClient) Init() {
	cl.status <- DHCP_CLIENT_UNINITIALIZED
	cl.status <- DHCP_CLIENT_INITIALIZED
	conn, err := net.Dial("udp", "255.255.255.255:67")
	defer conn.Close()
	if err != nil {
		fmt.Printf("Closing Connecting: %v\n", err)
		os.Exit(1)
	}

	ops := []DHCPOption{
		DHCP_MSG_TYPE_DISCOVER,
		DHCP_MAX_MSG_SIZE,
		DHCP_PARAM_REQ_LIST,
		DHCP_CLIENT_ID,
		DHCP_END,
	}

	msg := NewDiscoverMsg(cl.iface.HardwareAddr, ops)
	err = msg.WriteToConn(conn)
	if err != nil {
		fmt.Printf("Closing Connecting: %v\n", err)
		os.Exit(1)
	}
	cl.status <- DHCP_CLIENT_DISCOVERING

	listener, err := net.ListenUDP("udp", &net.UDPAddr{Port: 68})
	fail(err)

	status, err := bufio.NewReader(listener).ReadBytes(255)
	cl.status <- DHCP_CLIENT_OFFERED
	fmt.Println(status)
	cl.status <- DHCP_CLIENT_REQUESTING
	cl.status <- DHCP_CLIENT_ACCEPTACK
	cl.status <- DHCP_CLIENT_ACKED
}

func (cl *DHCPClient) Status() {
	select {
	case state := <-cl.status:
		switch state {
		case DHCP_CLIENT_UNINITIALIZED:
			fmt.Println("Client is Uninitialized...")
		case DHCP_CLIENT_INITIALIZED:
			fmt.Println("Initializing Client...")
		case DHCP_CLIENT_DISCOVERING:
			fmt.Println("Discovering DHCP Server...")
		case DHCP_CLIENT_OFFERED:
			fmt.Println("Offer from Server (<IP>, <SERVERIP>)...")
		case DHCP_CLIENT_REQUESTING:
			fmt.Println("Replying to Server with Request for <IP>...")
		case DHCP_CLIENT_ACCEPTACK:
			fmt.Println("Waiting for ACK from Server...")
		case DHCP_CLIENT_ACKED:
			fmt.Println("Server ACKED...")
		}
	default:
		//fmt.Println("Waiting...")
	}
}
