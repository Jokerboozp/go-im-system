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
	flag       int //当前client的模式
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	client.conn = conn
	return client

}

func (c *Client) Menu() bool {

	var flag int

	fmt.Println("1.公共聊天模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println("请输入范围合法的数字")
		return false
	}
}

var serverIp string
var serverPort int

// client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址,默认127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口,默认8888")
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.Menu() != true {

		}

		//根据不同的模式处理不同的业务
		switch c.flag {
		case 1:
			fmt.Println("公聊模式选择")
			break
		case 2:
			fmt.Println("私聊模式选择")
			break
		case 3:
			fmt.Println("更改用户名选择")
			break
		}
	}
}

func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("服务器链接失败")
		return
	}

	fmt.Println("服务器链接成功")

	client.Run()
}
