package main

import (
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

// 用户上线
func (this *User) Online() {
	this.server.mapLock.Lock() //加锁
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock() //释放锁
	this.server.BroadCast(this, "已上线")
}

// 用户下线
func (this *User) Outline() {
	this.server.mapLock.Lock() //加锁
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock() //释放锁
	this.server.BroadCast(this, "下线")
}

// 处理用户消息
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
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
	go user.ListenMessage()
	return user
}
