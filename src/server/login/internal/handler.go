package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"math/rand"
	"reflect"
	"server/msg"
	"server/game"
	"time"
)


func init() {
	// 向当前模块（game 模块）注册 Hello 消息的消息处理函数 handleHello
	handler(&msg.User{}, handleUser)
}


func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}
//加入房间
func handleUser(args []interface{}) {
	// 收到的消息
	m := args[0].(*msg.User)
	// 消息的发送者
	a := args[1].(gate.Agent)

	//log.Debug("login a1 %v", ginternal.Agents)
	var g int64 = rand.Int63n(2)

	rand.Seed(time.Now().UnixNano())

	m.Gid = g

	log.Debug("login uid %v", m.Uid)
	log.Debug("login gid %v", m.Gid)
    //存储房间相关信息
    //game.ChanRPC.Go("JoinGroup", a , g , m.Uid )
	skeleton.AsynCall(game.ChanRPC, "JoinGroup", a, g , m.Uid , func(err error) {
		if nil != err {
			log.Error("login failed:",err.Error())
			return
		}
	})

	a.SetUserData(m)
	// 给发送者回应一个 Hello 消息
	a.WriteMsg(m)
}

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}
