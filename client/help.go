package main

// 定义用户的结构体
// 使用mysql存储用户信息（注册，登录）（表的设计    1.用户id，用户名，发送信息总数量  2.用户id，日期，当天发送信息总数量  3.用户id，是否在线）
// 单向链表对用户进行排序，依据用户发送的消息数量进行降序排列
// 聊天室可以查到每天，发出消息最多的用户
// 聊天室服务器可以查看在线用户，离线用户，总用户

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// User 用户结构体
type User struct {
	UserName string
	ChatNum  int
}

// 定义一个全局对象db
var db *sql.DB

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

func getDate() string {
	return time.Now().Format("2006/01/02")
}

// 登录
func login(uname string) {
	// 查询数据库中的用户表，是否有该用户
	sqlStr := "select user_name, chat_num from user_info where user_name=?"
	var u User
	err := db.QueryRow(sqlStr, uname).Scan(&u.UserName, &u.ChatNum)
	if err != nil {
		fmt.Println("查无此人, 登录失败")
		os.Exit(1)
	}

	// 查看是否有使用者当天的记录，没有就添加
	sqlStr = "select user_name, chat_num from chat_date where user_name=? and date = ?"
	date := getDate()
	err = db.QueryRow(sqlStr, uname, date).Scan(&u.UserName, &u.ChatNum)
	if err != nil {
		sqlStr = "insert into chat_date(user_name, chat_num, date) values(?,?,?)"
		_, err1 := db.Exec(sqlStr, uname, 0, date)
		if err1 != nil {
			fmt.Printf("insert failed, err:%v\n", err)
			os.Exit(1)
		}
	}

	// 设置登录状态
	sqlStr = "update user_status set status=1 where user_name=?"
	_, err = db.Exec(sqlStr, uname)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("UserName:%s ChatNum:%d\n", u.UserName, u.ChatNum)
}

// 注册
func register(uname string) {
	// 查看用户名是否存在，存在则失败，不存在则创建用户到数据库
	// 查询数据库中的用户表，是否有该用户
	sqlStr := "select user_name, chat_num from user_info where user_name=?"
	var u User
	err := db.QueryRow(sqlStr, uname).Scan(&u.UserName, &u.ChatNum)
	if err != nil {
		sqlStr = "insert into user_info(user_name, chat_num) values (?,?)"
		_, err := db.Exec(sqlStr, uname, 0)
		if err != nil {
			fmt.Printf("insert failed, err:%v\n", err)
			return
		}
		// 状态表插入新用户
		sqlStr = "insert into user_status(user_name, status) values (?,0)"
		_, _ = db.Exec(sqlStr, uname)
		fmt.Println("注册成功")
	} else {
		fmt.Println("该用户已存在")
		os.Exit(1)
	}
}
