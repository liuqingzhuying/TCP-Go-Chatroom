package main

// 定义用户的结构体
// 使用mysql存储用户信息（注册，登录）（表的设计    1.用户id，用户名，密码，发送信息总数量  2.用户id，日期，当天发送信息总数量  3.用户id，是否在线）
// 单向链表对用户进行排序，依据用户发送的消息数量进行降序排列
// 聊天室可以查到每天，发出消息最多的用户
// 聊天室服务器可以查看在线用户，离线用户，总用户

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// User 用户结构体
type User struct {
	UserName string
	ChatNum  int
}

type userList []*User

// 排序规则重写
func (s userList) Len() int {
	return len(s)
}
func (s userList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s userList) Less(i, j int) bool {
	return s[i].ChatNum > s[j].ChatNum
}

func getAllChatData() userList {
	// 得到所有用户的所有聊天记录并排序
	var uList userList
	sql := "select user_name, chat_num from user_info"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("scan err", err)
	}
	for rows.Next() {
		var u User
		err := rows.Scan(&u.UserName, &u.ChatNum)
		if err != nil {
			fmt.Println("scan err", err)
			return uList
		}
		uList = append(uList, &u)
	}
	sort.Sort(uList)
	return uList
}

// 定义一个全局对象db
var db *sql.DB

func getDate() string {
	// 得到当前日期
	return time.Now().Format("2006/01/02")
}

// 定义一个初始化数据库的函数
func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:123456@tcp(127.0.0.1:3306)/chat_test?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

// Status 用户状态结构体
type Status struct {
	UserName   string
	UserStatus int
}

type statusList []*Status

// 得到所有的在线用户，离线用户，和总用户
func getUser() (statusList, statusList, statusList) {
	var allUser statusList
	var online statusList
	var offline statusList
	sql := "select * from user_status"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("scan err", err)
	}
	for rows.Next() {
		var status Status
		err := rows.Scan(&status.UserName, &status.UserStatus)
		if err != nil {
			fmt.Println("scan err", err)
			return allUser, online, offline
		}
		if status.UserStatus == 1 {
			online = append(online, &status)
		} else if status.UserStatus == 0 {
			offline = append(offline, &status)
		}
		allUser = append(allUser, &status)
	}
	return allUser, online, offline
}
