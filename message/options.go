package message

type DHCPMsgOption byte

const (
	OPTION_SUBNET_MASK    DHCPMsgOption = 1
	OPTION_ROUTER         DHCPMsgOption = 3
	OPTION_DNS            DHCPMsgOption = 6
	OPTION_DOMAIN_NAME    DHCPMsgOption = 15
	OPTION_LEASE_TIME     DHCPMsgOption = 51
	OPTION_MSG_TYPE       DHCPMsgOption = 53
	OPTION_SERVER_ID      DHCPMsgOption = 54
	OPTION_MAX_MSG_SIZE   DHCPMsgOption = 57
	OPTION_PARAM_REQ_LIST DHCPMsgOption = 55
	OPTION_CLIENT_ID      DHCPMsgOption = 61
	OPTION_END            DHCPMsgOption = 255
)

var (
	DHCP_MSG_TYPE_DISCOVER = DHCPOption{
		Option: OPTION_MSG_TYPE,
		Len:    1,
		Data:   []byte{1},
	}
	DHCP_MAX_MSG_SIZE = DHCPOption{
		Option: OPTION_MAX_MSG_SIZE,
		Len:    2,
		Data:   []byte{0x05, 0xdc},
	}
	DHCP_PARAM_REQ_LIST = DHCPOption{
		Option: OPTION_PARAM_REQ_LIST,
		Len:    3,
		Data:   []byte{1, 3, 6}, //Default to  Subnet Mask(1), Router IP(3), and DNS Server(s)(6)
	}
	DHCP_CLIENT_ID = DHCPOption{
		Option: OPTION_CLIENT_ID,
		Len:    5,
		Data:   []byte{1, byte('D'), byte('E'), byte('M'), byte('O')},
	}
	DHCP_END = DHCPOption{
		Option: OPTION_END,
	}
)

type DHCPOption struct {
	Option DHCPMsgOption
	Len    byte
	Data   []byte
}

func NewDHCPOption(code DHCPMsgOption, data []byte) DHCPOption {
	if len(data) > 255 {
		data = data[:255] //TODO: report this truncating up in some way so that it's not a surprise
	}
	op := DHCPOption{
		Option: code,
		Len:    byte(len(data)),
		Data:   data,
	}
	return op
}

func (op *DHCPOption) Bytes() []byte {
	b := []byte{}
	b = append(b, byte(op.Option), op.Len)
	b = append(b, op.Data...)
	return b
}

func OptionsFromBytes(ops []byte) ([]DHCPOption, error) {
	var options = []DHCPOption{}
	for i := 0; ; {
		opCode := ops[i]
		if opCode == 255 {
			options = append(options, DHCPOption{Option: DHCPMsgOption(opCode)})
			break
		} else {
			opLen := ops[i+1]
			opData := ops[i+2 : (i+2)+int(opLen)]
			options = append(options, DHCPOption{Option: DHCPMsgOption(opCode), Len: opLen, Data: opData})
			i = i + 2 + int(opLen)
		}
	}

	return options, nil
}
