package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	//链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1 公聊模式")
	fmt.Println("2 私聊模式")
	fmt.Println("3 更新用户名")
	fmt.Println("0 退出")
	// 接受用户输入的数据
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法字符")
		return false
	}
}

// 公聊模式
func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>聊天内容:")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("publicChat conn Write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>聊天内容:")
		fmt.Scanln(&chatMsg)
	}
}
func (client *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

// 私聊模式
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string
	client.SelectUser()
	fmt.Println(">>>>请输入聊天对象[用户名],exit退出")
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println(">>>>请输入聊天内容,exit退出")
		//阻止for循环，等待用户输入
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("publicChat conn Write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>请输入聊天内容:")
			fmt.Scanln(&chatMsg)
		}
		// exit会跳出for循环
		client.SelectUser()
		fmt.Println(">>>>请重新输入聊天对象[用户名],exit退出")
		fmt.Scanln(&remoteName)
	}
}

// 更新用户名
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>请输入用户名:")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return false
	}
	return true
}

// 处理server回应消息，直到显示标准输出
func (client *Client) DealRes() {
	//一旦client.conn有数据，就直接copy到stdout标准输出，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}
func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置默认服务器IP(127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置默认服务器PORT(8888)")
}
func main() {
	//解析命令行
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("客户端建立连接失败")
		return
	}
	// 单独开启goroutine处理server
	go client.DealRes()
	fmt.Println(">>>链接服务器成功...")
	//启动客户端业务
	client.Run()
}
