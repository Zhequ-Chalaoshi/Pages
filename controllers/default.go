package controllers

import (
	models "Pages/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func checkError(err error) {
	if err != nil {
		logs.Error(err)
	}
}

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["Path"] = this.Ctx.Request.RequestURI
	this.TplName = "index.html"
}

func (this *MainController) AddT() {
	this.TplName = "addt.html"
}

func (c *MainController) AddTPost() {
	Comment := c.GetString("content")
	if len(Comment) < 15 {
		c.TplName = "fail.html"
		return
	}

	res := models.WriteToFile(Comment)
	if res != true {
		c.TplName = "fail.html"
	} else {
		c.TplName = "addTsuccess.html"
	}

}

// ########################################################################################## 页面错误处理
type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error401() {
	c.TplName = "401.html"
}

func (c *ErrorController) Error403() {
	c.TplName = "403.html"
}

func (c *ErrorController) Error404() {
	c.TplName = "404.html"
}

func (c *ErrorController) Error500() {
	c.TplName = "500.html"
}

func (c *ErrorController) Error503() {
	c.TplName = "503.html"
}
