package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"reflect"
	"server/msg"
)

func init() {
	// 向当前模块（game 模块）注册 Hello 消息的消息处理函数 handleHello
	handler(&msg.Hello{}, handleHello)
	handler(&msg.REQUEST{}, handleRequest)

}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleHello(args []interface{}) {
	// 收到的 Hello 消息
	m := args[0].(*msg.Hello)
	// 消息的发送者
	a := args[1].(gate.Agent)

	// 输出收到的消息的内容
	log.Debug("hello %v", m.Name)

	// 给发送者回应一个 Hello 消息
	a.WriteMsg(&msg.Hello{
		Name: "client",
	})
}

//处理请求参数
func handleRequest(args []interface{}) {
	// 收到的 Hello 消息
	m := args[0].(*msg.REQUEST)
	// 消息的发送者
	a := args[1].(gate.Agent)

	u := a.UserData().(*msg.User)

	response := &msg.RESPONSE{}
	response.Code = 1
	response.Msg = "ok"
	response.Data = "xxx"

	// 输出收到的消息的内容
	log.Debug("handleAction %v", m)

	log.Debug("handleAction uids %v", Groups[u.Gid])
	// 给房间内人员发送消息
	for _,v := range Groups[u.Gid]{
		v.WriteMsg(response)
	}
}
