package msg

import (
	"github.com/name5566/leaf/network/json"
	"github.com/name5566/leaf/network/protobuf"
)

// 使用默认的 JSON 消息处理器（默认还提供了 protobuf 消息处理器）
var (
	Processor = json.NewProcessor()
	ProtobufProcessor = protobuf.NewProcessor()
	)
func init() {
	// 测试 这里我们注册了一个 JSON 消息 Hello
	Processor.Register(&Hello{})
	//登录授权
	Processor.Register(&User{})
	//请求
	Processor.Register(&REQUEST{})
	//响应
	Processor.Register(&RESPONSE{})
}

// 一个结构体定义了一个 JSON 消息的格式
// 消息名为 Hello
type Hello struct {
	Name string
}

//用户信息
type User struct {
	Uid int64
	Username string
	Rid int64  		//房间ID
	Token string 	//交互认证Token
	Gid int64  		//所属职业
	Action string   //动作
	Position string //位置
	Direction string//当前朝向
}
//请求参数
type REQUEST struct {
	Method string
	Body string
}

//返回值
type RESPONSE struct {
	Code int64
	Msg string
	Body string
}