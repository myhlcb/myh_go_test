package main

import (
	"fmt"
	"net"
	"strings"
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

// 给当前user对应客户端发消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 处理用户消息
func (this *User) DoMessage(msg string) {
	// 查询当前用户列表
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			OnlineMsg := fmt.Sprintf("[%s]%s:在线\n", user.Addr, user.Name)
			this.SendMsg(OnlineMsg)

		}
		this.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//消息格式rename|张三
		newName := strings.Split(msg, "|")[1]
		//判断name是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名被占用\n")

		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendMsg("您已经更新用户名:" + this.Name + "\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式 to|张三|消息内容
		//1 获取对方名字
		remoteName := strings.Split(msg, "|")[1]
		//判断name是否存在
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("当前用户不存在\n")
			return
		}
		//2 根据用户名得到对方User对象
		content := strings.Split(msg, "|")[2]
		remoteUser.SendMsg(this.Name + "对你说:" + content)
		//3 获取消息内容发送过去

	} else {
		this.server.BroadCast(this, msg)
	}

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
