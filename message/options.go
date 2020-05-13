package message

type OptionCode byte

const (
	OPTION_SUBNET_MASK    OptionCode = 1
	OPTION_ROUTER         OptionCode = 3
	OPTION_DNS            OptionCode = 6
	OPTION_DOMAIN_NAME    OptionCode = 15
	OPTION_LEASE_TIME     OptionCode = 51
	OPTION_REQUEST_IP     OptionCode = 50
	OPTION_MSG_TYPE       OptionCode = 53
	OPTION_SERVER_ID      OptionCode = 54
	OPTION_MAX_MSG_SIZE   OptionCode = 57
	OPTION_PARAM_REQ_LIST OptionCode = 55
	OPTION_CLIENT_ID      OptionCode = 61
	OPTION_END            OptionCode = 255
)

var (
	DefaultDiscoverOps = []Option{
		Option{
			Code: OPTION_MSG_TYPE,
			Len:  1,
			Data: []byte{1},
		},
		Option{
			Code: OPTION_MAX_MSG_SIZE,
			Len:  2,
			Data: []byte{0x05, 0xdc},
		},
		Option{
			Code: OPTION_PARAM_REQ_LIST,
			Len:  3,
			Data: []byte{1, 3, 6}, //Default to  Subnet Mask(1), Router IP(3), and DNS Server(s)(6)
		},
		Option{
			Code: OPTION_CLIENT_ID,
			Len:  5,
			Data: []byte{1, byte('D'), byte('E'), byte('M'), byte('O')},
		},
		Option{
			Code: OPTION_END,
		},
	}

	DefaultRequestOps = []Option{
		Option{
			Code: OPTION_MSG_TYPE,
			Len:  1,
			Data: []byte{3},
		},
		Option{
			Code: OPTION_MAX_MSG_SIZE,
			Len:  2,
			Data: []byte{0x05, 0xdc},
		},
		Option{
			Code: OPTION_CLIENT_ID,
			Len:  5,
			Data: []byte{1, byte('D'), byte('E'), byte('M'), byte('O')},
		},
		Option{
			Code: OPTION_END,
		},
	}
)

type Option struct {
	Code OptionCode
	Len  byte
	Data []byte
}

func NewOption(code OptionCode, data []byte) Option {
	if len(data) > 255 {
		data = data[:255] //TODO: report this truncating up in some way so that it's not a surprise
	}
	op := Option{
		Code: code,
		Len:  byte(len(data)),
		Data: data,
	}
	return op
}

func (op *Option) Bytes() []byte {
	b := []byte{}
	b = append(b, byte(op.Code), op.Len)
	b = append(b, op.Data...)
	return b
}

func OptionsFromBytes(ops []byte) ([]Option, error) {
	var options = []Option{}
	for i := 0; ; {
		opCode := ops[i]
		if opCode == 255 {
			options = append(options, Option{Code: OptionCode(opCode)})
			break
		} else {
			opLen := ops[i+1]
			opData := ops[i+2 : (i+2)+int(opLen)]
			options = append(options, Option{Code: OptionCode(opCode), Len: opLen, Data: opData})
			i = i + 2 + int(opLen)
		}
	}

	return options, nil
}

func FindOption(code OptionCode, ops []Option) *Option {
	for i := range ops {
		if ops[i].Code == code {
			return &ops[i]
		}
	}
	return nil
}
