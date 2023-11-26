package main

import (
	"fmt"
	"net"
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
}

// 用户工厂
func NewUser(conn net.Conn) *User {
	s := conn.RemoteAddr().String()
	u := &User{Username: s, Address: s, C: make(chan string), Conn: conn}
	go u.ListerMesage()
	return u
}

// 实时从私人消息中得到消息
func (this *User) ListerMesage() {
	for {
		// 阻塞，等待消息
		msg := <-this.C
		// 发送消息
		this.Conn.Write([]byte(msg))
		//服务端控制台打印
		fmt.Println(msg)
	}
}
