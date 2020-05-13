# Simple Toy DHCP Client

Simple client that prints out a list of interfaces to use, once one is selected it proceeds to broadcast a discover message and waits for a reply from a local dhcp server.
Handles the offer and makes a request, looking for a final ack before exiting.

Wireshark pcap is included.

# Example output from program
```
=== Ack -> Client ===
Message type: 2
Hardware Type: 1
Hardware Length: 6
XID: [23 23 43 23]
Client IP: [0 0 0 0]
Your IP: [192 168 0 5]
Server IP: [0 0 0 0]
Gateway IP: [0 0 0 0]
Client Hardware Address: [212 59 4 32 53 210 0 0 0 0 0 0 0 0 0 0]
Options: Option Name: Message Type
Code: 53
Len: 1
Data: [5]
Option Name: Server ID
Code: 54
Len: 4
Data: [192 168 0 1]
Option Name: Lease Time
Code: 51
Len: 4
Data: [0 1 81 128]
Option Name: Subnet Mask
Code: 1
Len: 4
Data: [255 255 255 0]
Option Name: Router
Code: 3
Len: 4
Data: [192 168 0 1]
Option Name: DNS Server(s)
Code: 6
Len: 4
Data: [192 168 0 1]
Option Name: END
Code: 255
Len: 0
Data: []
```

# Example output from Wireshark
```
Dynamic Host Configuration Protocol (ACK)
    Message type: Boot Reply (2)
    Hardware type: Ethernet (0x01)
    Hardware address length: 6
    Hops: 0
    Transaction ID: 0x17172b17
    Seconds elapsed: 0
    Bootp flags: 0x0000 (Unicast)
    Client IP address: 0.0.0.0 (0.0.0.0)
    Your (client) IP address: 192.168.0.5 (192.168.0.5)
    Next server IP address: 0.0.0.0 (0.0.0.0)
    Relay agent IP address: 0.0.0.0 (0.0.0.0)
    Client MAC address: Babylon.local (d4:3b:04:20:35:d2)
    Client hardware address padding: 00000000000000000000
    Server host name not given
    Boot file name not given
    Magic cookie: DHCP
    Option: (53) DHCP Message Type (ACK)
        Length: 1
        DHCP: ACK (5)
    Option: (54) DHCP Server Identifier (192.168.0.1)
        Length: 4
        DHCP Server Identifier: 192.168.0.1 (192.168.0.1)
    Option: (51) IP Address Lease Time
        Length: 4
        IP Address Lease Time: (86400s) 1 day
    Option: (1) Subnet Mask (255.255.255.0)
        Length: 4
        Subnet Mask: 255.255.255.0
    Option: (3) Router
        Length: 4
        Router: 192.168.0.1 (192.168.0.1)
    Option: (6) Domain Name Server
        Length: 4
        Domain Name Server: 192.168.0.1 (192.168.0.1)
    Option: (255) End
        Option End: 255
```

# Things Left TODO
	- Add more code documentation
	- Add some basic test cases

# Notice
Purely a fun project and not meant to be a serious attempt at implementing a DHCP client.
