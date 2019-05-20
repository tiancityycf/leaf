package internal

import (
	"encoding/json"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"reflect"
	"server/msg"
	"time"
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
var group = 0
//处理请求参数
func handleRequest(args []interface{}) {
	// 收到的 Hello 消息
	m := args[0].(*msg.REQUEST)
	// 消息的发送者
	a := args[1].(gate.Agent)

	u := a.UserData().(*msg.User)

	u.Position = m.Body

	log.Release("request",u)

	//response := &msg.RESPONSE{}
	//response.Code = 1
	//response.Msg = "ok"
	//response.Body = m.Body

	// 输出收到的消息的内容
	log.Debug("handleAction %v", m)

	log.Debug("handleAction uids %v", Groups[u.Gid])

	log.Release("group ",group)
	if(group == 0){
		go sendGroupMsg(u.Gid)
	}
	group = 1
	// 给房间内人员发送消息
	//for _,v := range Groups[u.Gid]{
	//	v.WriteMsg(response)
	//}
}

func sendGroupMsg(gid int64){
	tick := time.NewTicker(1000000000 * time.Nanosecond)
	//tick := time.NewTicker(16600000 * time.Nanosecond)

	response := &msg.RESPONSE{}

	response.Msg = "ok"

	for {
		select {
		//此处在等待channel中的信号，因此执行此段代码时会阻塞120秒
		case <-tick.C:
			for _,v := range Groups[gid] {
				if GroupUsers[gid] == nil {
					GroupUsers[gid] = make(map[int64]*msg.User)
				}
				user,ok := v.UserData().(*msg.User)
				if(ok) {
					GroupUsers[gid][user.Uid] = user
				}
			}
			b, _ := json.Marshal(GroupUsers[gid])
			for _,v := range Groups[gid]{
				//user := v.UserData().(*msg.User)
				response.Body = string(b)
				log.Release(response.Body)
				response.Code = 1
				v.WriteMsg(response)
			}
		}
	}
}
