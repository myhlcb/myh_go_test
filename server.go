package main

import (
	"fmt"
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
	user := NewUser(conn)
	this.mapLock.Lock() //加锁
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock() //释放锁
	this.BroadCast(user, "已上线")
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
