package main

import (
	"fmt"
	"io"
	"net"
	"sync" //加锁用
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex //加锁
	Message   chan string
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s:%s", user.Addr, user.Name, msg)
	this.Message <- sendMsg
}

// 处理用户登录
func (this *Server) Handle(conn net.Conn) {
	user := NewUser(conn, this)
	user.Online()
	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			//从当前conn中读取数据
			n, err := conn.Read(buf)
			//标识语用户下线
			if n == 0 {
				user.Outline()
				return
			}
			//非法操作
			if err != nil && err != io.EOF {
				fmt.Println("Conn READ err:", err)
				return
			}
			// 提取用户消息,去除\n
			msg := string(buf[:n-1])
			//消息处理(交给用户)
			user.DoMessage(msg)
		}
	}()
	//当前handle阻塞
	select {}
}
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("err happen:", err)
	}
	defer listener.Close()
	//启动监听message的goroutine
	go this.ListenMessage()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept err")
			continue
		}
		go this.Handle(conn)
	}

}
