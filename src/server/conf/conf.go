package conf

import (
	"log"
	"time"
)

var (
	// log conf
	LogFlag = log.LstdFlags
	//在日志中输出文件名和行号 可用的 LogFlag s见：https://golang.org/pkg/log/#pkg-constants
	//LogFlag = log.Lshortfile

	// gate conf
	PendingWriteNum        = 2000
	MaxMsgLen       uint32 = 4096
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 2
	LittleEndian           = false

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000
)
