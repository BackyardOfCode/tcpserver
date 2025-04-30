package server

import (
	"fmt"
	"net"
	"strings"
	"tcpserver/utils"
)

// User 对象
type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动监听消息
	go user.ListenMessage()
	return user
}

//消息监听

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

// 上线
func (this *User) Online() {
	this.server.mapLock.Lock()
	fmt.Println(this.Name, this)
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	//广播用户上线消息
	this.server.BroadCast(this, fmt.Sprintf("用户%s已经上线", this.Name))
}

// 下线
func (this *User) Offline() {
	this.server.mapLock.Lock()
	fmt.Println(this.Name, this)
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	//广播用户上线消息
	this.server.BroadCast(this, fmt.Sprintf("用户%s已下线", this.Name))

}

//Usersend message

func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}

// 处理消息发信息

func (this *User) DoMessage(msg string) {

	//用户发送以/开头命令为出发系统命令

	if strings.HasPrefix(msg, "/") {

		if strings.HasPrefix(msg, "/help") {
			this.SendMessage(utils.GetHelp())
			return
		}

		//查询所有用户
		if strings.HasPrefix(msg, "/all") {
			fmt.Println("查询用户在线信息")
			this.server.mapLock.Lock()
			for _, user := range this.server.OnlineMap {
				onlineMsg := "[地址端口为:" + user.Addr + "] 用户名为:【" + user.Name + "】:在线\n"
				this.SendMessage(onlineMsg)
			}

			this.server.mapLock.Unlock()
			return
		}

		//修改用户名

		if strings.HasPrefix(msg, "/rename:") {
			newName := strings.TrimPrefix(msg, "/rename:")

			_, ok := this.server.OnlineMap[newName]

			if ok {
				this.SendMessage("用户名已被占用")
			} else {
				this.server.mapLock.Lock()
				delete(this.server.OnlineMap, this.Name)
				this.server.OnlineMap[newName] = this
				this.Name = newName
				this.server.mapLock.Unlock()

				this.SendMessage("修改成功! 当前用户名为" + this.Name)
			}

			return
		}

		//私聊 消息格式 /msg/name:msgcontent
		if strings.HasPrefix(msg, "/msg/") {
			remoteName := strings.Split(msg, "/")[2]
			fmt.Println(remoteName)
			_, ok := this.server.OnlineMap[remoteName]

			if !ok {
				this.SendMessage("用户没有找到")
				return
			}

			msgContent := strings.Split(msg, "/")[3]
			if msgContent == "" {
				this.SendMessage("消息内容为空")
				return
			}

			this.server.OnlineMap[remoteName].SendMessage(fmt.Sprintf("来自%s的消息【%s】", this.Name, msgContent))

			return
		}
	}

	this.server.BroadCast(this, msg)
}
