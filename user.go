package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	//用户名
	Username string `json:"username,omitempty"`
	//地址
	Address string `json:"address,omitempty"`
	//channel
	C chan string `json:"c,omitempty"`
	//连接
	Conn net.Conn `json:"conn,omitempty"`
	//server

	UserServer *Server
}

// 用户工厂
func NewUser(conn net.Conn, UserServer *Server) *User {
	s := conn.RemoteAddr().String()
	u := &User{Username: s, Address: s, C: make(chan string), Conn: conn, UserServer: UserServer}
	go u.ListerMesage()
	return u
}

// 上线功能
func (this *User) Online() {
	this.UserServer.BroadCast(this, "上线了")
}

// 下线功能
func (this *User) Offline() {
	this.UserServer.BroadCast(this, "下线了")
}

// 发给自己
func (this *User) SendMy(msg string) {
	this.C <- msg
}

// 发送消息
func (this *User) SendMessage(msg string) {
	if msg == "why" {
		this.UserServer.mapLock.Lock()
		for _, user := range this.UserServer.OnlineMap {
			str := fmt.Sprintf("%s: 在线", user.Username)
			this.SendMy(str)
		}

		this.UserServer.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		name := strings.Split(msg, "|")[1]
		user := this.UserServer.OnlineMap[name]
		if user != nil {
			this.SendMy("用户名已存在")
		}
		delete(this.UserServer.OnlineMap, this.Username)
		this.Username = name
		this.UserServer.OnlineMap[name] = this
		this.SendMy("改名成功")
		this.UserServer.BroadCast(this, "改名成功")
	} else if len(msg) > 3 && msg[:3] == "to|" {
		name := strings.Split(msg, "|")[1]
		msg := strings.Split(msg, "|")[2]
		user := this.UserServer.OnlineMap[name]
		if user == nil {
			this.SendMy("你私信的用户不存在")
		}
		user.SendMy(msg)
	} else {
		this.UserServer.BroadCast(this, msg)
	}
}

// 实时从私人消息中得到消息
func (this *User) ListerMesage() {
	for {
		// 阻塞，等待消息
		msg := <-this.C
		// 得到消息
		this.Conn.Write([]byte(msg))
		//服务端控制台打印
		fmt.Println(msg)
	}
}
