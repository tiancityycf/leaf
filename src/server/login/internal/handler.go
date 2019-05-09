package internal

import (
	"context"
	"fmt"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2"
	"math/rand"
	"os"
	"reflect"
	"server/conf"
	"server/game"
	"server/msg"
	"time"
)

var session *mgo.Session
var client *mongo.Client
var collection *mongo.Collection
var insertOneRes *mongo.InsertOneResult


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

	//mgo 驱动操作数据库
	session,err := mgo.Dial(conf.Server.MgodbAddr)
	if err != nil {
		log.Fatal("dial mongodb error: %v", err)
	}

	db := session.DB("game")
	// load
	//err = db.C("users").Find(bson.M{"uid": m.Uid}).One(&m)
	err = db.C("users").Insert(m)
	if err != nil {
		// unknown error
		if err != mgo.ErrNotFound{
			log.Error("load acc %v data error: %v", m.Uid, err)
			return
		}

	}

	//mongo-go-driver驱动操作数据库
	client,err = mongo.Connect(getContext(),options.Client().ApplyURI(conf.Server.MgodbAddr))
	checkErr(err)
	//判断服务是否可用
	if err = client.Ping(getContext(), readpref.Primary()); err != nil {
		checkErr(err)
	}

	//选择数据库和集合
	collection = client.Database("game").Collection("user")

	//插入一条数据
	insertOneRes, err = collection.InsertOne(getContext(), m);
	checkErr(err)
	fmt.Printf("InsertOne插入的消息ID:%v\n", insertOneRes.InsertedID)

	a.SetUserData(m)
	// 给发送者回应一个 Hello 消息
	a.WriteMsg(m)
}

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	return
}

func checkErr(err error) {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("没有查到数据")
			os.Exit(0)
		} else {
			fmt.Println(err)
			os.Exit(0)
		}

	}
}

