package main

import (
	"bufio"
	// "encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {

	ints, err := net.Interfaces()
	fail(err)
	fmt.Print("Select Interface to Use:\n")
	for i := range ints {
		fmt.Printf("[%v] %v (%v)\n", i, ints[i].Name, ints[i].HardwareAddr)
	}
	fmt.Print("\n")

	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	fail(err)
	index, err := strconv.Atoi(string(input[0]))
	fail(err)

	// conn, err := net.Dial("udp", "255.255.255.255:67")
	// defer conn.Close()
	// if err != nil {
	// fmt.Printf("Closing Connecting: %v\n", err)
	// os.Exit(1)
	// }
	//
	// fmt.Println("Creating Msg")
	//
	// ops := []DHCPOption{
	// DHCP_MSG_TYPE_DISCOVER,
	// DHCP_MAX_MSG_SIZE,
	// DHCP_PARAM_REQ_LIST,
	// DHCP_CLIENT_ID,
	// DHCP_END,
	// }
	//
	// msg := NewDiscoverMsg(ints[index].HardwareAddr, ops)
	// fmt.Println("Writing to Conn")
	// err = msg.WriteToConn(conn)
	// if err != nil {
	// fmt.Printf("Closing Connecting: %v\n", err)
	// os.Exit(1)
	// }
	//
	// fmt.Println(net.Interfaces())
	//
	// listener, err := net.ListenUDP("udp", &net.UDPAddr{Port: 68})
	// fail(err)
	//
	// status, err := bufio.NewReader(listener).ReadBytes(255)
	// fmt.Println(status)

	client := NewClient(ints[index])
	go client.Init()

	for {
		client.Status()
	}

}

func fail(err error) {
	if err != nil {
		fmt.Printf("Failed to list Interfaces: %v", err)
		os.Exit(1)
	}
}
