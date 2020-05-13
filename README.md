# Simple Toy DHCP Client

Simple client that prints out a list of interfaces to use, once one is selected it proceeds to broadcast a discover message and waits for a reply from a local dhcp server.
Handles the offer and makes a request, looking for a final ack before exiting.

Wireshark pcap is included.

# Things Left TODO
	- Add more code documentation
	- Add some basic test cases
	- Retries
	- Timeouts

# Notice
Purely a fun project and not meant to be a serious attempt at implementing a DHCP client.
