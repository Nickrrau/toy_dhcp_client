package main

var (
	DHCP_MSG_TYPE_DISCOVER = DHCPOption{
		Option: 53,
		Len:    1,
		Data:   []byte{1},
	}
	DHCP_MAX_MSG_SIZE = DHCPOption{
		Option: 57,
		Len:    2,
		Data:   []byte{0x05, 0xdc},
	}
	DHCP_PARAM_REQ_LIST = DHCPOption{
		Option: 55,
		Len:    5,
		Data:   []byte{1, 2, 3, 6, 69},
	}
	DHCP_CLIENT_ID = DHCPOption{
		Option: 61,
		Len:    5,
		Data:   []byte{1, byte('D'), byte('E'), byte('M'), byte('O')},
	}
	DHCP_END = DHCPOption{
		Option: 255,
	}
)

type DHCPOption struct {
	Option byte
	Len    byte
	Data   []byte
}

func (op *DHCPOption) Bytes() []byte {
	b := []byte{}
	b = append(b, op.Option, op.Len)
	b = append(b, op.Data...)
	return b
}
