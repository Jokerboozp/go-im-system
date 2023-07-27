package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	client.conn = conn
	return client

}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("服务器链接失败")
		return
	}

	fmt.Println("服务器链接成功")

	select {}
}
