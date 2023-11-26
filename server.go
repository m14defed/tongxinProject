package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// 创建一个服务器
type Server struct {
	//服务器ip地址
	Ip        string
	port      string
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

// 服务器工厂
func NewServer(ip string, port string) *Server {
	return &Server{Ip: ip, port: port, OnlineMap: make(map[string]*User), Message: make(chan string)}
}

// 监听Messager实时广播消息
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		this.mapLock.RLock()
		for _, user := range this.OnlineMap {
			user.C <- msg

		}
		this.mapLock.RUnlock()

	}
}

// 向广播中心发送消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Username + ":" + msg + "]"
	this.Message <- sendMsg
}

// 连接后的的处理器
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("holle world")
	user := NewUser(conn, this)

	this.mapLock.Lock()
	this.OnlineMap[user.Username] = user
	this.mapLock.Unlock()
	//广播上线消息
	//this.BroadCast(user, "login success")
	user.Online()
	//广播用户的消息
	Duque := make([]byte, 2048)
	n, err := conn.Read(Duque)
	if n == 0 {
		fmt.Println("用户下线")
		//this.BroadCast(user, "user offline")
		user.Offline()
		return
	}
	if err != nil && err != io.EOF {
		fmt.Println("conn err", err)
	}
	msg := string(Duque[:n-1])
	//this.BroadCast(user, msg)
	user.SendMessage(msg)
	select {}

}

// 启动服务器
func (this *Server) Start() {
	//监听
	listen, err := net.Listen("tcp", fmt.Sprintf("%v:%v", this.Ip, this.port))
	if err != nil {
		fmt.Println("err", err)
	}
	//关闭
	defer listen.Close()
	//启动消息推发
	go this.ListenMessager()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		go this.Handler(conn)
	}
}
