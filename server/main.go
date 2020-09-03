package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
)

// var wg sync.WaitGroup

// ChatInfo 聊天信息结构体
type ChatInfo struct {
	UserName string
	ChatData string
}

// 处理函数
func process(conn net.Conn) {
	defer conn.Close() // 关闭连接
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}

		// 对用户发送的信息进行json解码，得到传输的结构体
		var chatInfo ChatInfo
		// movies2 := make([]Movie, 10)
		if err = json.Unmarshal(buf[:n], &chatInfo); err != nil {
			log.Fatalf("JSON unmarshling failed: %s", err)
		}
		// recvStr := string(buf[:n])
		fmt.Println(chatInfo.UserName, " : ", chatInfo.ChatData)

		// 用户发送的总信息数量加一
		sql := "update user_info set chat_num = chat_num+1 where user_name=?"
		db.Exec(sql, chatInfo.UserName)

		// 当天用户发送的总信息数量加一
		date := getDate()
		sql = "update chat_date set chat_num = chat_num+1 where user_name=? and date = ?"
		_, err = db.Exec(sql, chatInfo.UserName, date)
		if err != nil {
			fmt.Println("err :", err)
		}

		data, _ := json.Marshal(chatInfo)
		conn.Write(data) // 发送数据
	}
}

func server() {
	err := initDB() // 调用输出化数据库的函数
	if err != nil {
		fmt.Printf("init db failed,err:%v\n", err)
		return
	}
	defer db.Close()
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn) // 启动一个goroutine处理连接
	}
}

func main() {
	go server()
	// 创建一个默认的路由引擎
	r := gin.Default()
	// 当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
	r.GET("/hello", func(c *gin.Context) {
		uList := getAllChatData()
		maxUser := getDayMax()
		allUser, online, offline := getUser()
		// c.JSON：返回JSON格式的数据
		c.JSON(200, gin.H{
			"sort_by_send_info": uList,
			"today_max_send":    maxUser,
			"allUser":           allUser,
			"onlineUser":        online,
			"offlineUser":       offline,
		})
	})
	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run("0.0.0.0:8080")
}
