package utils

import (
	"fmt"
	"net"
)

func RandPort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		fmt.Println("")
	}
	listen, err := net.Listen("tcp", addr.String())
	if err != nil {
		panic(err)
	}
	defer listen.Close()
	port := listen.Addr().(*net.TCPAddr).Port
	fmt.Println(port)
	return port
}
