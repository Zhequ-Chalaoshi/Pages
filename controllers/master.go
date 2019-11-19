package controllers

import (
	models "Pages/models"
	"strconv"

	"github.com/astaxie/beego"
)

// ########################### master
type MasterIdController struct {
	beego.Controller
}

type MasterNameController struct {
	beego.Controller
}

func (c *MasterNameController) Get() {
	c.Data["Path"] = c.Ctx.Request.RequestURI
	c.TplName = "m_index.html"
}

// list: search item by name
func (this *MasterNameController) Post() {
	// this.TplName = "search.html"
	inputName := this.GetString("name")

	// 防止输入为空
	if inputName != "" {
		// inputName := this.Ctx.Input.Param(":name")
		res := models.GetMasterListByName(inputName)

		this.Data["list"] = res.Info
		this.TplName = "m_search.html"

		// this.Data["json"] = &res
		// this.ServeJSON()
	} else {
		this.TplName = "404.html"
	}

}

// detail: search detail by id
func (this *MasterIdController) Get() {
	input := this.Ctx.Input.Param(":id")
	inputId, err := strconv.Atoi(input)

	checkError(err)
	res := models.GetMasterDetailById(inputId)

	this.Data["info"] = res.Info
	this.Data["detail"] = res.Comment
	this.TplName = "m_result.html"

	// this.Data["json"] = &res
	// this.ServeJSON()
}

//  comment: insert item 处理教师新评论
func (this *MasterIdController) Post() {
	// 教师id
	inputId := this.Input().Get("id")
	Id, err1 := strconv.Atoi(inputId)
	checkError(err1)
	if err1 != nil {
		this.TplName = "fail.html"
		return
	}

	// 评分
	var Score [4]int
	inputScore := this.Input().Get("score1")
	t_score, err2 := strconv.Atoi(inputScore)
	checkError(err2)
	if err2 != nil {
		this.TplName = "fail.html"
		return
	}
	Score[0] = t_score

	inputScore = this.Input().Get("score2")
	t_score, err2 = strconv.Atoi(inputScore)
	checkError(err2)
	if err2 != nil {
		this.TplName = "fail.html"
		return
	}
	Score[1] = t_score

	inputScore = this.Input().Get("score3")
	t_score, err2 = strconv.Atoi(inputScore)
	checkError(err2)
	if err2 != nil {
		this.TplName = "fail.html"
		return
	}
	Score[2] = t_score

	inputScore = this.Input().Get("score4")
	t_score, err2 = strconv.Atoi(inputScore)
	checkError(err2)
	if err2 != nil {
		this.TplName = "fail.html"
		return
	}
	Score[3] = t_score

	// 具体评论内容
	Comment := this.GetString("comment")
	if len(Comment) < 15 {
		this.TplName = "fail.html"
		return
	}

	// 插入评论
	res := models.InsertCommentToMaster(Id, Score, Comment)

	// 数据返回
	// this.Data["json"] = &res
	// this.ServeJSON()
	this.Data["result"] = res.Msg
	this.TplName = "m_comment.html" //TODO: 设计该html文件

}

// update: 评论点赞功能
func (this *MasterIdController) Put() {
	// 教师id
	inputId := this.Input().Get("id")
	Id, err1 := strconv.Atoi(inputId)
	checkError(err1)
	if err1 != nil {
		return
	}

	// 评论编号
	inputNum := this.Input().Get("num")
	num, err2 := strconv.Atoi(inputNum)
	checkError(err2)
	if err2 != nil {
		return
	}

	// 赞or踩？
	isThumb := this.Input().Get("thumb")
	Thumb, err3 := strconv.Atoi(isThumb)
	checkError(err3)
	if err3 != nil {
		return
	}

	// 数据库处理
	res := models.UpdateThumbForMaster(Id, num, Thumb)

	// 数据返回
	// this.Data["json"] = &res
	// this.ServeJSON()
	this.Data["result"] = res.Msg
	this.TplName = "m_update.html" //TODO: 设计该html文件
}
