package controllers

import (
	models "Pages/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// ########################### bachelor
type BachelorsIdController struct {
	beego.Controller
}

type BachelorsNameController struct {
	beego.Controller
}

func (this *BachelorsNameController) Get() {
	this.Data["Path"] = this.Ctx.Request.RequestURI
	this.TplName = "b_index.html"
}

// list: search item by name
func (this *BachelorsNameController) Post() {
	// this.TplName = "search.html"
	inputName := this.GetString("name")

	// 防止输入为空
	if inputName != "" {
		// inputName := this.Ctx.Input.Param(":name")
		res := models.GetBachelorsListByName(inputName)
		this.Data["list"] = res.Info
		this.TplName = "b_search.html"
		// this.Data["json"] = &res
		// this.ServeJSON()

	} else {
		this.TplName = "404.html"
	}

}

// detail: search detail by id
func (this *BachelorsIdController) Get() {
	input := this.Ctx.Input.Param(":id")
	inputId, err := strconv.Atoi(input)
	checkError(err)
	res := models.GetBachelorsDetailById(inputId)

	this.Data["info"] = res.Info
	this.Data["course"] = res.Course
	this.Data["detail"] = res.Comment
	this.TplName = "b_result.html"

	// this.Data["json"] = &res
	// this.ServeJSON()
}

// comment: insert item 处理教师新评论
func (this *BachelorsIdController) Post() {
	// 教师id
	inputId := this.Input().Get("id")

	Id, err1 := strconv.Atoi(inputId)
	if err1 != nil {
		logs.Error(err1)
		this.TplName = "fail.html"
		return
	}

	// 评分
	inputScore := this.Input().Get("score")
	Score, err2 := strconv.Atoi(inputScore)
	if err2 != nil {
		logs.Error(err1)
		this.TplName = "fail.html"
		return
	}

	// 具体评论内容
	Comment := this.GetString("comment")
	if len(Comment) < 15 {
		this.TplName = "fail.html"
		return
	}

	// 插入评论
	res := models.InsertCommentToBachelors(Id, Score, Comment)
	// 数据返回
	// this.Data["json"] = &res
	// this.ServeJSON()
	this.Data["result"] = res.Msg
	this.TplName = "b_comment.html" //TODO: 设计该html文件

}

// update: 评论点赞功能
func (this *BachelorsIdController) Put() {
	// 教师id
	inputId := this.Input().Get("id")
	Id, err1 := strconv.Atoi(inputId)

	// 评论编号
	inputNum := this.Input().Get("num")
	num, err2 := strconv.Atoi(inputNum)

	// 赞or踩？
	isThumb := this.Input().Get("thumb")
	Thumb, err3 := strconv.Atoi(isThumb)

	if err1 != nil && err2 != nil && err3 != nil {
		// 数据库处理
		res := models.UpdateThumbForBachelors(Id, num, Thumb)

		// 数据返回
		// this.Data["json"] = &res
		// this.ServeJSON()
		this.Data["result"] = res.Msg
		this.TplName = "b_update.html" //TODO: 设计该html文件
	} else {
		logs.Error(err1, err2, err3)
	}

}
