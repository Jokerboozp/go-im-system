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

func (c *Client) UpdateName() bool {
	fmt.Println("请输入用户名")
	fmt.Scanln(&c.Name)

	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.Menu() != true {

		}

		//根据不同的模式处理不同的业务
		switch c.flag {
		case 1:
			c.PublicChat()
			fmt.Println(">>>>>>>>公聊模式选择<<<<<<<<")
			break
		case 2:
			fmt.Println(">>>>>>>>>私聊模式选择<<<<<<<<<")
			break
		case 3:
			fmt.Println(">>>>>>>更改用户名选择<<<<<")
			c.UpdateName()
			break
		}
	}
}

func (c *Client) PublicChat() {
	var chatMsg string
	//提示用户输入消息
	fmt.Println(">>>>>>>请输入聊天内容，exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println(err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>>>请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
	//发送给服务器
}

func (c *Client) DealResponse() {
	//一旦conn有数据，就直接copy到stdout标准输出中，永久阻塞监听
	io.Copy(os.Stdout, c.conn)
}

func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("服务器链接失败>>>>>>>>>>>")
		return
	}

	//单独开启一个goroutine去处理server的回执消息
	go client.DealResponse()

	fmt.Println("服务器链接成功>>>>>>>>>>>>")

	client.Run()
}
