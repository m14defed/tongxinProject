package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	port string
}

func NewServer(ip string, port string) *Server {
	return &Server{Ip: ip, port: port}
}
func (this *Server) Handler(coon net.Conn) {
	fmt.Println("holle world")
}

func (this *Server) Start() {
	//监听
	listen, err := net.Listen("tcp", fmt.Sprintf("%v:%v", this.Ip, this.port))
	if err != nil {
		fmt.Println("err", err)
	}
	//关闭
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		go this.Handler(conn)

	}
}
