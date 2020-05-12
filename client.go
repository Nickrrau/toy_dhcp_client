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
	DHCP_CLIENT_FATAL
)

type DHCPClient struct {
	state    ClientState
	iface    net.Interface
	ip       net.IP
	serverIP net.IP
	xid      []byte
	ops      []DHCPOption
}

func NewClient(iface net.Interface, ops []DHCPOption) *DHCPClient {
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

func (cl *DHCPClient) listen() (DHCPMsg, error) {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{Port: 68}) //TODO: Set Timeout
	defer listener.Close()
	fail(err)
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
			if MessageType(status[0]) == BOOT_REPLY && bytes.Compare(status[4:8], cl.xid) == 0 { //TODO: Move this check into parseMsg so that we can get rid of the duped code
				return cl.parseMsg(status)
			} else {
				return DHCPMsg{}, errors.New("Incorrect Message Type, Ignoring") //TODO: Make this more informative
			}
		case DHCP_CLIENT_REQUESTING:
			fmt.Println(status[0])
			if MessageType(status[0]) == BOOT_REPLY && bytes.Compare(status[4:8], cl.xid) == 0 { //TODO: Move this check into parseMsg so that we can get rid of the duped code
				return cl.parseMsg(status)
			} else {
				return DHCPMsg{}, errors.New("Incorrect Message Type, Ignoring") //TODO: Make this more informative
			}
		default:
			continue
		}

	}
}

func (cl *DHCPClient) discover() {
	conn, err := net.Dial("udp", "255.255.255.255:67")
	defer conn.Close()
	if err != nil {
		fmt.Printf("Closing Connecting: %v\n", err)
		os.Exit(1)
	}

	msg := NewDiscoverMsg(cl.iface.HardwareAddr, cl.ops)
	err = msg.WriteToConn(conn)
	if err != nil {
		fmt.Printf("Closing Connecting: %v\n", err)
		os.Exit(1)
	}

	cl.state = DHCP_CLIENT_DISCOVERING

}

func (cl *DHCPClient) parseMsg(data []byte) (DHCPMsg, error) {

	//TODO: Replace this with a method on DHCPMsg
	var msg DHCPMsg
	var err error

	msg.MsgType = MessageType(data[0])
	msg.HardwareType = HardwareType(data[1])
	msg.HardwareLength = data[2]
	msg.Hops = data[3]
	msg.XID = data[4:8]
	msg.ElapsedTime = data[8:10]
	msg.Flags = data[10:12]
	msg.ClientAddr = data[12:16]
	msg.YourAddr = data[16:20]
	msg.ServerAddr = data[20:24]
	msg.GatewayAddr = data[24:28]
	msg.ClientHardwareAddr = data[28:44]
	msg.ServerName = data[44:108]
	msg.File = data[108:236]
	msg.Magic = data[236:240]
	msg.RawBody = data[:240]

	msg.Options, err = OptionsFromBytes(data[240:])
	if err != nil {
		return msg, errors.New(fmt.Sprintf("Malformed Message:\n\tError = %v\n", err))
	}
	msg.RawOptions = data[240:]

	return msg, nil
}

func (cl *DHCPClient) reply() {

}

func (cl *DHCPClient) acceptAck() {

}

func (cl *DHCPClient) Run() {
	conn, err := net.Dial("udp", "255.255.255.255:67")
	defer conn.Close()
	if err != nil {
		fmt.Printf("Closing Connecting: %v\n", err)
		os.Exit(1)
	}

	msg := NewDiscoverMsg(cl.iface.HardwareAddr, cl.ops)
	err = msg.WriteToConn(conn)
	if err != nil {
		fmt.Printf("Closing Connecting: %v\n", err)
		os.Exit(1)
	}

	cl.state = DHCP_CLIENT_DISCOVERING

	for parsed, err := cl.listen(); ; {
		if err == nil {
			fmt.Println(parsed.String())
			break
		} //TODO: Wait before Retry
		//TODO: Fail on N Retries
	}
}

func (cl *DHCPClient) State() ClientState {
	return cl.state
}
