package models

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used sql driver
)

func init() {
	// 连接数据库
	maxIdle := 500
	maxConn := 4000
	err := orm.RegisterDataBase("default", "mysql", "Username:Password@tcp(ipAddr:3306)/Pages?charset=utf8", maxIdle, maxConn)
	// 参数1        数据库的别名，用来在 ORM 中切换数据库使用
	// 参数2        mysql
	// 参数3        对应的链接字符串 "账号:密码@tcp(ip:端口)/数据库"
	// 参数4(可选)  设置最大空闲连接
	// 参数5(可选)  设置最大数据库连接 (go >= 1.2)

	if err != nil {
		logs.Error("connect mysql err : ", err)
	} else {
		logs.Info("connect mysql success")
	}
	// register model 连接数据表
	orm.RegisterModel(new(BachelorsTInfo))
	orm.RegisterModel(new(BachelorsTCourses))
	orm.RegisterModel(new(BachelorsTComment))

	orm.RegisterModel(new(MasterTInfo))
	orm.RegisterModel(new(MasterTComment))

	// create table
	// orm.RunSyncdb("default", false, true)
	// orm.Debug = true
}

// WriteToFile 将反馈写入清单文件
func WriteToFile(filename string) bool {
	var filelist = "./logs/feedback.txt"

	t := time.Now()
	s := fmt.Sprintf("%02d.%02d.%2d  %2d:%2d:%2d\n",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	var write2String = "\n\n" + s + filename + "\n\n"
	var f *os.File
	var err1 error

	if checkFileIsExist(filelist) { //如果文件存在
		f, err1 = os.OpenFile(filelist, os.O_APPEND|os.O_WRONLY, os.ModeAppend) //打开文件
	} else {
		f, err1 = os.Create(filelist) //创建文件
		logs.Error("文件", filelist, "不存在！")
	}
	checkError(err1)

	_, err1 = io.WriteString(f, write2String) //写入文件(字符串)
	defer f.Close()

	checkError(err1)
	if err1 != nil {
		return false
	}
	return true
}

// 判断文件是否存在，存在返回true不存在返回false
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func checkError(err error) {
	if err != nil {
		logs.Error(err)
	}
}
