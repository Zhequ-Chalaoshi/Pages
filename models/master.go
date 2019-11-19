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

// MasterTInfo Model Struct
// 研究生教师信息表
type MasterTInfo struct {
	Id        int    `orm:"pk"`
	Name      string `orm:"size(40)"`
	College   string `orm:"size(40)"`
	Researchs string `orm:"size(400)"`
	Email     string `orm:"size(100)"`
	Content   string `orm:"size(2000)"`
	Rate      float32
	Rate1     float32
	Rate2     float32
	Rate3     float32
	Rate4     float32
}

type MasterTComment struct {
	Id      int `orm:"pk"`
	Num     int
	Content string `orm:"size(3000)"`
	Thumb   int
	Time    string `orm:"size(40)"`
}

type ReturnMasterData struct {
	Status  int              `json:"status"`
	Msg     string           `json:"msg"`
	Info    []MasterTInfo    `json:"info"`
	Comment []MasterTComment `json:"comment"`
}

// GetMasterListByName 按姓名查找
func GetMasterListByName(inputName string) ReturnMasterData {
	o := orm.NewOrm()
	info := o.QueryTable("MasterTInfo")
	var tmp_info []MasterTInfo

	// n1, err1 := info.Filter("name__contains", inputName).All(&tmp_info)

	cond := orm.NewCondition()
	cond1 := cond.And("name__contains", inputName).Or("college__contains", inputName)
	info = info.SetCond(cond1)
	n1, err1 := info.OrderBy("college", "-rate").All(&tmp_info)

	var res ReturnMasterData
	if err1 == nil {
		res.Status = 200
		res.Msg = "Success to search item!"
		res.Info = tmp_info
		res.Comment = nil
	} else {
		res.Status = 500
		res.Msg = "Fail to search item. " + err1.Error()
		res.Info = nil
		res.Comment = nil
	}
	if n1 == 0 {
		res.Status = 500
		res.Msg = "No such item"
		res.Info = nil
		res.Comment = nil
	}

	return res
}

// GetMasterDetailById 查找具体内容
func GetMasterDetailById(inputId int) ReturnMasterData {

	// o := orm.NewOrm()

	// info := o.QueryTable("BachelorsTInfo")
	// var tmp_info []BachelorsTInfo
	// n1, err1 := info.Filter("id", inputId).All(&tmp_info)
	tmp_info, n1, err1 := GetMasterInfoById(inputId)

	// comment := o.QueryTable("BachelorsTComment")
	// var tmp_comment []BachelorsTComment
	// n3, err3 := comment.Filter("id", inputId).All(&tmp_comment)
	tmp_comment, n3, err3 := GetMasterCommentById(inputId)

	var res ReturnMasterData
	if err1 == nil && err3 == nil {
		res.Status = 200
		res.Msg = "Success to search item!"
		res.Info = tmp_info
		res.Comment = tmp_comment
	} else {
		res.Status = 500
		res.Msg = "Fail to search item. " + err1.Error() + err3.Error()
		res.Info = nil
		res.Comment = nil
	}
	if n1 == 0 {
		res.Info = nil
	}
	if n3 == 0 {
		res.Comment = nil
	}
	if n1 == 0 && n3 == 0 {
		res.Status = 500
		res.Msg = "No such item"
		res.Info = nil
		res.Comment = nil
	}

	return res
}

// GetMasterDetailById 添加评论数据
func InsertCommentToMaster(Id int, Score [4]int, Comment string) ReturnMasterData {
	o := orm.NewOrm()
	var res ReturnMasterData

	// 新增评分处理
	tmp_info, _, err1 := GetMasterInfoById(Id)

	// 获取所有评论数
	tmp_c, total_comment_num, err3 := GetMasterCommentById(Id)

	// 计算分数
	tmp_info[0].Rate1 = (float32(total_comment_num)*tmp_info[0].Rate1 + float32(Score[0])) / float32(total_comment_num+1)
	tmp_info[0].Rate2 = (float32(total_comment_num)*tmp_info[0].Rate2 + float32(Score[1])) / float32(total_comment_num+1)
	tmp_info[0].Rate3 = (float32(total_comment_num)*tmp_info[0].Rate3 + float32(Score[2])) / float32(total_comment_num+1)
	tmp_info[0].Rate4 = (float32(total_comment_num)*tmp_info[0].Rate4 + float32(Score[3])) / float32(total_comment_num+1)

	tmp_info[0].Rate = (tmp_info[0].Rate1 + tmp_info[0].Rate2 + tmp_info[0].Rate3 + tmp_info[0].Rate4) / 4.0

	//	更改数据
	_, err2 := o.Update(&tmp_info[0])

	// 评论信息
	var tmp_comment MasterTComment

	t := time.Now()
	tmp_comment.Id = Id
	tmp_comment.Num = int(tmp_c[len(tmp_c)-1].Num + 1)
	tmp_comment.Content = Comment
	tmp_comment.Thumb = 0
	tmp_comment.Time = fmt.Sprintf("%02d.%02d", t.Year(), t.Month())
	// 插入
	_, err4 := o.Insert(&tmp_comment)

	writeString := "master comments add: " + strconv.Itoa(Id) + "\n" + strconv.Itoa(Score[0]) + "\n" + strconv.Itoa(Score[1]) + "\n" + strconv.Itoa(Score[2]) + "\n" + strconv.Itoa(Score[3]) + "\n" + Comment

	// 存储评论数据到文件中
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
	res.Comment = nil
	return res
}

// UpdateThumbForMaster 点赞功能
func UpdateThumbForMaster(inputId int, num int, thumb int) ReturnMasterData {
	o := orm.NewOrm()
	var res ReturnMasterData

	// 点赞
	var err error
	if thumb == 1 {
		_, err = o.QueryTable("MasterTComment").Filter("id", inputId).Filter("num", num).Update(orm.Params{
			"thumb": orm.ColValue(orm.ColAdd, 1),
		})
	} else {
		_, err = o.QueryTable("MasterTComment").Filter("id", inputId).Filter("num", num).Update(orm.Params{
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
	res.Comment = nil
	return res
}

func GetMasterInfoById(inputId int) ([]MasterTInfo, int64, error) {
	o := orm.NewOrm()
	info := o.QueryTable("MasterTInfo")
	var res []MasterTInfo
	n1, err1 := info.Filter("id", inputId).All(&res)
	return res, n1, err1
}

func GetMasterCommentById(inputId int) ([]MasterTComment, int64, error) {
	o := orm.NewOrm()
	info := o.QueryTable("MasterTComment")
	var res []MasterTComment
	n1, err1 := info.Filter("id", inputId).OrderBy("-thumb").All(&res)
	return res, n1, err1
}
