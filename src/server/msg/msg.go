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
	Password string
	Token string
	Gid int64
}
//请求参数
type REQUEST struct {
	Action string
	Data string
}

//返回值
type RESPONSE struct {
	Code int64
	Msg string
	Data string
}