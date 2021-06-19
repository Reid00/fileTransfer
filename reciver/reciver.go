package main

import (
	"fmt"
	"net"
	"os"
)

// 接收文件客户端

// 接收发送的文件内容，写到本地
func ReciverFile(conn net.Conn, filePath string) {

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("os.OpenFile error: ", err)
		return
	}
	defer f.Close()

	// 读取发送端发送的文件内容
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			fmt.Println("Reciver Done")
			return
		}
		if err != nil {
			fmt.Println("conn.Read error: ", err)
			return
		}
		// 读取的文件内容写到本地文件
		_, err = f.Write(buf[:n])
		if err != nil {
			fmt.Println("f.Write error: ", err)
			return
		}

	}
}

func main() {
	// 创建一个用于监听连接的socket
	listener, err := net.Listen("tcp", "127.0.0.1:8008")
	if err != nil {
		fmt.Println("net.Listen error: ", err)
		return
	}
	defer listener.Close()
	// 阻塞监听客户端连接
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("listener.Accept error: ", err)
		return
	}
	defer conn.Close()

	// 接收发送的文件名, most 1024 length
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read error: ", err)
		return
	}
	fileName := string(buf[:n])
	// 返回 回执给发送端
	_, err = conn.Write([]byte("ok"))
	if err != nil {
		fmt.Println("conn.Write error: ", err)
		return
	}

	// 接收发送端 发送的文件内容
	ReciverFile(conn, fileName)
}
