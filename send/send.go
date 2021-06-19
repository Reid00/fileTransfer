package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

// function: 发送文件端

// Get File Name from agrs
func GetFileName(args []string) string {
	filePath := args[1]

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	fileName := fileInfo.Name()
	return fileName
}

// Send file to client
func SendFile(conn net.Conn, filePath string) {
	// 1. 读取本地文件
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("os.Open error: ", err)
		return
	}
	defer f.Close()
	// 2. 发送给client
	buf := make([]byte, 1024)
	for {
		// 读文件内容
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Send done")
			} else {
				fmt.Println("f.Read error: ", err)
			}
			return
		}
		// 发送数据，读多少，写多少
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("conn.Write error", err)
		}
	}
}

func main() {

	// 1. 获取文件名称
	argList := os.Args
	if len(argList) != 2 {
		fmt.Println("right format is: send.exe abs file path")
		return
	}
	fileName := GetFileName(argList)
	// 2. 发送文件名给接收端
	// 2.1 创建用于通信的socket
	conn, err := net.Dial("tcp", "127.0.0.1:8008")
	if err != nil {
		fmt.Println("net.Dail error: ", err)
		return
	}
	defer conn.Close()

	// 2.2 发送文件名给接收端
	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("conn.Write error: ", err)
		return
	}

	// 2.3 接收 接收端返回的回执socket
	buf := make([]byte, 10)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read error: ", err)
		return
	}
	ret := string(buf[:n])
	// 判断回执是否成功
	if ret == "ok" {
		// 3. 发送文件给接收端
		SendFile(conn, os.Args[1])
	}
}
