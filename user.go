package main

import "net"

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

// 用户处理信息
func (u *User) DoMessage(msg string) {
	u.server.BoardCast(u, msg)
}
