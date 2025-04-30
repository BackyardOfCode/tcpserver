package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//	在线用户列表

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

//新建服务器

func NewServer(ip string, port int) *Server {
	s := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return s
}

// 启动服务器
func (this *Server) Start() {
	listner, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("net list err:", err)
	}
	fmt.Sprintf("Net Server Is Run : %s:%d", this.Ip, this.Port)
	defer listner.Close()
	go this.ListMessage()
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("list accept err : ", err)
			continue
		}

		go this.handler(conn)
	}
}

//监听Message 的go ， 有消息发送全部

func (this *Server) ListMessage() {
	for {
		msg := <-this.Message
		fmt.Println("receive Message :", msg)
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

//广播消息方法

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

// 服务器消息处理
func (this *Server) handler(conn net.Conn) {
	fmt.Println(conn)
	fmt.Println("tcp 链接成功")
	user := NewUser(conn, this)

	//广播用户上线消息
	user.Online()

	//接受客户发送的消息
	isLive := make(chan bool)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)

			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err : ", err)
				return
			}
			//去除\n
			msg := string(buf[:n])
			user.DoMessage(msg)
			isLive <- true
		}
	}()
	//当前handler 阻塞
	for {
		select {
		case <-isLive:
		//	当前用户是否活跃 应该重置定时器
		//不做任何事情 激活select 更细定时器

		case <-time.After(time.Second * 30):
			//	已经超时 讲
			user.SendMessage("超时被踢...")
			close(user.C)
			conn.Close()
			return
		}
	}

}
