package errno

// Errno 是一个自定义的错误封装类型
type Errno struct {
	State int    `json:"state"`
	Msg   string `json:"msg"`
}

func (e *Errno) Error() string {
	return e.Msg
}

func (e Errno) Add(s string) *Errno {
	e.Msg += ": " + s
	return &e
}

type ResponseErrno struct {
	State      int    `json:"state"`
	Msg        string `json:"msg"`
	HttpCode   int    `json:"http_code"`
	OriginBody string `json:"origin_body"`
}

func (e *ResponseErrno) Error() string {
	return e.Msg
}

func (e ResponseErrno) SetCode(code int, s string) *ResponseErrno {
	e.HttpCode = code
	e.OriginBody = s
	return &e
}

func (e ResponseErrno) Add(s string) *ResponseErrno {
	e.Msg += ": " + s
	return &e
}

var (
	NotSupportChainType   = &Errno{10001, "Not support this chain"}
	InvalidTx             = &Errno{10002, "Invalid Tx"}
	InvalidTxType         = &Errno{10003, "Invalid Tx type"}
	InvalidTypeAssert     = &Errno{20001, "Invalid type asset"}
	InvalidStringToBigNum = &Errno{20002, "Invalid string for big number"}
	TxFromNotSet          = &Errno{20002, "From of tx not set"}
	ProviderNotSet        = &Errno{20003, "Not set provider"}
	ParseTxError          = &Errno{20004, "Parse tx error"}
)
