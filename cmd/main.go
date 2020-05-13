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
	if err != nil {
		fmt.Printf("Failed to list Interfaces: %v", err)
		os.Exit(1)
	}
	fmt.Print("Select Interface to Use:\n")
	for i := range ints {
		fmt.Printf("[%v] %v (%v)\n", i, ints[i].Name, ints[i].HardwareAddr)
	}
	fmt.Print("\n")

	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to list Interfaces: %v", err)
		os.Exit(2)
	}
	index, err := strconv.Atoi(string(input[0]))
	if err != nil {
		fmt.Printf("Failed to list Interfaces: %v", err)
		os.Exit(3)
	}

	client := toy.NewClient(ints[index], DHCP_XID, message.DefaultDiscoverOps)

	err = client.Run()
	if err != nil {
		fmt.Printf("Client Failed\nLast State: %v\nError: %v\n", client.State(), err)
		os.Exit(4)
	}
}
