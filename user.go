package main

import (
	"fmt"
	"net"
)

type User struct {
	Username string      `json:"username,omitempty"`
	Address  string      `json:"address,omitempty"`
	C        chan string `json:"c,omitempty"`
	Conn     net.Conn    `json:"conn,omitempty"`
}

func NewUser(conn net.Conn) *User {
	s := conn.RemoteAddr().String()
	u := &User{Username: s, Address: s, C: make(chan string), Conn: conn}
	go u.ListerMesage()
	return u
}
func (this *User) ListerMesage() {
	for {
		// 阻塞，等待消息
		msg := <-this.C
		// 发送消息
		this.Conn.Write([]byte(msg))
		fmt.Println(msg)
	}
}
