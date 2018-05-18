package routers

import (
	"github.com/astaxie/beego"
	"loveHome/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetAreaInfo")
	beego.Router("/api/v1.0/session", &controllers.SessionController{},
		"get:GetSessionName;delete:DelSessionName")
	//登陆
	/*
		method: POST
				api/v1.0/sessions
	*/
	beego.Router("/api/v1.0/sessions", &controllers.UserController{}, "post:Login")
	beego.Router("/api/v1.0/houses/index", &controllers.HousesIndexController{},
		"get:GetHousesIndex")
	beego.Router("/api/v1.0/users", &controllers.UserController{}, "post:Reg")
	//文件上传业务
	beego.Router("api/v1.0/user/avatar", &controllers.UserController{}, "post:UploadAvatar")

}
