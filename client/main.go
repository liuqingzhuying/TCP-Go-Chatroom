package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

// ChatInfo 聊天信息结构体
type ChatInfo struct{
	UserName string
	ChatData string
}

var wg sync.WaitGroup

func setStatus(username string){
	sql := "update user_status set status=0 where user_name=?"
	_, err := db.Exec(sql,username)
	if err != nil{
		fmt.Println(err)
	}
}

// MessageSend 客户端发送信息 
func MessageSend(Name string,conn net.Conn){
	defer conn.Close()
	var input string
	for {
		reader:=bufio.NewReader(os.Stdin)
		data,_,_:=reader.ReadLine()
		// 获得用户输入信息
		input = string(data)

		// 判断用户输入的是否是EXIT退出指令
		if input == "EXIT"{
			setStatus(Name)
			// 还需要将用户的在线状态置位为0，传入数据库
			break
		}

		// 将信息传入结构体，得到用户发送的信息
		u := ChatInfo{
			UserName: Name,
			ChatData: input,
		}
		// 将结构体转换成json格式
		jsonData, err := json.Marshal(u)
		if err != nil {
			log.Fatalf("Json marshaling failed：%s", err)
			setStatus(Name)
			// 转换失败即发送失败，需要更改用户状态
			break
		}
		_, err = conn.Write(jsonData) // 发送数据
		if err != nil {
			setStatus(Name)
			// 发送失败，退出客户端并且更改用户状态
			break
		}

		// 接受数据
		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("recv failed, err:", err)
			setStatus(Name)
			// 接受失败，退出客户端并且更改用户状态
			break
		}
		fmt.Println(string(buf[:n]))
	}
	wg.Done()
}

// 客户端
func main() {
	err := initDB() // 调用输出化数据库的函数
	if err != nil {
		fmt.Printf("init db failed,err:%v\n", err)
		return
	}
	defer db.Close()
	wg.Add(1)
	
	fmt.Println("--------------------------------------")
	fmt.Println("1.登录请输入1")
	fmt.Println("2.注册请输入2")
	fmt.Println("--------------------------------------")
	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')
	inputInfo := strings.Trim(input, "\r\n")
	var username string
	switch inputInfo {
	case "1":
		fmt.Println("请输入登录用户名:")
		input, _ = inputReader.ReadString('\n')
		username = strings.Trim(input, "\r\n")
		login(username)
	case "2":
		fmt.Println("请输入注册用户名:")
		input, _ = inputReader.ReadString('\n')
		username = strings.Trim(input, "\r\n")
		register(username)
		login(username)
	}
	fmt.Println(inputInfo)

	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer conn.Close() // 关闭连接
	go MessageSend(username,conn)
	wg.Wait()
}
