package routers

import (
	"Pages/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// ------------------------------------------------------------------- 错误页面跳转
	beego.ErrorController(&controllers.ErrorController{})

	// ------------------------------------------------------------------- 主页
	beego.Router("/", &controllers.MainController{}, "*:Get")

	// ------------------------------------------------------------------- 本科生
	beego.Router("/bachelors", &controllers.BachelorsNameController{}, "get:Get")
	// 按名字及学院搜索
	beego.Router("/bachelorslist", &controllers.BachelorsNameController{}, "post:Post")
	// 查看教师具体内容
	beego.Router("/bachelorsdetail/:id:int", &controllers.BachelorsIdController{}, "get:Get")
	// 添加教师评分和评论
	beego.Router("/bachelorsdetail", &controllers.BachelorsIdController{}, "post:Post")
	// 点赞功能
	// beego.Router("/bachelorsdetail", &controllers.BachelorsIdController{}, "put:Put")

	// ------------------------------------------------------------------- 研究生
	beego.Router("/master", &controllers.MasterNameController{}, "get:Get")
	// 按名字及学院搜索
	beego.Router("/masterlist", &controllers.MasterNameController{}, "post:Post")
	// 查看教师具体内容
	beego.Router("/masterdetail/:id:int", &controllers.MasterIdController{}, "get:Get")
	// 添加教师评分和评论
	beego.Router("/masterdetail", &controllers.MasterIdController{}, "post:Post")
	// 点赞功能
	// beego.Router("/masterdetail", &controllers.MasterIdController{}, "put:Put")

	// ------------------------------------------------------------------- 添加反馈
	beego.Router("/addt", &controllers.MainController{}, "get:AddT")
	beego.Router("/addt", &controllers.MainController{}, "post:AddTPost")

}
