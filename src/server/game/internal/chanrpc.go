package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"server/msg"
)

//存储所有agents信息 广播用
var Agents = make(map[gate.Agent]struct{})
//存储房间相关用户信息 组播用
var Groups = make(map[int64]map[int64]gate.Agent)

func init() {
	log.Debug("game init")
	//这里的 NewAgent 和 CloseAgent 会被 LeafServer 的 gate 模块在连接建立和连接中断时调用
	//gate 模块这样调用 game 模块的 NewAgent ChanRPC
	//game.ChanRPC.Go("NewAgent", a)
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)

	//这里的 JoinGroup 和 LeaveGroup 会被 login模块在加入房间、离开房间、连接断开时调用
	//login 模块这样调用 game 模块的 JoinGroup ChanRPC
	//game.ChanRPC.Go("JoinGroup", a , g , m.Uid )
	skeleton.RegisterChanRPC("JoinGroup", rpcJoinGroup)
	skeleton.RegisterChanRPC("LeaveGroup", rpcLeaveGroup)
}
//建立连接
func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.Debug("gamerpcNewAgent %v", a)
	//_ = a
	Agents[a] = struct{}{}
	log.Debug("gamerpcNewAgent %v", Agents)

}
//连接断开
func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.Debug("gamerpcCloseAgent %v", a)
	delete(Agents, a)
	//连接断开时清除房间用户信息
	user := a.UserData().(*msg.User)
	g := user.Gid
	u := user.Uid
	delete(Groups[g], u)

}
//加入房间
func rpcJoinGroup(args []interface{}) {
	a := args[0].(gate.Agent)
	g := args[1].(int64)
	u := args[2].(int64)
	log.Debug("rpcJoinGroup %v", a)
	if Groups[g] == nil {
		Groups[g] = make(map[int64]gate.Agent)
	}
	Groups[g][u] =  a
	log.Debug("gamerpcNewAgent %v", Groups[g])

}
//离开房间
func rpcLeaveGroup(args []interface{}) {
	a := args[0].(gate.Agent)
	g := args[1].(int64)
	u := args[2].(int64)
	log.Debug("gamerpcCloseAgent %v", a)
	delete(Groups[g], u)
}