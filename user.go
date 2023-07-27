package main

import (
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

// 创建一个用户的api
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

// 监听当前user channel的方法，一旦有消息，就直接发送给对端客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))
	}
}

// 用户上线
func (u *User) Online() {
	//用户上线。将用户加入OnlineMap中
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	//广播当前用户上线消息
	u.server.BoardCast(u, "已上线")
}

// 用户下线
func (u *User) Offline() {
	//用户下线。将用户从OnlineMap删除
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	//广播当前用户上线消息
	u.server.BoardCast(u, "已下线")
}

func (u *User) SendMessage(msg string) {
	u.conn.Write([]byte(msg))
}

// 用户处理信息
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前在线用户都有哪些
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Name + "]:" + "在线...\n"
			u.SendMessage(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//消息格式：rename|张三
		newName := strings.Split(msg, "|")[1]

		//判断name是否存在
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMessage("当前用户名被使用")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMessage("你已经更新用户名：" + u.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式：to|张三|消息内容

		//获取对方的用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMessage("消息格式不正确")
		}
		//根据用户名得到对方的User对象
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMessage("该用户名不存在")
			return
		}
		//获取消息内容，通过对方的User对象将消息内容发送过去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMessage("消息为空")
			return
		}
		remoteUser.SendMessage(u.Name + "发送的:" + content + "\n")
	} else {
		u.server.BoardCast(u, msg)
	}
}
