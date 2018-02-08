package new

import (
	"fmt"
	"os"
	"path"

	"github.com/adolphlxm/atc-tool/commands"
	"github.com/adolphlxm/atc-tool/utils"
	"github.com/adolphlxm/atc/logs"
	"strings"
)


var CmdNew = &commands.Command{
	Usage: "new [appname]",
	Use:   "Create a Atc application",
	Options: `
The command "new" creates a folder named [appname] and generates the following structure:
├── conf [配置文件目录]
│   ├── app.ini [核心配置文件]
│   └── error.ini [错误码文件]
├── bin [可执行文件目录]
│── api [RESTFul API 目录]
│     ├── V1 [版本目录]
│     └── router.go [路由文件]
│── model [DB目录]
│── service [服务目录]
│── thrift [RPC]
│    ├── idl [存放Thrift IDL文件目录]
│    ├── gen-go [存放Thrift 生成 go文件目录]
│    ├── ...(.go)
│    └── router.go [路由文件]
└── atc.go [运行主程序文件]
`,
	Run: Run,
}

var atcgo = `package main
import (
	_ "{{.Appname}}/api"
	"github.com/adolphlxm/atc"
)

func main() {
	atc.Run()
}
`

var config = `
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
; 开发模式                                           ;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
[local]
; DEBUG模式
app.debug = true
; 项目名称
app.name = ATC-ram

; 错误码配置文件
error.file = ./conf/error.ini
; HTTP/Websocket
; 是否支持HTTP/Websocket通信
;   - true : 支持
;   - false: 不支持
http.support = true
; HTTP服务地址
;   - e.g. "","127.0.0.1","localhost"
http.addr = 127.0.0.1
; HTTP服务端口
;   - e.g. 80
http.port = 8080
; HTTP退出最多等待时间
;   - 单位:s
http.qtimeout = 30
; 请求超时时间
;   - 单位:s
http.readtimeout = 0
; 响应超时时间
;   - 单位:s
http.writetimeout = 0
; POST上传文件最大内存
; 默认值:1 << 26 64M
post.maxmemory = 67108864

; 支持访问域
;   - * 表示允许访问所有域
front.host = *

; Thrift-RPC
; thrift允许定义一个简单的定义文件中的数据类型和服务接口，
; 以作为输入文件，编译器生成代码用来方便地生成RPC客户端和服务器通信的无缝跨编程语言。
;
; Thrift-RPC通信开关
;   - true : 支持
;   - false: 不支持
thrift.support = true
; Thrift-DEBUG调试
;   - true : 打开,打开DEBUG模式后会输出ATC_logs为前缀的通信日志,方便调试时排查问题
;   - false: 关闭
thrift.debug = true
; Thrift服务地址
;   - e.g. "","127.0.0.1","localhost"
thrift.addr = 127.0.0.1
; Thrift服务端口
;   - e.g. 9090
thrift.port = 10011
thirft.secure = false
; Thrift传输格式(通信层)
;   - binary : 二进制编码格式进行数据传输
;   - bompact : 高效率的、密集的二进制编码格式进行数据传输(压缩)
;   - json : 使用JSON的数据编码协议进行数据传输
;   - [暂不支持]simplejson : 只提供JSON只写的协议,适用于通过脚本语言解析
thrift.protocol = binary
; Thrift数据传输方式(传输层)
;   - framed : 使用非阻塞式方式,按块的大小进行传输
;               以帧为传输单位，帧结构为：4个字节（int32_t）+传输字节串，头4个字节是存储后面字节串的长度，该字节串才是正确需要传输的数据
;   - [暂不支持]memorybuffer : 将内存用于I/O
;   - buffered : 对某个transport对象操作的数据进行buffer,即从buffer中读取数据进行传输,或将数据直接写入到buffer
thrift.transport = framed
; thriftRPC 退出最多等待时间
;   - 单位:s
thrift.qtimeout = 500
; Thrift客户端请求超时时间
;   - 单位:s, 0表示不限制
thrift.client.timeout = 10

; Log
;
;
log.support = true
; Log级别
;LevelFatal
;LevelError
;LevelWarn
;LevelNotice
;LevelInfo
;LevelTrace
;LevelDebug
log.level = LevelFatal
; Log输出
;   stdout : 控制台输出
;   file : 文件输出
log.output = stdout
; Log指定日志路径文件
;   写入file文件方式时需要填写该项
;   指定一个日志写入文件路径
log.dir = ./log/test.log
; Log日志文件最大尺寸,单位：字节, 默认1.7G
log.maxsize = 1887436800
; Log日志文件缓冲区，满了后会执行flush刷入磁盘, 默认256KB
log.buffersize = 262,144
; Log日志定时刷新时间, 单位: 秒
;   默认 : 10
log.flushinterval = 10

; 数据库
; 是否支持ORM数据库引擎
;   - true : 支持
;   - false: 不支持
orm.support = true
; ORM日志级别
;   LOG_UNKNOWN
;   LOG_OFF
;   LOG_WARNING
;   LOG_INFO
;   LOG_DEBUG
orm.log.level = LOG_DEBUG
; 连接池的空闲数大小
orm.maxidleconns =
; 最大打开的连接数
orm.maxopenconns =
; 定时监测数据库连接是否鲜活
;   单位: 秒
;   注意: 某些数据库有连接超时设置的，可以通过起一个定期Ping的Go程来保持连接鲜活。
;   默认值: 30
orm.pingtime = 180
; 数据库名
;   格式: 别名
;   e.g. test_w:写库,test_r:写库
orm.aliasnames = atc_ram
; 数据库配置
;   e.g. dbs.别名.类型 = 值
orm.atc_ram.driver = mysql
orm.atc_ram.host = 127.0.0.1:3306
orm.atc_ram.user = root
orm.atc_ram.password = 123456
orm.atc_ram.dbname = atc_ram
; 生产队列
;   - true : 支持
;   - false: 不支持
queue.publisher.support = false
queue.publisher.aliasnames = p1,p2
queue.publisher.p1.driver = redis
queue.publisher.p1.addrs = redis://127.0.0.1:6379
;   - redis
;   - nats
queue.publisher.drivername = redis
queue.publisher.addrs = redis://127.0.0.1:6379
; 消费队列
;   - true : 支持
;   - false: 不支持
queue.consumer.support = false
queue.consumer.aliasnames = c1,c2
queue.consumer.c1.driver = redis
queue.consumer.c1.addrs = redis://127.0.0.1:6379

; cacahe
;   格式: 别名
;   e.g. m1,r1
cache.aliasnames = mem,redis
cache.support = false
cache.mem.addrs = memcache://127.0.0.1:11211
cache.redis.addrs = redis://:123456@127.0.0.1:6379/0?maxIdle=10&maxActive=10&idleTimeout=3

; mongodb
;   格式: 别名
;   e.g. m1,r1
mgo.aliasnames = m
mgo.support = false
mgo.m.addrs = mongodb://127.0.0.1:27017
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
; 测试模式                                           ;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
[dev]
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
; 生产模式                                           ;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
[prod]
`

var restful_api = `package api

import (
	"github.com/adolphlxm/atc"
)

type ApiHandler struct {
	atc.Handler
}
func (this *ApiHandler) Get() {
	this.JSON(map[string]interface{}{"atc":"The ATC restful API runs successfully"})
}


type Api2Handler struct {
	atc.Handler
}
func (this *Api2Handler) Get() {
	userid := this.Ctx.ParamInt64("userid")
	this.JSON(map[string]interface{}{
		"atc":"The ATC restful API2 runs successfully",
		"userid":userid,
	})
}


type Api2TestHandler struct {
	atc.Handler
}
func (this *Api2TestHandler) Get() {
	this.JSON(map[string]interface{}{
		"atc":"The ATC restful API3 runs successfully",
	})
}

`
var restful_routers = `package api

import (
	"github.com/adolphlxm/atc"
)

func init() {
	v1 := atc.Route.Group("V1")
	{
		// V1版本过滤器, 根据路由规则加载。
		// 支持三种过滤器：
		//      1. EFORE_ROUTE                    //匹配路由之前
		//      2. BEFORE_HANDLER                 //匹配到路由后,执行Handler之前
		//      3. AFTER                          //执行完所有逻辑后
		//v1.AddFilter(atc.BEFORE_ROUTE, ".*", V1.AfterAuth)

		// V1版本测试
		v1.AddRouter("api", &ApiHandler{})
		v1.AddRouter("api2.{userid:[0-9]*}", &Api2Handler{})
		v1.AddRouter("api2.test", &Api2TestHandler{})
	}

}
`

func init() {
	commands.Register(CmdNew)
	//_, file, _, _ := runtime.Caller(1)
	//apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
}

//├── conf
//│   ├── app.ini
//│   └── error.ini
//├── bin
//├── src
//│   ├── api
//│         └── router.go
//│   └── model
//│   └── service
//│   └── thrift
//└── atc.go
func Run(cmd *commands.Command, args []string) int {

	if len(args) != 1 {
		fmt.Println("Usage: new [appname]")
		fmt.Println(commands.ErrUseError)
		logs.Error("Argument [appname] is missing")
		return 1
	}

	cmd.Flag.Parse(args)
	args = cmd.Flag.Args()

	apppath := args[0]
	packpath, err := utils.CheckEnv(apppath)
	if err != nil {
		logs.Error(err.Error())
		return 1
	}

	// Create appname
	os.Mkdir(apppath,0755)
	// Create conf
	os.Mkdir(path.Join(apppath, "conf"),0755)
	utils.WriteToFile(path.Join(apppath,"conf","app.ini"),config)
	utils.WriteToFile(path.Join(apppath,"conf","error.ini"),"")
	// Create bin
	os.Mkdir(path.Join(apppath, "bin"),0755)

	// Create api
	os.Mkdir(path.Join(apppath,"api"),0755)
	utils.WriteToFile(path.Join(apppath,"api","api.go"),restful_api)
	utils.WriteToFile(path.Join(apppath,"api","routers.go"),restful_routers)
	// Create src/model
	os.Mkdir(path.Join(apppath, "model"),0755)
	// Create src/service
	os.Mkdir(path.Join(apppath, "service"),0755)
	// Create src/thrift
	os.Mkdir(path.Join(apppath, "thrift"),0755)
	os.Mkdir(path.Join(apppath, "thrift","gen-go"),0755)
	os.Mkdir(path.Join(apppath, "thrift","idl"),0755)
	// Create atc.go
	utils.WriteToFile(path.Join(apppath, apppath + ".go"),strings.Replace(atcgo,"{{.Appname}}",packpath,-1))


	// lwgo
	// Create src/thrift
	os.Mkdir(path.Join(apppath, "vendor"),0755)
	logs.Tracef("'%s' project created successfully", apppath)
	return 2
}
