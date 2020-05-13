package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	toy "toy_dhcp_client"
	"toy_dhcp_client/message"
)

var (
	DHCP_XID = []byte{23, 23, 43, 23}
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

	ops := []message.DHCPOption{
		message.DHCP_MSG_TYPE_DISCOVER,
		message.DHCP_MAX_MSG_SIZE,
		message.DHCP_PARAM_REQ_LIST,
		message.DHCP_CLIENT_ID,
		message.DHCP_END,
	}

	client := toy.NewClient(ints[index], DHCP_XID, ops)

	client.Run()
}

func fail(err error) {
	if err != nil {
		fmt.Printf("Failed to list Interfaces: %v", err)
		os.Exit(1)
	}
}
