package msg

import (
	"github.com/name5566/leaf/network/json"
)

var Processor = json.NewProcessor()

func init() {
	Processor.Register(&S2C_Close{})
	Processor.Register(&C2S_Auth{})
	Processor.Register(&S2C_Auth{})
}

// Close
const (
	S2C_Close_LoginRepeated = 1
	S2C_Close_InnerError    = 2
)

type S2C_Close struct {
	Err int
}

// Auth
type C2S_Auth struct {
	AccID string
}

const (
	S2C_Auth_OK           = 0
	S2C_Auth_AccIDInvalid = 1
)

type S2C_Auth struct {
	Err int
}
