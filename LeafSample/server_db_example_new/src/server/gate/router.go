package gate

import (
	"server/login"
	"server/msg"
)

func init() {
	// login
	msg.Processor.SetRouter(&msg.C2S_Auth{}, login.ChanRPC)

	// game
}
