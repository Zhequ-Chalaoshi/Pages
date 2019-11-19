package models

// package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
)

// BachelorsTInfo Model Struct
// 本科生教师信息表
type BachelorsTInfo struct {
	Id       int    `orm:"pk"`
	College  string `orm:"size(100)"`
	ScoreNum int
	TName    string `orm:"size(100)"`
	AvrScore float32
}
type BachelorsTCourses struct {
	Id      int `orm:"pk"`
	Num     int
	Courses string `orm:"size(100)"`
	Scores  string `orm:"size(40)"`
}
type BachelorsTComment struct {
	Id      int `orm:"pk"`
	Num     int
	Content string `orm:"size(1000)"`
	Thumb   int
	Time    string `orm:"size(40)"`
}

type ReturnBachelorsData struct {
	Status  int                 `json:"status"`
	Msg     string              `json:"msg"`
	Info    []BachelorsTInfo    `json:"obj"`
	Course  []BachelorsTCourses `json:"course"`
	Comment []BachelorsTComment `json:"comment"`
}

// GetBachelorsListByName 查找列表
func GetBachelorsListByName(inputName string) ReturnBachelorsData {
	o := orm.NewOrm()

	info := o.QueryTable("BachelorsTInfo")
	var tmp_info []BachelorsTInfo

	// 按姓名、学院
	cond := orm.NewCondition()
	cond1 := cond.And("t_name__contains", inputName).Or("college__contains", inputName)
	info = info.SetCond(cond1)
	// 按学院，评分排序
	n1, err1 := info.OrderBy("college", "-avr_score").All(&tmp_info)

	// 按课程查找
	// course := o.QueryTable("BachelorsTCourses")
	// var tmp_course []BachelorsTCourses
	// n2, err2 := course.Filter("id", inputId).Distinct().All(&tmp_course)

	var res ReturnBachelorsData
	if err1 == nil {
		res.Status = 200
		res.Msg = "Success to search item!"
		res.Info = tmp_info
		res.Course = nil
		res.Comment = nil
	} else {
		res.Status = 500
		res.Msg = "Fail to search item. " + err1.Error()
		res.Info = nil
		res.Course = nil
		res.Comment = nil
	}
	if n1 == 0 {
		res.Status = 500
		res.Msg = "No such item"
		res.Info = nil
		res.Course = nil
		res.Comment = nil
	}

	return res
}

// GetBachelorsDetailById 查找具体内容
func GetBachelorsDetailById(inputId int) ReturnBachelorsData {

	// o := orm.NewOrm()

	// info := o.QueryTable("BachelorsTInfo")
	// var tmp_info []BachelorsTInfo
	// n1, err1 := info.Filter("id", inputId).All(&tmp_info)
	tmp_info, n1, err1 := GetBachInfoById(inputId)

	// course := o.QueryTable("BachelorsTCourses")
	// var tmp_course []BachelorsTCourses
	// n2, err2 := course.Filter("id", inputId).All(&tmp_course)
	tmp_course, n2, err2 := GetBachCourseById(inputId)

	// comment := o.QueryTable("BachelorsTComment")
	// var tmp_comment []BachelorsTComment
	// n3, err3 := comment.Filter("id", inputId).All(&tmp_comment)
	tmp_comment, n3, err3 := GetBachCommentById(inputId)
	var res ReturnBachelorsData
	if err1 == nil && err2 == nil && err3 == nil {
		res.Status = 200
		res.Msg = "Success to search item!"
		res.Info = tmp_info
		res.Course = tmp_course
		res.Comment = tmp_comment
	} else {
		res.Status = 500
		res.Msg = "Fail to search item. " + err1.Error() + err2.Error() + err3.Error()
		res.Info = nil
		res.Course = nil
		res.Comment = nil
	}
	if n1 == 0 {
		res.Info = nil
	}
	if n2 == 0 {
		res.Course = nil
	}
	if n3 == 0 {
		res.Comment = nil
	}
	if n1 == 0 && n2 == 0 && n3 == 0 {
		res.Status = 500
		res.Msg = "No such item"
		res.Info = nil
		res.Course = nil
		res.Comment = nil
	}

	return res
}

// InsertCommentToBachelors 添加评论数据
func InsertCommentToBachelors(Id int, Score int, Comment string) ReturnBachelorsData {
	o := orm.NewOrm()
	var res ReturnBachelorsData

	// 新增评分处理
	tmp_info, _, err1 := GetBachInfoById(Id)
	num := tmp_info[0].ScoreNum
	tmp_info[0].ScoreNum = num + 1
	tmp_info[0].AvrScore = (float32(num)*tmp_info[0].AvrScore + float32(Score)) / float32(num+1)

	//	更改数据
	_, err2 := o.Update(&tmp_info[0])

	// 评论信息
	var tmp_comment BachelorsTComment
	tmp_c, _, err3 := GetBachCommentById(Id)

	t := time.Now()
	tmp_comment.Id = Id

	tmp_comment.Num = int(tmp_c[len(tmp_c)-1].Num + 1) // 取最后一条评论的编号
	tmp_comment.Content = Comment
	tmp_comment.Thumb = 0
	tmp_comment.Time = fmt.Sprintf("%02d.%02d.%2d", t.Year(), t.Month(), t.Day())
	// 插入
	_, err4 := o.Insert(&tmp_comment)

	// 存储评论数据到文件中
	writeString := "bachelors comments add: " + strconv.Itoa(Id) + "\n" + strconv.Itoa(Score) + "\n" + Comment

	w := WriteToFile(writeString)
	if w == true {
		logs.Critical("OK to write to feedback.txt")
	}

	if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
		res.Status = 200
		res.Msg = "Success to add item!"
	} else {
		res.Status = 500
		res.Msg = "Fail to add item. " + err1.Error() + err2.Error() + err3.Error() + err4.Error()
	}
	res.Info = nil
	res.Course = nil
	res.Comment = nil
	return res
}

// UpdateThumbForBachelors 修改，点赞功能
func UpdateThumbForBachelors(inputId int, num int, thumb int) ReturnBachelorsData {
	o := orm.NewOrm()
	var res ReturnBachelorsData

	// 点赞
	var err error
	if thumb == 1 {
		_, err = o.QueryTable("BachelorsTComment").Filter("id", inputId).Filter("num", num).Update(orm.Params{
			"thumb": orm.ColValue(orm.ColAdd, 1),
		})
	} else { // 踩
		_, err = o.QueryTable("BachelorsTComment").Filter("id", inputId).Filter("num", num).Update(orm.Params{
			"thumb": orm.ColValue(orm.ColMinus, 1),
		})
	}

	if err == nil {
		res.Status = 200
		res.Msg = "Success to update item!"
	} else {
		res.Status = 500
		res.Msg = "Fail to update item. " + err.Error()
	}
	res.Info = nil
	res.Course = nil
	res.Comment = nil
	return res
}

func GetBachInfoById(inputId int) ([]BachelorsTInfo, int64, error) {
	o := orm.NewOrm()
	info := o.QueryTable("BachelorsTInfo")
	var res []BachelorsTInfo
	n1, err1 := info.Filter("id", inputId).All(&res)
	return res, n1, err1
}

func GetBachCourseById(inputId int) ([]BachelorsTCourses, int64, error) {
	o := orm.NewOrm()
	info := o.QueryTable("BachelorsTCourses")
	var res []BachelorsTCourses
	n1, err1 := info.Filter("id", inputId).All(&res)
	return res, n1, err1
}

func GetBachCommentById(inputId int) ([]BachelorsTComment, int64, error) {
	o := orm.NewOrm()
	info := o.QueryTable("BachelorsTComment")
	var res []BachelorsTComment
	n1, err1 := info.Filter("id", inputId).OrderBy("-thumb").All(&res)
	return res, n1, err1
}
