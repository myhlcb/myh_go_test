package main

import (
	"flag"
	"fmt"
	"net"
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
func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		switch client.flag {
		case 1:
			fmt.Println("公聊模式")
			break
		case 2:
			fmt.Println("私聊模式")
			break
		case 3:
			fmt.Println("更新用户名")
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
	fmt.Println(">>>链接服务器成功...")
	//启动客户端业务
	client.Run()
}
