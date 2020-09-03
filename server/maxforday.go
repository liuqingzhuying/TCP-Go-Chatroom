package main

import (
	"fmt"
	"sort"
)

// DateChat 用户每天聊天结构体
type DateChat struct {
	UserName string
	ChatNum  int
	Date     string
}

type dateList []*DateChat

func (s dateList) Len() int {
	return len(s)
}
func (s dateList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s dateList) Less(i, j int) bool {
	return s[i].ChatNum > s[j].ChatNum
}

func getDayMax() *DateChat {
	var dList dateList
	date := getDate()
	sql := "select user_name, chat_num,date from chat_date where date = ?"
	rows, err := db.Query(sql, date)
	if err != nil {
		fmt.Println("scan err", err)
	}
	for rows.Next() {
		var d DateChat
		err := rows.Scan(&d.UserName, &d.ChatNum, &d.Date)
		if err != nil {
			fmt.Println("scan err", err)
			return &DateChat{}
		}
		dList = append(dList, &d)
	}
	sort.Sort(dList)
	return dList[0]
}
